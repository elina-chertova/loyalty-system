package service

import (
	"errors"
	"github.com/elina-chertova/loyalty-system/internal/config"
	"github.com/elina-chertova/loyalty-system/internal/security"
	"gorm.io/gorm"
)

type UserAuth struct {
	U UserService
}

func NewUserAuth(model *UserModel) *UserAuth {
	return &UserAuth{U: model}
}

func (u *UserAuth) Register(login, password string, isAdmin bool) error {
	_, err := u.U.GetUserByName(login)
	if !errors.Is(err, gorm.ErrRecordNotFound) {
		return config.ErrorCreatingUser
	}

	pwd, err := security.HashPassword(password)
	if err != nil {
		return err
	}

	err = u.U.AddUser(login, pwd, isAdmin)
	if err != nil {
		return err
	}
	return nil
}
