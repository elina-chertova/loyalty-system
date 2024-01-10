// Package service provides functionalities for user authentication
// and management in the loyalty system.
package service

import (
	"errors"
	"fmt"

	"github.com/elina-chertova/loyalty-system/internal/db/userdb"
	"github.com/elina-chertova/loyalty-system/internal/security"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

// UserAuth handles operations related to user authentication, such as
// registration, login, and token management.
type UserAuth struct {
	userRep userdb.UserRepository
}

// NewUserAuth creates a new instance of UserAuth with the given UserRepository.
func NewUserAuth(model userdb.UserRepository) *UserAuth {
	return &UserAuth{userRep: model}
}

// Predefined errors for user authentication operations.
var (
	ErrorCreatingUser  = errors.New("user cannot be created")
	ErrorAddingUser    = errors.New("user cannot be added")
	ErrorFindingUser   = errors.New("user not found")
	ErrorPasswordCheck = errors.New("password is wrong")
)

// Register handles the user registration process. It checks if a user already exists,
// hashes the password, and adds the new user to the repository.
func (u *UserAuth) Register(login, password string, isAdmin bool) error {
	_, err := u.userRep.GetUserByName(login)
	if !errors.Is(err, gorm.ErrRecordNotFound) {
		return ErrorCreatingUser
	}

	pwd, err := security.HashPassword(password)
	if err != nil {
		return err
	}

	err = u.userRep.AddUser(login, pwd, isAdmin)
	if err != nil {
		return fmt.Errorf("%w: %v", ErrorAddingUser, err.Error())
	}
	return nil
}

// Login verifies user credentials. It checks if the user exists and if the
// provided password matches the stored hash.
// Returns true if authentication is successful, false otherwise.
func (u *UserAuth) Login(login, password string) (bool, error) {
	user, err := u.userRep.GetUserByName(login)
	if err != nil {
		return false, fmt.Errorf("%w: %v", ErrorFindingUser, err.Error())
	}

	isEqual := security.CheckPasswordHash(password, user.Password)
	if !isEqual {
		return isEqual, ErrorPasswordCheck
	}
	return isEqual, nil
}

// SetToken generates a new JWT token for a given user login.
// Returns the user ID and the generated token.
func (u *UserAuth) SetToken(login string) (uuid.UUID, string, error) {
	user, err := u.userRep.GetUserByName(login)
	if err != nil {
		return uuid.Nil, "", err
	}
	token, err := security.GenerateToken(user.ID)
	if err != nil {
		return uuid.Nil, "", err
	}
	return user.ID, token, nil
}
