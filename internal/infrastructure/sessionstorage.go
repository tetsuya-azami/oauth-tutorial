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
	sessionStore          = map[session.SessionID]SessionData{}
	ErrInvalidSessionID   = errors.New("sessionID is required")
	ErrInvalidSessionData = errors.New("invalid session data")
	ErrSessionNotFound    = errors.New("session not found")
)

func (s *SessionStorage) Save(sessionID session.SessionID, sessiondata *SessionData) error {
	if sessionID == "" {
		return ErrInvalidSessionID
	}
	if sessiondata == nil {
		return ErrInvalidSessionData
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

func (s *SessionStorage) Delete(sessionID session.SessionID) error {
	if sessionID == "" {
		return ErrInvalidSessionID
	}

	delete(sessionStore, sessionID)
	return nil
}
