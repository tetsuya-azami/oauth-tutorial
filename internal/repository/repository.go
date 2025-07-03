package repository

import "oauth-tutorial/internal/domain"

// AuthCodeRepository defines methods for storing and retrieving authorization codes
type AuthCodeRepository interface {
	Save(code *domain.AuthorizationCode) error
	FindByCode(code string) (*domain.AuthorizationCode, error)
	Delete(code string) error
}

// TokenRepository defines methods for storing and retrieving tokens
type TokenRepository interface {
	Save(token *domain.AccessToken) error
	FindByAccessToken(token string) (*domain.AccessToken, error)
}
