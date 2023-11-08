package main

import (
	"github.com/elina-chertova/loyalty-system/internal/auth/handlers"
	"github.com/elina-chertova/loyalty-system/internal/auth/service"
	"github.com/elina-chertova/loyalty-system/internal/db"
	"github.com/gin-gonic/gin"
)

func main() {
	if err := run(); err != nil {
		panic(err)
	}
}

func run() error {
	router := gin.Default()
	dbConn := db.Init()

	userService := service.NewUserModel(dbConn)
	authenticator := service.NewUserAuth(userService)
	handler := handlers.NewAuthHandler(authenticator)

	router.POST("/api/user/register", handler.RegisterHandler())
	err := router.Run()
	if err != nil {
		return err
	}
	return nil
}
