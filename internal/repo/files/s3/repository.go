package s3

import (
	"github.com/minio/minio-go/v7"
	_ "github.com/minio/minio-go/v7/pkg/credentials"
)

type Repository struct {
	client *minio.Client
	bucket string
}

func (r Repository) Close() error {
	return nil
}
