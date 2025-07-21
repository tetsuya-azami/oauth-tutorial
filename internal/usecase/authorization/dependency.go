package usecase

import (
	"oauth-tutorial/internal/domain"
	"oauth-tutorial/internal/session"
)

type IClientRepository interface {
	SelectByClientIDAndSecret(clientID string) (*domain.Client, error)
}

type IAuthParamSession interface {
	Save(sessionID session.SessionID, authParam *domain.AuthorizationCodeFlowParam) error
}

type ISessionIDGenerator interface {
	Generate() session.SessionID
}
