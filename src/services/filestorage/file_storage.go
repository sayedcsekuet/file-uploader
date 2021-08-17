package filestorage

import (
	"file-uploader/src/configs"
	"file-uploader/src/models"
	"fmt"
	"github.com/pkg/errors"
	logger "github.com/sirupsen/logrus"
	"gocloud.dev/blob"
	"io"
	"mime/multipart"
	"path/filepath"
	"strings"
)

type FileReader struct {
	*blob.Reader
	*models.File
}
type StorageService interface {
	Upload(file *multipart.FileHeader, data *models.File) error
	Delete(data *models.File) error
	Read(data *models.File) (*FileReader, error)
	Provider() string
}

func NewStorageService(fileStorage *configs.FileStorage) StorageService {
	return &storageService{fileStorage}
}

type storageService struct {
	*configs.FileStorage
}

const ChunkSize = 1024

func (fr *storageService) Provider() string {
	return fr.FileStorageConfig.Provider
}
func (fr *storageService) Upload(file *multipart.FileHeader, data *models.File) error {
	bucket, bucketName, prefix := fr.initBucket(data)
	defer bucket.Close()
	fileName := fmt.Sprintf("%s/%s%s", prefix, data.ID, filepath.Ext(file.Filename))
	data.BucketPath = fmt.Sprintf("%s/%s", bucketName, fileName)
	var writeOption *blob.WriterOptions
	if data.ExpiredAt.Valid {
		dMeta := data.ExpiredAt.Time.String()
		writeOption = &blob.WriterOptions{Metadata: map[string]string{"Expires": dMeta}}

	}
	w, err := bucket.NewWriter(fr.Ctx, fileName, writeOption)
	if err != nil {
		return err
	}
	err = fr.write(file, w)
	if err != nil {
		logger.Error(err, nil)
		return errors.New(fmt.Sprintf("Fail to upload the file to server!"))
	}
	return nil
}

func (fr *storageService) Delete(data *models.File) error {
	bucket, _, prefix := fr.initBucket(data)
	defer bucket.Close()
	err := bucket.Delete(fr.Ctx, prefix)
	if err != nil {
		return err
	}
	return nil
}
func (fr *storageService) Read(data *models.File) (*FileReader, error) {
	bucket, _, prefix := fr.initBucket(data)
	defer bucket.Close()
	r, err := bucket.NewReader(fr.Ctx, prefix, nil)
	if err != nil {
		return nil, err
	}
	return &FileReader{r, data}, nil
}

func (fr *storageService) DeleteFile(data *models.File) error {
	bucket, _, prefix := fr.initBucket(data)
	defer bucket.Close()
	err := bucket.Delete(fr.Ctx, prefix)
	if err != nil {
		return err
	}
	return nil
}

func (fr *storageService) initBucket(data *models.File) (*blob.Bucket, string, string) {
	bucketPaths := strings.Split(strings.TrimLeft(data.BucketPath, "/"), "/")
	bucketName := bucketPaths[0]
	prefix := ""
	if len(bucketPaths) > 1 {
		prefix = strings.Join(bucketPaths[1:], "/")
	}
	bucket := fr.FileStorage.Open(bucketName, data.Provider)
	return bucket, bucketName, prefix
}

func (fr *storageService) write(mFile *multipart.FileHeader, w *blob.Writer) error {
	file, err := mFile.Open()
	if err != nil {
		return err
	}
	defer file.Close()
	for {
		buf := make([]byte, ChunkSize)

		nr, err := file.Read(buf)
		if err == io.EOF {
			closeErr := w.Close()
			if closeErr != nil {
				return closeErr
			}
			return nil
		}
		if err != nil {
			return err
		}
		if nr > 0 {
			_, err = w.Write(buf[0:nr])
			if err != nil {
				return err
			}
		}
	}
}
