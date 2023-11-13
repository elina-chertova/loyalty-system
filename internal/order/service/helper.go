package service

import (
	"encoding/json"
	"github.com/elina-chertova/loyalty-system/internal/config"
	"github.com/levigross/grequests"
	"net/http"
)

type OrderLoyaltyFormat struct {
	Order   string  `json:"order"`
	Status  string  `json:"status"`
	Accrual float64 `json:"accrual,omitempty"`
}

func GetOrderInfo(orderID string) (OrderLoyaltyFormat, error) {
	response, err := grequests.Get("http://localhost:8080/api/orders/"+orderID, nil)
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
