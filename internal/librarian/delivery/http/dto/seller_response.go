package dto

import (
	"time"

	"github.com/google/uuid"

	"github.com/dinorain/pinjembuku/internal/models"
)

type LibrarianResponseDto struct {
	LibrarianID      uuid.UUID `json:"librarian_id"`
	Email         string    `json:"email"`
	FirstName     string    `json:"first_name"`
	LastName      string    `json:"last_name"`
	Avatar        *string   `json:"avatar"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}

func LibrarianResponseFromModel(librarian *models.Librarian) *LibrarianResponseDto {
	return &LibrarianResponseDto{
		LibrarianID:      librarian.LibrarianID,
		Email:         librarian.Email,
		FirstName:     librarian.FirstName,
		LastName:      librarian.LastName,
		Avatar:        librarian.Avatar,
		CreatedAt:     librarian.CreatedAt,
		UpdatedAt:     librarian.UpdatedAt,
	}
}
