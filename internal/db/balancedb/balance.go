package balancedb

import (
	"fmt"
	"github.com/elina-chertova/loyalty-system/internal/config"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type BalanceModel struct {
	DB *gorm.DB
}

func NewBalanceModel(db *gorm.DB) *BalanceModel {
	return &BalanceModel{DB: db}
}

type BalanceRepository interface {
	AddBalance(uuid.UUID, float64, float64) error
	GetBalanceByUserID(uuid.UUID) (Balance, error)
	UpdateBalance(uuid.UUID, float64, float64) error

	AddWithdrawFunds(uuid.UUID, string, float64) error
	GetOrdersWithdrawFunds() ([]string, error)
	GetWithdrawalByUserID(userID uuid.UUID) ([]Withdrawal, error)
}

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
		return fmt.Errorf("%w: %v", config.ErrorDownloadingBalance, result.Error)
	}
	return nil
}

func (balanceDB *BalanceModel) GetBalanceByUserID(userID uuid.UUID) (Balance, error) {
	var balance Balance
	result := balanceDB.DB.Order("updated_at desc").Where(&Balance{UserID: userID}).First(&balance)
	if result.Error != nil {
		return balance, result.Error
	}
	return balance, nil
}

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
		return fmt.Errorf("%w: %v", config.ErrorDownloadingWithdrawFunds, result.Error)
	}
	return nil
}

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

func (balanceDB *BalanceModel) GetWithdrawalByUserID(userID uuid.UUID) ([]Withdrawal, error) {
	var withdrawals []Withdrawal
	result := balanceDB.DB.Order("updated_at desc").Where(&Withdrawal{UserID: userID}).Find(&withdrawals)
	if result.Error != nil {
		return withdrawals, result.Error
	}
	return withdrawals, nil
}
