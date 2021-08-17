package helpers

import (
	mocks2 "file-uploader/src/mocks"
	"github.com/labstack/echo"
	"github.com/stretchr/testify/assert"
	"net/http"
	"testing"
)

func TestCollectFromContext(t *testing.T) {
	c, _, _ := mocks2.CreateTestMultipartFormEchoContext(http.MethodPost, echo.MIMETextHTML, `<virus/>`)
	collector := NewFileCollector()
	r, _ := collector.CollectFromContext(c)
	assert.Equal(t, "virus.dat", r[0].Filename)
	assert.Greater(t, r[0].Size, int64(0))
}

func TestEmptyCollectFromContext(t *testing.T) {
	c, _, _ := mocks2.CreateTestMultipartFormEchoContext(http.MethodPost, echo.MIMETextHTML, ``)
	collector := NewFileCollector()
	r, _ := collector.CollectFromContext(c)
	assert.Equal(t, int64(0), r[0].Size)
}
