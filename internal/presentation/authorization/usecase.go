package authorize

import (
	"oauth-tutorial/internal/domain"
)

type IAuthorizationFlow interface {
	Execute(param *domain.AuthorizationCodeFlowParam) error
}
