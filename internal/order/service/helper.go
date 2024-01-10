// Package service provides functionalities for processing and updating
// order statuses in the loyalty system.
package service

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/elina-chertova/loyalty-system/internal/config"
	"github.com/levigross/grequests"
)

// OrderLoyaltyFormat defines the format for loyalty data associated with an order.
type OrderLoyaltyFormat struct {
	Order   string  `json:"order"`
	Status  string  `json:"status"`
	Accrual float64 `json:"accrual,omitempty"`
}

// UpdateOrderStatus updates the status of orders based on the response
// from the accrual system. It queries the accrual system for each unprocessed order
// and updates the order's status accordingly in the local system.
// The function handles both successful accrual updates and cases where
// the accrual system returns no content, marking such orders as "INVALID".
func (ord *UserOrder) UpdateOrderStatus(accrualServerAddress string) error {
	endpoint := fmt.Sprintf(config.AccrualSystemAddress, accrualServerAddress)

	orders, err := ord.OrderRep.GetUnprocessedOrders()
	if err != nil {
		return err
	}
	for _, order := range orders {
		response, err := grequests.Get(endpoint+order.OrderID, nil)
		if err != nil {
			return err
		}
		if response.StatusCode == http.StatusOK {
			var orderLoyalty OrderLoyaltyFormat
			err := json.Unmarshal(response.Bytes(), &orderLoyalty)
			if err != nil {
				return err
			}

			if orderLoyalty.Status != "" {
				err := ord.OrderRep.UpdateOrderStatus(
					orderLoyalty.Order,
					orderLoyalty.Status,
					orderLoyalty.Accrual,
				)
				if err != nil {
					return err
				}
			}
			return nil
		} else if response.StatusCode == http.StatusNoContent {
			err := ord.OrderRep.UpdateOrderStatus(order.OrderID, "INVALID", 0.0)
			if err != nil {
				return err
			}
		}
	}
	return nil
}
