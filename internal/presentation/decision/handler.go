package decision

import (
	"errors"
	"net/http"
	"net/url"
	"oauth-tutorial/internal/presentation"
	"oauth-tutorial/internal/session"
	"oauth-tutorial/internal/usecase/decision"
	"oauth-tutorial/pkg/mylogger"
	"strconv"
)

type DecisionHandler struct {
	logger                   mylogger.Logger
	publishAuthorizationCode IPublishAuthorizationCodeUseCase
}

func NewDecisionHandler(logger mylogger.Logger, publishAuthorizationCode IPublishAuthorizationCodeUseCase) *DecisionHandler {
	return &DecisionHandler{logger: logger, publishAuthorizationCode: publishAuthorizationCode}
}

func (h *DecisionHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		h.logger.Info("Failed to parse form: %v", err)
		presentation.WriteJSONResponse(w, http.StatusBadRequest, ErrorResponse{Message: "パラメータの形式を確認してください"})
		return
	}

	input, err := h.convertParamToInput(r.Form, w, r)
	if err != nil {
		presentation.WriteJSONResponse(w, http.StatusBadRequest, ErrorResponse{Message: err.Error()})
		return
	}

	result, err := h.publishAuthorizationCode.Execute(input)
	if err != nil {
		switch {
		case errors.Is(err, decision.ErrSessionNotFound):
			presentation.WriteJSONResponse(w, http.StatusBadRequest, ErrorResponse{Message: err.Error()})
			return
		case errors.Is(err, decision.ErrAuthorizationDenied):
			presentation.WriteJSONResponse(w, http.StatusForbidden, ErrorResponse{Message: err.Error()})
			return
		case errors.Is(err, decision.ErrUnexpectedError):
			presentation.WriteJSONResponse(w, http.StatusInternalServerError, ErrorResponse{Message: err.Error()})
			return
		}
	}

	redirectUri := result.BaseRedirectUri() + "?code=" + result.AuthorizationCode() + "&state=" + result.State()
	http.Redirect(w, r, redirectUri, http.StatusSeeOther)
}

func (h *DecisionHandler) convertParamToInput(formValues url.Values, w http.ResponseWriter, r *http.Request) (*decision.PublishAuthorizationCodeInput, error) {
	approved, err := strconv.ParseBool(formValues.Get("approved"))
	if err != nil {
		h.logger.Info("Invalid 'approved' parameter: %v", err)
		return nil, errors.New("無効なリクエストです。もう一度初めからやり直してください")
	}
	sessionID, err := r.Cookie(session.SessionIDCookieName)
	if err != nil {
		h.logger.Info("SessionID cookie not found: %v", err)
		return nil, errors.New("セッションが見つかりません。もう一度初めからやり直してください")
	}

	param, err := decision.NewPublishAuthorizationCodeInput(session.SessionID(sessionID.Value), formValues.Get("login_id"), formValues.Get("password"), approved)
	if err != nil {
		h.logger.Info("Failed to create PublishAuthorizationCodeParam: %v", err)
		return nil, errors.New("無効なリクエストです。もう一度初めからやり直してください")
	}

	return param, nil
}
