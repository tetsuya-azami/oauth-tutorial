package decision

type PublishAuthorizationCodeOutput struct {
	baseRedirectUri   string
	authorizationCode string
	state             string
}

func NewPublishAuthorizationCodeOutput(baseRedirectUri, authorizationCode, state string) PublishAuthorizationCodeOutput {
	return PublishAuthorizationCodeOutput{
		baseRedirectUri:   baseRedirectUri,
		authorizationCode: authorizationCode,
		state:             state,
	}
}

func (r *PublishAuthorizationCodeOutput) BaseRedirectUri() string {
	return r.baseRedirectUri
}

func (r *PublishAuthorizationCodeOutput) AuthorizationCode() string {
	return r.authorizationCode
}

func (r *PublishAuthorizationCodeOutput) State() string {
	return r.state
}

type ErrPublishAuthorizationCode struct {
	err             error
	baseRedirectUri string
}

func NewErrPublishAuthorizationCode(err error, baseRedirectUri string) *ErrPublishAuthorizationCode {
	return &ErrPublishAuthorizationCode{
		err:             err,
		baseRedirectUri: baseRedirectUri,
	}
}

func (e *ErrPublishAuthorizationCode) Error() string {
	return e.err.Error()
}

func (e *ErrPublishAuthorizationCode) Unwrap() error {
	return e.err
}
