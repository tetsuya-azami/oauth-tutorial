package dto

import (
	"errors"
	"net/url"
)

// AuthorizeRequest は認可リクエストのDTOです。
type AuthorizeRequest struct {
	clientID     string
	redirectURI  string
	responseType string
	state        string
}

// NewAuthorizeRequestはパラメータのバリデーションを行い、AuthorizeRequestを生成します。
func NewAuthorizeRequest(values url.Values) (*AuthorizeRequest, error) {
	clientID := values.Get("client_id")
	redirectURI := values.Get("redirect_uri")
	responseType := values.Get("response_type")
	state := values.Get("state")

	if clientID == "" || redirectURI == "" || responseType != "code" {
		return nil, errors.New("invalid request parameters")
	}

	return &AuthorizeRequest{
		clientID:     clientID,
		redirectURI:  redirectURI,
		responseType: responseType,
		state:        state,
	}, nil
}

// ClientID は client_id を返します。
func (a *AuthorizeRequest) ClientID() string {
	return a.clientID
}

// RedirectURI は redirect_uri を返します。
func (a *AuthorizeRequest) RedirectURI() string {
	return a.redirectURI
}

// ResponseType は response_type を返します。
func (a *AuthorizeRequest) ResponseType() string {
	return a.responseType
}

// State は state を返します。
func (a *AuthorizeRequest) State() string {
	return a.state
}
