package user

import (
	"context"
	"github.com/minio/minio-go/v7"
	"github.com/mrizkisaputra/expenses-api/internal/user/model"
	"net/url"
	"time"
)

// MinIO AWS S3 (Simple Storage Service) Operations
type AWSUserRepository interface {
	PutObject(ctx context.Context, input *model.UserUploadInput) (*minio.UploadInfo, error)

	GetObject(ctx context.Context, bucketName, objectName string) (*minio.Object, error)

	RemoveObject(ctx context.Context, bucketName, objectName string) error

	PresignedGetObject(ctx context.Context, bucketName, objectName string, expiry time.Duration) (*url.URL, error)
}
