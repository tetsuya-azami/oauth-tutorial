package authorize

import (
	"errors"
	"oauth-tutorial/internal/domain"
	"oauth-tutorial/internal/infrastructure"
	"oauth-tutorial/internal/session"
	"oauth-tutorial/pkg/mylogger"
)

type AuthorizationCodeFlow struct {
	logger             mylogger.Logger
	clientRepository   IClientRepository
	sessionStore       ISessionStorage
	sessionIDGenerator ISessionIDGenerator
}

func NewAuthorizationCodeFlow(logger mylogger.Logger, cr IClientRepository, sessionIDGenerator ISessionIDGenerator, sessionStorage ISessionStorage) *AuthorizationCodeFlow {
	return &AuthorizationCodeFlow{
		logger:             logger,
		clientRepository:   cr,
		sessionIDGenerator: sessionIDGenerator,
		sessionStore:       sessionStorage,
	}
}

var (
	ErrClientNotFound     = errors.New("client not found")
	ErrUnExpected         = errors.New("unexpected error occurred")
	ErrInvalidRedirectURI = errors.New("invalid redirect URI")
	ErrServer             = errors.New("server error occurred")
)

func (c *AuthorizationCodeFlow) Execute(param *domain.AuthorizationCodeFlowParam) (session.SessionID, error) {
	cr := c.clientRepository
	client, err := cr.SelectByClientID(param.ClientID())
	if err != nil {
		switch {
		case errors.Is(err, infrastructure.ErrClientNotFound):
			c.logger.Info("client not found", "clientID", param.ClientID())
			return "", ErrClientNotFound
		default:
			c.logger.Error("unexpected error occured", "error", err)
			return "", ErrUnExpected
		}
	}

	if !client.ContainsRedirectURI(param.RedirectURI()) {
		c.logger.Info("invalid Redirect URI", "redirectURI", param.RedirectURI())
		return "", ErrInvalidRedirectURI
	}

	sessionID := c.sessionIDGenerator.Generate()

	err = c.sessionStore.Save(sessionID, infrastructure.NewSessionData(param, nil))
	if err != nil {
		switch {
		case errors.Is(err, infrastructure.ErrInvalidSessionID) || errors.Is(err, infrastructure.ErrInvalidSessionData):
			c.logger.Info("invalid session parameter", "error", err)
			return "", ErrServer
		default:
			c.logger.Error("unexpected error occured", "error", err)
			return "", ErrUnExpected
		}
	}

	return sessionID, nil
}
