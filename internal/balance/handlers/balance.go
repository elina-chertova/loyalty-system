package handlers

import (
	"encoding/json"
	"fmt"
	"github.com/elina-chertova/loyalty-system/internal/auth/handlers"
	"github.com/elina-chertova/loyalty-system/internal/balance/service"
	"github.com/elina-chertova/loyalty-system/internal/config"
	"github.com/elina-chertova/loyalty-system/internal/logger"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"go.uber.org/zap"
	"net/http"
)

type BalanceService interface {
	GetBalance(token string) (service.UserBalanceFormat, error)
	WithdrawFunds(token, order string, sum float64) (int, error)
	WithdrawalInfo(token string) ([]service.WithdrawalFormat, error)
	AddInitialBalance(userID uuid.UUID) error
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
		respondWithError(c, http.StatusInternalServerError, "Server error", err)
		return
	}

	c.Writer.WriteHeader(statusCode)
	c.Writer.Write(result)
}

func (balance *BalanceHandler) WithdrawalInfoHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		token, exists := c.Get("token")
		if !exists {
			respondWithError(
				c,
				http.StatusUnauthorized,
				"Token not found",
				config.ErrorTokenNotFound,
			)
			return
		}

		tokenStr := fmt.Sprintf("%v", token)
		withdrawalInfo, err := balance.Balance.WithdrawalInfo(tokenStr)
		if len(withdrawalInfo) == 0 {
			c.Writer.WriteHeader(http.StatusNoContent)
			return
		}

		if err != nil {
			respondWithError(c, http.StatusInternalServerError, "Error with WithdrawalInfo", err)
			return
		}

		respondWithJSON(c, http.StatusOK, withdrawalInfo)
	}
}

func (balance *BalanceHandler) GetBalanceHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		token, exists := c.Get("token")
		if !exists {
			respondWithError(
				c,
				http.StatusUnauthorized,
				"Token not found",
				config.ErrorTokenNotFound,
			)
			return
		}

		tokenStr := fmt.Sprintf("%v", token)
		userBalance, err := balance.Balance.GetBalance(tokenStr)
		if err != nil {
			respondWithError(c, http.StatusInternalServerError, "Error with GetBalance", err)
			return
		}

		respondWithJSON(c, http.StatusOK, userBalance)
	}
}

func (balance *BalanceHandler) RequestWithdrawFundsHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		var w withdraw
		if err := c.BindJSON(&w); err != nil {
			respondWithError(c, http.StatusBadRequest, "Check json input", err)
			return
		}

		token, exists := c.Get("token")
		if !exists {
			respondWithError(
				c,
				http.StatusUnauthorized,
				"Token not found",
				config.ErrorTokenNotFound,
			)
			return
		}

		tokenStr := fmt.Sprintf("%v", token)
		statusCode, err := balance.Balance.WithdrawFunds(tokenStr, w.Order, w.Sum)
		if err != nil {
			respondWithError(c, statusCode, "error in WithdrawFunds", err)
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

func respondWithError(c *gin.Context, statusCode int, message string, err error) {
	logger.Logger.Error(
		message,
		zap.String("endpoint", c.Request.URL.Path),
		zap.Error(err),
	)
	c.AbortWithStatusJSON(
		statusCode, handlers.Response{
			Message: message,
			Status:  http.StatusText(statusCode),
		},
	)
}
