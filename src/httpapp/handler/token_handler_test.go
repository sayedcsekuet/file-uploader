package handler

import (
	mocks2 "file-uploader/src/mocks"
	mocks3 "file-uploader/src/mocks/services/tokenservice"
	"github.com/labstack/echo"
	"github.com/stretchr/testify/assert"
	"net/http"
	"strings"
	"testing"
)

var tResponse = "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJhdXRob3JpemVkIjp0cnVlLCJleHAiOjE2MjEzNzYxMTksImZpbGVfaWQiOiJiNjk3MjRhMy04NDg4LTQwMmMtYTE2ZS1mYzBkYTBlYTY4MzIifQ.YSjK5MSoaUrtIK7_u0PVMd2UYWcDxSsWDUabOV69WwM"

func TestGenerateTokensIdsRequest(t *testing.T) {
	th, ts := createTokenHandler()
	c, _, _ := mocks2.CreateTestEchoContext("/tokens", http.MethodPost, echo.MIMEApplicationJSON, `{}`)
	ts.On("Generate", "b69724a3-8488-402c-a16e-fc0da0ea6832", "test", int64(1621376119)).Return(tResponse, nil)
	r := th.GenerateTokens(c)
	assert.Error(t, r)
	assert.Equal(t, "code=400, message=Key: 'TokenBody.Ids' Error:Field validation for 'Ids' failed on the 'required' tag", r.Error())
}
func TestGenerateTokensExpIsNotValidRequest(t *testing.T) {
	th, ts := createTokenHandler()
	c, _, _ := mocks2.CreateTestEchoContext("/tokens", http.MethodPost, echo.MIMEApplicationJSON, `{"ids":["b69724a3-8488-402c-a16e-fc0da0ea6832"],"expired_at":"sadfsdf"}`)
	ts.On("Generate", "b69724a3-8488-402c-a16e-fc0da0ea6832", "test", int64(1621376119)).Return(tResponse, nil)
	r := th.GenerateTokens(c)
	assert.Error(t, r)
	assert.Equal(t, "code=400, message=Expired at time is not valid!", r.Error())
}
func TestGenerateTokens(t *testing.T) {
	th, ts := createTokenHandler()
	c, _, res := mocks2.CreateTestEchoContext("/tokens", http.MethodPost, echo.MIMEApplicationJSON, `{"ids":["b69724a3-8488-402c-a16e-fc0da0ea6832"],"expired_at":"2021-05-19T00:15:19.913207+02:00"}`)
	ts.On("Generate", "b69724a3-8488-402c-a16e-fc0da0ea6832", "test", int64(1621376119)).Return(tResponse, nil)
	r := th.GenerateTokens(c)
	assert.NoError(t, r)
	assert.Contains(t, strings.TrimSpace(res.Body.String()), tResponse)
}
func createTokenHandler() (TokenHandler, *mocks3.TokenService) {
	ts := &mocks3.TokenService{}
	th := NewTokenHandler(ts)
	return th, ts
}
