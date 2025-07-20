package usecase

import (
	"oauth-tutorial/internal/domain"
)

type IClientRepository interface {
	SelectByClientIDAndSecret(clientID string) (*domain.Client, error)
}

type IAuthParamSession interface {
	Save(sessionID string, authParam *domain.AuthorizationCodeFlowParam) error
}
