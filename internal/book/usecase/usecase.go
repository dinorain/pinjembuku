package usecase

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/dinorain/pinjembuku/config"
	"github.com/dinorain/pinjembuku/internal/book"
	"github.com/dinorain/pinjembuku/internal/models"
	"github.com/dinorain/pinjembuku/pkg/logger"
	"github.com/dinorain/pinjembuku/pkg/utils"
)

const (
	bookByIdCacheDuration = 3600
)

// Book UseCase
type bookUseCase struct {
	cfg    *config.Config
	logger logger.Logger
}

var _ book.BookUseCase = (*bookUseCase)(nil)

// New Book UseCase
func NewBookUseCase(cfg *config.Config, logger logger.Logger) *bookUseCase {
	return &bookUseCase{cfg: cfg, logger: logger}
}

// FindAllBySubject find books by subject id
func (u *bookUseCase) FindAllBySubject(ctx context.Context, subject string, pagination *utils.Pagination) ([]models.Book, error) {
	resp, err := http.Get(fmt.Sprintf("https://openlibrary.org/subjects/%v.json?offset=%v&limit=%v", subject, pagination.GetOffset(), pagination.GetLimit()))
	if err != nil {
		u.logger.Errorf("bookUseCase.http.Get: %v", err)
		return nil, err
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)

	type openlibraryReponseDto struct {
		Works []models.Book `json:"works"`
	}

	var bookResponseDto openlibraryReponseDto
	_ = json.Unmarshal(body, &bookResponseDto)

	return bookResponseDto.Works, nil
}

// FindById find book by uuid
func (u *bookUseCase) FindByWork(ctx context.Context, bookKey string) (*models.Book, error) {
	resp, err := http.Get(fmt.Sprintf("https://openlibrary.org/%v.json", bookKey))
	if err != nil {
		u.logger.Errorf("bookUseCase.http.Get: %v", err)
		return nil, err
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)

	var book models.Book
	_ = json.Unmarshal(body, &book)

	return &book, nil
}
