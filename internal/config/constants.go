package config

import "time"

const (
	TableBalance    = "balances"
	TableWithdrawal = "withdrawals"
	TableOrder      = "orders"
	Processed       = "PROCESSED"
	TokenExp        = time.Minute * 10
	UpdateInterval  = 1 * time.Second

	AccrualSystemAddress = "%s/api/orders/"
)
