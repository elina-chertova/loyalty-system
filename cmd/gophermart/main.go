// @title Loyalty System
// @version 1.0
// @description Loyalty System description
// @host localhost:8081
// @BasePath /api/user
package main

import (
	"log"
	"os"
	"runtime"
	"time"

	_ "github.com/elina-chertova/loyalty-system/docs"
	"github.com/elina-chertova/loyalty-system/internal"
	"github.com/elina-chertova/loyalty-system/internal/auth/middleware"
	balService "github.com/elina-chertova/loyalty-system/internal/balance/service"
	"github.com/elina-chertova/loyalty-system/internal/config"
	"github.com/elina-chertova/loyalty-system/internal/db"
	handlersDB "github.com/elina-chertova/loyalty-system/internal/db/handlers"
	ordService "github.com/elina-chertova/loyalty-system/internal/order/service"
	"github.com/elina-chertova/loyalty-system/pkg/logger"
	"github.com/gin-contrib/gzip"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"go.uber.org/zap"
	p "net/http/pprof"
	"runtime/pprof"
)

func main() {
	f, err := os.Create("profiles/base.pprof")
	if err != nil {
		log.Fatal("could not create memory profile: ", err)
	}
	defer f.Close()

	if err := run(); err != nil {
		panic(err)
	}
	runtime.GC()
	if err := pprof.WriteHeapProfile(f); err != nil {
		log.Fatal("could not write memory profile: ", err)
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
	RegisterPprofRoutes(router)
	model := db.NewModels(dbConn)
	service := internal.NewServices(model)
	handler := internal.NewHandlers(service)

	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	router.GET("/api/user/ping", handlersDB.Ping(dbConn))

	router.POST("/api/user/register", handler.User.RegisterHandler())
	router.POST("/api/user/login", handler.User.LoginHandler())

	router.POST(
		"/api/user/orders",
		middleware.JWTAuth(),
		handler.Order.LoadOrderHandler(),
	)
	router.GET(
		"/api/user/orders",
		middleware.JWTAuth(),
		handler.Order.GetOrdersHandler(),
	)

	router.GET(
		"/api/user/balance",
		middleware.JWTAuth(),
		handler.Balance.GetBalanceHandler(),
	)

	router.POST(
		"/api/user/balance/withdraw",
		middleware.JWTAuth(),
		handler.Balance.RequestWithdrawFundsHandler(),
	)

	router.GET(
		"/api/user/withdrawals",
		middleware.JWTAuth(),
		handler.Balance.WithdrawalInfoHandler(),
	)

	go func() {
		for {
			updateOrderStatusLoop(service.Order, params.AccrualSystemAddress)
			updateBalanceLoop(service.Order, service.Balance)
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

func RegisterPprofRoutes(router *gin.Engine) {
	router.GET("/debug/pprof/", gin.WrapF(p.Index))
	router.GET("/debug/pprof/heap", gin.WrapF(p.Index))
	router.GET("/debug/pprof/goroutine", gin.WrapF(p.Index))
	router.GET("/debug/pprof/block", gin.WrapF(p.Index))
	router.GET("/debug/pprof/threadcreate", gin.WrapF(p.Index))
	router.GET("/debug/pprof/cmdline", gin.WrapF(p.Cmdline))
	router.GET("/debug/pprof/profile", gin.WrapF(p.Profile))
	router.GET("/debug/pprof/symbol", gin.WrapF(p.Symbol))
	router.GET("/debug/pprof/trace", gin.WrapF(p.Trace))
}
