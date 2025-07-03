package infrastructure

import (
	"errors"
	"oauth-tutorial/internal/domain"
	"sync"
)

type InMemoryTokenRepo struct {
	store map[string]*domain.AccessToken
	mu    sync.RWMutex
}

func NewInMemoryTokenRepo() *InMemoryTokenRepo {
	return &InMemoryTokenRepo{
		store: make(map[string]*domain.AccessToken),
	}
}

func (r *InMemoryTokenRepo) Save(token *domain.AccessToken) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.store[token.AccessToken] = token
	return nil
}

func (r *InMemoryTokenRepo) FindByAccessToken(token string) (*domain.AccessToken, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	t, ok := r.store[token]
	if !ok {
		return nil, errors.New("not found")
	}
	return t, nil
}
