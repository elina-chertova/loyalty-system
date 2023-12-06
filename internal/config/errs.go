package config

import "errors"

var (
	ErrorDownloadingBalance       = errors.New("balance cannot be created")
	ErrorDownloadingWithdrawFunds = errors.New("WithdrawFunds cannot be created")
	ErrorInsufficientFunds        = errors.New("insufficient funds")
)
