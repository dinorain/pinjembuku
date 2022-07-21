package repository

import (
	"context"
	"database/sql"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"

	"github.com/dinorain/pinjembuku/internal/models"
	"github.com/dinorain/pinjembuku/internal/librarian"
	"github.com/dinorain/pinjembuku/pkg/utils"
)

// Librarian repository
type LibrarianRepository struct {
	db *sqlx.DB
}

var _ librarian.LibrarianPGRepository = (*LibrarianRepository)(nil)

// Librarian repository constructor
func NewLibrarianPGRepository(db *sqlx.DB) *LibrarianRepository {
	return &LibrarianRepository{db: db}
}

// Create new librarian
func (r *LibrarianRepository) Create(ctx context.Context, librarian *models.Librarian) (*models.Librarian, error) {
	createdLibrarian := &models.Librarian{}
	if err := r.db.QueryRowxContext(
		ctx,
		createLibrarianQuery,
		librarian.FirstName,
		librarian.LastName,
		librarian.Email,
		librarian.Password,
		librarian.Avatar,
	).StructScan(createdLibrarian); err != nil {
		return nil, errors.Wrap(err, "LibrarianRepository.Create.QueryRowxContext")
	}

	return createdLibrarian, nil
}

// UpdateById update existing librarian
func (r *LibrarianRepository) UpdateById(ctx context.Context, librarian *models.Librarian) (*models.Librarian, error) {
	if res, err := r.db.ExecContext(
		ctx,
		updateByIdQuery,
		librarian.LibrarianID,
		librarian.FirstName,
		librarian.LastName,
		librarian.Email,
		librarian.Password,
		librarian.Avatar,
	); err != nil {
		return nil, errors.Wrap(err, "UpdateById.Update.ExecContext")
	} else {
		_, err := res.RowsAffected()
		if err != nil {
			return nil, errors.Wrap(err, "UpdateById.Update.RowsAffected")
		}
	}

	return librarian, nil
}

// FindAll Find librarians
func (r *LibrarianRepository) FindAll(ctx context.Context, pagination *utils.Pagination) ([]models.Librarian, error) {
	var librarians []models.Librarian
	if err := r.db.SelectContext(ctx, &librarians, findAllQuery, pagination.GetLimit(), pagination.GetOffset()); err != nil {
		return nil, errors.Wrap(err, "LibrarianRepository.FindById.SelectContext")
	}

	return librarians, nil
}

// FindByEmail Find by librarian email address
func (r *LibrarianRepository) FindByEmail(ctx context.Context, email string) (*models.Librarian, error) {
	librarian := &models.Librarian{}
	if err := r.db.GetContext(ctx, librarian, findByEmailQuery, email); err != nil {
		return nil, errors.Wrap(err, "FindByEmail.GetContext")
	}

	return librarian, nil
}

// FindById Find librarian by uuid
func (r *LibrarianRepository) FindById(ctx context.Context, librarianID uuid.UUID) (*models.Librarian, error) {
	librarian := &models.Librarian{}
	if err := r.db.GetContext(ctx, librarian, findByIdQuery, librarianID); err != nil {
		return nil, errors.Wrap(err, "LibrarianRepository.FindById.GetContext")
	}

	return librarian, nil
}

// DeleteById Find librarian by uuid
func (r *LibrarianRepository) DeleteById(ctx context.Context, librarianID uuid.UUID) error {
	if res, err := r.db.ExecContext(ctx, deleteByIdQuery, librarianID); err != nil {
		return errors.Wrap(err, "LibrarianRepository.DeleteById.ExecContext")
	} else {
		cnt, err := res.RowsAffected()
		if err != nil {
			return errors.Wrap(err, "LibrarianRepository.DeleteById.RowsAffected")
		} else if cnt == 0 {
			return sql.ErrNoRows
		}
	}

	return nil
}
