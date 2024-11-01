package repository

import (
	"context"
	"github.com/minio/minio-go/v7"
	"github.com/mrizkisaputra/expenses-api/internal/user"
	"github.com/mrizkisaputra/expenses-api/internal/user/model"
	"github.com/pkg/errors"
	"net/url"
	"time"
)

type awsUserRepository struct {
	s3Client *minio.Client
}

// AWS User Repository constructor
func NewAWSUserRepository(s3Client *minio.Client) user.AWSUserRepository {
	return &awsUserRepository{s3Client: s3Client}
}

// Upload file to MinIO S3
func (aws *awsUserRepository) PutObject(ctx context.Context, input *model.UserUploadInput) (*minio.UploadInfo, error) {
	opts := minio.PutObjectOptions{
		UserMetadata: map[string]string{"x-amz-acl": "public-read"},
		ContentType:  input.ContentType,
	}

	// s3 object operations
	uploadInfo, err := aws.s3Client.PutObject(ctx, input.BucketName, input.ObjectName, input.Object, input.ObjectSize, opts)
	if err != nil {
		return nil, errors.Wrap(err, "AWSUserRepository.PutObject.s3Client.PutObject")
	}
	return &uploadInfo, nil
}

// Download file from MinIO S3
func (aws *awsUserRepository) GetObject(ctx context.Context, bucketName, objectName string) (*minio.Object, error) {
	object, err := aws.s3Client.GetObject(ctx, bucketName, objectName, minio.GetObjectOptions{})
	defer object.Close()
	if err != nil {
		return nil, errors.Wrap(err, "AWSUserRepository.GetObject.s3Client.GetObject")
	}
	return object, nil
}

// Delete file from MinIO S3
func (aws *awsUserRepository) RemoveObject(ctx context.Context, bucketName, objectName string) error {
	opts := minio.RemoveObjectOptions{
		GovernanceBypass: true,
	}
	if err := aws.s3Client.RemoveObject(ctx, bucketName, objectName, opts); err != nil {
		return errors.Wrap(err, "AWSUserRepository.RemoveObject.s3Client.RemoveObject")
	}
	return nil
}

// Get temporary access url
func (aws *awsUserRepository) PresignedGetObject(ctx context.Context, bucketName, objectName string, expiry time.Duration) (*url.URL, error) {
	//Set request parameters for content-disposition.
	var reqParam = make(url.Values)

	presignedUrl, err := aws.s3Client.PresignedGetObject(ctx, bucketName, objectName, expiry, reqParam)
	if err != nil {
		return nil, errors.Wrap(err, "AWSUserRepository.PresignedGetObject.s3Client.PresignedGetObject")
	}

	return presignedUrl, nil
}
