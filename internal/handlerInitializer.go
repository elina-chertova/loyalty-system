package internal

import (
	handlersUser "github.com/elina-chertova/loyalty-system/internal/auth/handlers"
	handlersBal "github.com/elina-chertova/loyalty-system/internal/balance/handlers"
	handlersOrd "github.com/elina-chertova/loyalty-system/internal/order/handlers"
)

type handlers struct {
	User    *handlersUser.AuthHandler
	Order   *handlersOrd.OrderHandler
	Balance *handlersBal.BalanceHandler
}

func NewHandlers(s *services) *handlers {
	return &handlers{
		User:    handlersUser.NewAuthHandler(s.Balance, s.User),
		Order:   handlersOrd.NewOrderHandler(s.Order),
		Balance: handlersBal.NewBalanceHandler(s.Balance),
	}
}
