package userdb

import (
	"testing"

	"gorm.io/gorm"
)

func TestUserModel_AddUser(t *testing.T) {
	type fields struct {
		DB *gorm.DB
	}
	type args struct {
		name     string
		password string
		isAdmin  bool
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(
			tt.name, func(t *testing.T) {
				userDB := &UserModel{
					DB: tt.fields.DB,
				}
				if err := userDB.AddUser(
					tt.args.name,
					tt.args.password,
					tt.args.isAdmin,
				); (err != nil) != tt.wantErr {
					t.Errorf("AddUser() error = %v, wantErr %v", err, tt.wantErr)
				}
			},
		)
	}
}
