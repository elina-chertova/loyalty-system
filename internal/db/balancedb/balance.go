// Package balancedb provides data access functionalities for managing
// user balances and withdrawals in the loyalty system. It uses GORM
package balancedb

import (
	"errors"
	"fmt"

	"github.com/elina-chertova/loyalty-system/internal/config"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

// BalanceModel represents the model for managing balance and withdrawal data in the database.
type BalanceModel struct {
	DB *gorm.DB
}

// NewBalanceModel creates a new instance of BalanceModel with the given GORM DB instance.
// This function initializes the model for interacting with balance-related data.
func NewBalanceModel(db *gorm.DB) *BalanceModel {
	return &BalanceModel{DB: db}
}

// BalanceRepository defines the interface for balance data operations. It abstracts
// the methods to interact with user balances and withdrawals in the database.
type BalanceRepository interface {
	AddBalance(uuid.UUID, float64, float64) error
	GetBalanceByUserID(uuid.UUID) (Balance, error)
	UpdateBalance(uuid.UUID, float64, float64) error

	AddWithdrawFunds(uuid.UUID, string, float64) error
	GetOrdersWithdrawFunds() ([]string, error)
	GetWithdrawalByUserID(userID uuid.UUID) ([]Withdrawal, error)
}

// ErrorDownloadingBalance and ErrorDownloadingWithdrawFunds represent errors
// encountered during balance and withdrawal operations, respectively.
var (
	ErrorDownloadingBalance       = errors.New("balance cannot be created")
	ErrorDownloadingWithdrawFunds = errors.New("WithdrawFunds cannot be created")
)

// UpdateBalance updates the balance record in the database for a specific user.
// userID: Unique identifier of the user.
// current: Updated current balance to be set.
// withdrawn: Updated withdrawn amount to be set.
// Returns an error if the update operation fails.
func (balanceDB *BalanceModel) UpdateBalance(userID uuid.UUID, current, withdrawn float64) error {
	updateResult := balanceDB.DB.Table(config.TableBalance).Where("user_id = ?", userID).
		Updates(
			map[string]interface{}{
				"current":   current,
				"withdrawn": withdrawn,
			},
		)
	if updateResult.Error != nil {
		return updateResult.Error
	}
	return nil
}

// AddBalance adds a new balance entry to the database for the specified user.
// Returns an error if the balance cannot be created.
func (balanceDB *BalanceModel) AddBalance(
	userID uuid.UUID,
	current float64,
	withdrawn float64,
) error {
	result := balanceDB.DB.Create(
		&Balance{
			UserID:    userID,
			Current:   current,
			Withdrawn: withdrawn,
		},
	)
	if result.Error != nil {
		return fmt.Errorf("%w: %v", ErrorDownloadingBalance, result.Error)
	}
	return nil
}

// GetBalanceByUserID retrieves a user's balance by their user ID from the database.
// Returns the Balance object and an error if the balance is not found.
func (balanceDB *BalanceModel) GetBalanceByUserID(userID uuid.UUID) (Balance, error) {
	var balance Balance
	result := balanceDB.DB.Order("updated_at desc").Where(&Balance{UserID: userID}).First(&balance)
	if result.Error != nil {
		return balance, result.Error
	}
	return balance, nil
}

// AddWithdrawFunds adds a withdrawal record to the database for a specific order.
// userID: Unique identifier of the user.
// order: Identifier of the order for which the withdrawal is made.
// sum: Amount of funds to be withdrawn.
// Returns an error if the operation fails.
func (balanceDB *BalanceModel) AddWithdrawFunds(
	userID uuid.UUID,
	order string,
	sum float64,
) error {
	result := balanceDB.DB.Table(config.TableWithdrawal).Create(
		&Withdrawal{
			UserID: userID,
			Order:  order,
			Sum:    sum,
		},
	)
	if result.Error != nil {
		return fmt.Errorf("%w: %v", ErrorDownloadingWithdrawFunds, result.Error)
	}
	return nil
}

// GetOrdersWithdrawFunds retrieves a list of order IDs for which funds have been withdrawn.
// Returns a slice of order IDs and an error if retrieval fails.
func (balanceDB *BalanceModel) GetOrdersWithdrawFunds() ([]string, error) {
	var orders []string

	result := balanceDB.DB.Table(config.TableWithdrawal).Order("updated_at desc").Pluck(
		"order",
		&orders,
	).Find(&orders)
	if result.Error != nil {
		return orders, result.Error
	}

	return orders, nil
}

// GetWithdrawalByUserID retrieves all withdrawal records associated with a specific user.
// userID: Unique identifier of the user.
// Returns a slice of Withdrawal objects and an error if retrieval fails.
func (balanceDB *BalanceModel) GetWithdrawalByUserID(userID uuid.UUID) ([]Withdrawal, error) {
	var withdrawals []Withdrawal
	result := balanceDB.DB.Order("updated_at desc").Where(&Withdrawal{UserID: userID}).Find(&withdrawals)
	if result.Error != nil {
		return withdrawals, result.Error
	}
	return withdrawals, nil
}
