// Package orderdb provides data access functionalities for order management
// in the loyalty system. It uses GORM for database operations related to orders.
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

// OrderModel represents the model for order data and provides methods
// for interacting with the orders table in the database.
type OrderModel struct {
	DB *gorm.DB
}

// NewOrderModel creates a new instance of OrderModel with the given GORM DB instance.
func NewOrderModel(db *gorm.DB) *OrderModel {
	return &OrderModel{DB: db}
}

// OrderRepository defines the interface for order data operations. It abstracts
// the methods to interact with the orders in the database.
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

// ErrorDownloadingOrder represents an error encountered while creating an order.
var ErrorDownloadingOrder = errors.New("order cannot be created")

// AddOrder adds a new order to the database with the provided details.
// Returns an error if the order cannot be created.
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

// GetOrderByID retrieves an order by its ID from the database.
// Returns the Order object and an error if the order is not found.
func (orderDB *OrderModel) GetOrderByID(orderID string) (Order, error) {
	var order Order
	result := orderDB.DB.Where(&Order{OrderID: orderID}).First(&order)
	if result.Error != nil {
		return Order{}, result.Error
	}
	return order, nil
}

// GetOrderByUserID retrieves an order by user ID from the database.
// Returns the Order object and an error if the order is not found.
func (orderDB *OrderModel) GetOrderByUserID(userID uuid.UUID) ([]Order, error) {
	var orders []Order
	result := orderDB.DB.Order("updated_at desc").Where(&Order{UserID: userID}).Find(&orders)
	if result.Error != nil {
		return orders, result.Error
	}
	return orders, nil
}

// OrderAccrual represents the accrual data associated with an order.
type OrderAccrual struct {
	UserID     uuid.UUID `gorm:"column:user_id"`
	Order      string    `gorm:"column:order_id"`
	SumAccrual float64   `gorm:"column:accrual"`
}

// UserAccrual represents the total accrual for a user.
type UserAccrual struct {
	UserID     uuid.UUID `gorm:"column:user_id"`
	SumAccrual float64   `gorm:"column:total_accrual"`
}

// GetPreparedOrders retrieves prepared orders.
// Returns the OrderAccrual object and an error if the order is not found.
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

// GetTotalAccrualByUsers retrieves total accrual.
// Returns the UserAccrual object and an error if the rows are not found.
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

// SetProcessedOrders marks a batch of orders as processed in the database.
// order: A slice of OrderAccrual representing the orders to be marked as processed.
// The method updates the 'credited' field of each specified order to 'true'.
// Returns an error if the update operation fails.
func (orderDB *OrderModel) SetProcessedOrders(order []OrderAccrual) error {
	updateResult := orderDB.DB.Table(config.TableOrder).
		Where("order_id IN ? AND status = ?", getOrderIDs(order), config.Processed).
		Update("credited", true)

	if updateResult.Error != nil {
		return updateResult.Error
	}
	return nil
}

// GetOrderAccrual retrieves accrual information for each user based on their orders.
// It first fetches prepared orders, then calculates total accruals by users,
// and finally marks the processed orders as credited.
// Returns a slice of UserAccrual containing accrual information for each user, or an error.
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

// getOrderIDs extracts and returns a slice of order IDs from a slice of OrderAccrual.
// orders: A slice of OrderAccrual from which the order IDs are to be extracted.
func getOrderIDs(orders []OrderAccrual) []string {
	var orderIDs []string
	for _, o := range orders {
		orderIDs = append(orderIDs, o.Order)
	}
	return orderIDs
}

// GetUnprocessedOrders fetches orders that are either in 'PROCESSING' or 'NEW' status.
// Returns a slice of Order objects representing unprocessed orders, or an error if the fetch fails.
func (orderDB *OrderModel) GetUnprocessedOrders() ([]Order, error) {
	var orders []Order
	result := orderDB.DB.Where("status IN ?", []string{"PROCESSING", "NEW"}).First(&orders)
	if result.Error != nil {
		return []Order{}, result.Error
	}
	return orders, nil
}

// UpdateOrderStatus updates the status and accrual of an order identified by orderID.
// orderID: Identifier of the order to be updated.
// newStatus: New status to be set for the order.
// accrual: Accrual amount to be updated for the order.
// Returns an error if the update operation fails.
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
