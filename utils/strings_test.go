package utils

import (
	"testing"

	"github.com/stretchr/testify/suite"
)

type StringsTestSuite struct {
	suite.Suite
}

func (suite *StringsTestSuite) TestStringInSlice() {
	if !StringInSlice("a", []string{"a", "b", "c"}) {
		suite.T().Errorf("Cannot find string in slice")
	}
}

func TestStringsTestSuite(t *testing.T) {
	suite.Run(t, new(StringsTestSuite))
}
