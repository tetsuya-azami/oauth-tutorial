package infrastructure

import (
	"oauth-tutorial/internal/domain"
	"oauth-tutorial/internal/logger"
	"testing"
)

func Test_認可リクエストパラメータの保存(t *testing.T) {
	logger := logger.NewMyLogger()

	// テスト用のAuthorizationCodeFlowParamを作成
	validParam, err := domain.NewAuthorizationCodeFlowParam(
		logger,
		"code",
		"test-client",
		"https://example.com/callback",
		"read write",
		"test-state",
	)
	if err != nil {
		t.Fatalf("Failed to create valid AuthorizationCodeFlowParam: %v", err)
	}

	tests := []struct {
		name        string
		sessionID   string
		param       *domain.AuthorizationCodeFlowParam
		expectedErr error
		setupFunc   func(*AuthParamSession)
		checkFunc   func(*testing.T, *AuthParamSession, string)
	}{
		{
			name:        "正常ケース - 新しいセッションの保存",
			sessionID:   "test-session-id",
			param:       validParam,
			expectedErr: nil,
			setupFunc: func(session *AuthParamSession) {
				sessionStore = make(map[string]domain.AuthorizationCodeFlowParam)
			},
			checkFunc: func(t *testing.T, session *AuthParamSession, sessionID string) {
				// 正しいキーで保存されていること
				if _, exists := sessionStore[sessionID]; !exists {
					t.Error("authParam should be saved with correct sessionID")
				}
				if len(sessionStore) != 1 {
					t.Errorf("sessionStore length = %d, want 1", len(sessionStore))
				}
			},
		},
		{
			name:        "正常ケース - 既存セッションの上書き",
			sessionID:   "existing-session",
			param:       validParam,
			expectedErr: nil,
			setupFunc: func(session *AuthParamSession) {
				sessionStore = make(map[string]domain.AuthorizationCodeFlowParam)
				oldParam, _ := domain.NewAuthorizationCodeFlowParam(
					logger,
					"code",
					"old-client",
					"https://old.com/callback",
					"read",
					"old-state",
				)
				sessionStore["existing-session"] = *oldParam
			},
			checkFunc: func(t *testing.T, session *AuthParamSession, sessionID string) {
				// 上書きされていること
				if _, exists := sessionStore[sessionID]; !exists {
					t.Error("authParam should be saved with correct sessionID")
				}
				// 新しい値で上書きされていること
				saved := sessionStore[sessionID]
				if saved.ClientID() != validParam.ClientID() {
					t.Errorf("saved ClientID = %s, want %s", saved.ClientID(), validParam.ClientID())
				}
			},
		},
		{
			name:        "異常ケース - 空のセッションID",
			sessionID:   "",
			param:       validParam,
			expectedErr: ErrInvalidParameter,
			setupFunc: func(session *AuthParamSession) {
				sessionStore = make(map[string]domain.AuthorizationCodeFlowParam)
			},
			checkFunc: func(t *testing.T, session *AuthParamSession, sessionID string) {
				if _, exists := sessionStore[sessionID]; exists {
					t.Error("sessionStore should not have entry for empty sessionID")
				}
				if len(sessionStore) != 0 {
					t.Errorf("sessionStore length = %d, want 0", len(sessionStore))
				}
			},
		},
		{
			name:        "異常ケース - nilパラメータ",
			sessionID:   "test-session",
			param:       nil,
			expectedErr: ErrInvalidParameter,
			setupFunc: func(session *AuthParamSession) {
				sessionStore = make(map[string]domain.AuthorizationCodeFlowParam)
			},
			checkFunc: func(t *testing.T, session *AuthParamSession, sessionID string) {
				if _, exists := sessionStore[sessionID]; exists {
					t.Error("sessionStore should not have entry for nil param")
				}
				if len(sessionStore) != 0 {
					t.Errorf("sessionStore length = %d, want 0", len(sessionStore))
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			session := NewAuthParamSession()

			// given
			if tt.setupFunc != nil {
				tt.setupFunc(session)
			}

			// when
			err := session.Save(tt.sessionID, tt.param)

			// then
			if tt.expectedErr != nil {
				if err == nil {
					t.Errorf("Save() error = nil, want %v", tt.expectedErr)
				} else if err != tt.expectedErr {
					t.Errorf("Save() error = %v, want %v", err, tt.expectedErr)
				}
			} else {
				if err != nil {
					t.Errorf("Save() error = %v, want nil", err)
				}
			}

			// 追加のチェック
			if tt.checkFunc != nil {
				tt.checkFunc(t, session, tt.sessionID)
			}
		})
	}
}

func Test_認可リクエストパラメータの取得(t *testing.T) {
	session := NewAuthParamSession()
	logger := logger.NewMyLogger()

	// sessionStoreを初期化（他のテストの影響を避けるため）
	sessionStore = make(map[string]domain.AuthorizationCodeFlowParam)

	// テスト用のAuthorizationCodeFlowParamを作成・保存
	param, err := domain.NewAuthorizationCodeFlowParam(
		logger,
		"code",
		"test-client",
		"https://example.com/callback",
		"read write",
		"test-state",
	)
	if err != nil {
		t.Fatalf("Failed to create AuthorizationCodeFlowParam: %v", err)
	}

	sessionID := "test-session-id"
	err = session.Save(sessionID, param)
	if err != nil {
		t.Fatalf("Save() error = %v", err)
	}

	tests := []struct {
		name      string
		sessionID string
		wantErr   bool
		expectNil bool
	}{
		{
			name:      "正常系 - 既存セッションの取得",
			sessionID: sessionID,
			wantErr:   false,
			expectNil: false,
		},
		{
			name:      "異常系 - 存在しないセッションの取得",
			sessionID: "non-existing-session",
			wantErr:   true,
			expectNil: true,
		},
		{
			name:      "異常系 - 空のセッションIDでの取得",
			sessionID: "",
			wantErr:   true,
			expectNil: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := session.Get(tt.sessionID)

			if tt.wantErr {
				if err == nil {
					t.Errorf("Get() error = nil, wantErr %v", tt.wantErr)
				}
				if err != ErrSessionNotFound {
					t.Errorf("Get() error = %v, want %v", err, ErrSessionNotFound)
				}
			} else {
				if err != nil {
					t.Errorf("Get() error = %v, wantErr %v", err, tt.wantErr)
				}
			}

			if tt.expectNil {
				if result != nil {
					t.Errorf("Get() result = %v, want nil", result)
				}
			} else {
				if result == nil {
					t.Error("Get() result should not be nil")
					return
				}
				// 値が正しく取得できることを確認
				if result.ClientID() != param.ClientID() {
					t.Errorf("Get() result.ClientID() = %v, want %v", result.ClientID(), param.ClientID())
				}
				if result.RedirectURI() != param.RedirectURI() {
					t.Errorf("Get() result.RedirectURI() = %v, want %v", result.RedirectURI(), param.RedirectURI())
				}
				if result.State() != param.State() {
					t.Errorf("Get() result.State() = %v, want %v", result.State(), param.State())
				}
			}
		})
	}
}

