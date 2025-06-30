package domain

import "time"

type AuthorizationCode struct {
	value       string
	userID      string
	clientID    string
	scope       string
	redirectURI string
	expiresAt   int64
}

func NewAuthorizationCode(value, userID, clientID, scope, redirectURI string, now time.Time) *AuthorizationCode {
	expiresAt := now.Local().Add(10 * time.Minute).Unix() // デフォルトロケールの時間で10分後に設定
	return &AuthorizationCode{
		value:       value,
		userID:      userID,
		clientID:    clientID,
		scope:       scope,
		redirectURI: redirectURI,
		expiresAt:   expiresAt,
	}
}

func (a *AuthorizationCode) Value() string       { return a.value }
func (a *AuthorizationCode) UserID() string      { return a.userID }
func (a *AuthorizationCode) ClientID() string    { return a.clientID }
func (a *AuthorizationCode) RedirectURI() string { return a.redirectURI }
func (a *AuthorizationCode) ExpiresAt() int64    { return a.expiresAt }
func (a *AuthorizationCode) Scope() string       { return a.scope }
