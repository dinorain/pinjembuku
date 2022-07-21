//go:generate mockgen -source usecase.go -destination mock/usecase.go -package mock
package book

import (
	"context"

	"github.com/dinorain/pinjembuku/internal/models"
	"github.com/dinorain/pinjembuku/pkg/utils"
)

//  Book UseCase interface
type BookUseCase interface {
	FindAllBySubject(ctx context.Context, subject string, pagination *utils.Pagination) ([]models.Book, error)
	FindByWork(ctx context.Context, bookKey string) (*models.Book, error)
}
