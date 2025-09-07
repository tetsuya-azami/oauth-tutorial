package token

import (
	"net/http"
	"oauth-tutorial/internal/domain"
	"oauth-tutorial/internal/presentation"
	utoken "oauth-tutorial/internal/usecase/token"
	"oauth-tutorial/internal/usecase/token/authorizationcodeflow"
	"oauth-tutorial/internal/usecase/token/refreshtokenflow"
	"oauth-tutorial/pkg/mylogger"
	"strings"
)

type TokenHandler struct {
	logger mylogger.Logger
	pts    utoken.PublishTokenStrategy
}

func NewTokenHandler(logger mylogger.Logger, pts utoken.PublishTokenStrategy) *TokenHandler {
	return &TokenHandler{logger: logger, pts: pts}
}

func (h *TokenHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		h.logger.Info("formのParseに失敗しました。", "err", err)
		presentation.WriteJSONResponse(w, http.StatusBadRequest, NewErrorResponse(InvalidRequest, "リクエストが不正です。"))
		return
	}

	formGrantType := r.FormValue("grant_type")
	grantType, err := domain.ResolveGrantType(formGrantType)
	if err != nil {
		h.logger.Info("サポートされていないgrant_typeです。", "grant_type", formGrantType)
		presentation.WriteJSONResponse(w, http.StatusBadRequest, NewErrorResponse(UnsupportedGrantType, "サポートされていないgrant_typeです。"))
		return
	}

	// grantTypeに応じたトークン発行フローを解決する
	interactor, err := h.pts.ResolvePublishTokenFlow(grantType)
	if err != nil {
		h.logger.Error("トークン発行フローの解決に失敗しました。", "err", err)
		presentation.WriteJSONResponse(w, http.StatusInternalServerError, NewErrorResponse(ServerError, "サーバーエラーが発生しました。"))
		return
	}

	// client認証パラメータの解決(クライアント認証をしないクライアントのことを考慮し、この時点ではclient_id, client_secretの空値を許容する)
	clientID, clientSecret := resolveClientAuthenticationParameters(r)

	// トークン発行フローに応じてinputを解決する
	input := resolveInput(r, grantType, clientID, clientSecret)

	// トークン発行フローを実行する
	accessToken, refreshToken, err := interactor.Execute(input)
	if err != nil {
		handleError(w, err, h.logger)
		return
	}

	presentation.WriteJSONResponse(w, http.StatusOK, SuccessResponse{
		AccessToken:  accessToken.Value(),
		RefreshToken: refreshToken.Value(),
		TokenType:    "Bearer",
		ExpiresIn:    int(domain.AccessTokenDuration.Minutes()),
	})
}

func resolveClientAuthenticationParameters(r *http.Request) (clientID string, clientSecret string) {
	// Basic認証のケース
	authorizationHeader := r.Header.Get("Authorization")
	if strings.HasPrefix(authorizationHeader, "Basic ") {
		// Authorizationヘッダーからクレデンシャルを取得
		clientID, clientSecret, _ := r.BasicAuth()
		if clientID != "" && clientSecret != "" {
			return clientID, clientSecret
		}
	}

	// フォームパラメータのケース
	clientID = r.PostFormValue("client_id")
	clientSecret = r.PostFormValue("client_secret")
	if clientID != "" && clientSecret != "" {
		return clientID, clientSecret
	}

	// クライアント認証をしないクライアントのケース
	// クライアント認証をしない場合でも、client_idは必須なので、client_idのみ返す
	if clientID != "" {
		return clientID, ""
	}

	return "", ""
}

func resolveInput(r *http.Request, grantType domain.GrantType, clientID string, clientSecret string) any {
	redirectURI := r.FormValue("redirect_uri")
	if redirectURI == "" {
		return nil
	}

	// grantTypeに応じてinputを解決する
	switch grantType {
	case domain.GrantTypeAuthorizationCode:
		code := r.FormValue("code")
		if strings.TrimSpace(code) == "" {
			return nil
		}
		return authorizationcodeflow.NewAuthorizationCodeInput(clientID, clientSecret, code, redirectURI)
	case domain.GrantTypeRefreshToken:
		refreshToken := r.FormValue("refresh_token")
		if strings.TrimSpace(refreshToken) == "" {
			return nil
		}
		return refreshtokenflow.NewRefreshTokenInput(clientID, clientSecret, refreshToken, redirectURI)
	default:
		return nil
	}
}

// interactor毎にpresenterを用意した方が楽かも。
// 本当はパターンマッチしたいが。。。
func handleError(w http.ResponseWriter, err error, logger mylogger.Logger) {
	if err != nil {
		switch err {
		// 認可コードフローのエラーハンドリング
		case authorizationcodeflow.ErrInvalidInputType:
			presentation.WriteJSONResponse(w, http.StatusBadRequest, NewErrorResponse(InvalidRequest, "リクエストのパラメータが不正です。"))
			return
		case authorizationcodeflow.ErrClientNotFound:
			presentation.WriteJSONResponse(w, http.StatusUnauthorized, NewErrorResponse(InvalidClient, "該当するクライアントが見つかりません。"))
			return
		case authorizationcodeflow.ErrInvalidClientCredential:
			presentation.WriteJSONResponse(w, http.StatusUnauthorized, NewErrorResponse(InvalidClient, "該当するクライアントが見つかりません。"))
			return
		case authorizationcodeflow.ErrAuthorizationCodeNotFound:
			presentation.WriteJSONResponse(w, http.StatusBadRequest, NewErrorResponse(InvalidGrant, "codeが不正です。"))
			return
		case authorizationcodeflow.ErrInvalidClientID:
			presentation.WriteJSONResponse(w, http.StatusBadRequest, NewErrorResponse(InvalidGrant, "codeが不正です。"))
			return
		case authorizationcodeflow.ErrInvalidRedirectURI:
			presentation.WriteJSONResponse(w, http.StatusBadRequest, NewErrorResponse(InvalidGrant, "redirect_uriが不正です。"))
			return
		case authorizationcodeflow.ErrAuthorizationCodeExpired:
			presentation.WriteJSONResponse(w, http.StatusBadRequest, NewErrorResponse(InvalidGrant, "codeの有効期限が切れています。"))
			return
		case utoken.ErrNoMatchingStrategyFound:
			presentation.WriteJSONResponse(w, http.StatusBadRequest, NewErrorResponse(ServerError, "サーバーエラーが発生しました。"))
			return
		default:
			logger.Error("予期せぬエラーが起きました。", "err", err)
			presentation.WriteJSONResponse(w, http.StatusInternalServerError, NewErrorResponse(ServerError, "サーバーエラーが発生しました。"))
			return
		}
	}
}
