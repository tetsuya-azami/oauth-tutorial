package infrastructure

import (
	"errors"
	"oauth-tutorial/internal/domain"
)

type ClientRepository struct {
	clients map[string]*domain.Client
}

var ErrClientNotFound = errors.New("client not found")

func NewClientRepository() *ClientRepository {
	clients := map[string]*domain.Client{
		"iouobrnea": domain.ReconstructClient("iouobrnea", "client-1", "password", []string{"https://client.example.com/callback"}),
	}
	return &ClientRepository{clients: clients}
}

func (r *ClientRepository) SelectByClientIDAndSecret(clientID string) (*domain.Client, error) {
	client, ok := r.clients[clientID]
	if !ok {
		return nil, ErrClientNotFound
	}
	return client, nil
}
