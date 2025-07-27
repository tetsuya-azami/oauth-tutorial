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

func (sd *SessionData) AuthParam() *domain.AuthorizationCodeFlowParam {
	return sd.authParam
}

func (sd *SessionData) User() *domain.User {
	return sd.user
}
