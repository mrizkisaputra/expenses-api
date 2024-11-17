package aws

import (
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"github.com/mrizkisaputra/expenses-api/config"
	"github.com/pkg/errors"
)

func NewAWSClient(cfg *config.Config) (*minio.Client, error) {
	client, err := minio.New(cfg.AWS.Endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(cfg.AWS.MinioAccessKey, cfg.AWS.MinioSecretKey, ""),
		Secure: cfg.AWS.UseSSL,
	})
	if err != nil {
		return nil, errors.Wrap(err, "NewAWSClient.minio.New")
	}

	return client, err
}
