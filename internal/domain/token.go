package domain

import "oauth-tutorial/internal/crypt"

type AccessToken struct {
	accessToken string
	clientID    string
	userID      string
	scopes      []string
	expiresAt   int64
}

func NewAccessToken(clientID, userID string, scopes []string, expiresAt int64) *AccessToken {
	g := crypt.RandomGenerator{}
	at := g.GenerateURLSafeRandomString(32)
	return &AccessToken{
		accessToken: at,
		clientID:    clientID,
		userID:      userID,
		scopes:      scopes,
		expiresAt:   expiresAt,
	}
}

func (t *AccessToken) AccessToken() string { return t.accessToken }
func (t *AccessToken) ClientID() string    { return t.clientID }
func (t *AccessToken) UserID() string      { return t.userID }
func (t *AccessToken) Scopes() []string    { return t.scopes }
func (t *AccessToken) ExpiresAt() int64    { return t.expiresAt }
