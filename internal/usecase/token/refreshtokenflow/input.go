package refreshtokenflow

// refresh_token フロー用の入力
type RefreshTokenInput struct {
	clientID     string
	clientSecret string
	refreshToken string
	redirectURI  string
}

func NewRefreshTokenInput(clientID, clientSecret, refreshToken, redirectURI string) RefreshTokenInput {
	return RefreshTokenInput{
		clientID:     clientID,
		clientSecret: clientSecret,
		refreshToken: refreshToken,
		redirectURI:  redirectURI,
	}
}

func (i RefreshTokenInput) ClientID() string {
	return i.clientID
}
func (i RefreshTokenInput) ClientSecret() string {
	return i.clientSecret
}
func (i RefreshTokenInput) RefreshToken() string {
	return i.refreshToken
}
func (i RefreshTokenInput) RedirectURI() string {
	return i.redirectURI
}
