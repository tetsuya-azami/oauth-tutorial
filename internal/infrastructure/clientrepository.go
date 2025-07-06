package infrastructure

import (
	"errors"
	"oauth-tutorial/internal/domain"
)

type IClientRepository interface {
	SelectByClientIDAndSecret(clientID, clientSecret string) (*domain.Client, error)
}

type ClientRepository struct {
	clients map[ClientIdAndSecretPair]*domain.Client
}

type ClientIdAndSecretPair struct {
	ClientID     string
	ClientSecret string
}

var clients = map[ClientIdAndSecretPair]*domain.Client{
	{ClientID: "iouobrnea", ClientSecret: "password"}: domain.ReconstructClient("iouobrnea", "client-1", "password", []string{"https://client.example.com/callback"}),
}

var ErrClientNotFound = errors.New("client not found")

func (*ClientRepository) NewClientRepository() *ClientRepository {
	return &ClientRepository{clients: clients}
}

func (r *ClientRepository) SelectByClientIDAndSecret(clientID, clientSecret string) (*domain.Client, error) {
	client, ok := r.clients[ClientIdAndSecretPair{ClientID: clientID, ClientSecret: clientSecret}]
	if !ok {
		return nil, ErrClientNotFound
	}
	return client, nil
}
