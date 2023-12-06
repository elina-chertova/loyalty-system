package orderdb

import (
	"errors"
	"fmt"
	"github.com/elina-chertova/loyalty-system/internal/config"
	"github.com/elina-chertova/loyalty-system/pkg/logger"
	"github.com/google/uuid"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type OrderModel struct {
	DB *gorm.DB
}

func NewOrderModel(db *gorm.DB) *OrderModel {
	return &OrderModel{DB: db}
}

type OrderRepository interface {
	AddOrder(string, uuid.UUID, string, float64) error
	GetOrderByID(string) (Order, error)
	GetOrderByUserID(uuid.UUID) ([]Order, error)
	GetOrderAccrual() ([]UserAccrual, error)

	GetPreparedOrders() ([]OrderAccrual, error)
	SetProcessedOrders(order []OrderAccrual) error
	GetTotalAccrualByUsers() ([]UserAccrual, error)

	GetUnprocessedOrders() ([]Order, error)
	UpdateOrderStatus(orderID string, newStatus string, accrual float64) error
}

var ErrorDownloadingOrder = errors.New("order cannot be created")

func (orderDB *OrderModel) AddOrder(
	orderID string,
	userID uuid.UUID,
	status string,
	accrual float64,
) error {
	result := orderDB.DB.Create(
		&Order{
			OrderID:  orderID,
			UserID:   userID,
			Status:   status,
			Accrual:  accrual,
			Credited: false,
		},
	)
	if result.Error != nil {
		return fmt.Errorf("%w: %v", ErrorDownloadingOrder, result.Error)
	}
	return nil
}

func (orderDB *OrderModel) GetOrderByID(orderID string) (Order, error) {
	var order Order
	result := orderDB.DB.Where(&Order{OrderID: orderID}).First(&order)
	if result.Error != nil {
		return Order{}, result.Error
	}
	return order, nil
}

func (orderDB *OrderModel) GetOrderByUserID(userID uuid.UUID) ([]Order, error) {
	var orders []Order
	result := orderDB.DB.Order("updated_at desc").Where(&Order{UserID: userID}).Find(&orders)
	if result.Error != nil {
		return orders, result.Error
	}
	return orders, nil
}

type OrderAccrual struct {
	UserID     uuid.UUID `gorm:"column:user_id"`
	Order      string    `gorm:"column:order_id"`
	SumAccrual float64   `gorm:"column:accrual"`
}

type UserAccrual struct {
	UserID     uuid.UUID `gorm:"column:user_id"`
	SumAccrual float64   `gorm:"column:total_accrual"`
}

func (orderDB *OrderModel) GetPreparedOrders() ([]OrderAccrual, error) {
	var order []OrderAccrual

	result := orderDB.DB.Table(config.TableOrder).Select("user_id, order_id, accrual").Where(
		"credited = ? AND status = ?",
		false,
		config.Processed,
	).Find(&order)
	if result.Error != nil {
		return []OrderAccrual{}, result.Error
	}
	return order, nil
}

func (orderDB *OrderModel) GetTotalAccrualByUsers() ([]UserAccrual, error) {
	var userSum []UserAccrual
	resultUserAccrual := orderDB.DB.Table(config.TableOrder).Select("user_id, SUM(accrual) as total_accrual").Where(
		"credited = ? AND status = ?",
		false,
		config.Processed,
	).Group("user_id").Find(&userSum)
	if resultUserAccrual.Error != nil {
		return []UserAccrual{}, resultUserAccrual.Error
	}
	return userSum, nil
}

func (orderDB *OrderModel) SetProcessedOrders(order []OrderAccrual) error {
	updateResult := orderDB.DB.Table(config.TableOrder).
		Where("order_id IN ? AND status = ?", getOrderIDs(order), config.Processed).
		Update("credited", true)

	if updateResult.Error != nil {
		return updateResult.Error
	}
	return nil
}

func (orderDB *OrderModel) GetOrderAccrual() ([]UserAccrual, error) {
	order, err := orderDB.GetPreparedOrders()
	if err != nil {
		return []UserAccrual{}, err
	}

	userSum, err := orderDB.GetTotalAccrualByUsers()
	if err != nil {
		logger.Logger.Warn("Error in GetTotalAccrualByUsers", zap.Error(err))
		return []UserAccrual{}, err
	}
	err = orderDB.SetProcessedOrders(order)
	if err != nil {
		logger.Logger.Warn("Error in SetProcessedOrders", zap.Error(err))
		return []UserAccrual{}, err
	}

	return userSum, nil
}

func getOrderIDs(orders []OrderAccrual) []string {
	var orderIDs []string
	for _, o := range orders {
		orderIDs = append(orderIDs, o.Order)
	}
	return orderIDs
}

func (orderDB *OrderModel) GetUnprocessedOrders() ([]Order, error) {
	var orders []Order
	result := orderDB.DB.Where("status IN ?", []string{"PROCESSING", "NEW"}).First(&orders)
	if result.Error != nil {
		return []Order{}, result.Error
	}
	return orders, nil
}

func (orderDB *OrderModel) UpdateOrderStatus(
	orderID string,
	newStatus string,
	accrual float64,
) error {
	result := orderDB.DB.Model(&Order{}).Where("order_id = ?", orderID).Update(
		"status",
		newStatus,
	).Update("accrual", accrual)
	if result.Error != nil {
		return result.Error
	}
	return nil
}
