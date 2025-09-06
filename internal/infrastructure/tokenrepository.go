package infrastructure

import (
	"errors"
	"oauth-tutorial/internal/domain"
	"sync"
)

type ITokenRepository interface {
	Save(token *domain.AccessToken) error
	FindByAccessToken(token string) (*domain.AccessToken, error)
}

type TokenRepository struct {
	store map[string]*domain.AccessToken
	mu    sync.RWMutex
}

func NewTokenRespository() *TokenRepository {
	return &TokenRepository{
		store: make(map[string]*domain.AccessToken),
	}
}

func (r *TokenRepository) Save(token *domain.AccessToken) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.store[token.Value()] = token
}

func (r *TokenRepository) FindByAccessToken(token string) (*domain.AccessToken, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	t, ok := r.store[token]
	if !ok {
		return nil, errors.New("not found")
	}
	return t, nil
}
