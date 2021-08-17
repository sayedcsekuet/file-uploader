package fileservice

import (
	"file-uploader/src/helpers"
	mocks4 "file-uploader/src/mocks"
	mocks3 "file-uploader/src/mocks/repositories"
	mocks5 "file-uploader/src/mocks/services/filestorage"
	"file-uploader/src/mocks/services/scanner"
	mocks2 "file-uploader/src/mocks/services/tokenservice"
	"file-uploader/src/models"
	"file-uploader/src/repositories"
	"file-uploader/src/services/filestorage"
	"file-uploader/src/services/scanner"
	"github.com/labstack/echo"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"testing"
	"time"
)

var fileModel = models.NewFile("957ad83c-3e7f-494e-89b7-717cad82103d", "test", "test", "/fsd", "s3", models.MetaData{
	MimeType: "",
	Size:     0,
})

func TestList(t *testing.T) {
	testFileService, _, fr, _, _ := createHandler()
	fr.On("GetAllByOwner",
		"test",
		repositories.SearchParam{Name: "", CreatedDate: "", Offset: 0, Limit: 0},
	).Return([]*models.File{}, nil)
	// Assertions
	res, err := testFileService.List("test", repositories.SearchParam{Name: "", CreatedDate: "", Offset: 0, Limit: 0})
	assert.NoError(t, err)
	assert.Equal(t, 0, len(res))
}

func TestGetFileNotFound(t *testing.T) {
	testFileService, _, fr, _, _ := createHandler()
	fr.On("FindByOwnerAndId", "test", "test").Return(nil, errors.New("File not found!"))
	// Assertions
	_, err := testFileService.GetByOwner("test", "test")
	assert.Error(t, err)
	assert.Equal(t, "File not found!", err.Error())
}

func TestGet(t *testing.T) {
	testFileService, _, fr, _, _ := createHandler()
	resMock := &models.File{
		ID:         "",
		Name:       "",
		MetaData:   nil,
		OwnerID:    "",
		BucketPath: "",
		Provider:   "",
		CreatedAt:  time.Time{},
		ExpiredAt:  models.NullTime{},
	}
	fr.On("FindByOwnerAndId", "test", "test").Return(resMock, nil)
	// Assertions
	res, _ := testFileService.GetByOwner("test", "test")
	assert.Equal(t, res, res)
}

func TestDeleteFileNotFound(t *testing.T) {
	testFileService, _, fr, _, _ := createHandler()
	fr.On("FindByOwnerAndId", "test", "test").Return(nil, errors.New("File not found!"))
	// Assertions
	err := testFileService.Delete("test", "test")
	assert.Error(t, err)
	assert.Equal(t, "File not found!", err.Error())
}
func TestDelete(t *testing.T) {
	testFileService, _, fr, fss, _ := createHandler()
	fr.On("FindByOwnerAndId", "test", "test").Return(fileModel, nil)
	fr.On("Delete", "test", "test").Return(nil)
	fss.On("Delete", fileModel).Return(nil)
	// Assertions
	assert.NoError(t, testFileService.Delete("test", "test"))
}

func TestCreateFailOnVirusFile(t *testing.T) {
	testFileService, si, _, _, _ := createHandler()

	req := mocks4.NewTestFileHTTPRequest("POST", echo.MIMEOctetStream, `<virus/>`)
	req.Header["Content-Disposition"] = []string{"attachment;filename=virus.dat"}
	files, _ := helpers.NewFileCollector().Collect(req.MultipartForm)
	f, _ := files[0].Open()
	si.On("ScanForVirus", "virus.dat", f).Return(scanner.FileScanResult{
		Result: &scanner.Result{
			Status:      "",
			Virus:       true,
			Description: "",
		},
		Error:    false,
		FileName: "",
		Message:  "",
	})
	var filesData = map[string]*models.FileData{}
	filesData["1"] = &models.FileData{
		ID:         "",
		FileHeader: files[0],
	}
	result := testFileService.Create("test", &models.CreateFileData{
		Files:      filesData,
		BucketPath: "",
		ExpiredAt:  "",
	})
	// Assertions
	assert.Equal(t, true, result[0].Virus)
}

