package storage

import (
	"bytes"
	"context"
	"github.com/google/uuid"
	"path/filepath"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

type MinIOClient struct {
	client *minio.Client
}

// NewMinIOClient Initialize new MinIO client
func NewMinIOClient(endpoint, accessKey, secretKey string, useSSL bool) (*MinIOClient, error) {
	client, err := minio.New(endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(accessKey, secretKey, ""),
		Secure: useSSL,
	})
	if err != nil {
		return nil, err
	}

	return &MinIOClient{client: client}, nil
}

// SaveFile saves a file to a MinIO storage and returns an ID
func (minIoClient *MinIOClient) SaveFile(ctx context.Context, bucketName, fileName string, fileData []byte) (string, error) {
	// Checks if bucket exists. If it doesn't - creates it
	exists, err := minIoClient.client.BucketExists(ctx, bucketName)
	if err != nil {
		return "", err
	}
	if !exists {
		err = minIoClient.client.MakeBucket(ctx, bucketName, minio.MakeBucketOptions{})
		if err != nil {
			return "", err
		}
	}

	fileID := generateFileID(fileName)

	_, err = minIoClient.client.PutObject(
		ctx,
		bucketName,
		fileID,
		bytes.NewReader(fileData),
		int64(len(fileData)),
		minio.PutObjectOptions{},
	)
	return fileID, err
}

// Generates file id
func generateFileID(fileName string) string {
	return uuid.New().String() + filepath.Ext(fileName)
}
