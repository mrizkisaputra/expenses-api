package repository

import (
	"context"
	"github.com/alicebob/miniredis/v2"
	"github.com/google/uuid"
	"github.com/mrizkisaputra/expenses-api/internal/user"
	"github.com/mrizkisaputra/expenses-api/internal/user/model"
	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/require"
	"log"
	"testing"
	"time"
)

var key = "api-user:idtest"

func SetupRedis() user.UserRedisRepository {
	server, err := miniredis.Run()
	if err != nil {
		log.Fatalf("could not connet to redis server, got error %v", err)
	}

	redisClient := redis.NewClient(&redis.Options{
		Addr: server.Addr(),
	})

	repository := NewUserRedisRepository(redisClient)
	return repository
}

func TestUserRedisRepository_Set(t *testing.T) {
	repo := SetupRedis()

	expire := time.Second * 30
	testUser := &model.User{
		Email:    "test@test.com",
		Password: "secret",
		Information: model.Information{
			FirstName: "muhammat",
			LastName:  "saputra",
		},
	}

	err := repo.Set(context.Background(), key, expire, testUser)
	require.Nil(t, err)
	require.NoError(t, err)
}

func TestUserRedisRepository_Get(t *testing.T) {
	repo := SetupRedis()

	// cache data
	expire := time.Second * 30
	testUser := &model.User{
		Email:    "test@test.com",
		Password: "secret",
		Information: model.Information{
			FirstName: "muhammat",
			LastName:  "saputra",
		},
	}

	err := repo.Set(context.Background(), key, expire, testUser)
	require.Nil(t, err)
	require.NoError(t, err)

	// get cache data
	usr, err := repo.Get(context.Background(), key)
	require.Nil(t, err)
	require.NoError(t, err)
	require.NotNil(t, usr)
	require.Equal(t, testUser, usr)
	//fmt.Printf("%+v", usr)
}

func TestUserRedisRepository_Delete(t *testing.T) {
	repo := SetupRedis()

	key := uuid.New().String()
	err := repo.Delete(context.Background(), key)
	require.Nil(t, err)
	require.NoError(t, err)
}
