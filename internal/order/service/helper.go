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

func GetOrderInfo(orderID string, accrualServerAddress string) (OrderLoyaltyFormat, error) {
	endpoint := fmt.Sprintf(config.AccrualSystemAddress, accrualServerAddress)
	response, err := grequests.Get(endpoint+orderID, nil)
	if err != nil {
		return OrderLoyaltyFormat{}, err
	}

	if response.StatusCode == http.StatusOK {
		var orderLoyalty OrderLoyaltyFormat
		err := json.Unmarshal(response.Bytes(), &orderLoyalty)
		if err != nil {
			return OrderLoyaltyFormat{}, err
		}

		return orderLoyalty, nil
	}
	return OrderLoyaltyFormat{}, config.ErrorGettingOrder
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
				err := ord.OrderRep.UpdateOrderStatus(orderLoyalty.Order, orderLoyalty.Status)
				if err != nil {
					return err
				}
			}
			return nil
		}
	}
	return nil
}
