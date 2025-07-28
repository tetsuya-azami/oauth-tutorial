package decision

import (
	"errors"
	"oauth-tutorial/internal/session"
)

type PublishAuthorizationCodeParam struct {
	sessionId session.SessionID
	loginID   string
	password  string
	approved  bool
}

var (
	ErrEmptySessionID = errors.New("session ID cannot be empty")
	ErrEmptyLoginID   = errors.New("login ID cannot be empty")
	ErrEmptyPassword  = errors.New("password cannot be empty")
)

func NewPublishAuthorizationCodeParam(sessionId session.SessionID, loginID, password string, approved bool) (*PublishAuthorizationCodeParam, error) {
	if sessionId == "" {
		return nil, ErrEmptySessionID
	}
	if loginID == "" {
		return nil, ErrEmptyLoginID
	}
	if password == "" {
		return nil, ErrEmptyPassword
	}

	return &PublishAuthorizationCodeParam{sessionId: sessionId, approved: approved}, nil
}

func (p *PublishAuthorizationCodeParam) Approved() bool {
	return p.approved
}
