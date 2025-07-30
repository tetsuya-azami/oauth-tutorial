package decision

import (
	"errors"
	"oauth-tutorial/internal/domain"
	"oauth-tutorial/internal/infrastructure"
	"oauth-tutorial/pkg/mylogger"
	"time"
)

var (
	ErrSessionNotFound         = errors.New("session not found")
	ErrUnexpectedError         = errors.New("unexpected error occurred")
	ErrAuthorizationDenied     = errors.New("authorization denied by user")
	ErrInvalidLoginCredentials = errors.New("invalid login credentials")
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

func (uc *PublishAuthorizationCodeUseCase) Execute(param *PublishAuthorizationCodeParam) (PublishAuthorizationCodeResult, *ErrPublishAuthorizationCode) {
	session, err := uc.sessionStore.Get(param.sessionId)
	switch {
	case errors.Is(err, infrastructure.ErrSessionNotFound):
		uc.logger.Info("Session not found", err)
		return PublishAuthorizationCodeResult{}, &ErrPublishAuthorizationCode{
			err:             ErrSessionNotFound,
			baseRedirectUri: "",
		}
	case err != nil:
		uc.logger.Error("Unexpected error occurred", err)
		return PublishAuthorizationCodeResult{}, &ErrPublishAuthorizationCode{
			err:             ErrUnexpectedError,
			baseRedirectUri: "",
		}
	}
	if session == nil {
		uc.logger.Info("Session is nil")
		return PublishAuthorizationCodeResult{}, &ErrPublishAuthorizationCode{
			err:             ErrSessionNotFound,
			baseRedirectUri: "",
		}
	}

	if !param.approved {
		uc.logger.Info("Authorization denied by user")
		return PublishAuthorizationCodeResult{}, &ErrPublishAuthorizationCode{
			err:             ErrAuthorizationDenied,
			baseRedirectUri: session.AuthParam().RedirectURI(),
		}
	}

	user, err := uc.userRepository.SelectByLoginIDAndPassword(param.loginID, param.password)
	if err != nil {
		uc.logger.Info("Failed to select user by loginID and password", err)
		return PublishAuthorizationCodeResult{}, &ErrPublishAuthorizationCode{
			err:             ErrInvalidLoginCredentials,
			baseRedirectUri: session.AuthParam().RedirectURI(),
		}
	}
	if user == nil {
		uc.logger.Info("User not found for loginID: %s", param.loginID)
		return PublishAuthorizationCodeResult{}, &ErrPublishAuthorizationCode{
			err:             ErrInvalidLoginCredentials,
			baseRedirectUri: session.AuthParam().RedirectURI(),
		}
	}

	authorizationCode := domain.NewAuthorizationCode(uc.randomCodeGenerator, user.UserID(), session.AuthParam().ClientID(), session.AuthParam().Scopes(), session.AuthParam().RedirectURI(), time.Now())

	uc.authCodeRepository.Save(authorizationCode)

	return PublishAuthorizationCodeResult{
		baseRedirectUri:   session.AuthParam().RedirectURI(),
		authorizationCode: authorizationCode.Value(),
		state:             session.AuthParam().State(),
	}, nil
}
