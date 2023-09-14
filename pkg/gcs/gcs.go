/**
    @author: dongjs
    @date: 2023/9/12
    @description:
**/

package gcs

import (
	"cloud.google.com/go/storage"
	"context"
	"errors"
	"github.com/sirupsen/logrus"
	"google.golang.org/api/option"
	"io"
	"os"
	"time"
)

// gcsClient is the interface for the CloudStorageProvider.
type gcsClient struct {
	client *storage.Client
	bucket string
}

// NewGCSClient creates a new instance of the CloudStorageProvider.
// need to close by user
func NewGCSClient() (CloudStorageProvider, error) {
	opt := option.WithCredentialsJSON([]byte(os.Getenv("GOOGLE_APPLICATION_CREDENTIALS")))
	client, err := storage.NewClient(context.Background(), opt)
	if err != nil {
		return nil, errors.New("NewGCSClient: NewClient error: " + err.Error())
	}
	return &gcsClient{
		client: client,
		bucket: os.Getenv("GOOGLE_CLOUD_STORAGE_BUCKET"),
	}, nil
}

// DeleteFiles delete filePath(fileName) begin with filePathPrefix in google cloud storage
func (g gcsClient) DeleteFiles(filePathPrefix string) error {
	bucketHandle := g.client.Bucket(g.bucket)
	objects := bucketHandle.Objects(context.Background(), &storage.Query{Prefix: filePathPrefix, Delimiter: "/"})
	if objects != nil {
		if next, err := objects.Next(); err == nil && next != nil {
			err = bucketHandle.Object(next.Name).Delete(context.Background())
			if err != nil {
				return err
			}
		}
	}
	return nil
}

// UploadFile upload file to google cloud storage
func (g gcsClient) UploadFile(filePath string, reader io.Reader) (string, error) {
	bucketHandle := g.client.Bucket(g.bucket)
	objectHandle := bucketHandle.Object(filePath)
	wc := objectHandle.NewWriter(context.Background())
	if _, err := io.Copy(wc, reader); err != nil {
		logrus.Errorf("UploadFileToGCS: copy err: %+v", err)
		return "", err
	}
	if err := wc.Close(); err != nil {
		return "", err
	}
	return g.SignedURL(filePath)
}

// StorageClient get storageClient
func (g gcsClient) StorageClient() *storage.Client {
	return g.client
}

// SignedURL Signed URLs allow anyone to access to a restricted resource for a limited time
func (g gcsClient) SignedURL(filePath string) (string, error) {
	opts := &storage.SignedURLOptions{
		Scheme:  storage.SigningSchemeDefault,
		Method:  "GET",
		Expires: time.Now().Add(100 * 12 * 30 * 24 * time.Hour),
	}
	bucketHandle := g.client.Bucket(g.bucket)
	url, err := bucketHandle.SignedURL(filePath, opts)
	if err != nil {
		return "", err
	}
	return url, nil
}
