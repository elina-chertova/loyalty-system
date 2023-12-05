package internal

import (
	authService "github.com/elina-chertova/loyalty-system/internal/auth/service"
	balService "github.com/elina-chertova/loyalty-system/internal/balance/service"
	"github.com/elina-chertova/loyalty-system/internal/db"
	ordService "github.com/elina-chertova/loyalty-system/internal/order/service"
)

type services struct {
	User    *authService.UserAuth
	Order   *ordService.UserOrder
	Balance *balService.UserBalance
}

func NewServices(s *db.Models) *services {
	return &services{
		User:    authService.NewUserAuth(s.User),
		Order:   ordService.NewOrder(s.Order),
		Balance: balService.NewBalance(s.Balance),
	}
}
