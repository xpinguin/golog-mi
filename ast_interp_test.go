package mi

import (
	"fmt"
	"strings"
	"testing"
)

func TestParse2AST2Terms(t *testing.T) {
	trail := ParseToTerms("testdata/squares.go")
	fmt.Printf("%v\n", strings.Join(trail, ""))
}
