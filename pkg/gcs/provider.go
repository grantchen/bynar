/**
    @author: dongjs
    @date: 2023/9/12
    @description:
**/

package gcs

import (
	"io"

	"cloud.google.com/go/storage"
)

type CloudStorageProvider interface {
	// DeleteFiles delete file begin with filePathPrefix in google cloud storage
	DeleteFiles(filePathPrefix string) error
	// UploadFile upload file to google cloud storage
	UploadFile(filePath string, reader io.Reader) (string, error)
	// StorageClient get storageClient
	StorageClient() *storage.Client
	// SignedURL returns a URL for the specified object. Signed URLs allow anyone to access to a restricted resource for a limited time
	SignedURL(filePath string) (string, error)
}
