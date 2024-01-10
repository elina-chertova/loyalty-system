package handlers

import (
	"bytes"
	"encoding/json"
	"fmt"
	authService "github.com/elina-chertova/loyalty-system/internal/auth/service"
	balService "github.com/elina-chertova/loyalty-system/internal/balance/service"
	"github.com/elina-chertova/loyalty-system/internal/config"
	"github.com/elina-chertova/loyalty-system/internal/db"
	"github.com/elina-chertova/loyalty-system/internal/db/balancedb"
	"github.com/elina-chertova/loyalty-system/internal/db/userdb"
	"github.com/gin-gonic/gin"
	"math/rand"
	"net/http"
	"net/http/httptest"
	"time"
)

func ExampleAuthHandler_registerHandler() {
	router := gin.Default()
	params := config.NewServer()

	conn := db.Init(params.DatabaseDSN)
	udb := userdb.NewUserModel(conn)
	bdb := balancedb.NewBalanceModel(conn)
	b := balService.NewBalance(bdb)
	u := authService.NewUserAuth(udb)

	userHandler := NewAuthHandler(b, u)
	router.POST("/api/user/register", userHandler.RegisterHandler())

	userCredentials := map[string]string{
		"login":    randString(10),
		"password": randString(10),
	}
	bodyBytes, _ := json.Marshal(userCredentials)

	req, _ := http.NewRequest("POST", "/api/user/register", bytes.NewBuffer(bodyBytes))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	fmt.Printf("Response code: %d\n", w.Code)
	if w.Result().Cookies() != nil {
		for _, cookie := range w.Result().Cookies() {
			fmt.Printf("Cookie set: %s\n", cookie.Name)
		}
	}
	if authHeader := w.Result().Header.Get("Authorization"); authHeader != "" {
		fmt.Println("Authorization header set")
	}

	// Output:
	// Response code: 200
	// Cookie set: access_token
	// Authorization header set
}

func ExampleAuthHandler_loginHandler() {
	router := gin.Default()
	params := config.NewServer()

	conn := db.Init(params.DatabaseDSN)
	udb := userdb.NewUserModel(conn)
	bdb := balancedb.NewBalanceModel(conn)
	b := balService.NewBalance(bdb)
	u := authService.NewUserAuth(udb)
	userHandler := NewAuthHandler(b, u)

	router.POST("/api/user/login", userHandler.LoginHandler())

	userCredentials := map[string]string{
		"login":    randString(10),
		"password": randString(10),
	}
	bodyBytes, _ := json.Marshal(userCredentials)

	req, _ := http.NewRequest("POST", "/api/user/login", bytes.NewBuffer(bodyBytes))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	fmt.Printf("Response code: %d\n", w.Code)
	if w.Result().Cookies() != nil {
		for _, cookie := range w.Result().Cookies() {
			fmt.Printf("Cookie set: %s\n", cookie.Name)
		}
	}
	if authHeader := w.Result().Header.Get("Authorization"); authHeader != "" {
		fmt.Println("Authorization header set")
	}

	// Output for a successful authentication:
	// Response code: 200
	// Cookie set: access_token
	// Authorization header set
	//
	// Output for a failed authentication due to incorrect credentials (example):
	// Response code: 401
	//
	// Output for a failed authentication due to bad request format (example):
	// Response code: 400
	//
	// Output for a server error (example):
	// Response code: 500
}

func randString(length int) string {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	var seededRand = rand.New(rand.NewSource(time.Now().UnixNano()))
	b := make([]byte, length)
	for i := range b {
		b[i] = charset[seededRand.Intn(len(charset))]
	}
	return string(b)
}
