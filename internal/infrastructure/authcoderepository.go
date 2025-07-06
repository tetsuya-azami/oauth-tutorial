package infrastructure

import (
	"errors"
	"oauth-tutorial/internal/domain"
	"sync"
)

type IAuthCodeRepository interface {
	Save(code *domain.AuthorizationCode) error
	FindByCode(code string) (*domain.AuthorizationCode, error)
	Delete(code string) error
}

type AuthCodeRepository struct {
	authCodeStore map[string]*domain.AuthorizationCode
	mu            sync.RWMutex
}

func NewAuthCodeRepository() *AuthCodeRepository {
	return &AuthCodeRepository{
		authCodeStore: make(map[string]*domain.AuthorizationCode),
	}
}

func (r *AuthCodeRepository) Save(code *domain.AuthorizationCode) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.authCodeStore[code.Value()] = code
	return nil
}

func (r *AuthCodeRepository) FindByCode(code string) (*domain.AuthorizationCode, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	v, ok := r.authCodeStore[code]
	if !ok {
		return nil, errors.New("not found")
	}
	return v, nil
}

func (r *AuthCodeRepository) Delete(code string) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	delete(r.authCodeStore, code)
	return nil
}
