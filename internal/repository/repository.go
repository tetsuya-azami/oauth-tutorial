package repository

import "oauth-tutorial/internal/domain"

// AuthCodeRepository defines methods for storing and retrieving authorization codes
type AuthCodeRepository interface {
	Save(code *domain.AuthCode) error
	FindByCode(code string) (*domain.AuthCode, error)
	Delete(code string) error
}

// TokenRepository defines methods for storing and retrieving tokens
type TokenRepository interface {
	Save(token *domain.Token) error
	FindByAccessToken(token string) (*domain.Token, error)
}
