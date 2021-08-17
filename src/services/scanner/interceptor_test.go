package scanner

import (
	"file-uploader/src/helpers"
	mocks2 "file-uploader/src/mocks"
	"github.com/labstack/echo"
	"github.com/stretchr/testify/assert"
	"io"
	"testing"
)

type MockScanner struct {
	Scanner,
	address string
	debug bool
}

func (s *MockScanner) SetAddress(address string) {
	s.address = address
}

func (s *MockScanner) Address() string {
	return s.address
}

func (s *MockScanner) HasVirus(reader io.Reader) (bool, error) {
	panic("implement me")
}

func (s *MockScanner) Ping() error {
	panic("implement me")
}

func (s *MockScanner) Version() (string, error) {
	panic("implement me")
}

var result = new(Result)
var err = new(error)

func (s *MockScanner) Scan(io.Reader) (*Result, error) {
	return result, *err
}

var testScanInterceptor = NewScanInterceptor(new(MockScanner), "/temp/av.sock")

func TestNonMultipartRequest_VirusFound_Without_ContentDisposition(t *testing.T) {
	result.Virus = true
	req := mocks2.NewTestFileHTTPRequest("POST", echo.MIMEMultipartForm, `<virus/>`)
	files, _ := helpers.NewFileCollector().Collect(req.MultipartForm)
	rr := testScanInterceptor.ScanFiles(files)
	assert.Equal(t, true, rr.Files[0].Virus)
}

func TestNonMultipartRequest_VirusFound_With_ContentDisposition(t *testing.T) {
	result.Virus = true
	req := mocks2.NewTestFileHTTPRequest("POST", echo.MIMEOctetStream, `<virus/>`)
	req.Header["Content-Disposition"] = []string{"attachment;filename=virus.dat"}
	files, _ := helpers.NewFileCollector().Collect(req.MultipartForm)
	rr := testScanInterceptor.ScanFiles(files)
	assert.Equal(t, true, rr.Files[0].Virus)
	assert.Equal(t, "virus.dat", rr.Files[0].FileName)
}

func TestNonMultipartRequest_Clean(t *testing.T) {
	result.Virus = false
	req := mocks2.NewTestFileHTTPRequest("POST", echo.MIMEOctetStream, `<clean/>`)
	files, _ := helpers.NewFileCollector().Collect(req.MultipartForm)
	rr := testScanInterceptor.ScanFiles(files)
	assert.Equal(t, true, rr.Success)
}