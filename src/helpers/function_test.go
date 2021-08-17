package helpers

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestStringContains(t *testing.T) {
	assert.Equal(t, true, Contains([]string{"a", "b"}, "a"))
}

func TestStringNotContains(t *testing.T) {
	assert.Equal(t, false, Contains([]string{"a", "b"}, "d"))
}
