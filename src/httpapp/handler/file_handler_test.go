package handler

import (
	"encoding/json"
	"file-uploader/src/helpers"
	mocks2 "file-uploader/src/mocks"
	mocks3 "file-uploader/src/mocks/services/fileservice"
	"file-uploader/src/models"
	"file-uploader/src/repositories"
	"file-uploader/src/services/fileservice"
	"github.com/labstack/echo"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"net/http"
	"strings"
	"testing"
)

const body = `{"id":"991de811-1d97-40cc-9cc3-e7a38fc4913c","name":"test.vcf","meta_data":{"size":25919,"mime_type":"text/x-vcard"},"owner_id":"asdfjkasdf","bucket_path":"991de811-1d97-40cc-9cc3-e7a38fc4913c.vcf","provider":"s3","created_at":"2021-05-10T06:38:12Z","expired_at":{"Time":"2021-05-12T19:21:04Z","Valid":true}}`
const expected = `{"id":"991de811-1d97-40cc-9cc3-e7a38fc4913c","name":"test.vcf","meta_data":{"size":25919,"mime_type":"text/x-vcard"},"owner_id":"asdfjkasdf","bucket_path":"991de811-1d97-40cc-9cc3-e7a38fc4913c.vcf","provider":"s3","created_at":"2021-05-10T06:38:12Z","expired_at":"2021-05-12T19:21:04Z"}`

func TestList(t *testing.T) {
	testFileHandler, fs, _ := createHandler()
	fs.On("List",
		"test",
		repositories.SearchParam{Name: "", CreatedDate: "", Offset: 0, Limit: 0},
	).Return([]*models.File{}, nil)
	c, _, res := mocks2.CreateTestEchoContext("/files", http.MethodGet, echo.MIMETextHTML, "")
	// Assertions
	if assert.NoError(t, testFileHandler.List(c)) {
		assert.Equal(t, http.StatusOK, res.Code)
		assert.Equal(t, `[]`, strings.TrimSpace(res.Body.String()))
	}
}

func TestListError(t *testing.T) {
	testFileHandler, fs, _ := createHandler()
	fs.On("List",
		"test",
		repositories.SearchParam{Name: "", CreatedDate: "", Offset: 0, Limit: 0},
	).Return(nil, errors.New("Test error"))
	c, _, _ := mocks2.CreateTestEchoContext("/files", http.MethodGet, echo.MIMETextHTML, "")
	err := testFileHandler.List(c)
	// Assertions
	assert.Error(t, err)
	assert.Equal(t, "code=500, message=Internal server error", err.Error())
}

func TestGetFileNotFound(t *testing.T) {
	testFileHandler, fs, _ := createHandler()
	fs.On("GetByOwner",
		"991de811-1d97-40cc-9cc3-e7a38fc4913c",
		"test",
	).Return(nil, errors.New("File not found!"))
	c, _, _ := mocks2.CreateTestEchoContext("/files/991de811-1d97-40cc-9cc3-e7a38fc4913c", http.MethodGet, echo.MIMETextHTML, "")
	c.SetParamNames("id")
	c.SetParamValues("991de811-1d97-40cc-9cc3-e7a38fc4913c")
	response := testFileHandler.Get(c)
	// Assertions
	assert.Error(t, response)
}

func TestGet(t *testing.T) {
	testFileHandler, fs, _ := createHandler()
	file := &models.File{}
	json.Unmarshal([]byte(body), file)
	fs.On("GetByOwner",
		"991de811-1d97-40cc-9cc3-e7a38fc4913c",
		"test",
	).Return(file, nil)
	c, _, res := mocks2.CreateTestEchoContext("/files/991de811-1d97-40cc-9cc3-e7a38fc4913c", http.MethodGet, echo.MIMETextHTML, "")
	c.SetParamNames("id")
	c.SetParamValues("991de811-1d97-40cc-9cc3-e7a38fc4913c")
	// Assertions
	assert.NoError(t, testFileHandler.Get(c))
	assert.Equal(t, http.StatusOK, res.Code)
	assert.Equal(t, expected, strings.TrimSpace(res.Body.String()))
}

func TestDeleteFileNotFound(t *testing.T) {
	testFileHandler, fs, _ := createHandler()
	fs.On("Delete",
		"991de811-1d97-40cc-9cc3-e7a38fc4913c",
		"test",
	).Return(errors.New("File not found!"))
	c, _, _ := mocks2.CreateTestEchoContext("/files/991de811-1d97-40cc-9cc3-e7a38fc4913c", http.MethodDelete, echo.MIMETextHTML, "")
	c.SetParamNames("id")
	c.SetParamValues("991de811-1d97-40cc-9cc3-e7a38fc4913c")
	response := testFileHandler.Delete(c)
	// Assertions
	assert.Error(t, response)
}

func TestDelete(t *testing.T) {
	testFileHandler, fs, _ := createHandler()
	fs.On("Delete",
		"991de811-1d97-40cc-9cc3-e7a38fc4913c",
		"test",
	).Return(nil)
	c, _, res := mocks2.CreateTestEchoContext("/files/991de811-1d97-40cc-9cc3-e7a38fc4913c", http.MethodDelete, echo.MIMETextHTML, "")
	c.SetParamNames("id")
	c.SetParamValues("991de811-1d97-40cc-9cc3-e7a38fc4913c")
	// Assertions
	if assert.NoError(t, testFileHandler.Delete(c)) {
		assert.Equal(t, http.StatusOK, res.Code)
	}
}

func TestCreateFailOnVirusFile(t *testing.T) {
	testFileHandler, fs, fc := createHandler()
	c, _, res := mocks2.CreateTestMultipartFormEchoContext(http.MethodPost, echo.MIMETextHTML, `<virus/>`)
	fileCreateData, _ := fc.CollectFileCreateDataFromContext(c)
	response := []*fileservice.CreateFileResponse{{
		Id:       "",
		FileName: "virus.dat",
		Success:  false,
		Message:  "virus  found!",
		Virus:    true,
	}}
	fs.On("Create", "test", fileCreateData).Return(response)
	testFileHandler.Create(c)
	// Assertions
	assert.Equal(t, http.StatusMultiStatus, res.Code)
	assert.Equal(t, `[{"file_name":"virus.dat","success":false,"message":"virus  found!","virus":true,"meta_data":null}]`, strings.TrimSpace(res.Body.String()))

}

func TestCreateFile(t *testing.T) {
	testFileHandler, fs, fc := createHandler()
	c, _, res := mocks2.CreateTestMultipartFormEchoContext(http.MethodPost, echo.MIMETextHTML, `<virus/>`)
	fileCreateData, _ := fc.CollectFileCreateDataFromContext(c)
	response := []*fileservice.CreateFileResponse{{
		Id:       "",
		FileName: "virus.dat",
		Success:  true,
		Message:  "",
		Virus:    false,
	}}
	fs.On("Create", "test", fileCreateData).Return(response)
	testFileHandler.Create(c)
	// Assertions
	assert.Equal(t, http.StatusOK, res.Code)
	assert.Equal(t, `[{"file_name":"virus.dat","success":true,"message":"","virus":false,"meta_data":null}]`, strings.TrimSpace(res.Body.String()))
}

func createHandler() (FileHandler, *mocks3.FileService, helpers.FormCollector) {
	fs := &mocks3.FileService{}
	fc := helpers.NewFileCollector()
	h := NewFileHandler(fc, fs)
	return h, fs, fc
}
