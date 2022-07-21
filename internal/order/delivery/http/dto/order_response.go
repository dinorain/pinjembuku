package dto

import (
	"time"

	"github.com/google/uuid"

	"github.com/dinorain/pinjembuku/internal/models"
)

type OrderResponseDto struct {
	OrderID        uuid.UUID        `json:"order_id"`
	UserID         uuid.UUID        `json:"user_id"`
	LibrarianID    *uuid.UUID       `json:"librarian_id"`
	Item           models.OrderItem `json:"item"`
	Status         string           `json:"status" db:"status"`
	PickupSchedule time.Time        `json:"pickup_schedule,omitempty"`
	CreatedAt      time.Time        `json:"created_at,omitempty"`
	UpdatedAt      time.Time        `json:"updated_at,omitempty"`
}

func OrderResponseFromModel(order *models.Order) *OrderResponseDto {
	return &OrderResponseDto{
		OrderID:        order.OrderID,
		UserID:         order.UserID,
		LibrarianID:    order.LibrarianID,
		Item:           order.Item,
		Status:         order.Status,
		PickupSchedule: order.PickupSchedule,
		CreatedAt:      order.CreatedAt,
		UpdatedAt:      order.UpdatedAt,
	}
}
