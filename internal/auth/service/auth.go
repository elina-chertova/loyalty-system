package service

import (
	"errors"
	"fmt"
	"github.com/elina-chertova/loyalty-system/internal/db/userdb"
	"github.com/elina-chertova/loyalty-system/internal/security"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type UserAuth struct {
	userRep userdb.UserRepository
}

func NewUserAuth(model userdb.UserRepository) *UserAuth {
	return &UserAuth{userRep: model}
}

var (
	ErrorCreatingUser  = errors.New("user cannot be created")
	ErrorAddingUser    = errors.New("user cannot be added")
	ErrorFindingUser   = errors.New("user not found")
	ErrorPasswordCheck = errors.New("password is wrong")
)

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
