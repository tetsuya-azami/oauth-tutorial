package decision

type PublishAuthorizationCodeResult struct {
	baseRedirectUri   string
	authorizationCode string
	state             string
}

func (r *PublishAuthorizationCodeResult) BaseRedirectUri() string {
	return r.baseRedirectUri
}

func (r *PublishAuthorizationCodeResult) AuthorizationCode() string {
	return r.authorizationCode
}

func (r *PublishAuthorizationCodeResult) State() string {
	return r.state
}

type ErrPublishAuthorizationCode struct {
	err             error
	baseRedirectUri string
}

func (e *ErrPublishAuthorizationCode) Error() string {
	return e.err.Error()
}
