package mocks

import (
	"bytes"
	"io"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"strings"
)

func NewTestFileHTTPRequest(method string, contentType string, content string) *http.Request {
	fileReader := strings.NewReader(content)
	b := bytes.Buffer{} // buffer to write the request payload into
	fw := multipart.NewWriter(&b)
	fFile, _ := fw.CreateFormFile("file", "virus.dat")
	io.Copy(fFile, fileReader)
	fw.Close()
	req := httptest.NewRequest(method, "http://clammit/scan", &b)
	req.Header = map[string][]string{
		"Content-Type":    []string{contentType},
		"X-Forwarded-For": []string{"kermit"},
	}
	req.Header.Set("Content-Type", fw.FormDataContentType())
	req.ParseMultipartForm(32 << 20)
	return req
}
