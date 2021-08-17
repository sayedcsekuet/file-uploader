package helpers

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

//Validates the URL is OK (fatal error if not) and returns it
func TestEmptyStringIsNotValidUrl(t *testing.T) {
	assert.Equal(t, false, IsValidUrl(""))
}

func TestIsNotValidUrl(t *testing.T) {
	assert.Equal(t, false, IsValidUrl("34unix:://sadhfkjsd//sdaf.sock"))
}
func TestIsValidUrl(t *testing.T) {
	assert.Equal(t, true, IsValidUrl("unix:://sadhfkjsd//sdaf.sock"))
}
