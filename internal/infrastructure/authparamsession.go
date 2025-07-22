package infrastructure

import (
	"errors"
	"oauth-tutorial/internal/domain"
	"oauth-tutorial/internal/session"
)

type AuthParamSession struct{}

func NewAuthParamSession() *AuthParamSession {
	return &AuthParamSession{}
}

var (
	sessionStore        = map[session.SessionID]domain.AuthorizationCodeFlowParam{}
	ErrInvalidParameter = errors.New("sessionID is required")
	ErrSessionNotFound  = errors.New("authParam not found")
)

func (s *AuthParamSession) Save(sessionID session.SessionID, authParam *domain.AuthorizationCodeFlowParam) error {
	if sessionID == "" {
		return ErrInvalidParameter
	}
	if authParam == nil {
		return ErrInvalidParameter
	}

	sessionStore[sessionID] = *authParam
	return nil
}

func (s *AuthParamSession) Get(sessionID session.SessionID) (*domain.AuthorizationCodeFlowParam, error) {
	authParam, ok := sessionStore[sessionID]
	if !ok {
		return nil, ErrSessionNotFound
	}
	return &authParam, nil
}

func (s *AuthParamSession) Clear(sessionID session.SessionID) {
	delete(sessionStore, sessionID)
}
