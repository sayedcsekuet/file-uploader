package handler

import (
	"file-uploader/src/helpers"
	mocks4 "file-uploader/src/mocks"
	mocks5 "file-uploader/src/mocks/helpers"
	mocks2 "file-uploader/src/mocks/services/scanner"
	"file-uploader/src/services/scanner"
	"github.com/labstack/echo"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"net/http"
	"testing"
)

func TestScanOnError(t *testing.T) {
	handler, si, mfc := createScanHandler()
	c, _, _ := mocks4.CreateTestMultipartFormEchoContext(http.MethodPost, echo.MIMETextHTML, `<virus/>`)
	mfc.On("CollectFromContext", c).Return(nil, errors.New("Files not found!"))
	si.On("ScanFiles", nil).Return(&scanner.ScanResult{
		Success: true,
		Files:   nil,
		Message: "",
	})
	err := handler.ScanFiles(c)
	assert.Error(t, err)
	assert.Equal(t, "code=400, message=Files not found!", err.Error())
}
func TestScan(t *testing.T) {
	handler, si, mfc := createScanHandler()
	c, _, _ := mocks4.CreateTestMultipartFormEchoContext(http.MethodPost, echo.MIMETextHTML, `<virus/>`)
	fc := helpers.NewFileCollector()
	files, err := fc.CollectFromContext(c)
	mfc.On("CollectFromContext", c).Return(files, err)
	si.On("ScanFiles", files).Return(&scanner.ScanResult{
		Success: true,
		Files:   nil,
		Message: "",
	})
	assert.NoError(t, handler.ScanFiles(c))
	assert.Equal(t, http.StatusOK, c.Response().Status)
}

func TestHealth(t *testing.T) {
	handler, si, _ := createScanHandler()
	c, _, _ := mocks4.CreateTestMultipartFormEchoContext(http.MethodPost, echo.MIMETextHTML, `<virus/>`)
	s := &mocks2.Scanner{}
	si.On("Scanner").Return(s)
	s.On("Ping").Return(nil)
	assert.NoError(t, handler.Health(c))
}
func createScanHandler() (ScanHandler, *mocks2.ScanInterceptor, *mocks5.FormCollector) {
	si := &mocks2.ScanInterceptor{}
	fc := &mocks5.FormCollector{}
	handler := NewScanHandler(si, fc)
	return handler, si, fc
}
