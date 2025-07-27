package authorize

import (
	"errors"
	"net/http"
	"oauth-tutorial/internal/domain"
	"oauth-tutorial/internal/logger"
	"oauth-tutorial/internal/presentation"
	usecase "oauth-tutorial/internal/usecase/authorization"
)

type AuthorizeHandler struct {
	logger logger.MyLogger
	// TODO: response_typeによってflowを変える。(認可コードフロー以外にも対応)(factory patternを検討)
	authorizationFlow IAuthorizationFlow
}

func NewAuthorizeHandler(logger logger.MyLogger, clientGetter IAuthorizationFlow) *AuthorizeHandler {
	return &AuthorizeHandler{logger: logger, authorizationFlow: clientGetter}
}

func (h *AuthorizeHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	queries := r.URL.Query()

	// クエリパラメータを変数に格納
	responseType := queries.Get("response_type")
	clientID := queries.Get("client_id")
	redirectURI := queries.Get("redirect_uri")
	state := queries.Get("state")
	scope := queries.Get("scope")

	param, err := domain.NewAuthorizationCodeFlowParam(h.logger, responseType, clientID, redirectURI, scope, state)
	if err != nil {
		presentation.WriteJSONResponse(w, http.StatusBadRequest, ErrorResponse{Message: err.Error()})
		return
	}

	sessionID, err := h.authorizationFlow.Execute(param)
	if err != nil {
		switch {
		case errors.Is(err, usecase.ErrClientNotFound):
			h.logger.Info("Client not found", "clientID", clientID)
			presentation.WriteJSONResponse(w, http.StatusBadRequest, ErrorResponse{Message: err.Error()})
		case errors.Is(err, usecase.ErrInvalidRedirectURI):
			h.logger.Info("Invalid redirect URI", "redirectURI", redirectURI)
			presentation.WriteJSONResponse(w, http.StatusBadRequest, ErrorResponse{Message: err.Error()})
		case errors.Is(err, usecase.ErrUnExpected):
			h.logger.Error("Unexpected error occurred", "error", err)
			presentation.WriteJSONResponse(w, http.StatusInternalServerError, ErrorResponse{Message: err.Error()})
		case errors.Is(err, usecase.ErrServer):
			h.logger.Error("Server error occurred", "error", err)
			presentation.WriteJSONResponse(w, http.StatusInternalServerError, ErrorResponse{Message: err.Error()})
		default:
			h.logger.Error("Unexpected error occurred while getting client", "error", err)
			presentation.WriteJSONResponse(w, http.StatusInternalServerError, ErrorResponse{Message: err.Error()})
		}
		return
	}

	h.logger.Info("Client authorized successfully")

	// 本来は認証画面を表示するが、ここではOKのレスポンスを返すだけとする
	http.SetCookie(w, &http.Cookie{
		Name:     "SESSION_ID",
		Value:    string(sessionID),
		Path:     "/",
		HttpOnly: true,
		Secure:   true,
	})
	presentation.WriteJSONResponse(w, http.StatusOK, SuccessResponse{Message: "OK"})
}
