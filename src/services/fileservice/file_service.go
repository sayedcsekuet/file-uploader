package fileservice

import (
	"database/sql"
	"file-uploader/src/models"
	"file-uploader/src/repositories"
	"file-uploader/src/services/filestorage"
	"file-uploader/src/services/scanner"
	"file-uploader/src/services/tokenservice"
	"fmt"
	"github.com/araddon/dateparse"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"github.com/pkg/errors"
	logger "github.com/sirupsen/logrus"
	"gorm.io/datatypes"
)

type FileService interface {
	List(ownerId string, search repositories.SearchParam) ([]*models.File, error)
	Create(ownerId string, createFileData *models.CreateFileData) []*CreateFileResponse
	ScanAndSave(cFileData *models.FileData, ownerId, bucketPath, expiredAt string) *CreateFileResponse
	GetByOwner(id, ownerId string) (*models.File, error)
	Delete(id, ownerId string) error
	FileReader(id, ownerId, token string) (*filestorage.FileReader, error)
	DeleteExpiredFiles() error
}
type fileService struct {
	scanner.ScanInterceptor
	repositories.FileRepository
	filestorage.StorageService
	tokenservice.TokenService
}

type CreateFileResponse struct {
	Id       string         `json:"id,omitempty"`
	FileName string         `json:"file_name"`
	Success  bool           `json:"success"`
	Message  string         `json:"message"`
	Virus    bool           `json:"virus"`
	MetaData datatypes.JSON `json:"meta_data"`
}

func NewFileService(si scanner.ScanInterceptor, fr repositories.FileRepository, ss filestorage.StorageService, ts tokenservice.TokenService) FileService {
	return &fileService{
		ScanInterceptor: si,
		FileRepository:  fr,
		StorageService:  ss,
		TokenService:    ts,
	}
}
func (fs *fileService) DeleteExpiredFiles() error {
	files, err := fs.FileRepository.GetExpiredFiles()
	if err != nil {
		return err
	}
	var ids []string
	for _, f := range files {
		err = fs.StorageService.Delete(f)
		if err != nil {
			logger.Error(err, nil)
		}
		ids = append(ids, f.ID)
	}
	if len(ids) == 0 {
		return nil
	}
	err = fs.FileRepository.DeleteAll(ids)
	if err != nil {
		logger.Error(err, nil)
		return err
	}
	return nil
}

func (fs *fileService) List(ownerId string, search repositories.SearchParam) ([]*models.File, error) {
	return fs.FileRepository.GetAllByOwner(ownerId, search)
}
func (fs *fileService) ScanAndSave(cFileData *models.FileData, ownerId, bucketPath, expiredAt string) *CreateFileResponse {
	fileName := cFileData.Filename
	openedFile, err := cFileData.FileHeader.Open()
	if err != nil {
		return fs.prepareResponse(cFileData.ID, fileName, err, false)
	}
	var scanResult scanner.FileScanResult
	scanResult = fs.ScanInterceptor.ScanForVirus(fileName, openedFile)
	_ = openedFile.Close()
	if scanResult.Error {
		return fs.prepareResponse(cFileData.ID, fileName, errors.New(scanResult.Message), false)
	}
	if scanResult.Virus {
		return fs.prepareResponse(cFileData.ID, fileName, errors.New(fmt.Sprintf("Virus found with definition %s", scanResult.Description)), scanResult.Virus)
	}
	file, err := fs.prepareFileData(ownerId, bucketPath, expiredAt, cFileData)
	if err != nil {
		return fs.prepareResponse(cFileData.ID, fileName, err, false)
	}
	err = fs.StorageService.Upload(cFileData.FileHeader, file)
	if err != nil {
		return fs.prepareResponse(cFileData.ID, fileName, err, false)
	}
	_, err = fs.FileRepository.Create(file)
	if err != nil {
		e := fs.StorageService.Delete(file)
		logger.Error(e, nil)
		return fs.prepareResponse(cFileData.ID, fileName, err, false)
	}
	fileRes := fs.prepareResponse(file.ID, fileName, err, false)
	fileRes.MetaData = file.MetaData
	return fileRes
}

func (fs *fileService) Create(ownerId string, createFileData *models.CreateFileData) []*CreateFileResponse {
	var responses []*CreateFileResponse
	for i := range createFileData.Files {
		responses = append(responses, fs.ScanAndSave(createFileData.Files[i], ownerId, createFileData.BucketPath, createFileData.ExpiredAt))
	}
	return responses
}

func (fs *fileService) GetByOwner(id, ownerId string) (*models.File, error) {
	return fs.FileRepository.FindByOwnerAndId(id, ownerId)
}

func (fs *fileService) Delete(id, ownerId string) error {
	fileData, err := fs.FileRepository.FindByOwnerAndId(id, ownerId)
	if err != nil {
		return err
	}
	err = fs.FileRepository.Delete(id, ownerId)
	if err != nil {
		return err
	}
	return fs.StorageService.Delete(fileData)
}
func (fs *fileService) FileReader(id, ownerId, token string) (*filestorage.FileReader, error) {
	if token != "" {
		fileData, err := fs.FileRepository.Get(id)
		if err != nil {
			return nil, err
		}
		err = fs.TokenService.Verify(token, fileData.ID, fileData.OwnerID)
		if err != nil {
			return nil, err
		}
		return fs.StorageService.Read(fileData)
	}
	fileData, err := fs.FileRepository.FindByOwnerAndId(id, ownerId)
	if err != nil {
		return nil, err
	}
	return fs.StorageService.Read(fileData)
}
func (fs *fileService) prepareFileData(ownerId, bucketPath, expiredAt string, file *models.FileData) (*models.File, error) {
	id := uuid.NewString()
	if file.ID != "" {
		id = file.ID
	}
	fileData := models.NewFile(
		id,
		file.Filename,
		ownerId,
		bucketPath,
		fs.StorageService.Provider(),
		models.MetaData{MimeType: file.Header.Get("Content-Type"), Size: file.Size},
	)
	if expiredAt != "" {
		t, err := dateparse.ParseAny(expiredAt)
		if err != nil {
			return nil, err
		}
		fileData.ExpiredAt = models.NullTime{NullTime: sql.NullTime{
			Time:  t,
			Valid: true,
		}}
	}

	err := validator.New().Struct(fileData)
	if err != nil {
		return nil, err
	}
	return fileData, nil
}
func (fs *fileService) prepareResponse(id string, fileName string, err error, virus bool) *CreateFileResponse {
	ur := CreateFileResponse{
		Id:       id,
		FileName: fileName,
		Success:  true,
		Message:  "",
		Virus:    virus,
		MetaData: nil,
	}
	if err != nil {
		ur.Success = false
		ur.Message = err.Error()
	}
	return &ur
}
