package userdb

import (
	"gorm.io/gorm"
)

type UserModel struct {
	DB *gorm.DB
}

func NewUserModel(db *gorm.DB) *UserModel {
	return &UserModel{DB: db}
}

type UserRepository interface {
	AddUser(string, string, bool) error
	GetUserByName(string) (User, error)
}

func (userDB *UserModel) AddUser(name, password string, isAdmin bool) error {
	result := userDB.DB.Create(&User{Name: name, Password: password, IsAdmin: isAdmin})
	if result.Error != nil {
		return result.Error
	}
	return nil
}

func (userDB *UserModel) GetUserByName(name string) (User, error) {
	var u User
	result := userDB.DB.Where(&User{Name: name}).First(&u)
	if result.Error != nil {
		return User{}, result.Error
	}
	return u, nil
}
