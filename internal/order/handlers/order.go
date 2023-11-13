package handlers

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/elina-chertova/loyalty-system/internal/auth/handlers"
	"github.com/elina-chertova/loyalty-system/internal/config"
	"github.com/elina-chertova/loyalty-system/internal/order/service"
	"github.com/gin-gonic/gin"
	"net/http"
)

type OrderService interface {
	LoadOrder(token string, orderID string) (int, error)
	GetOrders(token string) ([]service.UserOrderFormat, error)
}

type OrderHandler struct {
	Order OrderService
}

func NewOrderHandler(orderAuth OrderService) *OrderHandler {
	return &OrderHandler{Order: orderAuth}
}

func (order *OrderHandler) LoadOrderHandler() gin.HandlerFunc {
	return func(c *gin.Context) {

		if err := c.Request.ParseForm(); err != nil {
			c.Status(http.StatusBadRequest)
			c.Abort()
			return
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
		orderNumber := c.Param("orderNumber")

		tokenStr := fmt.Sprintf("%v", token)
		statusCode, err := order.Order.LoadOrder(tokenStr, orderNumber)
		if err != nil {
			handleLoadOrderError(c, err)
			return
		}

		handleLoadOrderSuccess(c, statusCode)
	}
}

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
		if err != nil {
			c.AbortWithStatusJSON(
				http.StatusInternalServerError, handlers.Response{
					Message: err.Error(),
					Status:  "Server error",
				},
			)
			return
		}

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
		c.JSON(
			http.StatusUnprocessableEntity, handlers.Response{
				Message: err.Error(),
				Status:  "Order number is incorrect",
			},
		)
		c.Abort()
	case errors.Is(err, config.ErrorOrderBelongsAnotherUser):
		c.JSON(
			http.StatusConflict, handlers.Response{
				Message: err.Error(),
				Status:  "Conflict with other user",
			},
		)
		c.Abort()
	default:
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
