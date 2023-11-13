package orderdb

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
	"time"
)

type Order struct {
	gorm.Model
	OrderID   string    `json:"id" gorm:"unique_index"`
	UserID    uuid.UUID `json:"user_id"`
	Status    string    `json:"status"`
	Accrual   float64   `json:"accrual"`
	Credited  bool      `json:"credited"`
	CreatedAt time.Time `json:"created_at"`
}
