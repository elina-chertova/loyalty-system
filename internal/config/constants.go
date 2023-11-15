package config

import "time"

const (
	DSN             = "postgres://app:123qwe@localhost:5432/loyalty_db"
	DBName          = "loyalty_db"
	TableBalance    = "balances"
	TableWithdrawal = "withdrawals"
	TableOrder      = "orders"
	Processed       = "PROCESSED"
	TokenExp        = time.Minute * 10
)
