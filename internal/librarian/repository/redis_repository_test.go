package repository

import (
	"context"
	"log"
	"testing"

	"github.com/alicebob/miniredis"
	"github.com/go-redis/redis/v8"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"

	"github.com/dinorain/pinjembuku/internal/models"
)

func SetupRedis() *librarianRedisRepo {
	mr, err := miniredis.Run()
	if err != nil {
		log.Fatal(err)
	}
	client := redis.NewClient(&redis.Options{
		Addr: mr.Addr(),
	})

	librarianRedisRepository := NewLibrarianRedisRepo(client, nil)
	return librarianRedisRepository
}

func TestLibrarianRedisRepo_SetLibrarianCtx(t *testing.T) {
	t.Parallel()

	redisRepo := SetupRedis()

	t.Run("SetLibrarianCtx", func(t *testing.T) {
		librarian := &models.Librarian{
			LibrarianID: uuid.New(),
		}

		err := redisRepo.SetLibrarianCtx(context.Background(), redisRepo.createKey(librarian.LibrarianID.String()), 10, librarian)
		require.NoError(t, err)
	})
}

func TestLibrarianRedisRepo_GetByIdCtx(t *testing.T) {
	t.Parallel()

	redisRepo := SetupRedis()

	t.Run("GetByIdCtx", func(t *testing.T) {
		librarian := &models.Librarian{
			LibrarianID: uuid.New(),
		}

		err := redisRepo.SetLibrarianCtx(context.Background(), redisRepo.createKey(librarian.LibrarianID.String()), 10, librarian)
		require.NoError(t, err)

		librarian, err = redisRepo.GetByIdCtx(context.Background(), redisRepo.createKey(librarian.LibrarianID.String()))
		require.NoError(t, err)
		require.NotNil(t, librarian)
	})
}

func TestLibrarianRedisRepo_DeleteLibrarianCtx(t *testing.T) {
	t.Parallel()

	redisRepo := SetupRedis()

	t.Run("DeleteLibrarianCtx", func(t *testing.T) {
		librarian := &models.Librarian{
			LibrarianID: uuid.New(),
		}

		err := redisRepo.DeleteLibrarianCtx(context.Background(), redisRepo.createKey(librarian.LibrarianID.String()))
		require.NoError(t, err)
	})
}
