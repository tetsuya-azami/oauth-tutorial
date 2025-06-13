package infrastructure

import (
	"errors"
	"oauth-tutorial/internal/domain"
	"sync"
)

type InMemoryAuthCodeRepo struct {
	store map[string]*domain.AuthCode
	mu    sync.RWMutex
}

func NewInMemoryAuthCodeRepo() *InMemoryAuthCodeRepo {
	return &InMemoryAuthCodeRepo{
		store: make(map[string]*domain.AuthCode),
	}
}

func (r *InMemoryAuthCodeRepo) Save(code *domain.AuthCode) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.store[code.Code] = code
	return nil
}

func (r *InMemoryAuthCodeRepo) FindByCode(code string) (*domain.AuthCode, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	c, ok := r.store[code]
	if !ok {
		return nil, errors.New("not found")
	}
	return c, nil
}

func (r *InMemoryAuthCodeRepo) Delete(code string) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	delete(r.store, code)
	return nil
}
