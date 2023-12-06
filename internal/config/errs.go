package config

import "errors"

var (
	ErrorParseClaims              = errors.New("couldn't parse claims")
	ErrorTokenExpired             = errors.New("token expired")
	ErrorTokenNotFound            = errors.New("token not found")
	ErrorDownloadingOrder         = errors.New("order cannot be created")
	ErrorAddingOrder              = errors.New("order cannot be added")
	ErrorOrderBelongsAnotherUser  = errors.New("order belongs to another user")
	ErrorNotValidOrderNumber      = errors.New("order number is not valid")
	ErrorSystem                   = errors.New("error in loyality system")
	ErrorDownloadingBalance       = errors.New("balance cannot be created")
	ErrorDownloadingWithdrawFunds = errors.New("WithdrawFunds cannot be created")
	ErrorInsufficientFunds        = errors.New("insufficient funds")
)
