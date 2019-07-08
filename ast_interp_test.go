package mi

import (
	"fmt"
	"strings"
	"testing"
)

func testParse2AST2Terms(t *testing.T) {
	trail := ParseFileToTerms("testdata/squares.go")
	fmt.Printf("%v\n", strings.Join(trail, ""))
}

func TestTermsForProgram(t *testing.T) {
	ProgramToTerms(
		"github.com/xpinguin/golog-mi",
		"testdata/sample/generated")
}
