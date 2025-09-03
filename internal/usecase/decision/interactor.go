package decision

import (
	"errors"
	"oauth-tutorial/internal/domain"
	"oauth-tutorial/internal/infrastructure"
	"oauth-tutorial/pkg/mylogger"
	"time"
)

var (
	ErrSessionNotFound           = errors.New("session not found")
	ErrUnexpectedSessionGetError = errors.New("unexpected error occurred while getting session")
	ErrAuthorizationDenied       = errors.New("authorization denied by user")
	ErrInvalidLoginCredentials   = errors.New("invalid login credentials")
)

type PublishAuthorizationCodeUseCase struct {
	logger              mylogger.Logger
	randomCodeGenerator IRandomCodeGenerator
	sessionStore        ISessionStorage
	userRepository      IUserRepository
	authCodeRepository  IAuthorizationCodeRepository
}

func NewPublishAuthorizationCodeUseCase(logger mylogger.Logger, randomCodeGenerator IRandomCodeGenerator, sessionStore ISessionStorage, userRepository IUserRepository, authCodeRepository IAuthorizationCodeRepository) *PublishAuthorizationCodeUseCase {
	return &PublishAuthorizationCodeUseCase{
		logger:              logger,
		randomCodeGenerator: randomCodeGenerator,
		sessionStore:        sessionStore,
		userRepository:      userRepository,
		authCodeRepository:  authCodeRepository,
	}
}

func (uc *PublishAuthorizationCodeUseCase) Execute(param *PublishAuthorizationCodeInput) (PublishAuthorizationCodeOutput, error) {
	session, err := uc.sessionStore.Get(param.sessionId)
	switch {
	case errors.Is(err, infrastructure.ErrSessionNotFound):
		uc.logger.Info("Session not found", err)
		return PublishAuthorizationCodeOutput{}, &ErrPublishAuthorizationCode{
			err:             ErrSessionNotFound,
			baseRedirectUri: "",
			state:           "",
		}
	case err != nil:
		uc.logger.Error("Unexpected error occurred", err)
		return PublishAuthorizationCodeOutput{}, &ErrPublishAuthorizationCode{
			err:             ErrUnexpectedSessionGetError,
			baseRedirectUri: "",
			state:           "",
		}
	}
	if session == nil {
		uc.logger.Info("Session is nil")
		return PublishAuthorizationCodeOutput{}, &ErrPublishAuthorizationCode{
			err:             ErrSessionNotFound,
			baseRedirectUri: "",
			state:           "",
		}
	}

	if !param.approved {
		uc.logger.Info("Authorization denied by user")
		return PublishAuthorizationCodeOutput{}, &ErrPublishAuthorizationCode{
			err:             ErrAuthorizationDenied,
			baseRedirectUri: session.AuthParam().RedirectURI(),
			state:           session.AuthParam().State(),
		}
	}

	user, err := uc.userRepository.SelectByLoginIDAndPassword(param.loginID, param.password)
	if err != nil {
		uc.logger.Info("Failed to select user by loginID and password", err)
		return PublishAuthorizationCodeOutput{}, &ErrPublishAuthorizationCode{
			err:             ErrInvalidLoginCredentials,
			baseRedirectUri: "",
			state:           "",
		}
	}
	if user == nil {
		uc.logger.Info("Failed to select user by loginID and password", err)
		return PublishAuthorizationCodeOutput{}, &ErrPublishAuthorizationCode{
			err:             ErrInvalidLoginCredentials,
			baseRedirectUri: "",
			state:           "",
		}
	}

	authorizationCode := domain.NewAuthorizationCode(uc.randomCodeGenerator, user.UserID(), session.AuthParam().ClientID(), session.AuthParam().Scopes(), session.AuthParam().RedirectURI(), time.Now())

	uc.authCodeRepository.Save(authorizationCode)

	return NewPublishAuthorizationCodeOutput(
		session.AuthParam().RedirectURI(),
		authorizationCode.Value(),
		session.AuthParam().State(),
	), nil
}
