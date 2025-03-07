package s3

import (
	"bytes"
	"context"
	"errors"
	"io"
	"time"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"

	"github.com/MagicNetLab/ya-practicum-diplom/internal/config"
	"github.com/MagicNetLab/ya-practicum-diplom/internal/logger"
)

// Client интерфейс клиента S3
type Client interface {
	GetObject(ctx context.Context, name string) ([]byte, error)
	PutObject(ctx context.Context, name string, obj io.Reader, size int64) error
	RemoveObject(ctx context.Context, name string) error
}

// New конструктор клиента S3
func New(cnf config.S3Configurator) (Client, error) {
	client, err := minio.New(cnf.GetEndpoint(), &minio.Options{
		Creds:  credentials.NewStaticV4(cnf.GetAccessKey(), cnf.GetSecretKey(), ""),
		Secure: false,
	})
	if err != nil {
		logger.Error("Failed to create minio client", logger.StrArg("err", err.Error()))
		return nil, errors.New("failed to create minio client")
	}
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	exists, err := client.BucketExists(ctx, cnf.GetBucket())
	if err != nil {
		logger.Error("Failed to check bucket existence", logger.StrArg("err", err.Error()))
		return nil, errors.New("failed to check bucket existence")
	}
	if !exists {
		logger.Debug("bucket hasn't been found")
		err = client.MakeBucket(ctx, cnf.GetBucket(), minio.MakeBucketOptions{})
		if err != nil {
			logger.Error("Failed to create bucket", logger.StrArg("err", err.Error()))
			return nil, errors.New("failed to create bucket")
		}
	}

	return &S3Client{client: client, bucket: cnf.GetBucket()}, nil
}

// S3Client клиент хранилища S3
type S3Client struct {
	client *minio.Client
	bucket string
}

// GetObject получение объекта из хранилища
func (c *S3Client) GetObject(ctx context.Context, name string) ([]byte, error) {
	obj, err := c.client.GetObject(ctx, c.bucket, name, minio.GetObjectOptions{})
	if err != nil {
		logger.Error("Failed to get object", logger.StrArg("err", err.Error()))
		return nil, errors.New("failed to get object")
	}

	defer func() {
		if err = obj.Close(); err != nil {
			logger.Error("fialed to close object", logger.StrArg("err", err.Error()))
		}
	}()

	buf := new(bytes.Buffer)
	if _, err = buf.ReadFrom(obj); err != nil {
		return nil, errors.New("failed to read object")
	}
	return buf.Bytes(), nil
}

// PutObject сохранение объекта в хранилище
func (c *S3Client) PutObject(ctx context.Context, name string, obj io.Reader, size int64) error {
	if _, err := c.client.PutObject(
		ctx,
		c.bucket,
		name,
		obj,
		size,
		minio.PutObjectOptions{},
	); err != nil {
		logger.Error("Failed to put object", logger.StrArg("err", err.Error()))
		return errors.New("failed to put object")
	}

	return nil
}

// RemoveObject удаление объекта из хранилища
func (c *S3Client) RemoveObject(ctx context.Context, name string) error {
	err := c.client.RemoveObject(ctx, c.bucket, name, minio.RemoveObjectOptions{})
	if err != nil {
		logger.Error("Failed to remove object", logger.StrArg("err", err.Error()))
		return errors.New("failed to remove object")
	}

	return nil
}
