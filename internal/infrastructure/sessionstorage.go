package infrastructure

import (
	"errors"
	"oauth-tutorial/internal/session"
)

type SessionStorage struct{}

func NewSessionStorage() *SessionStorage {
	return &SessionStorage{}
}

var (
	sessionStore        = map[session.SessionID]SessionData{}
	ErrInvalidParameter = errors.New("sessionID is required")
	ErrSessionNotFound  = errors.New("authParam not found")
)

func (s *SessionStorage) Save(sessionID session.SessionID, sessiondata *SessionData) error {
	if sessionID == "" {
		return ErrInvalidParameter
	}

	sessionStore[sessionID] = *sessiondata
	return nil
}

func (s *SessionStorage) Get(sessionID session.SessionID) (*SessionData, error) {
	sessionData, ok := sessionStore[sessionID]
	if !ok {
		return nil, ErrSessionNotFound
	}
	return &sessionData, nil
}
