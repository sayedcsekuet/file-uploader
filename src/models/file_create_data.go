package models

import "mime/multipart"

type FileData struct {
	ID                    string `form:"id"`
	*multipart.FileHeader `form:"-"`
}
type CreateFileData struct {
	Files      map[string]*FileData `form:"files"`
	BucketPath string               `form:"bucket_path"`
	ExpiredAt  string               `form:"expired_at"`
}
