package decision

import (
	"errors"
	"oauth-tutorial/internal/session"
)

type PublishAuthorizationCodeInput struct {
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

func NewPublishAuthorizationCodeInput(sessionId session.SessionID, loginID, password string, approved bool) (*PublishAuthorizationCodeInput, error) {
	if sessionId == "" {
		return nil, ErrEmptySessionID
	}
	if loginID == "" {
		return nil, ErrEmptyLoginID
	}
	if password == "" {
		return nil, ErrEmptyPassword
	}

	return &PublishAuthorizationCodeInput{sessionId: sessionId, loginID: loginID, password: password, approved: approved}, nil
}

func (p *PublishAuthorizationCodeInput) Approved() bool {
	return p.approved
}
