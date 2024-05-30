// Package service provides functionalities for managing user orders
// in the loyalty system.
package service

import (
	"errors"
	"fmt"
	"time"

	"github.com/elina-chertova/loyalty-system/internal/db/orderdb"
	"github.com/elina-chertova/loyalty-system/internal/order/utils"
	"github.com/elina-chertova/loyalty-system/internal/security"
	"gorm.io/gorm"
)

// UserOrder struct handles operations related to user orders.
type UserOrder struct {
	OrderRep orderdb.OrderRepository
}

// NewOrder creates a new instance of UserOrder with the given OrderRepository.
func NewOrder(model orderdb.OrderRepository) *UserOrder {
	return &UserOrder{OrderRep: model}
}

// Predefined errors for order operations.
var (
	ErrorAddingOrder             = errors.New("order cannot be added")
	ErrorOrderBelongsAnotherUser = errors.New("order belongs to another user")
	ErrorNotValidOrderNumber     = errors.New("order number is not valid")
)

// LoadOrder handles the loading of an order. It checks the validity of the
// order number, associates it with a user, and updates the order status.
// It returns a LoadOrderResult indicating the outcome.
func (ord *UserOrder) LoadOrder(token string, orderID string) (*LoadOrderResult, error) {
	if !utils.IsLuhnValid(orderID) {
		return nil, ErrorNotValidOrderNumber
	}
	userID, err := security.GetUserIDFromToken(token)
	if err != nil {
		return nil, err
	}

	order, err := ord.OrderRep.GetOrderByID(orderID)
	if !errors.Is(err, gorm.ErrRecordNotFound) && err != nil {
		return nil, err
	}

	if (order == orderdb.Order{}) {
		err = ord.OrderRep.AddOrder(orderID, userID, "NEW", 0.0)
		if err != nil {
			return nil, fmt.Errorf("%w: %v", ErrorAddingOrder, err.Error())
		}
		return &LoadOrderResult{Status: StatusAccepted}, nil
	}

	if order.UserID != userID {
		return nil, ErrorOrderBelongsAnotherUser
	}
	return &LoadOrderResult{Status: StatusOK}, nil
}

// LoadOrderResult represents the status of an order after attempting to load it.
type LoadOrderResult struct {
	Status string
}

// StatusAccepted and StatusOK are constants representing the possible states
// of an order after being processed by LoadOrder.
const (
	StatusAccepted = "Accepted"
	StatusOK       = "OK"
)

// GetOrders retrieves the orders associated with a user token. It returns a
// slice of UserOrderFormat with details of each order.
func (ord *UserOrder) GetOrders(token string) ([]UserOrderFormat, error) {
	userID, err := security.GetUserIDFromToken(token)
	if err != nil {
		return []UserOrderFormat{}, err
	}

	orders, err := ord.OrderRep.GetOrderByUserID(userID)
	if !errors.Is(err, gorm.ErrRecordNotFound) && err != nil {
		return []UserOrderFormat{}, err
	}

	newOrders := make([]UserOrderFormat, 0, len(orders))
	for _, originalOrder := range orders {
		reducedOrder := ConvertToUserOrderFormat(originalOrder)
		newOrders = append(newOrders, *reducedOrder)
	}
	return newOrders, nil
}

// UserOrderFormat defines the format for representing user orders.
type UserOrderFormat struct {
	Number     string    `json:"number"`
	Status     string    `json:"status"`
	Accrual    *float64  `json:"accrual,omitempty"`
	UploadedAt time.Time `json:"uploaded_at"`
}

// ConvertToUserOrderFormat converts an orderdb.Order to UserOrderFormat
// for external representation.
func ConvertToUserOrderFormat(originalOrder orderdb.Order) *UserOrderFormat {
	if originalOrder.Accrual == 0 {
		return &UserOrderFormat{
			Number:     originalOrder.OrderID,
			Status:     originalOrder.Status,
			UploadedAt: originalOrder.UpdatedAt,
		}
	}
	return &UserOrderFormat{
		Number:     originalOrder.OrderID,
		Status:     originalOrder.Status,
		Accrual:    &originalOrder.Accrual,
		UploadedAt: originalOrder.UpdatedAt,
	}
}
