package service

import (
	"errors"
	"fmt"
	"github.com/elina-chertova/loyalty-system/internal/config"
	"github.com/elina-chertova/loyalty-system/internal/db/orderdb"
	"github.com/elina-chertova/loyalty-system/internal/order/utils"
	"github.com/elina-chertova/loyalty-system/internal/security"
	"gorm.io/gorm"
	"time"
)

type UserOrder struct {
	OrderRep orderdb.OrderRepository
}

func NewOrder(model orderdb.OrderRepository) *UserOrder {
	return &UserOrder{OrderRep: model}
}

func (ord *UserOrder) LoadOrder(token string, orderID string) (*LoadOrderResult, error) {
	if !utils.IsLuhnValid(orderID) {
		return nil, config.ErrorNotValidOrderNumber
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
			return nil, fmt.Errorf("%w: %v", config.ErrorAddingOrder, err.Error())
		}
		return &LoadOrderResult{Status: StatusAccepted}, nil
	}

	if order.UserID != userID {
		return nil, config.ErrorOrderBelongsAnotherUser
	}
	return &LoadOrderResult{Status: StatusOK}, nil
}

type LoadOrderResult struct {
	Status string
}

const (
	StatusAccepted = "Accepted"
	StatusOK       = "OK"
)

func (ord *UserOrder) GetOrders(token string) ([]UserOrderFormat, error) {
	userID, err := security.GetUserIDFromToken(token)
	if err != nil {
		return []UserOrderFormat{}, err
	}

	orders, err := ord.OrderRep.GetOrderByUserID(userID)
	if !errors.Is(err, gorm.ErrRecordNotFound) && err != nil {
		return []UserOrderFormat{}, err
	}

	var newOrders []UserOrderFormat

	for _, originalOrder := range orders {
		reducedOrder := ConvertToUserOrderFormat(originalOrder)
		newOrders = append(newOrders, *reducedOrder)
	}
	return newOrders, nil
}

type UserOrderFormat struct {
	Number     string    `json:"number"`
	Status     string    `json:"status"`
	Accrual    *float64  `json:"accrual,omitempty"`
	UploadedAt time.Time `json:"uploaded_at"`
}

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
