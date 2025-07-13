package presentation

import "github.com/google/uuid"

type SessionID string

func GenerateSessionID() string {
	return uuid.New().String()
}
