package refreshtokenflow

import (
	"errors"
	"fmt"
	"oauth-tutorial/internal/domain"
	tokenport "oauth-tutorial/internal/usecase/token/port"
	"oauth-tutorial/pkg/mylogger"
)

var ErrInvalidInputType = errors.New("invalid input type")

type RefreshTokenFlow struct {
	logger mylogger.Logger
	cr     tokenport.IClientRepository
	tr     tokenport.ITokenRepository
}

func NewRefreshTokenFlow(logger mylogger.Logger, cr tokenport.IClientRepository, tr tokenport.ITokenRepository) *RefreshTokenFlow {
	return &RefreshTokenFlow{
		logger: logger,
		cr:     cr,
		tr:     tr,
	}
}

func (r *RefreshTokenFlow) Execute(input any) (*domain.AccessToken, *domain.RefreshToken, error) {
	rti, ok := input.(RefreshTokenInput)
	if !ok {
		return nil, nil, ErrInvalidInputType
	}
	// unusedエラー防止
	fmt.Print(rti)
	// 実装予定
	return nil, nil, nil
}
