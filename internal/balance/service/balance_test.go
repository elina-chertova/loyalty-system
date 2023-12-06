package service

import (
	"errors"
	"fmt"
	"github.com/elina-chertova/loyalty-system/internal/db/balancedb"
	"github.com/elina-chertova/loyalty-system/internal/security"
	"github.com/google/uuid"
	"gorm.io/gorm"
	"testing"
	"time"
)

type MockBalanceRepository struct{}

var (
	ErrorDownloadingBalance       = errors.New("balance cannot be created")
	ErrorDownloadingWithdrawFunds = errors.New("WithdrawFunds cannot be created")
)

func (m *MockBalanceRepository) AddBalance(
	userID uuid.UUID,
	current float64,
	withdrawn float64,
) error {
	uuidIDTest, _ := uuid.Parse("69359037-9599-48e7-b8f2-48393c019135")
	if userID == uuidIDTest {
		return fmt.Errorf("%w: %v", ErrorDownloadingBalance, gorm.ErrRecordNotFound)
	}
	return nil
}

func (m *MockBalanceRepository) AddWithdrawFunds(
	userID uuid.UUID,
	order string,
	sum float64,
) error {
	uuidIDTest, _ := uuid.Parse("69359037-9599-48e7-b8f2-48393c019135")
	if userID == uuidIDTest {
		return fmt.Errorf("%w: %v", ErrorDownloadingWithdrawFunds, gorm.ErrRecordNotFound)
	}
	return nil
}

func (m *MockBalanceRepository) GetBalanceByUserID(userID uuid.UUID) (balancedb.Balance, error) {
	uuidIDTest1, _ := uuid.Parse("69359037-9599-48e7-b8f2-48393c019135")

	if uuidIDTest1 == userID {
		return balancedb.Balance{
			UserID:    uuidIDTest1,
			Current:   545.6,
			Withdrawn: 53,
			UpdatedAt: time.Now(),
		}, nil
	}
	return balancedb.Balance{}, errors.New("null balance")
}

func (m *MockBalanceRepository) UpdateBalance(userID uuid.UUID, current, withdrawn float64) error {
	return nil
}

func (m *MockBalanceRepository) GetOrdersWithdrawFunds() ([]string, error) {
	var orders []string
	return orders, nil
}

func (m *MockBalanceRepository) GetWithdrawalByUserID(userID uuid.UUID) (
	[]balancedb.Withdrawal,
	error,
) {
	var withdrawals []balancedb.Withdrawal
	return withdrawals, nil
}

func TestUserBalance_GetBalance(t *testing.T) {
	uuidRight, _ := uuid.Parse("69359037-9599-48e7-b8f2-48393c019135")
	tokenRight, err := security.GenerateToken(uuidRight)
	if err != nil {
		return
	}
	uuidWrong, _ := uuid.Parse("bc5b369e-2b4d-4958-8560-c8fbcc5c8188")
	tokenWrong, err := security.GenerateToken(uuidWrong)
	if err != nil {
		return
	}
	type args struct {
		token string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name:    "Balance is okey",
			args:    args{token: tokenRight},
			wantErr: false,
		},
		{
			name:    "Balance is empty",
			args:    args{token: tokenWrong},
			wantErr: true,
		},
	}

	rep := &MockBalanceRepository{}
	userBalance := NewBalance(rep)

	for _, tt := range tests {
		t.Run(
			tt.name, func(t *testing.T) {
				_, err := userBalance.GetBalance(tt.args.token)
				if (err != nil) != tt.wantErr {
					t.Errorf("WithdrawFunds() error = %v, wantErr %v", err, tt.wantErr)
					return
				}
			},
		)
	}
}
