package infrastructure

import (
	"errors"
	"oauth-tutorial/internal/domain"
)

type AuthParamSession struct{}

func NewAuthParamSession() *AuthParamSession {
	return &AuthParamSession{}
}

var (
	sessionStore        = map[string]domain.AuthorizationCodeFlowParam{}
	ErrInvalidParameter = errors.New("sessionID is required")
	ErrSessionNotFound  = errors.New("authParam not found")
)

// TODO: 型で制約をかける
func (s *AuthParamSession) Save(sessionID string, authParam *domain.AuthorizationCodeFlowParam) error {
	if sessionID == "" {
		return ErrInvalidParameter
	}
	if authParam == nil {
		return ErrInvalidParameter
	}

	sessionStore[sessionID] = *authParam
	return nil
}

func (s *AuthParamSession) Get(sessionID string) (*domain.AuthorizationCodeFlowParam, error) {
	authParam, ok := sessionStore[sessionID]
	if !ok {
		return nil, ErrSessionNotFound
	}
	return &authParam, nil
}

func (s *AuthParamSession) Clear(sessionID string) {
	delete(sessionStore, sessionID)
}
