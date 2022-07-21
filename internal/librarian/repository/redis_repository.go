package repository

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/go-redis/redis/v8"

	"github.com/dinorain/pinjembuku/internal/models"
	"github.com/dinorain/pinjembuku/internal/librarian"
	"github.com/dinorain/pinjembuku/pkg/logger"
)

// Librarian redis repository
type librarianRedisRepo struct {
	redisClient *redis.Client
	basePrefix  string
	logger      logger.Logger
}

var _ librarian.LibrarianRedisRepository = (*librarianRedisRepo)(nil)

// Librarian redis repository constructor
func NewLibrarianRedisRepo(redisClient *redis.Client, logger logger.Logger) *librarianRedisRepo {
	return &librarianRedisRepo{redisClient: redisClient, basePrefix: "librarian:", logger: logger}
}

// Get librarian by id
func (r *librarianRedisRepo) GetByIdCtx(ctx context.Context, key string) (*models.Librarian, error) {
	librarianBytes, err := r.redisClient.Get(ctx, r.createKey(key)).Bytes()
	if err != nil {
		return nil, err
	}
	librarian := &models.Librarian{}
	if err = json.Unmarshal(librarianBytes, librarian); err != nil {
		return nil, err
	}

	return librarian, nil
}

// Cache librarian with duration in seconds
func (r *librarianRedisRepo) SetLibrarianCtx(ctx context.Context, key string, seconds int, librarian *models.Librarian) error {
	librarianBytes, err := json.Marshal(librarian)
	if err != nil {
		return err
	}

	return r.redisClient.Set(ctx, r.createKey(key), librarianBytes, time.Second*time.Duration(seconds)).Err()
}

// Delete librarian by key
func (r *librarianRedisRepo) DeleteLibrarianCtx(ctx context.Context, key string) error {
	return r.redisClient.Del(ctx, r.createKey(key)).Err()
}

func (r *librarianRedisRepo) createKey(value string) string {
	return fmt.Sprintf("%s: %s", r.basePrefix, value)
}
