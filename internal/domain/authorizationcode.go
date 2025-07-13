package domain

import (
	"time"
)

// TODO: PKCEあとで実装する
type AuthorizationCode struct {
	value       string
	userID      string
	clientID    string
	scopes      []string
	redirectURI string
	expiresAt   int64
}

const (
	AUTHORIZATION_CODE_DURATION = 10 * time.Minute
)

var (
	SUPPORTED_SCOPES = []string{"read", "write"}
)

type RandomGenerator interface {
	GenerateURLSafeRandomString(n int) string
}

func NewAuthorizationCode(randomGenerator RandomGenerator, userID string, clientID string, scopes []string, redirectURI string, now time.Time) *AuthorizationCode {
	expiresAt := now.Local().Add(AUTHORIZATION_CODE_DURATION).Unix()
	v := randomGenerator.GenerateURLSafeRandomString(32)
	return &AuthorizationCode{
		value:       v,
		userID:      userID,
		clientID:    clientID,
		scopes:      scopes,
		redirectURI: redirectURI,
		expiresAt:   expiresAt,
	}
}

func IsValidScopes(scopes []string) bool {
	for _, scope := range scopes {
		valid := false
		for _, supportedScope := range SUPPORTED_SCOPES {
			if scope == supportedScope {
				valid = true
				break
			}
		}
		if !valid {
			return false
		}
	}
	return true
}

func (a *AuthorizationCode) Value() string       { return a.value }
func (a *AuthorizationCode) UserID() string      { return a.userID }
func (a *AuthorizationCode) ClientID() string    { return a.clientID }
func (a *AuthorizationCode) RedirectURI() string { return a.redirectURI }
func (a *AuthorizationCode) ExpiresAt() int64    { return a.expiresAt }
func (a *AuthorizationCode) Scopes() []string    { return a.scopes }
