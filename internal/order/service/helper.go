package service

import (
	"encoding/json"
	"fmt"
	"github.com/elina-chertova/loyalty-system/internal/config"
	"github.com/levigross/grequests"
	"net/http"
)

type OrderLoyaltyFormat struct {
	Order   string  `json:"order"`
	Status  string  `json:"status"`
	Accrual float64 `json:"accrual,omitempty"`
}

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
