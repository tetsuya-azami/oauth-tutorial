package domain

import (
	"errors"
	"oauth-tutorial/pkg/mylogger"
	"strings"
)

type AuthorizationCodeFlowParam struct {
	responseType ResponseType
	clientID     string
	redirectURI  string
	scopes       []string
	state        string
}

func NewAuthorizationCodeFlowParam(logger mylogger.Logger, responseType string, clientID string, redirectURI string, scope string, state string) (*AuthorizationCodeFlowParam, error) {
	// 認可フローによって処理が異なるケースを想定。現状、認可コードフローのみサポートのため、取得した値を使用してはいない
	rt, err := GetResponseType(responseType)
	if err != nil {
		logger.Info("Invalid response_type", "error", err)
		return &AuthorizationCodeFlowParam{}, err
	}

	if clientID == "" {
		logger.Info("client_id is empty")
		return &AuthorizationCodeFlowParam{}, errors.New("client_id is required")
	}

	if redirectURI == "" {
		logger.Info("redirect_uri is empty")
		return &AuthorizationCodeFlowParam{}, errors.New("redirect_uri is required")
	}

	if scope == "" {
		logger.Info("scope is empty")
		return &AuthorizationCodeFlowParam{}, errors.New("scope is required")
	}

	scopes := strings.Split(scope, " ")
	if !IsValidScopes(scopes) {
		logger.Info("Invalid scopes", "scopes", scopes, "supportedScopes", SUPPORTED_SCOPES)
		return &AuthorizationCodeFlowParam{}, errors.New("invalid scopes. Supported scopes are: " + strings.Join(SUPPORTED_SCOPES, ", "))
	}

	if state == "" {
		logger.Info("state is empty")
		return &AuthorizationCodeFlowParam{}, errors.New("state is required")
	}

	return &AuthorizationCodeFlowParam{
		responseType: rt,
		clientID:     clientID,
		redirectURI:  redirectURI,
		scopes:       scopes,
		state:        state,
	}, nil
}

func (p AuthorizationCodeFlowParam) ResponseType() ResponseType {
	return p.responseType
}

func (p AuthorizationCodeFlowParam) ClientID() string {
	return p.clientID
}

func (p AuthorizationCodeFlowParam) RedirectURI() string {
	return p.redirectURI
}

func (p AuthorizationCodeFlowParam) Scopes() []string {
	return p.scopes
}

func (p AuthorizationCodeFlowParam) State() string {
	return p.state
}
