package s3

import (
	"context"
	"time"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"

	"github.com/MagicNetLab/ya-practicum-diplom/internal/conf"
	"github.com/MagicNetLab/ya-practicum-diplom/internal/logger"
)

// New возвращает репозиторий с подключенным s3 хранилищем
func New(cnf conf.Configurator) (Repository, error) {
	r := Repository{}

	// подключение к s3 хранилищу
	client, err := minio.New(cnf.S3Endpoint(), &minio.Options{
		Creds:  credentials.NewStaticV4(cnf.S3SecretKey(), cnf.S3AccessKey(), ""),
		Secure: false,
	})
	if err != nil {
		logger.Error("failed create minio client", logger.StrArg("error", err.Error()))
		return r, err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	bucketExists, err := client.BucketExists(ctx, cnf.S3Bucket())
	if err != nil {
		logger.Error("failed check bucket exists", logger.StrArg("error", err.Error()))
		return r, err
	}
	if !bucketExists {
		err = client.MakeBucket(ctx, cnf.S3Bucket(), minio.MakeBucketOptions{})
		if err != nil {
			logger.Error("failed create s3 bucket", logger.StrArg("error", err.Error()))
			return r, err
		}
		r.bucket = cnf.S3Bucket()
	}

	r.client = client

	return r, nil
}
