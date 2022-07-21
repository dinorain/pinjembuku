package repository

import (
	"context"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/require"

	"github.com/dinorain/pinjembuku/internal/models"
	"github.com/dinorain/pinjembuku/pkg/utils"
)

func TestLibrarianRepository_Create(t *testing.T) {
	t.Parallel()

	db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	require.NoError(t, err)
	defer db.Close()

	sqlxDB := sqlx.NewDb(db, "sqlmock")
	defer sqlxDB.Close()

	librarianPGRepository := NewLibrarianPGRepository(sqlxDB)

	columns := []string{"librarian_id", "first_name", "last_name", "email", "password", "avatar", "created_at", "updated_at"}
	librarianUUID := uuid.New()
	mockLibrarian := &models.Librarian{
		LibrarianID:      librarianUUID,
		Email:         "email@gmail.com",
		FirstName:     "FirstName",
		LastName:      "LastName",
		Avatar:        nil,
		Password:      "123456",
	}

	rows := sqlmock.NewRows(columns).AddRow(
		librarianUUID,
		mockLibrarian.FirstName,
		mockLibrarian.LastName,
		mockLibrarian.Email,
		mockLibrarian.Password,
		mockLibrarian.Avatar,
		time.Now(),
		time.Now(),
	)

	mock.ExpectQuery(createLibrarianQuery).WithArgs(
		mockLibrarian.FirstName,
		mockLibrarian.LastName,
		mockLibrarian.Email,
		mockLibrarian.Password,
		mockLibrarian.Avatar,
	).WillReturnRows(rows)

	createdLibrarian, err := librarianPGRepository.Create(context.Background(), mockLibrarian)
	require.NoError(t, err)
	require.NotNil(t, createdLibrarian)
}

func TestLibrarianRepository_FindByEmail(t *testing.T) {
	t.Parallel()

	db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	require.NoError(t, err)
	defer db.Close()

	sqlxDB := sqlx.NewDb(db, "sqlmock")
	defer sqlxDB.Close()

	librarianPGRepository := NewLibrarianPGRepository(sqlxDB)

	columns := []string{"librarian_id", "first_name", "last_name", "email", "password", "avatar", "created_at", "updated_at"}
	librarianUUID := uuid.New()
	mockLibrarian := &models.Librarian{
		LibrarianID:      librarianUUID,
		Email:         "email@gmail.com",
		FirstName:     "FirstName",
		LastName:      "LastName",
		Avatar:        nil,
		Password:      "123456",
	}

	rows := sqlmock.NewRows(columns).AddRow(
		librarianUUID,
		mockLibrarian.FirstName,
		mockLibrarian.LastName,
		mockLibrarian.Email,
		mockLibrarian.Password,
		mockLibrarian.Avatar,
		time.Now(),
		time.Now(),
	)

	mock.ExpectQuery(findByEmailQuery).WithArgs(mockLibrarian.Email).WillReturnRows(rows)

	foundLibrarian, err := librarianPGRepository.FindByEmail(context.Background(), mockLibrarian.Email)
	require.NoError(t, err)
	require.NotNil(t, foundLibrarian)
	require.Equal(t, foundLibrarian.Email, mockLibrarian.Email)
}

func TestLibrarianRepository_FindAll(t *testing.T) {
	t.Parallel()

	db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	require.NoError(t, err)
	defer db.Close()

	sqlxDB := sqlx.NewDb(db, "sqlmock")
	defer sqlxDB.Close()

	librarianPGRepository := NewLibrarianPGRepository(sqlxDB)

	columns := []string{"librarian_id", "first_name", "last_name", "email", "password", "avatar", "created_at", "updated_at"}
	librarianUUID := uuid.New()
	mockLibrarian := &models.Librarian{
		LibrarianID:      librarianUUID,
		Email:         "email@gmail.com",
		FirstName:     "FirstName",
		LastName:      "LastName",
		Avatar:        nil,
		Password:      "123456",
	}

	rows := sqlmock.NewRows(columns).AddRow(
		librarianUUID,
		mockLibrarian.FirstName,
		mockLibrarian.LastName,
		mockLibrarian.Email,
		mockLibrarian.Password,
		mockLibrarian.Avatar,
		time.Now(),
		time.Now(),
	)

	size := 10
	mock.ExpectQuery(findAllQuery).WithArgs(size, 0).WillReturnRows(rows)
	foundLibrarians, err := librarianPGRepository.FindAll(context.Background(), utils.NewPaginationQuery(size, 1))
	require.NoError(t, err)
	require.NotNil(t, foundLibrarians)
	require.Equal(t, len(foundLibrarians), 1)

	mock.ExpectQuery(findAllQuery).WithArgs(size, 10).WillReturnRows(rows)
	foundLibrarians, err = librarianPGRepository.FindAll(context.Background(), utils.NewPaginationQuery(size, 2))
	require.NoError(t, err)
	require.Nil(t, foundLibrarians)
}

