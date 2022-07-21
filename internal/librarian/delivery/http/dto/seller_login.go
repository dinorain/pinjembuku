package dto

import (
	"github.com/google/uuid"
)

type LibrarianLoginRequestDto struct {
	Email    string `json:"email" validate:"required,lte=60,email"`
	Password string `json:"password" validate:"required"`
}

type LibrarianLoginResponseDto struct {
	LibrarianID uuid.UUID                      `json:"user_id" validate:"required"`
	Tokens   *LibrarianRefreshTokenResponseDto `json:"tokens" validate:"required"`
}
