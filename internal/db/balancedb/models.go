package balancedb

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Balance struct {
	gorm.Model
	UserID    uuid.UUID `json:"user_id"`
	Current   float64   `json:"current"`
	Withdrawn float64   `json:"withdrawn"`
	UpdatedAt time.Time `json:"updated_at"`
}

type Withdrawal struct {
	gorm.Model
	UserID    uuid.UUID `json:"user_id"`
	Order     string    `json:"order" gorm:"unique_index"`
	Sum       float64   `json:"sum"`
	UpdatedAt time.Time `json:"updated_at"`
}
