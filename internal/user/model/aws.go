package model

import "io"

// MinIO AWS User Upload Input
type UserUploadInput struct {
	Object      io.Reader
	ObjectName  string
	ObjectSize  int64
	BucketName  string
	ContentType string
}
