//go:generate mockgen -source usecase.go -destination mock/usecase.go -package mock
package librarian

import (
	"context"

	"github.com/google/uuid"

	"github.com/dinorain/pinjembuku/internal/models"
	"github.com/dinorain/pinjembuku/pkg/utils"
)

//  Librarian UseCase interface
type LibrarianUseCase interface {
	Register(ctx context.Context, librarian *models.Librarian) (*models.Librarian, error)
	Login(ctx context.Context, email string, password string) (*models.Librarian, error)
	FindAll(ctx context.Context, pagination *utils.Pagination) ([]models.Librarian, error)
	FindByEmail(ctx context.Context, email string) (*models.Librarian, error)
	FindById(ctx context.Context, librarianID uuid.UUID) (*models.Librarian, error)
	CachedFindById(ctx context.Context, librarianID uuid.UUID) (*models.Librarian, error)
	UpdateById(ctx context.Context, librarian *models.Librarian) (*models.Librarian, error)
	DeleteById(ctx context.Context, librarianID uuid.UUID) error
	GenerateTokenPair(librarian *models.Librarian, sessionID string) (access string, refresh string, err error)
}
