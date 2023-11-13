package main

import (
	handlersUser "github.com/elina-chertova/loyalty-system/internal/auth/handlers"
	"github.com/elina-chertova/loyalty-system/internal/auth/middleware"
	authService "github.com/elina-chertova/loyalty-system/internal/auth/service"
	handlersBal "github.com/elina-chertova/loyalty-system/internal/balance/handlers"
	balService "github.com/elina-chertova/loyalty-system/internal/balance/service"
	"github.com/elina-chertova/loyalty-system/internal/db"
	"github.com/elina-chertova/loyalty-system/internal/db/balancedb"
	"github.com/elina-chertova/loyalty-system/internal/db/orderdb"
	"github.com/elina-chertova/loyalty-system/internal/db/userdb"
	handlersOrd "github.com/elina-chertova/loyalty-system/internal/order/handlers"
	ordService "github.com/elina-chertova/loyalty-system/internal/order/service"
	"github.com/gin-gonic/gin"
	"time"
)

func main() {
	if err := run(); err != nil {
		panic(err)
	}
}

func run() error {
	router := gin.Default()
	dbConn := db.Init()

	userService := userdb.NewUserModel(dbConn)
	orderService := orderdb.NewOrderModel(dbConn)
	balanceService := balancedb.NewBalanceModel(dbConn)

	authenticator := authService.NewUserAuth(userService)
	order := ordService.NewOrder(orderService)
	balance := balService.NewBalance(balanceService)

	handlerAuth := handlersUser.NewAuthHandler(authenticator)
	handlerOrder := handlersOrd.NewOrderHandler(order)
	handlerBalance := handlersBal.NewBalanceHandler(balance)

	router.POST("/api/user/register", handlerAuth.RegisterHandler())
	router.POST("/api/user/login", handlerAuth.LoginHandler())

	router.POST(
		"/api/user/orders/:orderNumber",
		middleware.JWTAuth(),
		handlerOrder.LoadOrderHandler(),
	)
	router.GET(
		"/api/user/orders",
		middleware.JWTAuth(),
		handlerOrder.GetOrdersHandler(),
	)

	router.GET(
		"/api/user/balance",
		middleware.JWTAuth(),
		handlerBalance.GetBalanceHandler(),
	)

	router.POST(
		"/api/user/balance/withdraw",
		middleware.JWTAuth(),
		handlerBalance.RequestWithdrawFundsHandler(),
	)

	router.GET(
		"/api/user/withdrawals",
		middleware.JWTAuth(),
		handlerBalance.WithdrawalInfoHandler(),
	)

	go func() {
		for {
			err := balance.UpdateBalance(order)
			if err != nil {
				return
			}
			time.Sleep(5 * time.Second)
		}
	}()

	err := router.Run("localhost:8081")
	if err != nil {
		return err
	}
	return nil
}
