package decision

import (
	"oauth-tutorial/internal/domain"
	"oauth-tutorial/internal/infrastructure"
	"oauth-tutorial/internal/session"
)

type IRandomCodeGenerator interface {
	GenerateURLSafeRandomString(n int) string
}

type ISessionStorage interface {
	Get(sessionID session.SessionID) (*infrastructure.SessionData, error)
}

type IUserRepository interface {
	SelectByLoginIDAndPassword(loginID, password string) (*domain.User, error)
}

type IAuthorizationCodeRepository interface {
	Save(code *domain.AuthorizationCode)
}
