package infrastructure

import (
	"errors"
	"oauth-tutorial/internal/domain"
	"sync"
)

type InMemoryTokenRepo struct {
	store map[string]*domain.Token
	mu    sync.RWMutex
}

func NewInMemoryTokenRepo() *InMemoryTokenRepo {
	return &InMemoryTokenRepo{
		store: make(map[string]*domain.Token),
	}
}

func (r *InMemoryTokenRepo) Save(token *domain.Token) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.store[token.AccessToken] = token
	return nil
}

func (r *InMemoryTokenRepo) FindByAccessToken(token string) (*domain.Token, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	t, ok := r.store[token]
	if !ok {
		return nil, errors.New("not found")
	}
	return t, nil
}
