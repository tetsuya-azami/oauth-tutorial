package decision

import (
	"oauth-tutorial/internal/usecase/decision"
)

type IPublishAuthorizationCodeUseCase interface {
	Execute(param *decision.PublishAuthorizationCodeInput) (decision.PublishAuthorizationCodeOutput, *decision.ErrPublishAuthorizationCode)
}
