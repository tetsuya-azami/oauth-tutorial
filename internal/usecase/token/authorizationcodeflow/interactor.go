package authorizationcodeflow

import (
	"errors"
	"oauth-tutorial/internal/domain"
	tokenport "oauth-tutorial/internal/usecase/token/port"
	"oauth-tutorial/pkg/mylogger"
	"time"
)

var (
	ErrInvalidInputType          = errors.New("invalid input type")
	ErrClientNotFound            = errors.New("client not found")
	ErrInvalidClientCredential   = errors.New("invalid client credentials")
	ErrAuthorizationCodeNotFound = errors.New("authorization code not found")
	ErrInvalidClientID           = errors.New("invalid client ID")
	ErrInvalidRedirectURI        = errors.New("invalid redirect URI")
	ErrAuthorizationCodeExpired  = errors.New("authorization code expired")
)

type AuthorizationCodeFlow struct {
	logger mylogger.Logger
	cr     tokenport.IClientRepository
	ar     tokenport.IAuthorizationCodeRepository
	tr     tokenport.ITokenRepository
}

func NewAuthorizationCodeFlow(logger mylogger.Logger, cr tokenport.IClientRepository, ar tokenport.IAuthorizationCodeRepository, tr tokenport.ITokenRepository) *AuthorizationCodeFlow {
	return &AuthorizationCodeFlow{
		logger: logger,
		cr:     cr,
		ar:     ar,
		tr:     tr,
	}
}

// Token発行処理
func (i *AuthorizationCodeFlow) Execute(input any) (*domain.AccessToken, *domain.RefreshToken, error) {
	// TODO: 時刻のinjectの仕方考える
	now := time.Now()
	ai, ok := input.(AuthorizationCodeInput)
	if !ok {
		i.logger.Info("inputとinteractorの不整合です。", "input", input)
		return nil, nil, ErrInvalidInputType
	}

	// パブリッククライアントの場合クライアント認証をしないため、ここではClientIDのみでClient情報を取得する
	client, err := i.cr.FindByID(ai.ClientID())
	if err != nil {
		i.logger.Info("client_idに該当するClientが存在しません。", "err", err, "client_id", ai.ClientID())
		return nil, nil, ErrClientNotFound
	}

	// コンフィデンシャルクライアントはClient認証
	// TODO: ブルートフォース攻撃対策
	if client.ClientType() == domain.ConfidentialClient {
		if ai.ClientSecret() != client.Secret() {
			i.logger.Info("client認証に失敗しました。", "client_id", ai.ClientID())
			return nil, nil, ErrInvalidClientCredential
		}
	}

	// 認可コード情報取得
	authCode, err := i.ar.FindByCode(ai.Code())
	if err != nil {
		i.logger.Info("codeに該当する認可コードが存在しません。", "err", err, "code", ai.Code())
		return nil, nil, ErrAuthorizationCodeNotFound
	}

	// 認可コードとToken発行リクエストの検証
	err = i.isExchangeable(authCode, ai, now, i.logger)
	if err != nil {
		return nil, nil, err
	}

	// Token発行
	token := domain.NewAccessToken(ai.ClientID(), authCode.UserID(), authCode.Scopes(), now)
	// Token登録
	i.tr.Save(token)

	// RefreshToken発行・登録
	refreshToken := domain.NewRefreshToken(ai.ClientID(), now)
	i.tr.SaveRefreshToken(refreshToken, token)

	// 認可コード削除
	i.ar.Delete(authCode.Value())

	return token, refreshToken, nil
}

func (*AuthorizationCodeFlow) isExchangeable(authCode *domain.AuthorizationCode, ai AuthorizationCodeInput, now time.Time, logger mylogger.Logger) error {
	if authCode.IsExpired(now) {
		logger.Info("認可コードの有効期限が切れています。", "input.client_id", ai.ClientID(), "authCode.client_id", authCode.ClientID())
		return ErrAuthorizationCodeExpired
	}
	if ai.ClientID() != authCode.ClientID() {
		logger.Info("リクエストのclient_idが認可コードのclient_idと一致しません。", "input.client_id", ai.ClientID(), "authCode.client_id", authCode.ClientID())
		return ErrInvalidClientID
	}
	if ai.RedirectURI() != authCode.RedirectURI() {
		logger.Info("リクエストのredirect_uriが認可コードのredirect_uriと一致しません。", "input.redirect_uri", ai.RedirectURI(), "authCode.redirect_uri", authCode.RedirectURI())
		return ErrInvalidRedirectURI
	}

	return nil
}
