package usecase

import (
	"context"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/golang-jwt/jwt"
	"github.com/google/uuid"
	"github.com/pkg/errors"

	"github.com/dinorain/pinjembuku/config"
	"github.com/dinorain/pinjembuku/internal/models"
	"github.com/dinorain/pinjembuku/internal/librarian"
	"github.com/dinorain/pinjembuku/pkg/grpc_errors"
	"github.com/dinorain/pinjembuku/pkg/logger"
	"github.com/dinorain/pinjembuku/pkg/utils"
)

const (
	librarianByIdCacheDuration = 3600
)

// Librarian UseCase
type librarianUseCase struct {
	cfg          *config.Config
	logger       logger.Logger
	librarianPgRepo librarian.LibrarianPGRepository
	redisRepo    librarian.LibrarianRedisRepository
}

var _ librarian.LibrarianUseCase = (*librarianUseCase)(nil)

// New Librarian UseCase
func NewLibrarianUseCase(cfg *config.Config, logger logger.Logger, librarianRepo librarian.LibrarianPGRepository, redisRepo librarian.LibrarianRedisRepository) *librarianUseCase {
	return &librarianUseCase{cfg: cfg, logger: logger, librarianPgRepo: librarianRepo, redisRepo: redisRepo}
}

// Register new librarian
func (u *librarianUseCase) Register(ctx context.Context, librarian *models.Librarian) (*models.Librarian, error) {
	existsLibrarian, err := u.librarianPgRepo.FindByEmail(ctx, librarian.Email)
	if existsLibrarian != nil || err == nil {
		return nil, grpc_errors.ErrEmailExists
	}

	return u.librarianPgRepo.Create(ctx, librarian)
}

// FindAll find librarians
func (u *librarianUseCase) FindAll(ctx context.Context, pagination *utils.Pagination) ([]models.Librarian, error) {
	librarians, err := u.librarianPgRepo.FindAll(ctx, pagination)
	if err != nil {
		return nil, errors.Wrap(err, "librarianPgRepo.FindAll")
	}

	return librarians, nil
}

// FindByEmail find librarian by email address
func (u *librarianUseCase) FindByEmail(ctx context.Context, email string) (*models.Librarian, error) {
	findByEmail, err := u.librarianPgRepo.FindByEmail(ctx, email)
	if err != nil {
		return nil, errors.Wrap(err, "librarianPgRepo.FindByEmail")
	}

	findByEmail.SanitizePassword()

	return findByEmail, nil
}

// FindById find librarian by uuid
func (u *librarianUseCase) FindById(ctx context.Context, librarianID uuid.UUID) (*models.Librarian, error) {
	foundLibrarian, err := u.librarianPgRepo.FindById(ctx, librarianID)
	if err != nil {
		return nil, errors.Wrap(err, "librarianPgRepo.FindById")
	}

	return foundLibrarian, nil
}

// CachedFindById find librarian by uuid from cache
func (u *librarianUseCase) CachedFindById(ctx context.Context, librarianID uuid.UUID) (*models.Librarian, error) {
	cachedLibrarian, err := u.redisRepo.GetByIdCtx(ctx, librarianID.String())
	if err != nil && !errors.Is(err, redis.Nil) {
		u.logger.Errorf("redisRepo.GetByIdCtx", err)
	}
	if cachedLibrarian != nil {
		return cachedLibrarian, nil
	}

	foundLibrarian, err := u.librarianPgRepo.FindById(ctx, librarianID)
	if err != nil {
		return nil, errors.Wrap(err, "librarianPgRepo.FindById")
	}

	if err := u.redisRepo.SetLibrarianCtx(ctx, foundLibrarian.LibrarianID.String(), librarianByIdCacheDuration, foundLibrarian); err != nil {
		u.logger.Errorf("redisRepo.SetLibrarianCtx", err)
	}

	return foundLibrarian, nil
}

// UpdateById update librarian by uuid
func (u *librarianUseCase) UpdateById(ctx context.Context, librarian *models.Librarian) (*models.Librarian, error) {
	updatedLibrarian, err := u.librarianPgRepo.UpdateById(ctx, librarian)
	if err != nil {
		return nil, errors.Wrap(err, "librarianPgRepo.UpdateById")
	}

	if err := u.redisRepo.SetLibrarianCtx(ctx, updatedLibrarian.LibrarianID.String(), librarianByIdCacheDuration, updatedLibrarian); err != nil {
		u.logger.Errorf("redisRepo.SetLibrarianCtx", err)
	}

	updatedLibrarian.SanitizePassword()

	return updatedLibrarian, nil
}

// DeleteById delete librarian by uuid
func (u *librarianUseCase) DeleteById(ctx context.Context, librarianID uuid.UUID) error {
	err := u.librarianPgRepo.DeleteById(ctx, librarianID)
	if err != nil {
		return errors.Wrap(err, "librarianPgRepo.DeleteById")
	}

	if err := u.redisRepo.DeleteLibrarianCtx(ctx, librarianID.String()); err != nil {
		u.logger.Errorf("redisRepo.DeleteLibrarianCtx", err)
	}

	return nil
}

// Login librarian with email and password
func (u *librarianUseCase) Login(ctx context.Context, email string, password string) (*models.Librarian, error) {
	foundLibrarian, err := u.librarianPgRepo.FindByEmail(ctx, email)
	if err != nil {
		return nil, errors.Wrap(err, "librarianPgRepo.FindByEmail")
	}

	if err := foundLibrarian.ComparePasswords(password); err != nil {
		return nil, errors.Wrap(err, "librarian.ComparePasswords")
	}

	return foundLibrarian, err
}

func (u *librarianUseCase) GenerateTokenPair(librarian *models.Librarian, sessionID string) (access string, refresh string, err error) {
	token := jwt.New(jwt.SigningMethodHS256)

	claims := token.Claims.(jwt.MapClaims)
	claims["session_id"] = sessionID
	claims["librarian_id"] = librarian.LibrarianID
	claims["email"] = librarian.Email
	claims["exp"] = time.Now().Add(time.Minute * 15).Unix()

	access, err = token.SignedString([]byte(u.cfg.Server.JwtSecretKey))
	if err != nil {
		return "", "", err
	}

	refreshToken := jwt.New(jwt.SigningMethodHS256)
	rtClaims := refreshToken.Claims.(jwt.MapClaims)
	rtClaims["session_id"] = sessionID
	rtClaims["exp"] = time.Now().Add(time.Hour * 24).Unix()

	refresh, err = refreshToken.SignedString([]byte(u.cfg.Server.JwtSecretKey))
	if err != nil {
		return "", "", err
	}

	return access, refresh, nil
}
