//go:generate mockgen -source pg_repository.go -destination mock/pg_repository.go -package mock
package librarian

import (
	"context"

	"github.com/google/uuid"

	"github.com/dinorain/pinjembuku/internal/models"
	"github.com/dinorain/pinjembuku/pkg/utils"
)

// Librarian pg repository
type LibrarianPGRepository interface {
	Create(ctx context.Context, user *models.Librarian) (*models.Librarian, error)
	FindAll(ctx context.Context, pagination *utils.Pagination) ([]models.Librarian, error)
	FindByEmail(ctx context.Context, email string) (*models.Librarian, error)
	FindById(ctx context.Context, userID uuid.UUID) (*models.Librarian, error)
	UpdateById(ctx context.Context, user *models.Librarian) (*models.Librarian, error)
	DeleteById(ctx context.Context, userID uuid.UUID) error
}
