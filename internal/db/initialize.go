package db

import (
	"github.com/elina-chertova/loyalty-system/internal/db/balancedb"
	"github.com/elina-chertova/loyalty-system/internal/db/orderdb"
	"github.com/elina-chertova/loyalty-system/internal/db/userdb"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"log"
)

func Init(databaseDSN string) *gorm.DB {
	db, err := gorm.Open(postgres.Open(databaseDSN), &gorm.Config{})
	if err != nil {
		log.Fatalln(err)
	}

	err = db.AutoMigrate(
		&userdb.User{},
		&orderdb.Order{},
		&balancedb.Balance{},
		&balancedb.Withdrawal{},
	)
	if err != nil {
		log.Fatalln(err)
	}

	return db
}
