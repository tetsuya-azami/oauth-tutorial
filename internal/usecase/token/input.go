package token

type PublishTokenInput struct {
	clientID     string
	clientSecret string
	code         string
	redirectURI  string
}

func NewPublishTokenInput(clientID, clientSecret, grantType, code, redirectURI string) PublishTokenInput {
	return PublishTokenInput{
		clientID:     clientID,
		clientSecret: clientSecret,
		code:         code,
		redirectURI:  redirectURI,
	}
}

func (i PublishTokenInput) ClientID() string {
	return i.clientID
}
func (i PublishTokenInput) ClientSecret() string {
	return i.clientSecret
}
func (i PublishTokenInput) Code() string {
	return i.code
}
func (i PublishTokenInput) RedirectURI() string {
	return i.redirectURI
}
