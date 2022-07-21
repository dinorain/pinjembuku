package usecase

import (
	"context"
	"database/sql"
	"testing"

	"github.com/go-redis/redis/v8"
	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"

	"github.com/dinorain/pinjembuku/config"
	"github.com/dinorain/pinjembuku/internal/models"
	"github.com/dinorain/pinjembuku/internal/librarian/mock"
	"github.com/dinorain/pinjembuku/pkg/logger"
)

func TestLibrarianUseCase_Register(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	librarianPGRepository := mock.NewMockLibrarianPGRepository(ctrl)
	librarianRedisRepository := mock.NewMockLibrarianRedisRepository(ctrl)
	apiLogger := logger.NewAppLogger(nil)

	cfg := &config.Config{Server: config.ServerConfig{JwtSecretKey: "secret123"}}
	librarianUC := NewLibrarianUseCase(cfg, apiLogger, librarianPGRepository, librarianRedisRepository)

	librarianID := uuid.New()
	mockLibrarian := &models.Librarian{
		Email:         "email@gmail.com",
		FirstName:     "FirstName",
		LastName:      "LastName",
		Avatar:        nil,
		Password:      "123456",
	}

	ctx := context.Background()

	librarianPGRepository.EXPECT().FindByEmail(gomock.Any(), mockLibrarian.Email).Return(nil, sql.ErrNoRows)

	librarianPGRepository.EXPECT().Create(gomock.Any(), mockLibrarian).Return(&models.Librarian{
		LibrarianID:      librarianID,
		Email:         "email@gmail.com",
		FirstName:     "FirstName",
		LastName:      "LastName",
		Avatar:        nil,
		Password:      "123456",
	}, nil)

	createdLibrarian, err := librarianUC.Register(ctx, mockLibrarian)
	require.NoError(t, err)
	require.NotNil(t, createdLibrarian)
	require.Equal(t, createdLibrarian.LibrarianID, librarianID)
}

func TestLibrarianUseCase_FindByEmail(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	librarianPGRepository := mock.NewMockLibrarianPGRepository(ctrl)
	librarianRedisRepository := mock.NewMockLibrarianRedisRepository(ctrl)
	apiLogger := logger.NewAppLogger(nil)

	cfg := &config.Config{Server: config.ServerConfig{JwtSecretKey: "secret123"}}
	librarianUC := NewLibrarianUseCase(cfg, apiLogger, librarianPGRepository, librarianRedisRepository)

	librarianID := uuid.New()
	mockLibrarian := &models.Librarian{
		LibrarianID:      librarianID,
		Email:         "email@gmail.com",
		FirstName:     "FirstName",
		LastName:      "LastName",
		Avatar:        nil,
		Password:      "123456",
	}

	ctx := context.Background()

	librarianPGRepository.EXPECT().FindByEmail(gomock.Any(), mockLibrarian.Email).Return(mockLibrarian, nil)

	librarian, err := librarianUC.FindByEmail(ctx, mockLibrarian.Email)
	require.NoError(t, err)
	require.NotNil(t, librarian)
	require.Equal(t, librarian.Email, mockLibrarian.Email)
}

func TestLibrarianUseCase_Login(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	librarianPGRepository := mock.NewMockLibrarianPGRepository(ctrl)
	librarianRedisRepository := mock.NewMockLibrarianRedisRepository(ctrl)
	apiLogger := logger.NewAppLogger(nil)

	cfg := &config.Config{Server: config.ServerConfig{JwtSecretKey: "secret123"}}
	librarianUC := NewLibrarianUseCase(cfg, apiLogger, librarianPGRepository, librarianRedisRepository)

	librarianID := uuid.New()
	mockLibrarian := &models.Librarian{
		LibrarianID:      librarianID,
		Email:         "email@gmail.com",
		FirstName:     "FirstName",
		LastName:      "LastName",
		Avatar:        nil,
		Password:      "123456",
	}

	ctx := context.Background()

	librarianPGRepository.EXPECT().FindByEmail(gomock.Any(), mockLibrarian.Email).Return(mockLibrarian, nil)
	_, err := librarianUC.Login(ctx, mockLibrarian.Email, mockLibrarian.Password)
	require.NotNil(t, err)
}

