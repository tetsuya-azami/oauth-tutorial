package decision

import (
	"oauth-tutorial/internal/domain"
	inf_dto "oauth-tutorial/internal/infrastructure/dto"
	"oauth-tutorial/internal/session"
)

type IRandomCodeGenerator interface {
	GenerateURLSafeRandomString(n int) string
}

type ISessionStorage interface {
	Get(sessionID session.SessionID) (*inf_dto.SessionData, error)
}

type IUserRepository interface {
	SelectByLoginIDAndPassword(loginID, password string) (*domain.User, error)
}

type IAuthorizationCodeRepository interface {
	Save(code *domain.AuthorizationCode)
}
