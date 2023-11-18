// @title Loyalty System
// @version 1.0
// @description Loyalty System description
// @host localhost:8081
// @BasePath /api/user
package main

import (
	_ "github.com/elina-chertova/loyalty-system/docs"
	handlersUser "github.com/elina-chertova/loyalty-system/internal/auth/handlers"
	"github.com/elina-chertova/loyalty-system/internal/auth/middleware"
	authService "github.com/elina-chertova/loyalty-system/internal/auth/service"
	handlersBal "github.com/elina-chertova/loyalty-system/internal/balance/handlers"
	balService "github.com/elina-chertova/loyalty-system/internal/balance/service"
	"github.com/elina-chertova/loyalty-system/internal/config"
	"github.com/elina-chertova/loyalty-system/internal/db"
	"github.com/elina-chertova/loyalty-system/internal/db/balancedb"
	handlersDB "github.com/elina-chertova/loyalty-system/internal/db/handlers"
	"github.com/elina-chertova/loyalty-system/internal/db/orderdb"
	"github.com/elina-chertova/loyalty-system/internal/db/userdb"
	handlersOrd "github.com/elina-chertova/loyalty-system/internal/order/handlers"
	ordService "github.com/elina-chertova/loyalty-system/internal/order/service"
	"github.com/elina-chertova/loyalty-system/pkg/logger"
	"github.com/gin-contrib/gzip"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"go.uber.org/zap"
	"gorm.io/gorm"
	"time"
)

func main() {
	if err := run(); err != nil {
		panic(err)
	}
}

func run() error {
	params := config.NewServer()
	config.LoadEnv()
	dbConn := db.Init(params.DatabaseDSN)
	err := logger.InitLogger()
	if err != nil {
		panic(err)
	}

	router := routerInit()
	model := NewModels(dbConn)
	service := NewServices(model)
	handler := NewHandlers(service)

	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	router.GET("/api/user/ping", handlersDB.Ping(dbConn))

	router.POST("/api/user/register", handler.user.RegisterHandler())
	router.POST("/api/user/login", handler.user.LoginHandler())

	router.POST(
		"/api/user/orders",
		middleware.JWTAuth(),
		handler.order.LoadOrderHandler(params.AccrualSystemAddress),
	)
	router.GET(
		"/api/user/orders",
		middleware.JWTAuth(),
		handler.order.GetOrdersHandler(),
	)

	router.GET(
		"/api/user/balance",
		middleware.JWTAuth(),
		handler.balance.GetBalanceHandler(),
	)

	router.POST(
		"/api/user/balance/withdraw",
		middleware.JWTAuth(),
		handler.balance.RequestWithdrawFundsHandler(),
	)

	router.GET(
		"/api/user/withdrawals",
		middleware.JWTAuth(),
		handler.balance.WithdrawalInfoHandler(),
	)

	go func() {
		for {
			updateOrderStatusLoop(service.order, params.AccrualSystemAddress)
			updateBalanceLoop(service.order, service.balance)
			time.Sleep(config.UpdateInterval)
		}

	}()

	err = router.Run(params.Address)
	if err != nil {
		return err
	}
	return nil
}

func updateOrderStatusLoop(order *ordService.UserOrder, accrualServerAddress string) {
	err := order.UpdateOrderStatus(accrualServerAddress)
	if err != nil {
		logger.Logger.Warn("Order status has not been updated", zap.Error(err))
	}
}

func updateBalanceLoop(order *ordService.UserOrder, balance *balService.UserBalance) {
	err := balance.UpdateBalance(order)
	if err != nil {
		logger.Logger.Warn("Balance has not been updated", zap.Error(err))
		return
	}
}

func routerInit() *gin.Engine {
	router := gin.Default()
	router.Use(logger.GinLogger(logger.Logger))
	router.Use(gzip.Gzip(gzip.DefaultCompression))
	return router
}

type models struct {
	user    *userdb.UserModel
	order   *orderdb.OrderModel
	balance *balancedb.BalanceModel
}

func NewModels(conn *gorm.DB) *models {
	return &models{
		user:    userdb.NewUserModel(conn),
		order:   orderdb.NewOrderModel(conn),
		balance: balancedb.NewBalanceModel(conn),
	}
}

type services struct {
	user    *authService.UserAuth
	order   *ordService.UserOrder
	balance *balService.UserBalance
}

func NewServices(s *models) *services {
	return &services{
		user:    authService.NewUserAuth(s.user),
		order:   ordService.NewOrder(s.order),
		balance: balService.NewBalance(s.balance),
	}
}

type handlers struct {
	user    *handlersUser.AuthHandler
	order   *handlersOrd.OrderHandler
	balance *handlersBal.BalanceHandler
}

func NewHandlers(s *services) *handlers {
	return &handlers{
		user:    handlersUser.NewAuthHandler(s.balance, s.user),
		order:   handlersOrd.NewOrderHandler(s.order),
		balance: handlersBal.NewBalanceHandler(s.balance),
	}
}
