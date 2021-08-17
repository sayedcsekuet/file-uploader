package tokenservice

import (
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestGenerate(t *testing.T) {
	s := NewTokenService()
	_, err := s.Generate("test", "test", time.Now().Add(1).Unix())
	assert.NoError(t, err)
}

func TestTokenIsValid(t *testing.T) {
	s := NewTokenService()
	token, _ := s.Generate("test", "test", time.Now().Add(1).Unix())
	err := s.Verify(token, "test", "test")
	assert.NoError(t, err)
}
func TestTokenSigIsNotValid(t *testing.T) {
	s := NewTokenService()
	token, _ := s.Generate("test", "test", time.Now().Add(1).Unix())
	err := s.Verify(token, "test", "test1")
	assert.Error(t, err)
	assert.Equal(t, "signature is invalid", err.Error())
}
func TestTokenIsExpired(t *testing.T) {
	s := NewTokenService()
	token, _ := s.Generate("test", "test", time.Now().Add(time.Hour*(-1)).Unix())
	err := s.Verify(token, "test", "test")
	assert.Error(t, err)
	assert.Equal(t, "Token is expired", err.Error())
}
