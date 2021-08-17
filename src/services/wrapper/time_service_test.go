package wrapper

import (
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func Test_TimeService_Now(t *testing.T) {
	s := NewTimeService()
	now := s.Now().Format("Mon Jan 2 15:04 2006")
	expected := time.Now().Format("Mon Jan 2 15:04 2006")
	assert.Equal(t, expected, now)
}

func Test_TimeService_NowNotEqual(t *testing.T) {
	s := NewTimeService()
	now := s.Now().Format("Mon Jan 2 15:04 2006")
	expected := time.Now().Format("Mon Jan 2 15:0 2006")
	assert.NotEqual(t, expected, now)
}
