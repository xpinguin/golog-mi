package mi

import (
	"fmt"
	"os"
	"strings"
	"testing"

	"golang.org/x/tools/go/ssa"

	"github.com/iancoleman/strcase"
	"github.com/mndrix/golog/read"
	"github.com/mndrix/golog/term"
	"github.com/xpinguin/golog-mi/util"

	"github.com/leveltravel/dynpkg/util/meta"
)

const SFF = `
package main

func Square(A int) (R int) {
	factor := 1
	ind := A
	for ind != 0 {
		R = R + factor
		factor = factor + 2
		ind = ind - 1
	}
	return
}`

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

func TestTrivialInterpreter(t *testing.T) {
	ssapkg, pkg := util.BuildPkg("sff", SFF)
	if pkg.Name() != "main" {
		t.Errorf("pkg.Name() = %s, want main", pkg.Name())
	}

	testFuncs(t, ssapkg, "Square")
}
