package service

import (
	"errors"
	"github.com/elina-chertova/loyalty-system/internal/db/orderdb"
	"github.com/elina-chertova/loyalty-system/internal/security"
	"github.com/google/uuid"
	"gorm.io/gorm"
	"testing"
)

type MockOrderRepository struct {
	Orders map[string]orderdb.Order
}

func NewMockOrderRepository() *MockOrderRepository {
	return &MockOrderRepository{
		Orders: make(map[string]orderdb.Order),
	}
}

func (m *MockOrderRepository) GetOrderAccrual() ([]orderdb.UserAccrual, error) {
	return []orderdb.UserAccrual{}, nil
}

func (m *MockOrderRepository) GetPreparedOrders() ([]orderdb.OrderAccrual, error) {
	return []orderdb.OrderAccrual{}, nil
}

func (m *MockOrderRepository) SetProcessedOrders(order []orderdb.OrderAccrual) error {
	return nil
}

func (m *MockOrderRepository) GetTotalAccrualByUsers() ([]orderdb.UserAccrual, error) {
	return []orderdb.UserAccrual{}, nil
}

func (m *MockOrderRepository) GetUnprocessedOrders() ([]orderdb.Order, error) {
	return []orderdb.Order{}, nil
}

func (m *MockOrderRepository) UpdateOrderStatus(
	orderID string,
	newStatus string,
	accrual float64,
) error {
	return nil
}

func (m *MockOrderRepository) GetOrderByID(orderID string) (orderdb.Order, error) {
	if order, exists := m.Orders[orderID]; exists {
		return order, nil
	}
	return orderdb.Order{}, gorm.ErrRecordNotFound
}

func (m *MockOrderRepository) AddOrder(
	orderID string,
	userID uuid.UUID,
	status string,
	accrual float64,
) error {
	if _, exists := m.Orders[orderID]; exists {
		return errors.New("order already exists")
	}
	m.Orders[orderID] = orderdb.Order{
		OrderID: orderID,
		UserID:  userID,
		Status:  status,
		Accrual: accrual,
	}
	return nil
}

func (m *MockOrderRepository) GetOrderByUserID(userID uuid.UUID) ([]orderdb.Order, error) {
	var orders []orderdb.Order
	for _, order := range m.Orders {
		if order.UserID == userID {
			orders = append(orders, order)
		}
	}
	return orders, nil
}

func TestUserOrder_LoadOrder(t *testing.T) {
	mockRepo := &MockOrderRepository{Orders: make(map[string]orderdb.Order)}
	userOrder := NewOrder(mockRepo)

	token, _ := security.GenerateToken(uuid.New())
	orderID := "6231543915765652"

	result, err := userOrder.LoadOrder(token, orderID)
	if err != nil {
		t.Errorf("LoadOrder() error = %v", err)
	} else if result.Status != StatusAccepted {
		t.Errorf("Expected status %s, got %s", StatusAccepted, result.Status)
	}
}

func TestUserOrder_GetOrders(t *testing.T) {
	mockRepo := NewMockOrderRepository()
	userOrder := NewOrder(mockRepo)

	token, _ := security.GenerateToken(uuid.New())

	_, err := userOrder.GetOrders(token)
	if err != nil {
		t.Errorf("GetOrders() error = %v", err)
	}
}

func BenchmarkUserOrder_LoadOrder(b *testing.B) {
	mockRepo := NewMockOrderRepository()
	userOrder := NewOrder(mockRepo)

	token, _ := security.GenerateToken(uuid.New())
	orderID := "6231543915765652"

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = userOrder.LoadOrder(token, orderID)
	}
}

func BenchmarkUserOrder_GetOrders(b *testing.B) {
	mockRepo := NewMockOrderRepository()
	userOrder := NewOrder(mockRepo)

	token, _ := security.GenerateToken(uuid.New())

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = userOrder.GetOrders(token)
	}
}
