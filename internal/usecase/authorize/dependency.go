package authorize

import (
	"oauth-tutorial/internal/domain"
	"oauth-tutorial/internal/infrastructure"
	"oauth-tutorial/internal/session"
)

type IClientRepository interface {
	SelectByClientID(clientID string) (*domain.Client, error)
}

type ISessionStorage interface {
	Save(sessionID session.SessionID, sessionData *infrastructure.SessionData) error
}

type ISessionIDGenerator interface {
	Generate() session.SessionID
}
