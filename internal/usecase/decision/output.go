package decision

type PublishAuthorizationCodeOutput struct {
	baseRedirectUri   string
	authorizationCode string
	state             string
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

func (e *ErrPublishAuthorizationCode) Error() string {
	return e.err.Error()
}
