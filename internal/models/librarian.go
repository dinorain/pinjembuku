package models

import (
	"strings"
	"time"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

// Librarian model
type Librarian struct {
	LibrarianID uuid.UUID `json:"librarian_id" db:"librarian_id"`
	Email       string    `json:"email" db:"email"`
	FirstName   string    `json:"first_name" db:"first_name"`
	LastName    string    `json:"last_name" db:"last_name"`
	Avatar      *string   `json:"avatar" db:"avatar"`
	Password    string    `json:"-" db:"password"`
	CreatedAt   time.Time `json:"created_at,omitempty" db:"created_at"`
	UpdatedAt   time.Time `json:"updated_at,omitempty" db:"updated_at"`
}

func (s *Librarian) SanitizePassword() {
	s.Password = ""
}

func (s *Librarian) HashPassword() error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(s.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	s.Password = string(hashedPassword)
	return nil
}

func (s *Librarian) ComparePasswords(password string) error {
	return bcrypt.CompareHashAndPassword([]byte(s.Password), []byte(password))
}

func (s *Librarian) PrepareCreate() error {
	s.Email = strings.ToLower(strings.TrimSpace(s.Email))
	s.Password = strings.TrimSpace(s.Password)

	if err := s.HashPassword(); err != nil {
		return err
	}
	return nil
}

// Get avatar string
func (s *Librarian) GetAvatar() string {
	if s.Avatar == nil {
		return ""
	}
	return *s.Avatar
}
