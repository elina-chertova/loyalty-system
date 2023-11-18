package config

import (
	"flag"
	"os"
)

type Settings struct {
	Address              string
	DatabaseDSN          string
	AccrualSystemAddress string
}

func ParseServerFlags(s *Settings) {
	flag.StringVar(&s.Address, "a", "", "address and port to run server")
	flag.StringVar(
		&s.AccrualSystemAddress,
		"r",
		"",
		"address and port to run server Accrual System",
	)

	flag.StringVar(
		&s.DatabaseDSN,
		"d",
		"postgres://app:123qwe@localhost:5432/loyalty_db",
		"Database DSN. Ex: postgres://app:123qwe@localhost:5432/loyalty_db",
	)
	flag.Parse()
	if envRunAddr := os.Getenv("RUN_ADDRESS"); envRunAddr != "" {
		s.Address = envRunAddr
	}
	if envRunAccrual := os.Getenv("ACCRUAL_SYSTEM_ADDRESS"); envRunAccrual != "" {
		s.AccrualSystemAddress = envRunAccrual
	}
	if envDSN := os.Getenv("DATABASE_URI"); envDSN != "" {
		s.DatabaseDSN = envDSN
	}

}

func NewServer() *Settings {
	s := &Settings{}
	ParseServerFlags(s)
	return s
}
