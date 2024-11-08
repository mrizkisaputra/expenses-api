package user

import (
	"context"
	"github.com/minio/minio-go/v7"
	"github.com/mrizkisaputra/expenses-api/internal/user/model"
	"net/url"
	"time"
)

// UserPostgresRepository defines methods the services layer expects.
// any repositories it interacts with to implement
type UserPostgresRepository interface {
	Create(ctx context.Context, entity *model.User) (*model.User, error)

	Update(ctx context.Context, entity *model.User) (*model.User, error)

	FindByEmail(ctx context.Context, entity *model.User) (*model.User, error)

	FindById(ctx context.Context, entity *model.User) (*model.User, error)

	FindAlreadyExistByEmail(ctx context.Context, entity *model.User) (int64, error)
}

// UserRedisRepository defines methods the services layer expects.
// any repositories it interacts with to implement.
type UserRedisRepository interface {
	Set(ctx context.Context, key string, expiration time.Duration, value *model.User) error

	Get(ctx context.Context, key string) (*model.User, error)

	Delete(ctx context.Context, key string) error
}

// MinIO AWS S3 (Simple Storage Service) Operations.
// AWSUserRepository defines methods the services layer expects.
// any repositories it interacts with to implement.
type AWSUserRepository interface {
	PutObject(ctx context.Context, input *model.UserUploadInput) (*minio.UploadInfo, error)

	GetObject(ctx context.Context, bucketName, objectName string) (*minio.Object, error)

	RemoveObject(ctx context.Context, bucketName, objectName string) error

	PresignedGetObject(ctx context.Context, bucketName, objectName string, expiry time.Duration) (*url.URL, error)
}
