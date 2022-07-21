package dto

import (
	"time"

	"github.com/google/uuid"
)

type OrderCreateRequestDto struct {
	BookKey        string    `json:"key" validate:"required"`
	PickupSchedule time.Time `json:"pickup_schedule" validate:"required"`
}

type OrderCreateResponseDto struct {
	OrderID uuid.UUID `json:"order_id" validate:"required"`
}
