package domain

import (
	"oauth-tutorial/pkg/mycrypto"
	"time"
)

type AccessToken struct {
	value     string
	clientID  string
	userID    string
	scopes    []string
	expiresAt int64
}

const (
	ACCESS_TOKEN_DURATION = 24 * time.Hour // Access token valid for 24 hours
)

func NewAccessToken(clientID, userID string, scopes []string, now time.Time) *AccessToken {
	g := mycrypto.RandomGenerator{}
	expiresAt := now.Local().Add(ACCESS_TOKEN_DURATION).Unix()
	at := g.GenerateURLSafeRandomString(32)
	return &AccessToken{
		value:     at,
		clientID:  clientID,
		userID:    userID,
		scopes:    scopes,
		expiresAt: expiresAt,
	}
}

func (t *AccessToken) Value() string    { return t.value }
func (t *AccessToken) ClientID() string { return t.clientID }
func (t *AccessToken) UserID() string   { return t.userID }
func (t *AccessToken) Scopes() []string { return t.scopes }
func (t *AccessToken) ExpiresAt() int64 { return t.expiresAt }
