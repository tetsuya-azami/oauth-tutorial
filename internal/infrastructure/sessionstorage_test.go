package infrastructure

import (
	"oauth-tutorial/internal/domain"
	"oauth-tutorial/internal/infrastructure/dto"
	"oauth-tutorial/internal/session"
	"oauth-tutorial/pkg/mylogger"
	"testing"
)

func Test_認可リクエストパラメータの保存(t *testing.T) {
	logger := mylogger.NewLogger()

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
		sessionID   session.SessionID
		sessiondata *dto.SessionData
		expectedErr error
		setupFunc   func(*SessionStorage)
		checkFunc   func(*testing.T, *SessionStorage, session.SessionID)
	}{
		{
			name:        "正常ケース - 新しいセッションの保存",
			sessionID:   session.SessionID("test-session-id"),
			sessiondata: dto.NewSessionData(validParam, nil),
			expectedErr: nil,
			setupFunc: func(ss *SessionStorage) {
				sessionStore = make(map[session.SessionID]dto.SessionData)
			},
			checkFunc: func(t *testing.T, ss *SessionStorage, sessionID session.SessionID) {
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
			sessiondata: dto.NewSessionData(validParam, nil),
			expectedErr: nil,
			setupFunc: func(ss *SessionStorage) {
				sessionStore = make(map[session.SessionID]dto.SessionData)
				oldParam, _ := domain.NewAuthorizationCodeFlowParam(
					logger,
					"code",
					"old-client",
					"https://old.com/callback",
					"read",
					"old-state",
				)
				sessionStore[session.SessionID("existing-session")] = *dto.NewSessionData(oldParam, nil)
			},
			checkFunc: func(t *testing.T, ss *SessionStorage, sessionID session.SessionID) {
				// 上書きされていること
				if _, exists := sessionStore[sessionID]; !exists {
					t.Error("authParam should be saved with correct sessionID")
				}
				// 新しい値で上書きされていること
				saved := sessionStore[sessionID]
				if saved.AuthParam().ClientID() != validParam.ClientID() {
					t.Errorf("saved ClientID = %s, want %s", saved.AuthParam().ClientID(), validParam.ClientID())
				}
			},
		},
		{
			name:        "異常ケース - 空のセッションID",
			sessionID:   session.SessionID(""),
			sessiondata: dto.NewSessionData(validParam, nil),
			expectedErr: ErrInvalidSessionID,
			setupFunc: func(ss *SessionStorage) {
				sessionStore = make(map[session.SessionID]dto.SessionData)
			},
			checkFunc: func(t *testing.T, ss *SessionStorage, sessionID session.SessionID) {
				if _, exists := sessionStore[sessionID]; exists {
					t.Error("sessionStore should not have entry for empty sessionID")
				}
				if len(sessionStore) != 0 {
					t.Errorf("sessionStore length = %d, want 0", len(sessionStore))
				}
			},
		},
		{
			name:        "異常ケース - 空のセッションデータ",
			sessionID:   session.SessionID("test-session-id"),
			sessiondata: nil,
			expectedErr: ErrInvalidSessionData,
			setupFunc: func(ss *SessionStorage) {
				sessionStore = make(map[session.SessionID]dto.SessionData)
			},
			checkFunc: func(t *testing.T, ss *SessionStorage, sessionID session.SessionID) {
				if _, exists := sessionStore[sessionID]; exists {
					t.Error("sessionStore should not have entry for empty sessionID")
				}
				if len(sessionStore) != 0 {
					t.Errorf("sessionStore length = %d, want 0", len(sessionStore))
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			session := NewSessionStorage()

			// given
			if tt.setupFunc != nil {
				tt.setupFunc(session)
			}

			// when
			err := session.Save(tt.sessionID, tt.sessiondata)

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
	ss := NewSessionStorage()
	logger := mylogger.NewLogger()

	// sessionStoreを初期化（他のテストの影響を避けるため）
	sessionStore = make(map[session.SessionID]dto.SessionData)

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

	sessionID := session.SessionID("test-session-id")
	err = ss.Save(sessionID, dto.NewSessionData(param, nil))
	if err != nil {
		t.Fatalf("Save() error = %v", err)
	}

	tests := []struct {
		name      string
		sessionID session.SessionID
		wantErr   bool
		expectNil bool
	}{
		{
			name:      "正常系 - 既存セッションの取得",
			sessionID: session.SessionID("test-session-id"),
			wantErr:   false,
			expectNil: false,
		},
		{
			name:      "異常系 - 存在しないセッションの取得",
			sessionID: session.SessionID("non-existing-session"),
			wantErr:   true,
			expectNil: true,
		},
		{
			name:      "異常系 - 空のセッションIDでの取得",
			sessionID: session.SessionID(""),
			wantErr:   true,
			expectNil: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := ss.Get(tt.sessionID)

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
				if result.AuthParam().ClientID() != param.ClientID() {
					t.Errorf("Get() result.ClientID() = %v, want %v", result.AuthParam().ClientID(), param.ClientID())
				}
				if result.AuthParam().RedirectURI() != param.RedirectURI() {
					t.Errorf("Get() result.RedirectURI() = %v, want %v", result.AuthParam().RedirectURI(), param.RedirectURI())
				}
				if result.AuthParam().State() != param.State() {
					t.Errorf("Get() result.State() = %v, want %v", result.AuthParam().State(), param.State())
				}
			}
		})
	}
}

func Test_認可リクエストパラメータの削除(t *testing.T) {
	logger := mylogger.NewLogger()

	tests := []struct {
		name      string
		setupFunc func(*SessionStorage)
		sessionID session.SessionID
	}{
		{
			name: "正常系 - 既存セッションの削除",
			setupFunc: func(ss *SessionStorage) {
				sessionStore = make(map[session.SessionID]dto.SessionData)
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
				sessionStore[session.SessionID("test-session-id")] = *dto.NewSessionData(param, nil)
			},
			sessionID: session.SessionID("test-session-id"),
		},
		{
			name: "異常系 - 存在しないセッションの削除",
			setupFunc: func(ss *SessionStorage) {
				sessionStore = make(map[session.SessionID]dto.SessionData)
			},
			sessionID: session.SessionID("non-existing-session"),
		},
		{
			name: "異常系 - 空のセッションIDの削除",
			setupFunc: func(ss *SessionStorage) {
				sessionStore = make(map[session.SessionID]dto.SessionData)
			},
			sessionID: session.SessionID(""),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			session := NewSessionStorage()
			tt.setupFunc(session)

			session.Delete(tt.sessionID)
			afterLen := len(sessionStore)

			if afterLen != 0 {
				t.Errorf("after Clear: sessionStore length = %d, want 0", afterLen)
			}
		})
	}
}

func Test_セッションの削除(t *testing.T) {
	logger := mylogger.NewMockLogger()

	tests := []struct {
		name        string
		setupFunc   func(*SessionStorage)
		sessionID   session.SessionID
		expectExist bool
		expectErr   error
	}{
		{
			name: "正常系 - 既存セッションの削除",
			setupFunc: func(ss *SessionStorage) {
				sessionStore = make(map[session.SessionID]dto.SessionData)
				param, _ := domain.NewAuthorizationCodeFlowParam(
					logger,
					"code",
					"test-client",
					"https://example.com/callback",
					"read write",
					"test-state",
				)
				sessionStore[session.SessionID("delete-session-id")] = *dto.NewSessionData(param, nil)
			},
			sessionID:   session.SessionID("delete-session-id"),
			expectExist: false,
			expectErr:   nil,
		},
		{
			name: "異常系 - 存在しないセッションの削除",
			setupFunc: func(ss *SessionStorage) {
				sessionStore = make(map[session.SessionID]dto.SessionData)
			},
			sessionID:   session.SessionID("not-exist-session"),
			expectExist: false,
			expectErr:   nil,
		},
		{
			name: "異常系 - 空のセッションIDの削除",
			setupFunc: func(ss *SessionStorage) {
				sessionStore = make(map[session.SessionID]dto.SessionData)
			},
			sessionID:   session.SessionID(""),
			expectExist: false,
			expectErr:   ErrInvalidSessionID,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ss := NewSessionStorage()
			tt.setupFunc(ss)

			err := ss.Delete(tt.sessionID)

			if _, exists := sessionStore[tt.sessionID]; exists != tt.expectExist {
				t.Errorf("sessionStore existence after Delete = %v, want %v", exists, tt.expectExist)
			}
			if err != nil && err != tt.expectErr {
				t.Errorf("Delete() error = %v, want %v", err, tt.expectErr)
			}
			if err == nil && tt.expectErr != nil {
				t.Errorf("Delete() error = nil, want %v", tt.expectErr)
			}
		})
	}
}
