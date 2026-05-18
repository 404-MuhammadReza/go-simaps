package config

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

func MinioConfig(cfg *Config) *minio.Client {
	client, err := minio.New(cfg.MinioEndpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(cfg.MinioAccessKey, cfg.MinioSecretKey, ""),
		Secure: cfg.MinioUseSSL,
	})
	if err != nil { log.Fatal("Can't Connect to MinIO: ", err) }

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	exist, err := client.BucketExists(ctx, cfg.MinioBucket)
	if err != nil { log.Fatal("Can't checking bucket existence: ", err) }

	if !exist {
		err = client.MakeBucket(ctx, cfg.MinioBucket, minio.MakeBucketOptions{})
		if err != nil { log.Fatal("Can't create bucket: ", err) }
	}

	fmt.Println("Berhasil Connect Minio")
	return client
}