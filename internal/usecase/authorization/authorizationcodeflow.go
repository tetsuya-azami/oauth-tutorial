package usecase

import (
	"errors"
	"log/slog"
	"oauth-tutorial/internal/domain"
	"oauth-tutorial/internal/infrastructure"
	"oauth-tutorial/internal/presentation"
)

type AuthorizationCodeFlow struct {
	logger           *slog.Logger
	clientRepository IClientRepository
	sessionStore     IAuthParamSession
}

func NewAuthorizationCodeFlow(logger *slog.Logger, cr IClientRepository, sessionStore IAuthParamSession) *AuthorizationCodeFlow {
	return &AuthorizationCodeFlow{
		logger:           logger,
		clientRepository: cr,
		sessionStore:     sessionStore,
	}
}

var (
	ErrClientNotFound     = errors.New("client not found")
	ErrUnExpected         = errors.New("unexpected error occurred")
	ErrInvalidRedirectURI = errors.New("invalid redirect URI")
)

func (c *AuthorizationCodeFlow) Execute(param *domain.AuthorizationCodeFlowParam) error {
	cr := c.clientRepository
	client, err := cr.SelectByClientIDAndSecret(param.ClientID())
	if err != nil {
		switch {
		case errors.Is(err, infrastructure.ErrClientNotFound):
			return ErrClientNotFound
		default:
			return ErrUnExpected
		}
	}

	if !client.ContainsRedirectURI(param.RedirectURI()) {
		c.logger.Info("invalid Redirect URI", "redirectURI", param.RedirectURI())
		return ErrInvalidRedirectURI
	}

	sessionID := presentation.GenerateSessionID()
	c.sessionStore.Save(sessionID, param)

	return nil
}
