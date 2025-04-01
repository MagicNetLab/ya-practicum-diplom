package s3

import (
	"bytes"
	"context"
	"testing"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"

	cm "github.com/MagicNetLab/ya-practicum-diplom/internal/config/mocks"
	"github.com/stretchr/testify/assert"
)

// TestNew проверяет создание нового клиента S3.
func TestNew(t *testing.T) {
	mockCnf := new(cm.AppConfig)
	mockCnf.On("S3Endpoint").Return("localhost:9000")
	mockCnf.On("S3AccessKey").Return("minioadmin")
	mockCnf.On("S3SecretKey").Return("minioadmin")
	mockCnf.On("S3Bucket").Return("testbucket")
	mockCnf.On("IsValid").Return(true)

	s3Client, err := minio.New(mockCnf.S3Endpoint(), &minio.Options{
		Creds:  credentials.NewStaticV4(mockCnf.S3AccessKey(), mockCnf.S3SecretKey(), ""),
		Secure: false,
	})

	assert.NoError(t, err)
	t.Cleanup(func() {
		_ = s3Client.RemoveBucket(context.Background(), "testbucket")
	})

	t.Run("Bucket отсутствует", func(t *testing.T) {
		client, err := New(mockCnf)
		assert.NoError(t, err)
		assert.NotNil(t, client)
	})

	t.Run("Bucket уже есть", func(t *testing.T) {
		client, err := New(mockCnf)
		assert.NoError(t, err)
		assert.NotNil(t, client)
	})
}

// TestClient_GetObject проверяет получение объекта из S3.
func TestClient_GetObject(t *testing.T) {
	testbucket := "testbucket"
	testobjectName := "testobject.txt"

	mockCnf := new(cm.AppConfig)
	mockCnf.On("S3Endpoint").Return("localhost:9000")
	mockCnf.On("S3AccessKey").Return("minioadmin")
	mockCnf.On("S3SecretKey").Return("minioadmin")
	mockCnf.On("S3Bucket").Return(testbucket)
	mockCnf.On("IsValid").Return(true)

	s3Client, err := minio.New(mockCnf.S3Endpoint(), &minio.Options{
		Creds:  credentials.NewStaticV4(mockCnf.S3AccessKey(), mockCnf.S3SecretKey(), ""),
		Secure: false,
	})
	assert.NoError(t, err)

	err = s3Client.MakeBucket(context.Background(), testbucket, minio.MakeBucketOptions{})
	assert.NoError(t, err)
	_, err = s3Client.PutObject(context.Background(), testbucket, testobjectName, bytes.NewReader([]byte("test data")), int64(9), minio.PutObjectOptions{})
	assert.NoError(t, err)

	t.Cleanup(func() {
		_ = s3Client.RemoveObject(context.Background(), testbucket, testobjectName, minio.RemoveObjectOptions{})
		_ = s3Client.RemoveBucket(context.Background(), testbucket)
	})

	service, err := New(mockCnf)
	assert.NoError(t, err)

	t.Run("Успешное получение объекта", func(t *testing.T) {
		obj, err := service.GetObject(context.Background(), testobjectName)
		assert.NoError(t, err)
		assert.NotNil(t, obj)
	})

	t.Run("Попытка получить несуществующий объект", func(t *testing.T) {
		obj, err := service.GetObject(context.Background(), "not-found-object.txt")
		assert.Error(t, err)
		assert.Nil(t, obj)
	})

}

// TestClient_PutObject проверяет сохранение объекта в S3.
func TestClient_PutObject(t *testing.T) {
	testbucket := "testbucket"
	testobjectName := "/files/testobject.txt"

	mockCnf := new(cm.AppConfig)
	mockCnf.On("S3Endpoint").Return("localhost:9000")
	mockCnf.On("S3AccessKey").Return("minioadmin")
	mockCnf.On("S3SecretKey").Return("minioadmin")
	mockCnf.On("S3Bucket").Return(testbucket)
	mockCnf.On("IsValid").Return(true)

	s3Client, err := minio.New(mockCnf.S3Endpoint(), &minio.Options{
		Creds:  credentials.NewStaticV4(mockCnf.S3AccessKey(), mockCnf.S3SecretKey(), ""),
		Secure: false,
	})
	assert.NoError(t, err)

	err = s3Client.MakeBucket(context.Background(), testbucket, minio.MakeBucketOptions{})
	assert.NoError(t, err)

	t.Cleanup(func() {
		_ = s3Client.RemoveObject(context.Background(), testbucket, testobjectName, minio.RemoveObjectOptions{})
		_ = s3Client.RemoveBucket(context.Background(), testbucket)
	})

	service, err := New(mockCnf)
	assert.NoError(t, err)

	t.Run("Успешное сохранение объекта", func(t *testing.T) {
		err := service.PutObject(
			context.Background(),
			testobjectName,
			bytes.NewReader([]byte("test data")),
			int64(9),
		)
		assert.NoError(t, err)

		r, err := s3Client.GetObjectACL(context.Background(), testbucket, testobjectName)
		assert.NoError(t, err)
		assert.NotNil(t, r)
	})
}

// TestClient_RemoveObject проверяет удаление объекта из S3.
func TestClient_RemoveObject(t *testing.T) {
	testbucket := "testbucket"
	testobjectName := "testobject.txt"

	mockCnf := new(cm.AppConfig)
	mockCnf.On("S3Endpoint").Return("localhost:9000")
	mockCnf.On("S3AccessKey").Return("minioadmin")
	mockCnf.On("S3SecretKey").Return("minioadmin")
	mockCnf.On("S3Bucket").Return(testbucket)
	mockCnf.On("IsValid").Return(true)

	s3Client, err := minio.New(mockCnf.S3Endpoint(), &minio.Options{
		Creds:  credentials.NewStaticV4(mockCnf.S3AccessKey(), mockCnf.S3SecretKey(), ""),
		Secure: false,
	})
	assert.NoError(t, err)

	err = s3Client.MakeBucket(context.Background(), testbucket, minio.MakeBucketOptions{})
	assert.NoError(t, err)
	_, err = s3Client.PutObject(context.Background(), testbucket, testobjectName, bytes.NewReader([]byte("test data")), int64(9), minio.PutObjectOptions{})
	assert.NoError(t, err)

	t.Cleanup(func() {
		_ = s3Client.RemoveObject(context.Background(), testbucket, testobjectName, minio.RemoveObjectOptions{})
		_ = s3Client.RemoveBucket(context.Background(), testbucket)
	})

	service, err := New(mockCnf)
	assert.NoError(t, err)

	t.Run("Успешное удаление объекта", func(t *testing.T) {
		err := service.RemoveObject(context.Background(), testobjectName)
		assert.NoError(t, err)
		r, err := s3Client.GetObjectACL(context.Background(), testbucket, testobjectName)
		assert.Error(t, err)
		assert.Nil(t, r)
	})

}
