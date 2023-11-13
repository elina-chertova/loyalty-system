package handlers

import (
	"encoding/json"
	"fmt"
	"github.com/elina-chertova/loyalty-system/internal/auth/handlers"
	"github.com/elina-chertova/loyalty-system/internal/balance/service"
	"github.com/gin-gonic/gin"
	"net/http"
)

type BalanceService interface {
	GetBalance(token string) (service.UserBalanceFormat, error)
	WithdrawFunds(token, order string, sum float64) (int, error)
	WithdrawalInfo(token string) ([]service.WithdrawalFormat, error)
}

type BalanceHandler struct {
	Balance BalanceService
}

func NewBalanceHandler(b BalanceService) *BalanceHandler {
	return &BalanceHandler{Balance: b}
}

type withdraw struct {
	Order string  `json:"order"`
	Sum   float64 `json:"sum"`
}

func respondWithJSON(c *gin.Context, statusCode int, data interface{}) {
	result, err := json.MarshalIndent(data, "", "    ")
	if err != nil {
		respondWithError(c, http.StatusInternalServerError, "Server error")
		return
	}

	c.Writer.WriteHeader(statusCode)
	c.Writer.Write(result)
}

func (balance *BalanceHandler) WithdrawalInfoHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		token, exists := c.Get("token")
		if !exists {
			respondWithError(c, http.StatusUnauthorized, "Token not found")
			return
		}

		tokenStr := fmt.Sprintf("%v", token)
		withdrawalInfo, err := balance.Balance.WithdrawalInfo(tokenStr)
		if len(withdrawalInfo) == 0 {
			c.Writer.WriteHeader(http.StatusNoContent)
			return
		}

		if err != nil {
			respondWithError(c, http.StatusInternalServerError, err.Error())
			return
		}

		respondWithJSON(c, http.StatusOK, withdrawalInfo)
	}
}

func (balance *BalanceHandler) GetBalanceHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		token, exists := c.Get("token")
		if !exists {
			respondWithError(c, http.StatusUnauthorized, "Token not found")
			return
		}

		tokenStr := fmt.Sprintf("%v", token)
		userBalance, err := balance.Balance.GetBalance(tokenStr)
		if err != nil {
			respondWithError(c, http.StatusInternalServerError, err.Error())
			return
		}

		respondWithJSON(c, http.StatusOK, userBalance)
	}
}

func (balance *BalanceHandler) RequestWithdrawFundsHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		var w withdraw
		if err := c.BindJSON(&w); err != nil {
			respondWithError(c, http.StatusBadRequest, "Check json input")
			return
		}

		token, exists := c.Get("token")
		if !exists {
			respondWithError(c, http.StatusUnauthorized, "Token not found")
			return
		}

		tokenStr := fmt.Sprintf("%v", token)
		statusCode, err := balance.Balance.WithdrawFunds(tokenStr, w.Order, w.Sum)
		if err != nil {
			respondWithError(c, statusCode, err.Error())
			return
		}

		respondWithJSON(
			c, http.StatusOK, handlers.Response{
				Message: "Funds have been withdrawn",
				Status:  "OK",
			},
		)
	}
}

func respondWithError(c *gin.Context, statusCode int, message string) {
	c.AbortWithStatusJSON(
		statusCode, handlers.Response{
			Message: message,
			Status:  http.StatusText(statusCode),
		},
	)
}
