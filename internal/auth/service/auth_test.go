package service

import (
	"errors"
	"fmt"
	"github.com/elina-chertova/loyalty-system/internal/config"
	"github.com/elina-chertova/loyalty-system/internal/db/userdb"
	"github.com/elina-chertova/loyalty-system/internal/security"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
	"testing"
)

type MockUserRepository struct{}

func (m *MockUserRepository) GetUserByName(login string) (userdb.User, error) {
	if login == "existingUser" {
		uuidID, err := uuid.Parse("69359037-9599-48e7-b8f2-48393c019135")
		if err != nil {
			return userdb.User{}, err
		}

		pass, _ := security.HashPassword("hashedPassword")
		return userdb.User{
			ID:       uuidID,
			Name:     "existingUser",
			Password: pass,
		}, nil
	}
	return userdb.User{}, gorm.ErrRecordNotFound
}

func (m *MockUserRepository) AddUser(login, password string, isAdmin bool) error {
	if login == "existingUser" {
		return errors.New("user already exists")
	}
	return nil
}

func TestUserAuth_Register(t *testing.T) {

	type args struct {
		login    string
		password string
		isAdmin  bool
	}
	tests := []struct {
		name    string
		args    args
		wantErr error
	}{
		{
			name:    "User are registered",
			args:    args{login: "name", password: "password", isAdmin: false},
			wantErr: nil,
		},
		{
			name:    "User are exists",
			args:    args{login: "existingUser", password: "hashedPassword", isAdmin: false},
			wantErr: config.ErrorCreatingUser,
		},
	}

	rep := &MockUserRepository{}
	userAuth := NewUserAuth(rep)

	for _, tt := range tests {
		t.Run(
			tt.name, func(t *testing.T) {
				if err := userAuth.Register(
					tt.args.login,
					tt.args.password,
					false,
				); err != tt.wantErr {
					assert.Equal(t, err, tt.wantErr)
				}
			},
		)
	}

}

func TestUserAuth_Login(t *testing.T) {
	type args struct {
		login    string
		password string
	}
	tests := []struct {
		name    string
		args    args
		wantErr error
	}{
		{
			name:    "User are login",
			args:    args{login: "name", password: "password"},
			wantErr: fmt.Errorf("%w: %v", config.ErrorFindingUser, gorm.ErrRecordNotFound),
		},
		{
			name:    "User are exists",
			args:    args{login: "existingUser", password: "hashedPassword"},
			wantErr: nil,
		},
		{
			name:    "Password is wrong",
			args:    args{login: "existingUser", password: "hashedPassword123"},
			wantErr: config.ErrorPasswordCheck,
		},
	}

	rep := &MockUserRepository{}
	userAuth := NewUserAuth(rep)

	for _, tt := range tests {
		t.Run(
			tt.name, func(t *testing.T) {
				if _, err := userAuth.Login(
					tt.args.login,
					tt.args.password,
				); err != tt.wantErr {
					assert.Equal(
						t,
						err,
						tt.wantErr,
					)
				}
			},
		)
	}

}

func TestUserAuth_SetToken(t *testing.T) {
	type args struct {
		login string
	}
	tests := []struct {
		name    string
		args    args
		wantErr error
	}{
		{
			name:    "User are login",
			args:    args{login: "name"},
			wantErr: gorm.ErrRecordNotFound,
		},
		{
			name:    "User are exists",
			args:    args{login: "existingUser"},
			wantErr: config.ErrorCreatingUser,
		},
	}

	rep := &MockUserRepository{}
	userAuth := NewUserAuth(rep)

	for _, tt := range tests {
		t.Run(
			tt.name, func(t *testing.T) {
				if _, _, err := userAuth.SetToken(tt.args.login); err != tt.wantErr {
					assert.Equal(t, config.ErrorCreatingUser, tt.wantErr)
				}
			},
		)
	}

}
