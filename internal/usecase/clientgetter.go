package usecase

import (
	"errors"
	"oauth-tutorial/internal/domain"
	"oauth-tutorial/internal/infrastructure"
)

type IClientGetter interface {
	GetClient(clientID string, secret string) (*domain.Client, error)
}

type ClientGetter struct {
	clientRepository infrastructure.IClientRepository
}

func NewClientGetter(cr infrastructure.IClientRepository) *ClientGetter {
	return &ClientGetter{
		clientRepository: cr,
	}
}

var ErrClientNotFound = errors.New("client not found")
var ErrUnExpected = errors.New("unexpected error occurred")

func (c *ClientGetter) GetClient(clientID string, secret string) (*domain.Client, error) {
	cr := c.clientRepository
	client, err := cr.SelectByClientIDAndSecret(clientID, secret)
	if err != nil {
		switch {
		case errors.Is(err, infrastructure.ErrClientNotFound):
			return nil, ErrClientNotFound
		default:
			return nil, ErrUnExpected
		}
	}

	return client, nil
}
