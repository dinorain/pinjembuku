//go:generate mockgen -source redis_repository.go -destination mock/redis_repository.go -package mock
package librarian

import (
	"context"

	"github.com/dinorain/pinjembuku/internal/models"
)

// Librarian Redis repository interface
type LibrarianRedisRepository interface {
	GetByIdCtx(ctx context.Context, key string) (*models.Librarian, error)
	SetLibrarianCtx(ctx context.Context, key string, seconds int, user *models.Librarian) error
	DeleteLibrarianCtx(ctx context.Context, key string) error
}
