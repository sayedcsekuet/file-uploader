package errors

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestNewKnownf(t *testing.T) {
	args := []*Argument{
		{
			Key:   "arg_1",
			Value: "str",
		},
		{
			Key:   "arg_2",
			Value: 5,
		},
		{
			Key:   "arg_3",
			Value: []int{1, 2, 3},
		},
	}
	err := NewKnownf(400, "Some error [%s], [%d], %v", args)

	assert.IsType(t, Known{}, err)
	assert.Equal(t, 400, err.Code())
	expectedArgs := map[string]interface{}{
		"arg_1": "str",
		"arg_2": 5,
		"arg_3": []int{1, 2, 3},
	}
	assert.Equal(t, expectedArgs, err.Args())
	assert.EqualError(t, err, "Some error [str], [5], [1 2 3]")
}
