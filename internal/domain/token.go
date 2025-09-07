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

type RefreshToken struct {
	value     string
	expiresAt int64
}

const (
	AccessTokenDuration  = 24 * time.Hour      // Access token valid for 24 hours
	RefreshTokenDuration = 60 * 24 * time.Hour // Refresh token valid for 60 days
)

func NewAccessToken(clientID, userID string, scopes []string, now time.Time) *AccessToken {
	// TODO: generatorのinjectの仕方考える
	g := mycrypto.RandomGenerator{}
	expiresAt := now.Local().Add(AccessTokenDuration).Unix()
	v := g.GenerateURLSafeRandomString(32)
	return &AccessToken{
		value:     v,
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

func NewRefreshToken(clientID string, now time.Time) *RefreshToken {
	// TODO: generatorのinjectの仕方考える
	g := mycrypto.RandomGenerator{}
	expiresAt := now.Local().Add(RefreshTokenDuration).Unix()
	v := g.GenerateURLSafeRandomString(32)
	return &RefreshToken{
		value:     v,
		expiresAt: expiresAt,
	}
}

func (t *RefreshToken) Value() string    { return t.value }
func (t *RefreshToken) ExpiresAt() int64 { return t.expiresAt }
