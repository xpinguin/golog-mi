package mi

import (
	"fmt"
	"io/ioutil"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/xpinguin/golog-mi/util"
)

func testParse2AST2Terms(t *testing.T) {
	trail := ParseToTerms("testdata/squares.go")
	fmt.Printf("%v\n", strings.Join(trail, ""))
}

func TestTermsForProgram(t *testing.T) {
	prg := util.LoadProgram("github.com/leveltravel/storage")
	fset := prg.Fset

	trails := map[string][]string{}
	for _, pkgInfo := range prg.AllPackages {
		for _, f := range pkgInfo.Files {
			var name string
			if f.Name != nil {
				name = f.Name.Name
			} else {
				// failsafe naming :P
				name = fmt.Sprintf("%#v", f.Pos())
			}
			if _, ok := trails[name]; ok {
				name += "_" + time.Now().String()
			}
			trails[name] = FileToTerms(f, fset)
		}
	}

	for fname, trail := range trails {
		ioutil.WriteFile(
			"testdata/sample/generated/"+fname+".pl",
			[]byte(strings.Join(trail, "")),
			os.ModePerm)
	}
}
