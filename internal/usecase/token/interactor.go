package token

import (
	"errors"
	"oauth-tutorial/internal/domain"
	"oauth-tutorial/pkg/mylogger"
	"time"
)

var (
	ErrClientNotFound            = errors.New("client not found")
	ErrInvalidClientCredential   = errors.New("invalid client credentials")
	ErrAuthorizationCodeNotFound = errors.New("authorization code not found")
	ErrInvalidClientID           = errors.New("invalid client ID")
	ErrInvalidRedirectURI        = errors.New("invalid redirect URI")
	ErrAuthorizationCodeExpired  = errors.New("authorization code expired")
)

type PublishTokenUsecase struct {
	logger mylogger.Logger
	cr     IClientRepository
	ar     IAuthorizationCodeRepository
	tr     ITokenRepository
}

func NewPublishTokenUsecase(logger mylogger.Logger, cr IClientRepository, ar IAuthorizationCodeRepository, tr ITokenRepository) *PublishTokenUsecase {
	return &PublishTokenUsecase{
		logger: logger,
		cr:     cr,
		ar:     ar,
		tr:     tr,
	}
}

// Token発行処理
func (i *PublishTokenUsecase) Execute(input PublishTokenInput) (string, error) {
	now := time.Now()

	// パブリッククライアントの場合クライアント認証をしないため、ここではClientIDのみでClient情報を取得する
	client, err := i.cr.FindByID(input.ClientID())
	if err != nil {
		return "", ErrClientNotFound
	}

	// コンフィデンシャルクライアントはClient認証
	if client.ClientType() == domain.ConfidentialClient {
		if input.ClientSecret() != client.Secret() {
			return "", ErrInvalidClientCredential
		}
	}

	// 認可コード情報取得
	authCode, err := i.ar.FindByCode(input.Code())
	if err != nil {
		return "", ErrAuthorizationCodeNotFound
	}

	// 認可コードとToken発行リクエストの検証
	err = i.isExchangeable(authCode, input, now)
	if err != nil {
		return "", err
	}

	// Token発行
	token := domain.NewAccessToken(input.clientID, authCode.UserID(), authCode.Scopes(), now)
	// Token登録
	i.tr.Save(token)

	// 認可コード削除
	i.ar.Delete(authCode.Value())

	return token.Value(), nil
}

func (*PublishTokenUsecase) isExchangeable(authCode *domain.AuthorizationCode, input PublishTokenInput, now time.Time) error {
	if authCode.IsExpired(now) {
		return ErrAuthorizationCodeExpired
	}
	if input.ClientID() != authCode.ClientID() {
		return ErrInvalidClientID
	}
	if input.RedirectURI() != authCode.RedirectURI() {
		return ErrInvalidRedirectURI
	}

	return nil
}
