package usecase

import (
	"errors"
	"oauth-tutorial/internal/domain"
	"oauth-tutorial/internal/infrastructure"
	"oauth-tutorial/internal/session"
	"oauth-tutorial/internal/test"
	"testing"
)

type MockClientRepository struct {
	client *domain.Client
	err    error
}

func NewMockClientRepository(client *domain.Client, err error) *MockClientRepository {
	return &MockClientRepository{
		client: client,
		err:    err,
	}
}

func (m *MockClientRepository) SelectByClientID(clientID string) (*domain.Client, error) {
	if m.err != nil {
		return nil, m.err
	}
	return m.client, nil
}

type MockSessionIdGenerator struct {
	id session.SessionID
}

func NewMockSessionIdGenerator(id session.SessionID) *MockSessionIdGenerator {
	return &MockSessionIdGenerator{
		id: id,
	}
}

func (m *MockSessionIdGenerator) Generate() session.SessionID {
	return m.id
}

type MockAuthParamSession struct {
	err error
}

func NewMockAuthParamSession(err error) *MockAuthParamSession {
	return &MockAuthParamSession{
		err: err,
	}
}

func (m *MockAuthParamSession) Save(sessionID session.SessionID, authParam *domain.AuthorizationCodeFlowParam) error {
	if m.err != nil {
		return m.err
	}
	return nil
}

func Test_認可コードフローユースケース(t *testing.T) {
	validClient := domain.ReconstructClient(
		"test-client",
		"Test Client",
		"test-secret",
		[]string{"https://example.com/callback", "https://app.example.com/auth"},
	)

	logger := test.NewMockLogger()
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

	invalidRedirectParam, err := domain.NewAuthorizationCodeFlowParam(
		logger,
		"code",
		"test-client",
		"https://malicious.com/callback",
		"read write",
		"test-state",
	)
	if err != nil {
		t.Fatalf("Failed to create invalid redirect AuthorizationCodeFlowParam: %v", err)
	}

	tests := []struct {
		name        string
		param       *domain.AuthorizationCodeFlowParam
		setupFunc   func() *AuthorizationCodeFlow
		wantErr     bool
		expectedErr error
	}{
		{
			name:  "正常ケース - 認可フロー成功",
			param: validParam,
			setupFunc: func() *AuthorizationCodeFlow {
				cr := NewMockClientRepository(validClient, nil)
				aps := NewMockAuthParamSession(nil)
				sig := NewMockSessionIdGenerator("test-session-id")
				return NewAuthorizationCodeFlow(logger, cr, sig, aps)
			},
			wantErr:     false,
			expectedErr: nil,
		},
		{
			name:  "異常ケース - クライアントが見つからない",
			param: validParam,
			setupFunc: func() *AuthorizationCodeFlow {
				clientRepo := NewMockClientRepository(nil, infrastructure.ErrClientNotFound)
				sessionStore := NewMockAuthParamSession(nil)
				sig := NewMockSessionIdGenerator("test-session-id")
				return NewAuthorizationCodeFlow(logger, clientRepo, sig, sessionStore)
			},
			wantErr:     true,
			expectedErr: ErrClientNotFound,
		},
		{
			name:  "異常ケース - client取得でデータベースエラー",
			param: validParam,
			setupFunc: func() *AuthorizationCodeFlow {
				clientRepo := NewMockClientRepository(nil, errors.New("database error"))
				sessionStore := NewMockAuthParamSession(nil)
				sig := NewMockSessionIdGenerator("test-session-id")
				return NewAuthorizationCodeFlow(logger, clientRepo, sig, sessionStore)
			},
			wantErr:     true,
			expectedErr: ErrUnExpected,
		},
		{
			name:  "異常ケース - 無効なリダイレクトURI",
			param: invalidRedirectParam,
			setupFunc: func() *AuthorizationCodeFlow {
				clientRepo := NewMockClientRepository(validClient, nil)
				sessionStore := NewMockAuthParamSession(nil)
				sig := NewMockSessionIdGenerator("test-session-id")
				return NewAuthorizationCodeFlow(logger, clientRepo, sig, sessionStore)
			},
			wantErr:     true,
			expectedErr: ErrInvalidRedirectURI,
		},
		{
			name:  "異常ケース - セッション保存エラー",
			param: validParam,
			setupFunc: func() *AuthorizationCodeFlow {
				clientRepo := NewMockClientRepository(validClient, nil)
				sessionStore := NewMockAuthParamSession(infrastructure.ErrInvalidParameter)
				sig := NewMockSessionIdGenerator("test-session-id")
				return NewAuthorizationCodeFlow(logger, clientRepo, sig, sessionStore)
			},
			wantErr:     true,
			expectedErr: ErrServer,
		},
		{
			name:  "異常ケース - セッション保存で予期しないエラー",
			param: validParam,
			setupFunc: func() *AuthorizationCodeFlow {
				clientRepo := NewMockClientRepository(validClient, nil)
				sessionStore := NewMockAuthParamSession(errors.New("unexpected error"))
				sig := NewMockSessionIdGenerator("test-session-id")
				return NewAuthorizationCodeFlow(logger, clientRepo, sig, sessionStore)
			},
			wantErr:     true,
			expectedErr: ErrUnExpected,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// given
			flow := tt.setupFunc()

			// when
			sessionID, err := flow.Execute(tt.param)

			// then
			if tt.wantErr {
				if err == nil {
					t.Errorf("Execute() error = nil, wantErr %v", tt.wantErr)
					return
				}
				if !errors.Is(err, tt.expectedErr) {
					t.Errorf("Execute() error = %v, want %v", err, tt.expectedErr)
				}
			} else {
				if err != nil {
					t.Errorf("Execute() error = %v, wantErr %v", err, tt.wantErr)
				}
				if sessionID != "test-session-id" {
					t.Errorf("Execute() sessionID = %v, want %v", sessionID, "test-session-id")
				}
			}
		})
	}
}
