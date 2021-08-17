package errors

import (
	"errors"
	"github.com/stretchr/testify/suite"
	"testing"
)

type collectorTestSuite struct {
	suite.Suite
}

func TestCollectorTestSuite(t *testing.T) {
	suite.Run(t, new(collectorTestSuite))
}

func (suite *collectorTestSuite) Test_Add_And_Errors() {
	c := NewCollector()
	err1 := errors.New("error 1")
	err2 := errors.New("error 2")
	c.Add(err1)
	c.Add(nil)
	c.Add(err2)

	suite.Equal([]error{err1, nil, err2}, c.Errors())
}

func (suite *collectorTestSuite) Test_HasErrors() {
	suite.Run("True", func() {
		c := NewCollector()
		err1 := errors.New("error 1")
		c.Add(err1)

		suite.True(c.HasErrors())

		c.Add(nil)

		suite.True(c.HasErrors())
	})

	suite.Run("False", func() {
		c := NewCollector()

		suite.False(c.HasErrors())

		c.Add(nil)
		c.Add(nil)

		suite.False(c.HasErrors())
	})
}

func (suite *collectorTestSuite) Test_Error() {
	suite.Run("One error", func() {
		c := NewCollector()
		err1 := errors.New("error 1")
		c.Add(err1)

		suite.EqualError(c, "error 1")
	})

	suite.Run("Multiple errors", func() {
		c := NewCollector()
		c.Add(errors.New("error 1"))
		c.Add(errors.New("error 2"))
		c.Add(errors.New("error 3"))

		suite.EqualError(c, "error 1 | error 2 | error 3")
	})

	suite.Run("Multiple errors with nils", func() {
		c := NewCollector()
		c.Add(errors.New("error 1"))
		c.Add(nil)
		c.Add(errors.New("error 2"))

		suite.EqualError(c, "error 1 | error 2")
	})

	suite.Run("Only nils", func() {
		c := NewCollector()
		c.Add(nil)
		c.Add(nil)

		suite.EqualError(c, "")
	})

	suite.Run("Empty", func() {
		c := NewCollector()

		suite.EqualError(c, "")
	})
}
