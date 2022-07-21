package models

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"time"

	"github.com/google/uuid"
)

const (
	OrderStatusPending  = "pending"
	OrderStatusAccepted = "accepted"
)

// Order model
type Order struct {
	OrderID        uuid.UUID  `json:"order_id" db:"order_id"`
	UserID         uuid.UUID  `json:"user_id" db:"user_id"`
	LibrarianID    *uuid.UUID `json:"librarian_id" db:"librarian_id"`
	Item           OrderItem  `json:"item" db:"item"`
	Status         string     `json:"status" db:"status"`
	PickupSchedule time.Time  `json:"pickup_schedule,omitempty" db:"created_at"`
	CreatedAt      time.Time  `json:"created_at,omitempty" db:"created_at"`
	UpdatedAt      time.Time  `json:"updated_at,omitempty" db:"updated_at"`
}

type OrderItem Book

func (o *OrderItem) Scan(value interface{}) error {
	val, ok := value.([]byte)
	if !ok {
		return fmt.Errorf("unable to scan")
	}
	var item OrderItem
	if err := json.Unmarshal(val, &item); err != nil {
		return fmt.Errorf("json.Unmarshal %v", value)
	}
	*o = item
	return nil
}

func (o OrderItem) Value() (driver.Value, error) {
	valueJson, _ := json.Marshal(o)
	return string(valueJson), nil
}
