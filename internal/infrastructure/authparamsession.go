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
	sessionStore       = map[string]domain.AuthorizationCodeFlowParam{}
	ErrSessionNotFound = errors.New("authParam not found")
)

func (s *AuthParamSession) Save(sessionID string, authParam *domain.AuthorizationCodeFlowParam) error {
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

func (s *AuthParamSession) Clear(sessionID string) error {
	delete(sessionStore, sessionID)
	return nil
}
