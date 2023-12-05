package db

import (
	"github.com/elina-chertova/loyalty-system/internal/db/balancedb"
	"github.com/elina-chertova/loyalty-system/internal/db/orderdb"
	"github.com/elina-chertova/loyalty-system/internal/db/userdb"
	"gorm.io/gorm"
)

type Models struct {
	User    *userdb.UserModel
	Order   *orderdb.OrderModel
	Balance *balancedb.BalanceModel
}

func NewModels(conn *gorm.DB) *Models {
	return &Models{
		User:    userdb.NewUserModel(conn),
		Order:   orderdb.NewOrderModel(conn),
		Balance: balancedb.NewBalanceModel(conn),
	}
}