func TestLibrarianUseCase_FindAll(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	librarianPGRepository := mock.NewMockLibrarianPGRepository(ctrl)
	librarianRedisRepository := mock.NewMockLibrarianRedisRepository(ctrl)
	apiLogger := logger.NewAppLogger(nil)

	cfg := &config.Config{}
	librarianUC := NewLibrarianUseCase(cfg, apiLogger, librarianPGRepository, librarianRedisRepository)

	librarianID := uuid.New()
	mockLibrarian := &models.Librarian{
		LibrarianID:      librarianID,
		Email:         "email@gmail.com",
		FirstName:     "FirstName",
		LastName:      "LastName",
		Avatar:        nil,
		Password:      "123456",
	}

	ctx := context.Background()

	librarianPGRepository.EXPECT().FindAll(gomock.Any(), nil).AnyTimes().Return(append([]models.Librarian{}, *mockLibrarian), nil)

	librarians, err := librarianUC.FindAll(ctx, nil)
	require.NoError(t, err)
	require.NotNil(t, librarians)
	require.Equal(t, len(librarians), 1)
}

func TestLibrarianUseCase_FindById(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	librarianPGRepository := mock.NewMockLibrarianPGRepository(ctrl)
	librarianRedisRepository := mock.NewMockLibrarianRedisRepository(ctrl)
	apiLogger := logger.NewAppLogger(nil)

	cfg := &config.Config{}
	librarianUC := NewLibrarianUseCase(cfg, apiLogger, librarianPGRepository, librarianRedisRepository)

	librarianID := uuid.New()
	mockLibrarian := &models.Librarian{
		LibrarianID:      librarianID,
		Email:         "email@gmail.com",
		FirstName:     "FirstName",
		LastName:      "LastName",
		Avatar:        nil,
		Password:      "123456",
	}

	ctx := context.Background()

	librarianRedisRepository.EXPECT().GetByIdCtx(gomock.Any(), mockLibrarian.LibrarianID.String()).AnyTimes().Return(nil, redis.Nil)
	librarianPGRepository.EXPECT().FindById(gomock.Any(), mockLibrarian.LibrarianID).Return(mockLibrarian, nil)

	librarian, err := librarianUC.FindById(ctx, mockLibrarian.LibrarianID)
	require.NoError(t, err)
	require.NotNil(t, librarian)
	require.Equal(t, librarian.LibrarianID, mockLibrarian.LibrarianID)

	librarianRedisRepository.EXPECT().GetByIdCtx(gomock.Any(), mockLibrarian.LibrarianID.String()).AnyTimes().Return(nil, redis.Nil)
}

func TestLibrarianUseCase_CachedFindById(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	librarianPGRepository := mock.NewMockLibrarianPGRepository(ctrl)
	librarianRedisRepository := mock.NewMockLibrarianRedisRepository(ctrl)
	apiLogger := logger.NewAppLogger(nil)

	cfg := &config.Config{}
	librarianUC := NewLibrarianUseCase(cfg, apiLogger, librarianPGRepository, librarianRedisRepository)

	librarianID := uuid.New()
	mockLibrarian := &models.Librarian{
		LibrarianID:      librarianID,
		Email:         "email@gmail.com",
		FirstName:     "FirstName",
		LastName:      "LastName",
		Avatar:        nil,
		Password:      "123456",
	}

	ctx := context.Background()

	librarianRedisRepository.EXPECT().GetByIdCtx(gomock.Any(), mockLibrarian.LibrarianID.String()).AnyTimes().Return(nil, redis.Nil)
	librarianPGRepository.EXPECT().FindById(gomock.Any(), mockLibrarian.LibrarianID).Return(mockLibrarian, nil)
	librarianRedisRepository.EXPECT().SetLibrarianCtx(gomock.Any(), mockLibrarian.LibrarianID.String(), 3600, mockLibrarian).AnyTimes().Return(nil)

	librarian, err := librarianUC.CachedFindById(ctx, mockLibrarian.LibrarianID)
	require.NoError(t, err)
	require.NotNil(t, librarian)
	require.Equal(t, librarian.LibrarianID, mockLibrarian.LibrarianID)
}

