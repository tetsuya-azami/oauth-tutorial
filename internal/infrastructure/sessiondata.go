package infrastructure

import "oauth-tutorial/internal/domain"

type SessionData struct {
	authParam *domain.AuthorizationCodeFlowParam
	user      *domain.User
}

func NewSessionData(authParam *domain.AuthorizationCodeFlowParam, user *domain.User) *SessionData {
	return &SessionData{
		authParam: authParam,
		user:      user,
	}
}