func Test_認可リクエストパラメータの削除(t *testing.T) {
	logger := logger.NewMyLogger()

	tests := []struct {
		name      string
		setupFunc func(*AuthParamSession)
		sessionID string
	}{
		{
			name: "正常系 - 既存セッションの削除",
			setupFunc: func(session *AuthParamSession) {
				sessionStore = make(map[string]domain.AuthorizationCodeFlowParam)
				param, err := domain.NewAuthorizationCodeFlowParam(
					logger,
					"code",
					"test-client",
					"https://example.com/callback",
					"read write",
					"test-state",
				)
				if err != nil {
					t.Fatalf("Failed to create AuthorizationCodeFlowParam: %v", err)
				}
				sessionStore["test-session-id"] = *param
			},
			sessionID: "test-session-id",
		},
		{
			name: "異常系 - 存在しないセッションの削除",
			setupFunc: func(session *AuthParamSession) {
				sessionStore = make(map[string]domain.AuthorizationCodeFlowParam)
			},
			sessionID: "non-existing-session",
		},
		{
			name: "異常系 - 空のセッションIDの削除",
			setupFunc: func(session *AuthParamSession) {
				sessionStore = make(map[string]domain.AuthorizationCodeFlowParam)
			},
			sessionID: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			session := NewAuthParamSession()
			tt.setupFunc(session)

			session.Clear(tt.sessionID)
			afterLen := len(sessionStore)

			if afterLen != 0 {
				t.Errorf("after Clear: sessionStore length = %d, want 0", afterLen)
			}
		})
	}
}

func TestAuthParamSession_MultipleSession(t *testing.T) {
	session := NewAuthParamSession()
	logger := logger.NewMyLogger()

	// sessionStoreを初期化
	sessionStore = make(map[string]domain.AuthorizationCodeFlowParam)

	// 複数のセッションを保存
	sessions := []struct {
		sessionID string
		clientID  string
		state     string
	}{
		{"session1", "client1", "state1"},
		{"session2", "client2", "state2"},
		{"session3", "client3", "state3"},
	}

	params := make([]*domain.AuthorizationCodeFlowParam, len(sessions))

	// 複数のセッションを保存
	for i, s := range sessions {
		param, err := domain.NewAuthorizationCodeFlowParam(
			logger,
			"code",
			s.clientID,
			"https://example.com/callback",
			"read write",
			s.state,
		)
		if err != nil {
			t.Fatalf("Failed to create AuthorizationCodeFlowParam: %v", err)
		}

		err = session.Save(s.sessionID, param)
		if err != nil {
			t.Fatalf("Save() error = %v", err)
		}

		params[i] = param
	}

	// 全てのセッションが正しく取得できることを確認
	for i, s := range sessions {
		result, err := session.Get(s.sessionID)
		if err != nil {
			t.Errorf("Get() for session %s error = %v", s.sessionID, err)
		}
		if result.ClientID() != params[i].ClientID() {
			t.Errorf("Get() result.ClientID() = %v, want %v", result.ClientID(), params[i].ClientID())
		}
		if result.State() != params[i].State() {
			t.Errorf("Get() result.State() = %v, want %v", result.State(), params[i].State())
		}
	}
}
