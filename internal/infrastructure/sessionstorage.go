package infrastructure

import (
	"errors"
	"oauth-tutorial/internal/domain"
	"oauth-tutorial/internal/session"
)

type SessionStorage struct{}

func NewSessionStorage() *SessionStorage {
	return &SessionStorage{}
}

var (
	sessionStore        = map[session.SessionID]domain.AuthorizationCodeFlowParam{}
	ErrInvalidParameter = errors.New("sessionID is required")
	ErrSessionNotFound  = errors.New("authParam not found")
)

func (s *SessionStorage) Save(sessionID session.SessionID, authParam *domain.AuthorizationCodeFlowParam) error {
	if sessionID == "" {
		return ErrInvalidParameter
	}
	if authParam == nil {
		return ErrInvalidParameter
	}

	sessionStore[sessionID] = *authParam
	return nil
}

func (s *SessionStorage) Get(sessionID session.SessionID) (*domain.AuthorizationCodeFlowParam, error) {
	authParam, ok := sessionStore[sessionID]
	if !ok {
		return nil, ErrSessionNotFound
	}
	return &authParam, nil
}

func (s *SessionStorage) Clear(sessionID session.SessionID) {
	delete(sessionStore, sessionID)
}
