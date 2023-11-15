package config

import "time"

const (
	TableBalance    = "balances"
	TableWithdrawal = "withdrawals"
	TableOrder      = "orders"
	Processed       = "PROCESSED"
	TokenExp        = time.Minute * 10

	AccrualSystemAddress = "http://%s/api/orders/"
)
