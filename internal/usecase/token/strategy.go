package token

import (
	"errors"
	"oauth-tutorial/internal/domain"
	"oauth-tutorial/internal/usecase/token/authorizationcodeflow"
	tokenport "oauth-tutorial/internal/usecase/token/port"
	"oauth-tutorial/internal/usecase/token/refreshtokenflow"
	"oauth-tutorial/pkg/mylogger"
)

var ErrNoMatchingStrategyFound = errors.New("no matching strategy found")

type PublishTokenStrategy struct {
	logger mylogger.Logger
	cr     tokenport.IClientRepository
	ar     tokenport.IAuthorizationCodeRepository
	tr     tokenport.ITokenRepository
}

func NewPublishTokenStrategy(logger mylogger.Logger, cr tokenport.IClientRepository, ar tokenport.IAuthorizationCodeRepository, tr tokenport.ITokenRepository) *PublishTokenStrategy {
	return &PublishTokenStrategy{
		logger: logger,
		cr:     cr,
		ar:     ar,
		tr:     tr,
	}
}

type UsecaseInput = any

type Usecase interface {
	Execute(input UsecaseInput) (*domain.AccessToken, *domain.RefreshToken, error)
}

func (s *PublishTokenStrategy) ResolvePublishTokenFlow(grantType domain.GrantType) (Usecase, error) {
	switch grantType {
	case domain.GrantTypeAuthorizationCode:
		return authorizationcodeflow.NewAuthorizationCodeFlow(s.logger, s.cr, s.ar, s.tr), nil
	case domain.GrantTypeRefreshToken:
		return refreshtokenflow.NewRefreshTokenFlow(s.logger, s.cr, s.tr), nil
	default:
		s.logger.Error("enumでサポートしているgrant_typeがinteractorで実装されていません。")
		return nil, ErrNoMatchingStrategyFound
	}
}
