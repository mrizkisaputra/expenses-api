package mock

import (
	"context"
	"github.com/minio/minio-go/v7"
	"github.com/mrizkisaputra/expenses-api/internal/user/model"
	"github.com/stretchr/testify/mock"
	"net/url"
	"time"
)

type AwsRepositoryMock struct {
	mock.Mock
}

func (a *AwsRepositoryMock) PutObject(ctx context.Context, input *model.UserUploadInput) (*minio.UploadInfo, error) {
	args := a.Called(ctx, input)
	if args.Get(0) != nil {
		return args.Get(0).(*minio.UploadInfo), args.Error(1)
	}
	return nil, args.Error(1)
}

func (a *AwsRepositoryMock) GetObject(ctx context.Context, bucketName, objectName string) (*minio.Object, error) {
	//TODO implement me
	panic("implement me")
}

func (a *AwsRepositoryMock) RemoveObject(ctx context.Context, bucketName, objectName string) error {
	//TODO implement me
	panic("implement me")
}

func (a *AwsRepositoryMock) PresignedGetObject(ctx context.Context, bucketName, objectName string, expiry time.Duration) (*url.URL, error) {
	//TODO implement me
	panic("implement me")
}
