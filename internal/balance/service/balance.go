// Package service provides functionalities for managing user balances
// and withdrawals in the loyalty system.
package service

import (
	"errors"
	"fmt"
	"time"

	"github.com/elina-chertova/loyalty-system/internal/db/balancedb"
	"github.com/elina-chertova/loyalty-system/internal/order/service"
	"github.com/elina-chertova/loyalty-system/internal/order/utils"
	"github.com/elina-chertova/loyalty-system/internal/security"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

// UserBalance handles operations related to user balances.
type UserBalance struct {
	balanceRep balancedb.BalanceRepository
}

// NewBalance creates a new instance of UserBalance with the given BalanceRepository.
func NewBalance(model balancedb.BalanceRepository) *UserBalance {
	return &UserBalance{balanceRep: model}
}

// Predefined errors for balance operations.
var (
	ErrorNotValidOrderNumber = errors.New("order number is not valid")
	ErrorSystem              = errors.New("error in loyality system")
	ErrorInsufficientFunds   = errors.New("insufficient funds")
)

// AddInitialBalance sets the initial balance for a given user ID.
func (bal *UserBalance) AddInitialBalance(userID uuid.UUID) error {
	err := bal.balanceRep.AddBalance(userID, 0.0, 0.0)
	if err != nil {
		return err
	}
	return nil
}

// WithdrawFunds processes a withdrawal request for a user identified by a token.
// It verifies the validity of the order number and checks if sufficient funds are available.
func (bal *UserBalance) WithdrawFunds(token, order string, sum float64) error {
	userID, err := security.GetUserIDFromToken(token)
	if err != nil {
		return fmt.Errorf("%w; %v", ErrorSystem, err)
	}

	if !utils.IsLuhnValid(order) {
		return ErrorNotValidOrderNumber
	}

	balance, err := bal.balanceRep.GetBalanceByUserID(userID)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return err
	}
	if err != nil {
		return fmt.Errorf("%w; %v", ErrorSystem, err)
	}

	current := balance.Current - sum
	withdrawn := balance.Withdrawn + sum
	if current < 0 {
		return ErrorInsufficientFunds
	}

	err = bal.balanceRep.AddWithdrawFunds(userID, order, sum)
	if err != nil {
		return fmt.Errorf("%w; %v", ErrorSystem, err)
	}

	err = bal.balanceRep.UpdateBalance(userID, current, withdrawn)
	if err != nil {
		return fmt.Errorf("%w; %v", ErrorSystem, err)
	}

	return nil
}

// GetBalance retrieves the current balance and withdrawn amount for a user identified by a token.
func (bal *UserBalance) GetBalance(token string) (UserBalanceFormat, error) {
	userID, err := security.GetUserIDFromToken(token)
	if err != nil {
		return UserBalanceFormat{}, err
	}
	balance, err := bal.balanceRep.GetBalanceByUserID(userID)

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return UserBalanceFormat{}, nil
	} else if err != nil {
		return UserBalanceFormat{}, err
	}

	return *ConvertToUserBalanceFormat(balance), nil
}

// UserBalanceFormat defines the format for representing user balances.
type UserBalanceFormat struct {
	Current   float64 `json:"current"`
	Withdrawn float64 `json:"withdrawn"`
}

// ConvertToUserBalanceFormat converts a balancedb.Balance to UserBalanceFormat
// for external representation.
func ConvertToUserBalanceFormat(originalBalance balancedb.Balance) *UserBalanceFormat {
	return &UserBalanceFormat{
		Current:   originalBalance.Current,
		Withdrawn: originalBalance.Withdrawn,
	}
}

// WithdrawalFormat defines the format for representing user withdrawals.
type WithdrawalFormat struct {
	Order       string    `json:"order" gorm:"unique_index"`
	Sum         float64   `json:"sum"`
	ProcessedAt time.Time `json:"processed_at"`
}

// WithdrawalInfo retrieves the withdrawal information for a user identified by a token.
func (bal *UserBalance) WithdrawalInfo(token string) ([]WithdrawalFormat, error) {
	userID, err := security.GetUserIDFromToken(token)
	if err != nil {
		return nil, err
	}
	withdrawals, err := bal.balanceRep.GetWithdrawalByUserID(userID)
	if err != nil {
		return nil, err
	}

	newWithdrawals := make([]WithdrawalFormat, 0, len(withdrawals))
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

// UpdateBalance updates the balance of users based on their recent order accruals.
func (bal *UserBalance) UpdateBalance(ord *service.UserOrder) error {
	orderAccrual, err := ord.OrderRep.GetOrderAccrual()
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return err
	}

	for _, rows := range orderAccrual {
		balance, err := bal.balanceRep.GetBalanceByUserID(rows.UserID)
		switch e := err; {
		case errors.Is(e, gorm.ErrRecordNotFound):
			err := bal.balanceRep.AddBalance(rows.UserID, rows.SumAccrual, 0.0)
			if err != nil {
				return err
			}
			continue
		case !errors.Is(e, gorm.ErrRecordNotFound) && e != nil:
			return e
		}

		current := balance.Current + rows.SumAccrual

		if err = bal.balanceRep.UpdateBalance(
			rows.UserID,
			current,
			balance.Withdrawn,
		); err != nil {
			return err
		}

	}

	return nil
}
