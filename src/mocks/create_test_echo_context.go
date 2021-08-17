package mocks

import (
	"file-uploader/src/validators"
	"github.com/labstack/echo"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
)

func CreateTestEchoContext(path, method, mimeType, body string) (echo.Context, *http.Request, *httptest.ResponseRecorder) {
	e := echo.New()
	var ioReader io.Reader
	if body != "" {
		ioReader = strings.NewReader(body)
	}
	req := httptest.NewRequest(method, path, ioReader)
	req.Header.Set(echo.HeaderContentType, mimeType)
	req.Header.Set("x-api-key", "test")
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	e.Validator = validators.NewAppValidator()
	return c, req, rec
}

func CreateTestMultipartFormEchoContext(method, mimeType string, body string) (echo.Context, *http.Request, *httptest.ResponseRecorder) {
	e := echo.New()
	req := NewTestFileHTTPRequest(method, mimeType, body)
	req.Header.Set("x-api-key", "test")
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	e.Validator = validators.NewAppValidator()
	return c, req, rec
}
