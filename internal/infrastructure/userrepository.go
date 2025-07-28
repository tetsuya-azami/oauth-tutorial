package infrastructure

import (
	"errors"
	"oauth-tutorial/internal/domain"
)

type LoginPasswordPair struct {
	LoginID  string
	Password string
}

var ErrUserNotFound = errors.New("user not found")

var users = map[LoginPasswordPair]*domain.User{
	{LoginID: "test-user@example.com", Password: "password"}: domain.ReconstructUser("IU7ewbuvey", "test-user@example.com", "password"),
}

type UserRepository struct {
}

func NewUserRepository() *UserRepository {
	return &UserRepository{}
}

func (r *UserRepository) SelectByLoginIDAndPassword(loginID, password string) (*domain.User, error) {
	user, ok := users[LoginPasswordPair{LoginID: loginID, Password: password}]
	if !ok {
		return nil, ErrUserNotFound
	}
	return user, nil
}
