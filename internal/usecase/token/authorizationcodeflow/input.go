package authorizationcodeflow

// AuthorizationCodeInput は authorization_code フロー専用の入力
type AuthorizationCodeInput struct {
	clientID     string
	clientSecret string
	code         string
	redirectURI  string
}

func NewAuthorizationCodeInput(clientID, clientSecret, code, redirectURI string) AuthorizationCodeInput {
	return AuthorizationCodeInput{
		clientID:     clientID,
		clientSecret: clientSecret,
		code:         code,
		redirectURI:  redirectURI,
	}
}

func (i AuthorizationCodeInput) ClientID() string {
	return i.clientID
}
func (i AuthorizationCodeInput) ClientSecret() string {
	return i.clientSecret
}
func (i AuthorizationCodeInput) Code() string {
	return i.code
}
func (i AuthorizationCodeInput) RedirectURI() string {
	return i.redirectURI
}
