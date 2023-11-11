package config

import "errors"

var (
	ErrorCreatingUser  = errors.New("user cannot be created")
	ErrorAddingUser    = errors.New("user cannot be added")
	ErrorFindingUser   = errors.New("user not found")
	ErrorPasswordCheck = errors.New("password is wrong")
	ErrorParseClaims   = errors.New("couldn't parse claims")
	ErrorTokenExpired  = errors.New("token expired")
)
