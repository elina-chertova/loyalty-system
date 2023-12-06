package service

import (
	"errors"
	"fmt"
	"github.com/elina-chertova/loyalty-system/internal/config"
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

func (u *UserAuth) Register(login, password string, isAdmin bool) error {
	_, err := u.userRep.GetUserByName(login)
	if !errors.Is(err, gorm.ErrRecordNotFound) {
		return config.ErrorCreatingUser
	}

	pwd, err := security.HashPassword(password)
	if err != nil {
		return err
	}

	err = u.userRep.AddUser(login, pwd, isAdmin)
	if err != nil {
		return fmt.Errorf("%w: %v", config.ErrorAddingUser, err.Error())
	}
	return nil
}

func (u *UserAuth) Login(login, password string) (bool, error) {
	user, err := u.userRep.GetUserByName(login)
	if err != nil {
		return false, fmt.Errorf("%w: %v", config.ErrorFindingUser, err.Error())
	}

	isEqual := security.CheckPasswordHash(password, user.Password)
	if !isEqual {
		return isEqual, config.ErrorPasswordCheck
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