func TestLibrarianUseCase_UpdateById(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	librarianPGRepository := mock.NewMockLibrarianPGRepository(ctrl)
	librarianRedisRepository := mock.NewMockLibrarianRedisRepository(ctrl)
	apiLogger := logger.NewAppLogger(nil)

	cfg := &config.Config{}
	librarianUC := NewLibrarianUseCase(cfg, apiLogger, librarianPGRepository, librarianRedisRepository)

	librarianID := uuid.New()
	mockLibrarian := &models.Librarian{
		LibrarianID:      librarianID,
		Email:         "email@gmail.com",
		FirstName:     "FirstName",
		LastName:      "LastName",
		Avatar:        nil,
		Password:      "123456",
	}

	ctx := context.Background()

	librarianPGRepository.EXPECT().UpdateById(gomock.Any(), mockLibrarian).Return(mockLibrarian, nil)
	librarianRedisRepository.EXPECT().SetLibrarianCtx(gomock.Any(), mockLibrarian.LibrarianID.String(), 3600, mockLibrarian).AnyTimes().Return(nil)

	librarian, err := librarianUC.UpdateById(ctx, mockLibrarian)
	require.NoError(t, err)
	require.NotNil(t, librarian)
	require.Equal(t, librarian.LibrarianID, mockLibrarian.LibrarianID)
}

func TestLibrarianUseCase_DeleteById(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	librarianPGRepository := mock.NewMockLibrarianPGRepository(ctrl)
	librarianRedisRepository := mock.NewMockLibrarianRedisRepository(ctrl)
	apiLogger := logger.NewAppLogger(nil)

	cfg := &config.Config{}
	librarianUC := NewLibrarianUseCase(cfg, apiLogger, librarianPGRepository, librarianRedisRepository)

	librarianID := uuid.New()
	mockLibrarian := &models.Librarian{
		LibrarianID:      librarianID,
		Email:         "email@gmail.com",
		FirstName:     "FirstName",
		LastName:      "LastName",
		Avatar:        nil,
		Password:      "123456",
	}

	ctx := context.Background()

	librarianPGRepository.EXPECT().DeleteById(gomock.Any(), mockLibrarian.LibrarianID).Return(nil)
	librarianRedisRepository.EXPECT().DeleteLibrarianCtx(gomock.Any(), mockLibrarian.LibrarianID.String()).AnyTimes().Return(nil)

	err := librarianUC.DeleteById(ctx, mockLibrarian.LibrarianID)
	require.NoError(t, err)

	librarianPGRepository.EXPECT().FindById(gomock.Any(), mockLibrarian.LibrarianID).AnyTimes().Return(nil, nil)
	librarianRedisRepository.EXPECT().GetByIdCtx(gomock.Any(), mockLibrarian.LibrarianID.String()).AnyTimes().Return(nil, redis.Nil)
}

func TestLibrarianUseCase_GenerateTokenPair(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	librarianPGRepository := mock.NewMockLibrarianPGRepository(ctrl)
	librarianRedisRepository := mock.NewMockLibrarianRedisRepository(ctrl)
	apiLogger := logger.NewAppLogger(nil)

	cfg := &config.Config{}
	librarianUC := NewLibrarianUseCase(cfg, apiLogger, librarianPGRepository, librarianRedisRepository)

	librarianID := uuid.New()
	mockLibrarian := &models.Librarian{
		LibrarianID:      librarianID,
		Email:         "email@gmail.com",
		FirstName:     "FirstName",
		LastName:      "LastName",
		Avatar:        nil,
		Password:      "123456",
	}

	at, rt, err := librarianUC.GenerateTokenPair(mockLibrarian, mockLibrarian.LibrarianID.String())
	require.NoError(t, err)
	require.NotEqual(t, at, "")
	require.NotEqual(t, rt, "")
}
