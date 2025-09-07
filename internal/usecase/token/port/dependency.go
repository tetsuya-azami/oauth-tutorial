package tokenport

import "oauth-tutorial/internal/domain"

type IClientRepository interface {
	FindByID(clientID string) (*domain.Client, error)
}

type ITokenRepository interface {
	Save(token *domain.AccessToken)
	SaveRefreshToken(token *domain.RefreshToken, accessToken *domain.AccessToken)
	FindByAccessToken(token string) (*domain.AccessToken, error)
}

type IAuthorizationCodeRepository interface {
	FindByCode(code string) (*domain.AuthorizationCode, error)
	Delete(code string)
}