func TestCreateFile(t *testing.T) {
	testFileService, si, fr, fss, _ := createHandler()

	req := mocks4.NewTestFileHTTPRequest("POST", echo.MIMEOctetStream, `<clean/>`)
	req.Header["Content-Disposition"] = []string{"attachment;filename=virus.dat"}
	files, _ := helpers.NewFileCollector().Collect(req.MultipartForm)
	f, _ := files[0].Open()
	defer f.Close()
	si.On("ScanForVirus", "virus.dat", f).Return(scanner.FileScanResult{
		Result: &scanner.Result{
			Status:      "",
			Virus:       false,
			Description: "",
		},
		Error:    false,
		FileName: "",
		Message:  "",
	})
	var filesData = map[string]*models.FileData{}
	filesData["1"] = &models.FileData{
		ID:         "",
		FileHeader: files[0],
	}

	fss.On("Provider").Return("s3")
	fss.On("Upload", files[0], mock.MatchedBy(func(file *models.File) bool {
		return true
	})).Return(nil)
	fr.On("Create", mock.MatchedBy(func(file *models.File) bool {
		return true
	})).Return(fileModel, nil)
	result := testFileService.Create("test", &models.CreateFileData{
		Files:      filesData,
		BucketPath: "/fsd",
		ExpiredAt:  "",
	})
	// Assertions
	assert.Equal(t, false, result[0].Virus)
	assert.Equal(t, true, result[0].Success)

}
func TestFileReaderWithEmptyToken(t *testing.T) {
	testFileService, _, fr, fss, _ := createHandler()
	fr.On("FindByOwnerAndId", "957ad83c-3e7f-494e-89b7-717cad82103d", "test").Return(fileModel, nil)
	fileReader := &filestorage.FileReader{
		Reader: nil,
		File:   nil,
	}
	fss.On("Read", fileModel).Return(fileReader, nil)
	_, err := testFileService.FileReader("957ad83c-3e7f-494e-89b7-717cad82103d", "test", "")
	assert.NoError(t, err)
}
func TestFileReaderWithToken(t *testing.T) {
	testFileService, _, fr, fss, ts := createHandler()
	fr.On("Get", "957ad83c-3e7f-494e-89b7-717cad82103d").Return(fileModel, nil)
	fileReader := &filestorage.FileReader{
		Reader: nil,
		File:   nil,
	}
	fss.On("Read", fileModel).Return(fileReader, nil)
	ts.On("Verify", "test", "957ad83c-3e7f-494e-89b7-717cad82103d", "test").Return(nil)
	_, err := testFileService.FileReader("957ad83c-3e7f-494e-89b7-717cad82103d", "test", "test")
	assert.NoError(t, err)
}
func TestFileReaderWithTokenInvalid(t *testing.T) {
	testFileService, _, fr, fss, ts := createHandler()
	fr.On("Get", "957ad83c-3e7f-494e-89b7-717cad82103d").Return(fileModel, nil)
	fileReader := &filestorage.FileReader{
		Reader: nil,
		File:   nil,
	}
	fss.On("Read", fileModel).Return(fileReader, nil)
	ts.On("Verify", "test", "957ad83c-3e7f-494e-89b7-717cad82103d", "test").Return(errors.New("Token is not valid!"))
	_, err := testFileService.FileReader("957ad83c-3e7f-494e-89b7-717cad82103d", "test", "test")

	assert.Error(t, err)
	assert.Equal(t, "Token is not valid!", err.Error())
}

func TestDeleteExpiredFilesDbError(t *testing.T) {
	testFileService, _, fr, _, _ := createHandler()
	fr.On("GetExpiredFiles").Return(nil, errors.New("fail!"))
	err := testFileService.DeleteExpiredFiles()
	assert.Error(t, err)
}

func TestFailToDeleteExpiredFilesFromDb(t *testing.T) {
	testFileService, _, fr, ss, _ := createHandler()
	fr.On("GetExpiredFiles").Return([]*models.File{fileModel}, nil)
	fr.On("DeleteAll", []string{"957ad83c-3e7f-494e-89b7-717cad82103d"}).Return(errors.New("fail!"))
	ss.On("Delete", mock.MatchedBy(func(file *models.File) bool {
		return true
	})).Return(nil)
	err := testFileService.DeleteExpiredFiles()
	assert.Error(t, err)
}
func TestDeleteExpiredFilesSuccess(t *testing.T) {
	testFileService, _, fr, ss, _ := createHandler()
	fr.On("GetExpiredFiles").Return([]*models.File{fileModel}, nil)
	fr.On("DeleteAll", []string{"957ad83c-3e7f-494e-89b7-717cad82103d"}).Return(nil)
	ss.On("Delete", mock.MatchedBy(func(file *models.File) bool {
		return true
	})).Return(nil)
	err := testFileService.DeleteExpiredFiles()
	assert.NoError(t, err)
}
func createHandler() (FileService, *mocks.ScanInterceptor, *mocks3.FileRepository, *mocks5.StorageService, *mocks2.TokenService) {
	si := &mocks.ScanInterceptor{}
	fr := &mocks3.FileRepository{}
	fss := &mocks5.StorageService{}
	ts := &mocks2.TokenService{}
	h := NewFileService(si, fr, fss, ts)
	return h, si, fr, fss, ts
}
