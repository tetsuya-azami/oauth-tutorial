package authorize

import (
	"oauth-tutorial/internal/domain"
	inf_dto "oauth-tutorial/internal/infrastructure/dto"
	"oauth-tutorial/internal/session"
)

type IClientRepository interface {
	SelectByClientID(clientID string) (*domain.Client, error)
}

type ISessionStorage interface {
	Save(sessionID session.SessionID, sessionData *inf_dto.SessionData) error
}

type ISessionIDGenerator interface {
	Generate() session.SessionID
}
