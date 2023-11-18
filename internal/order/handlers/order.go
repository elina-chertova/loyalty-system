package handlers

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/elina-chertova/loyalty-system/internal/auth/handlers"
	"github.com/elina-chertova/loyalty-system/internal/config"
	"github.com/elina-chertova/loyalty-system/internal/order/service"
	"github.com/elina-chertova/loyalty-system/pkg/logger"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"net/http"
)

type OrderService interface {
	LoadOrder(token string, orderID string, accrualServerAddress string) (int, error)
	GetOrders(token string) ([]service.UserOrderFormat, error)
}

type OrderHandler struct {
	Order OrderService
}

func NewOrderHandler(orderAuth OrderService) *OrderHandler {
	return &OrderHandler{Order: orderAuth}
}

// LoadOrderHandler @Load Order Number
// @Description Load Order Number
// @ID load-order
// @Tags Order
// @Accept json
// @Produce json
// @Param order_id path string true "Order ID"
// @Success 200 {object} Response
// @Success 202 {object} Response
// @Failure 400 {object} Response
// @Failure 401 {object} Response
// @Failure 409 {object} Response
// @Failure 422 {object} Response
// @Failure 500 {object} Response
// @Router /orders [post]
func (order *OrderHandler) LoadOrderHandler(accrualServerAddress string) gin.HandlerFunc {
	return func(c *gin.Context) {
		var orderNumber string
		if body, err := c.GetRawData(); err == nil {
			orderNumber = string(body)
		}
		token, exists := c.Get("token")
		if !exists {
			c.JSON(
				http.StatusUnauthorized,
				handlers.Response{
					Message: "Token not found",
					Status:  "Unauthorized",
				},
			)
			return
		}
		fmt.Println("orderNumber", orderNumber)
		tokenStr := fmt.Sprintf("%v", token)
		statusCode, err := order.Order.LoadOrder(tokenStr, orderNumber, accrualServerAddress)
		if err != nil {
			handleLoadOrderError(c, err)
			return
		}

		handleLoadOrderSuccess(c, statusCode)
	}
}

// GetOrdersHandler @Get User Orders
// @Description Get User Orders
// @ID get-orders
// @Tags Order
// @Accept json
// @Produce json
// @Success 200 {object} []service.UserOrderFormat
// @Success 204 {object} Response
// @Failure 401 {object} Response
// @Failure 500 {object} Response
// @Router /orders [get]
func (order *OrderHandler) GetOrdersHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		token, exists := c.Get("token")
		if !exists {
			c.JSON(
				http.StatusUnauthorized,
				handlers.Response{
					Message: "Token not found",
					Status:  "Unauthorized",
				},
			)
			return
		}

		tokenStr := fmt.Sprintf("%v", token)
		orders, err := order.Order.GetOrders(tokenStr)
		if err != nil {
			c.AbortWithStatusJSON(
				http.StatusInternalServerError, handlers.Response{
					Message: err.Error(),
					Status:  "Server error",
				},
			)
			return
		}

		userOrders, err := json.MarshalIndent(orders, "", "    ")
		fmt.Println("userOrders", string(userOrders))
		if err != nil {
			c.AbortWithStatusJSON(
				http.StatusInternalServerError, handlers.Response{
					Message: err.Error(),
					Status:  "Server error",
				},
			)
			return
		}
		c.Writer.Header().Set("Content-Type", "application/json")
		if len(orders) == 0 {
			c.Writer.WriteHeader(http.StatusNoContent)
			c.Writer.Write(userOrders)
		} else {
			c.Writer.WriteHeader(http.StatusOK)
			c.Writer.Write(userOrders)
		}

	}
}

func handleLoadOrderError(c *gin.Context, err error) {
	switch {
	case errors.Is(err, config.ErrorNotValidOrderNumber):
		logger.Logger.Error(
			"Order number is incorrect",
			zap.String("endpoint", c.Request.URL.Path),
			zap.Error(err),
		)
		c.JSON(
			http.StatusUnprocessableEntity, handlers.Response{
				Message: err.Error(),
				Status:  "Order number is incorrect",
			},
		)
		c.Abort()
	case errors.Is(err, config.ErrorOrderBelongsAnotherUser):
		logger.Logger.Error(
			"Conflict with other user",
			zap.String("endpoint", c.Request.URL.Path),
			zap.Error(err),
		)
		c.JSON(
			http.StatusConflict, handlers.Response{
				Message: err.Error(),
				Status:  "Conflict with other user",
			},
		)
		c.Abort()
	default:
		logger.Logger.Error(
			"Server error",
			zap.String("endpoint", c.Request.URL.Path),
			zap.Error(err),
		)
		c.AbortWithStatusJSON(
			http.StatusInternalServerError, handlers.Response{
				Message: err.Error(),
				Status:  "Server error",
			},
		)
	}
}

func handleLoadOrderSuccess(c *gin.Context, statusCode int) {
	var message, status string

	switch statusCode {
	case http.StatusOK:
		message, status = "Order is already loaded", "OK"
	case http.StatusAccepted:
		message, status = "Order is loaded", "OK"
	default:
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	c.IndentedJSON(
		statusCode, handlers.Response{
			Message: message,
			Status:  status,
		},
	)
}
