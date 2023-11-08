package service

import (
	"fmt"
	"github.com/elina-chertova/loyalty-system/internal/config"
	"github.com/elina-chertova/loyalty-system/internal/db/user"
	"gorm.io/gorm"
)

type UserModel struct {
	DB *gorm.DB
}

func NewUserModel(db *gorm.DB) *UserModel {
	return &UserModel{DB: db}
}

type UserService interface {
	AddUser(string, string, bool) error
	GetUserByName(string) (user.User, error)
}

func (userDB *UserModel) AddUser(name string, password string, isAdmin bool) error {
	result := userDB.DB.Create(&user.User{Name: name, Password: password, IsAdmin: isAdmin})
	if result.Error != nil {
		return fmt.Errorf("%w: %v", config.ErrorCreatingUser, result.Error)
	}
	return nil
}

func (userDB *UserModel) GetUserByName(name string) (user.User, error) {
	var u user.User
	result := userDB.DB.Where(&user.User{Name: name}).First(&u)
	if result.Error != nil {
		return user.User{}, result.Error
	}
	return u, nil
}
