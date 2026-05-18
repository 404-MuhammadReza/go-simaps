package repository

import (
	"io"
	"time"
	"context"

	"github.com/minio/minio-go/v7"
)

type StorageRepository interface {
	EnsureBucket(ctx context.Context, bucketName string) error
	GetObject(ctx context.Context, objectName string) (*minio.Object, error)
	GetPresignedURL(ctx context.Context, objectName string, expiry time.Duration) (string, error)
	Upload(ctx context.Context, objectName string, reader io.Reader, objectSize int64, contentType string) (*minio.UploadInfo, error)
	Delete(ctx context.Context, objectName string) error
}

type storageRepository struct {
	client *minio.Client
	bucket string
}

func NewStorageRepository(client *minio.Client, bucket string) StorageRepository {
	return &storageRepository{ client, bucket }
}

func (r *storageRepository) EnsureBucket(ctx context.Context, bucketName string) error {
	exist, err := r.client.BucketExists(ctx, bucketName)
	if err != nil { return err }

	if !exist {
		err = r.client.MakeBucket(ctx, bucketName, minio.MakeBucketOptions{})
		if err != nil { return err }
	}

	return nil
}

func (r *storageRepository) GetObject(ctx context.Context, objectName string) (*minio.Object, error) {
	object, err := r.client.GetObject(ctx, r.bucket, objectName, minio.GetObjectOptions{})
	if err != nil { return nil, err }

	return object, nil
}

func (r *storageRepository) GetPresignedURL(ctx context.Context, objectName string, expiry time.Duration) (string, error) {
	url, err := r.client.PresignedGetObject(ctx, r.bucket, objectName, expiry, nil)
	if err != nil { return "", err }

	return url.String(), nil
}

func (r *storageRepository) Upload(ctx context.Context, objectName string, reader io.Reader, objectSize int64, contentType string) (*minio.UploadInfo, error) {
	option := minio.PutObjectOptions{ ContentType: contentType }
	info, err := r.client.PutObject(ctx, r.bucket, objectName, reader, objectSize, option)
	if  err != nil { return nil, err }

	return &info, err
}

func (r *storageRepository) Delete(ctx context.Context, objectName string) error {
	err := r.client.RemoveObject(ctx, r.bucket, objectName, minio.RemoveObjectOptions{})
	if err != nil { return err }

	return nil
}