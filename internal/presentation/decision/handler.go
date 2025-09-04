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
		h.logger.Info("Failed to parse form", "err", err)
		presentation.WriteJSONResponse(w, http.StatusBadRequest, ErrorResponse{Message: "パラメータの形式を確認してください"})
		return
	}

	input, err := h.convertParamToInput(r.Form, r)
	if err != nil {
		presentation.WriteJSONResponse(w, http.StatusBadRequest, ErrorResponse{Message: err.Error()})
		return
	}

	result, err := h.publishAuthorizationCode.Execute(input)
	if err != nil {
		var errPac *decision.ErrPublishAuthorizationCode
		if errors.As(err, &errPac) {
			switch {
			case errors.Is(errPac, decision.ErrSessionNotFound):
				// sessionが見つからない場合、redirectURIを取得できないため、JSONでエラーを返す
				presentation.WriteJSONResponse(w, http.StatusBadRequest, ErrorResponse{Message: err.Error()})
				return
			case errors.Is(errPac, decision.ErrUnexpectedSessionGetError):
				// session取得時に予期しないエラーが起きた場合、redirectURIを取得できないため、JSONでエラーを返す
				presentation.WriteJSONResponse(w, http.StatusBadRequest, ErrorResponse{Message: err.Error()})
				return
			case errors.Is(errPac, decision.ErrAuthorizationDenied):
				redirectUri := errPac.BaseRedirectUri() + "?error=access_denied&error_description=" + errPac.Error() + "&state=" + errPac.State()
				http.Redirect(w, r, redirectUri, http.StatusSeeOther)
				return
			case errors.Is(errPac, decision.ErrInvalidLoginCredentials):
				// クレデンシャルが異なる場合、リダイレクトせずにフロントでの再入力を促すためJSONでエラーを返す
				presentation.WriteJSONResponse(w, http.StatusUnauthorized, ErrorResponse{Message: errPac.Error()})
				return
			}
		}
	}

	redirectUri := result.BaseRedirectUri() + "?code=" + result.AuthorizationCode() + "&state=" + result.State()
	http.Redirect(w, r, redirectUri, http.StatusSeeOther)
}

func (h *DecisionHandler) convertParamToInput(formValues url.Values, r *http.Request) (*decision.PublishAuthorizationCodeInput, error) {
	approved, err := strconv.ParseBool(formValues.Get("approved"))
	if err != nil {
		h.logger.Info("Invalid 'approved' parameter", "err", err)
		return nil, errors.New("無効なリクエストです。もう一度初めからやり直してください")
	}
	sessionID, err := r.Cookie(session.SessionIDCookieName)
	if err != nil {
		h.logger.Info("SessionID cookie not found", "err", err)
		return nil, errors.New("セッションが見つかりません。もう一度初めからやり直してください")
	}

	input, err := decision.NewPublishAuthorizationCodeInput(session.SessionID(sessionID.Value), formValues.Get("login_id"), formValues.Get("password"), approved)
	if err != nil {
		h.logger.Info("Failed to create param", "err", err)
		return nil, errors.New("無効なリクエストです。もう一度初めからやり直してください")
	}

	return input, nil
}
