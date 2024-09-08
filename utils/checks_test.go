package utils

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type BoxTestSuite struct {
	suite.Suite
}

func (suite *BoxTestSuite) TestBoxWidth() {
	assert.Equal(suite.T(), box_width(), 76)
}

func (suite *BoxTestSuite) TestBoxWidthU() {
	assert.Equal(suite.T(), box_widthu(), uint(76))
}

func (suite *BoxTestSuite) TestBoxIndent() {
	assert.Equal(suite.T(), box_indent(), 2)
}

func  (suite *BoxTestSuite) TestInnerBoxWidth() {
	assert.Equal(suite.T(), 72, inner_box_width())
}

func (suite *BoxTestSuite) TestGetInnerString() {
	retv := get_inner_string("abcde")
	expected := "                                 abcde"
	assert.Equal(suite.T(), retv, expected)
}

func (suite *BoxTestSuite) TestBoxHeader() {
	retv := box_header()
	expected := ""
	expected += fmt.Sprintf("  +--------------------------------------------------------------------------+\n")
	expected += fmt.Sprintf("  |                                                                          |")
	assert.Equal(suite.T(), expected, retv)
}

func (suite *BoxTestSuite) TestBoxFooter() {
	retv := box_footer()
	expected := ""
	expected += fmt.Sprintf("  |                                                                          |\n")
	expected += fmt.Sprintf("  +--------------------------------------------------------------------------+")
	assert.Equal(suite.T(), expected, retv)
}

func (suite *BoxTestSuite) TestCenterText() {
	retv := center_text("abcde")
	expected := ""
	expected += fmt.Sprintf("  |                                  abcde                                   |")
	assert.Equal(suite.T(), expected, retv)
}


func TestBoxTestSuite(t *testing.T) {
	suite.Run(t, new(BoxTestSuite))
}


// func ErrorExit(msg interface{}) {
// func Abort(format string, input ...interface{}) {
// func box_footer() string {
// func center_text(instr string) string {
