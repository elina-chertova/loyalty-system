// Package userdb provides data access functionalities for user management
// in the loyalty system. It uses GORM for database operations related to users.
package userdb

import (
	"errors"
	"fmt"

	"gorm.io/gorm"
)

// UserModel represents the model for user data and provides methods
// for interacting with the users table in the database.
type UserModel struct {
	DB *gorm.DB
}

// NewUserModel creates a new instance of UserModel with the given GORM DB instance.
func NewUserModel(db *gorm.DB) *UserModel {
	return &UserModel{DB: db}
}

// UserRepository defines the interface for user data operations. It abstracts
// the methods to interact with the users in the database.
type UserRepository interface {
	AddUser(string, string, bool) error
	GetUserByName(string) (User, error)
}

// AddUser adds a new user to the database with the provided name, password, and admin status.
// Returns an error if the user cannot be created.
func (userDB *UserModel) AddUser(name, password string, isAdmin bool) error {
	result := userDB.DB.Create(&User{Name: name, Password: password, IsAdmin: isAdmin})
	if result.Error != nil {
		return fmt.Errorf("%w: %v", errors.New("error during creating user"), result.Error)
	}
	return nil
}

// GetUserByName retrieves a user by their name from the database.
// Returns the User object and an error if the user is not found.
func (userDB *UserModel) GetUserByName(name string) (User, error) {
	var u User
	result := userDB.DB.Where(&User{Name: name}).First(&u)
	if result.Error != nil {
		return User{}, result.Error
	}
	return u, nil
}
