package service

import (
	"errors"
	"fmt"
	"testing"
	"time"

	"github.com/elina-chertova/loyalty-system/internal/db/balancedb"
	"github.com/elina-chertova/loyalty-system/internal/security"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type MockBalanceRepository struct{}

var (
	ErrorDownloadingBalance       = errors.New("balance cannot be created")
	ErrorDownloadingWithdrawFunds = errors.New("WithdrawFunds cannot be created")
)

func TestUserBalance_AddInitialBalance(t *testing.T) {
	rep := &MockBalanceRepository{}
	userBalance := NewBalance(rep)

	uuidTest, _ := uuid.Parse("69359037-9599-48e7-b8f2-48393c019133") // Use a different UUID

	err := userBalance.AddInitialBalance(uuidTest)
	if err != nil {
		t.Errorf("AddInitialBalance() error = %v", err)
	}
}

func TestUserBalance_WithdrawFunds(t *testing.T) {
	rep := &MockBalanceRepository{}
	userBalance := NewBalance(rep)
	uuidIDTest, _ := uuid.Parse("69359037-9599-48e7-b8f2-48393c019135")
	token, _ := security.GenerateToken(uuidIDTest)
	order := "6231543915765652"
	sum := 50.0

	err := userBalance.WithdrawFunds(token, order, sum)
	if err != nil {
		t.Errorf("WithdrawFunds() error = %v", err)
	}
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

func BenchmarkUserBalance_AddInitialBalance(b *testing.B) {
	rep := &MockBalanceRepository{}
	userBalance := NewBalance(rep)

	uuidTest, _ := uuid.Parse("69359037-9599-48e7-b8f2-48393c019133")

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = userBalance.AddInitialBalance(uuidTest)
	}
}

func BenchmarkUserBalance_WithdrawFunds(b *testing.B) {
	rep := &MockBalanceRepository{}
	userBalance := NewBalance(rep)

	uuidTest, _ := uuid.Parse("69359037-9599-48e7-b8f2-48393c019135")
	token, _ := security.GenerateToken(uuidTest)
	order := "6231543915765652"
	sum := 50.0

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = userBalance.WithdrawFunds(token, order, sum)
	}
}

func BenchmarkUserBalance_GetBalance(b *testing.B) {
	rep := &MockBalanceRepository{}
	userBalance := NewBalance(rep)

	uuidRight, _ := uuid.Parse("69359037-9599-48e7-b8f2-48393c019135")
	tokenRight, _ := security.GenerateToken(uuidRight)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = userBalance.GetBalance(tokenRight)
	}
}

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
	if userID != uuidIDTest {
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
