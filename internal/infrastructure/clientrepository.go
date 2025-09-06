package infrastructure

import (
	"errors"
	"oauth-tutorial/internal/domain"
)

type ClientRepository struct {
	clients map[domain.ClientID]*domain.Client
}

var ErrClientNotFound = errors.New("client not found")

func NewClientRepository() *ClientRepository {
	clients := map[domain.ClientID]*domain.Client{
		"iouobrnea": domain.ReconstructClient(domain.ClientID("iouobrnea"), "client-1", domain.ConfidentialClient, "password", []string{"https://client.example.com/callback"}),
	}
	return &ClientRepository{clients: clients}
}

func (r *ClientRepository) SelectByClientID(clientID domain.ClientID) (*domain.Client, error) {
	client, ok := r.clients[clientID]
	if !ok {
		return nil, ErrClientNotFound
	}
	return client, nil
}

func (r *ClientRepository) SelectByClientIDAndClientSecret(clientID domain.ClientID, clientSecret string) (*domain.Client, error) {
	client, ok := r.clients[clientID]
	if !ok {
		return nil, ErrClientNotFound
	}
	if client.Secret() != clientSecret {
		return nil, ErrClientNotFound
	}

	return client, nil
}
