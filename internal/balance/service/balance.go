package service

import (
	"errors"
	"github.com/elina-chertova/loyalty-system/internal/config"
	"github.com/elina-chertova/loyalty-system/internal/db/balancedb"
	"github.com/elina-chertova/loyalty-system/internal/order/service"
	"github.com/elina-chertova/loyalty-system/internal/order/utils"
	"github.com/elina-chertova/loyalty-system/internal/security"
	"gorm.io/gorm"
	"net/http"
	"time"
)

type UserBalance struct {
	BalanceRep balancedb.BalanceRepository
}

func NewBalance(model balancedb.BalanceRepository) *UserBalance {
	return &UserBalance{BalanceRep: model}
}

func (bal *UserBalance) WithdrawFunds(token, order string, sum float64) (int, error) {
	userID, err := security.GetUserIDFromToken(token)
	if err != nil {
		return http.StatusInternalServerError, err
	}

	if !utils.IsLuhnValid(order) {
		return http.StatusUnprocessableEntity, config.ErrorNotValidOrderNumber
	}

	balance, err := bal.BalanceRep.GetBalanceByUserID(userID)
	if err != nil {
		return http.StatusInternalServerError, err
	}

	current := balance.Current - sum
	withdrawn := balance.Withdrawn + sum
	if current < 0 {
		return http.StatusPaymentRequired, config.ErrorInsufficientFunds
	}

	//orders, err := bal.balanceRep.GetOrdersWithdrawFunds()

	err = bal.BalanceRep.AddWithdrawFunds(userID, order, sum)
	if err != nil {
		return http.StatusInternalServerError, err
	}

	err = bal.BalanceRep.UpdateBalance(userID, current, withdrawn)
	if err != nil {
		return http.StatusInternalServerError, err
	}

	return http.StatusOK, nil
}

func (bal *UserBalance) GetBalance(token string) (UserBalanceFormat, error) {
	userID, err := security.GetUserIDFromToken(token)
	if err != nil {
		return UserBalanceFormat{}, err
	}
	balance, err := bal.BalanceRep.GetBalanceByUserID(userID)

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return UserBalanceFormat{}, nil
	} else if err != nil {
		return UserBalanceFormat{}, err
	}

	return *ConvertToUserBalanceFormat(balance), nil
}

type UserBalanceFormat struct {
	Current   float64 `json:"current"`
	Withdrawn float64 `json:"withdrawn"`
}

func ConvertToUserBalanceFormat(originalBalance balancedb.Balance) *UserBalanceFormat {
	return &UserBalanceFormat{
		Current:   originalBalance.Current,
		Withdrawn: originalBalance.Withdrawn,
	}
}

type WithdrawalFormat struct {
	Order       string    `json:"order" gorm:"unique_index"`
	Sum         float64   `json:"sum"`
	ProcessedAt time.Time `json:"processed_at"`
}

func (bal *UserBalance) WithdrawalInfo(token string) ([]WithdrawalFormat, error) {
	userID, err := security.GetUserIDFromToken(token)
	if err != nil {
		return nil, err
	}
	withdrawals, err := bal.BalanceRep.GetWithdrawalByUserID(userID)
	if err != nil {
		return nil, err
	}

	var newWithdrawals []WithdrawalFormat
	for _, w := range withdrawals {
		newWithdrawals = append(
			newWithdrawals, WithdrawalFormat{
				Order:       w.Order,
				Sum:         w.Sum,
				ProcessedAt: w.UpdatedAt,
			},
		)
	}

	return newWithdrawals, nil
}

func (bal *UserBalance) UpdateBalance(ord *service.UserOrder) error {
	orderAccrual, err := ord.OrderRep.GetOrderAccrual()
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return err
	}

	for _, rows := range orderAccrual {
		balance, err := bal.BalanceRep.GetBalanceByUserID(rows.UserID)
		switch e := err; {
		case errors.Is(e, gorm.ErrRecordNotFound):
			err := bal.BalanceRep.AddBalance(rows.UserID, rows.SumAccrual, 0.0)
			if err != nil {
				return err
			}
			continue
		case !errors.Is(e, gorm.ErrRecordNotFound) && e != nil:
			return e
		}

		current := balance.Current + rows.SumAccrual
		err = bal.BalanceRep.AddBalance(
			rows.UserID,
			current,
			balance.Withdrawn,
		)
		if err != nil {
			return err
		}

	}

	return nil
}