package authorize

import (
	"oauth-tutorial/internal/domain"
	"oauth-tutorial/internal/session"
)

type IAuthorizationFlow interface {
	Execute(param *domain.AuthorizationCodeFlowParam) (session.SessionID, error)
}
