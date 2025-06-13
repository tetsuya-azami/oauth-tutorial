package usecase

import (
	"oauth-tutorial/internal/repository"
)

type AuthUsecase struct {
	AuthCodeRepo repository.AuthCodeRepository
	TokenRepo    repository.TokenRepository
}

func NewAuthUsecase(acr repository.AuthCodeRepository, tr repository.TokenRepository) *AuthUsecase {
	return &AuthUsecase{
		AuthCodeRepo: acr,
		TokenRepo:    tr,
	}
}

// TODO: 認可コード発行、トークン発行などのビジネスロジックを実装
