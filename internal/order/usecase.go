//go:generate mockgen -source usecase.go -destination mock/usecase.go -package mock
package order

import (
	"context"

	"github.com/google/uuid"

	"github.com/dinorain/pinjembuku/internal/models"
	"github.com/dinorain/pinjembuku/pkg/utils"
)

//  Order UseCase interface
type OrderUseCase interface {
	Create(ctx context.Context, order *models.Order) (*models.Order, error)
	FindAll(ctx context.Context, pagination *utils.Pagination) ([]models.Order, error)
	FindAllByUserId(ctx context.Context, userID uuid.UUID, pagination *utils.Pagination) ([]models.Order, error)
	FindAllByLibrarianId(ctx context.Context, librarianID uuid.UUID, pagination *utils.Pagination) ([]models.Order, error)
	FindAllByUserIdLibrarianId(ctx context.Context, userID uuid.UUID, librarianID uuid.UUID, pagination *utils.Pagination) ([]models.Order, error)
	FindById(ctx context.Context, orderID uuid.UUID) (*models.Order, error)
	CachedFindById(ctx context.Context, orderID uuid.UUID) (*models.Order, error)
	UpdateById(ctx context.Context, order *models.Order) (*models.Order, error)
	DeleteById(ctx context.Context, orderID uuid.UUID) error
}
