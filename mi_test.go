package mi

import (
	"bytes"
	"fmt"
	"go/build"
	"go/types"
	"log"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"golang.org/x/tools/go/loader"
	"golang.org/x/tools/go/ssa/interp"
	"golang.org/x/tools/go/ssa/ssautil"

	"golang.org/x/tools/go/ssa"

	"github.com/iancoleman/strcase"
	"github.com/mndrix/golog/read"
	"github.com/mndrix/golog/term"
	"github.com/xpinguin/golog-mi/util"

	"github.com/leveltravel/dynpkg/util/meta"
)

/////////////////////////////
//////////////////////////// (FROM ssa.interp/interp_test.go)
func run(t *testing.T, input string) bool {
	t.Logf("Input: %s\n", input)

	start := time.Now()

	ctx := build.Default    // copy
	ctx.GOROOT = "testdata" // fake goroot
	/*ctx.GOOS = "linux"
	ctx.GOARCH = "amd64"*/

	conf := loader.Config{Build: &ctx}
	if _, err := conf.FromArgs([]string{input}, true); err != nil {
		t.Errorf("FromArgs(%s) failed: %s", input, err)
		return false
	}

	conf.Import("runtime")

	// Print a helpful hint if we don't make it to the end.
	var hint string
	defer func() {
		if hint != "" {
			fmt.Println("FAIL")
			fmt.Println(hint)
		} else {
			fmt.Println("PASS")
		}

		interp.CapturedOutput = nil
	}()

	hint = fmt.Sprintf("To dump SSA representation, run:\n%% go build golang.org/x/tools/cmd/ssadump && ./ssadump -test -build=CFP %s\n", input)

	iprog, err := conf.Load()
	if err != nil {
		t.Errorf("conf.Load(%s) failed: %s", input, err)
		return false
	}

	prog := ssautil.CreateProgram(iprog, ssa.SanityCheckFunctions)
	prog.Build()

	mainPkg := prog.Package(iprog.Created[0].Pkg)
	if mainPkg == nil {
		t.Fatalf("not a main package: %s", input)
	}

	interp.CapturedOutput = new(bytes.Buffer)

	hint = fmt.Sprintf("To trace execution, run:\n%% go build golang.org/x/tools/cmd/ssadump && ./ssadump -build=C -test -run --interp=T %s\n", input)
	exitCode := interp.Interpret(mainPkg, 0, &types.StdSizes{WordSize: 8, MaxAlign: 8}, input, []string{})
	if exitCode != 0 {
		t.Fatalf("interpreting %s: exit code was %d", input, exitCode)
	}
	// $GOROOT/test tests use this convention:
	if strings.Contains(interp.CapturedOutput.String(), "BUG") {
		t.Fatalf("interpreting %s: exited zero but output contained 'BUG'", input)
	}

	hint = "" // call off the hounds

	if false {
		t.Log(input, time.Since(start)) // test profiling
	}

	return true
}

/////////////////////////////
////////////////////////////

const SFF = `
package main

import (
	"runtime" // neccessary for the SSA-interp
	"fmt"
)

func Square(A int) (R int) {
	factor := 1
	ind := A
	for ind != 0 {
		R = R + factor
		factor = factor + 2
		ind = ind - 1
	}
	return
}

func init() {
	fmt.Println("from init()", runtime.GOARCH)
}

func main() {
	fmt.Println("from within")
}
`

//////
func testFuncs(t *testing.T, pkg *ssa.Package, funcRef ...string) {
	for _, ref := range funcRef {
		var f *ssa.Function

		refParts := strings.Split(ref, ".")
		var typeName, funcName string
		if len(refParts) > 1 {
			typeName = refParts[len(refParts)-2]
		}
		funcName = refParts[len(refParts)-1]

		for f = range util.FindMethod(typeName, funcName, pkg) {
			{

				// TODO: funcinfo
				funcInfo := struct {
					Name, Pkg, SrcLoc string
					Unk               map[string]string
				}{}
				funcInfo.Unk = map[string]string{}
				for _, lineFields := range util.FuncStrFields(f) {
					if lineFields[0] != "#" || len(lineFields) < 3 {
						continue
					}

					v := strings.Join(lineFields[2:], ", ")
					switch k := lineFields[1]; k {
					case "Name:":
						funcInfo.Name = v
					case "Package:":
						funcInfo.Pkg = v
					case "Location:":
						funcInfo.SrcLoc = v
					default:
						funcInfo.Unk[k] = v
					}
				}
				////
				unks := []Term{}
				for k, v := range funcInfo.Unk {
					k = strcase.ToSnake(strings.TrimRight(k, ":"))
					pair, err := read.Term(
						term.NewAtom(k).String() +
							"=" +
							term.NewAtom(v).String() +
							".")
					if err != nil {
						fmt.Println("ERR: failed to construct pair: ", err)
						continue
					}
					unks = append(unks, pair)
				}
				fmt.Println(
					meta.Fntr_("func_info",
						util.RelMethodName(f, f.Pkg),
						term.SliceToList([]Term{
							meta.Fntr_("blocks", len(f.Blocks)),
							meta.Fntr_("name", funcInfo.Name),
							meta.Fntr_("pkg", funcInfo.Pkg),
							meta.Fntr_("src", funcInfo.SrcLoc),
						}),
						term.SliceToList(unks),
					).String() + ".")
			}
			f.WriteTo(os.Stdout)
			fmt.Print("%% ==================================\n")
			////
			_ = meta.FunctionTerm(f, nil)
			//fmt.Print("==================================\n")
			fmt.Print("%% ==================================\n\n")
		}

		if f == nil {
			t.Fatal("function not found")
		}
	}
}

func testTrivialInterpreter(t *testing.T) {
	ssapkg, pkg := util.BuildPkg("sff", SFF)
	if pkg.Name() != "main" {
		t.Errorf("pkg.Name() = %s, want main", pkg.Name())
	}

	testFuncs(t, ssapkg, "Square")
}

func TestSSAInterp(t *testing.T) {
	cwd, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}

	if !run(t, filepath.Join(cwd, "testdata", "sff.go")) {
		t.Fail()
	}
}
