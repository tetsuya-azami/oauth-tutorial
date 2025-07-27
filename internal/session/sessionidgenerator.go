package session

import (
	"github.com/google/uuid"
)

const SessionIDCookieName = "SESSION_ID"

type SessionID string
type SessionIDGenerator struct{}

func NewSessionIDGenerator() *SessionIDGenerator {
	return &SessionIDGenerator{}
}

func (g *SessionIDGenerator) Generate() SessionID {
	return SessionID(uuid.New().String())
}