func TestLibrarianRepository_FindById(t *testing.T) {
	t.Parallel()

	db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	require.NoError(t, err)
	defer db.Close()

	sqlxDB := sqlx.NewDb(db, "sqlmock")
	defer sqlxDB.Close()

	librarianPGRepository := NewLibrarianPGRepository(sqlxDB)

	columns := []string{"librarian_id", "first_name", "last_name", "email", "password", "avatar", "created_at", "updated_at"}
	librarianUUID := uuid.New()
	mockLibrarian := &models.Librarian{
		LibrarianID:      librarianUUID,
		Email:         "email@gmail.com",
		FirstName:     "FirstName",
		LastName:      "LastName",
		Avatar:        nil,
		Password:      "123456",
	}

	rows := sqlmock.NewRows(columns).AddRow(
		librarianUUID,
		mockLibrarian.FirstName,
		mockLibrarian.LastName,
		mockLibrarian.Email,
		mockLibrarian.Password,
		mockLibrarian.Avatar,
		time.Now(),
		time.Now(),
	)

	mock.ExpectQuery(findByIdQuery).WithArgs(mockLibrarian.LibrarianID).WillReturnRows(rows)

	foundLibrarian, err := librarianPGRepository.FindById(context.Background(), mockLibrarian.LibrarianID)
	require.NoError(t, err)
	require.NotNil(t, foundLibrarian)
	require.Equal(t, foundLibrarian.LibrarianID, mockLibrarian.LibrarianID)
}

func TestLibrarianRepository_UpdateById(t *testing.T) {
	t.Parallel()

	db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	require.NoError(t, err)
	defer db.Close()

	sqlxDB := sqlx.NewDb(db, "sqlmock")
	defer sqlxDB.Close()

	librarianPGRepository := NewLibrarianPGRepository(sqlxDB)

	columns := []string{"librarian_id", "first_name", "last_name", "email", "password", "avatar", "created_at", "updated_at"}
	librarianUUID := uuid.New()
	mockLibrarian := &models.Librarian{
		LibrarianID:      librarianUUID,
		Email:         "email@gmail.com",
		FirstName:     "FirstName",
		LastName:      "LastName",
		Avatar:        nil,
		Password:      "123456",
	}

	_ = sqlmock.NewRows(columns).AddRow(
		librarianUUID,
		mockLibrarian.FirstName,
		mockLibrarian.LastName,
		mockLibrarian.Email,
		mockLibrarian.Password,
		mockLibrarian.Avatar,
		time.Now(),
		time.Now(),
	)

	mockLibrarian.FirstName = "FirstNameChanged"
	mock.ExpectExec(updateByIdQuery).WithArgs(
		mockLibrarian.LibrarianID,
		mockLibrarian.FirstName,
		mockLibrarian.LastName,
		mockLibrarian.Email,
		mockLibrarian.Password,
		mockLibrarian.Avatar,
	).WillReturnResult(sqlmock.NewResult(0, 1))

	updatedLibrarian, err := librarianPGRepository.UpdateById(context.Background(), mockLibrarian)
	require.NoError(t, err)
	require.NotNil(t, mockLibrarian)
	require.Equal(t, updatedLibrarian.FirstName, mockLibrarian.FirstName)
	require.Equal(t, updatedLibrarian.LibrarianID, mockLibrarian.LibrarianID)
}

func TestLibrarianRepository_DeleteById(t *testing.T) {
	t.Parallel()

	db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	require.NoError(t, err)
	defer db.Close()

	sqlxDB := sqlx.NewDb(db, "sqlmock")
	defer sqlxDB.Close()

	librarianPGRepository := NewLibrarianPGRepository(sqlxDB)

	columns := []string{"librarian_id", "first_name", "last_name", "email", "password", "avatar", "created_at", "updated_at"}
	librarianUUID := uuid.New()
	mockLibrarian := &models.Librarian{
		LibrarianID:      librarianUUID,
		Email:         "email@gmail.com",
		FirstName:     "FirstName",
		LastName:      "LastName",
		Avatar:        nil,
		Password:      "123456",
	}

	_ = sqlmock.NewRows(columns).AddRow(
		librarianUUID,
		mockLibrarian.FirstName,
		mockLibrarian.LastName,
		mockLibrarian.Email,
		mockLibrarian.Password,
		mockLibrarian.Avatar,
		time.Now(),
		time.Now(),
	)

	mock.ExpectExec(deleteByIdQuery).WithArgs(mockLibrarian.LibrarianID).WillReturnResult(sqlmock.NewResult(0, 1))

	err = librarianPGRepository.DeleteById(context.Background(), mockLibrarian.LibrarianID)
	require.NoError(t, err)
	require.NotNil(t, mockLibrarian)
}
