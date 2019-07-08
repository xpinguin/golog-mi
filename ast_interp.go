package mi

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"io/ioutil"
	"os"
	"reflect"
	"strings"
	"time"

	"github.com/iancoleman/strcase"
	"github.com/mndrix/golog/term"
	"github.com/xpinguin/golog-mi/util"
)

func FileToTerms(f *ast.File, fset *token.FileSet) (termsTrail []string) {
	if f == nil {
		panic("nil AST File to traverse")
	}
	if fset == nil {
		fset = token.NewFileSet()
	}

	depth := 0 // increase - [, decrease - ]
	ctr := 0   // preorder counter (nil's are ignored)
	trail := []string{}

	ast.Inspect(f, func(n ast.Node) bool {
		if n == nil {
			depth--
			if strings.TrimSpace(trail[len(trail)-1]) == "," {
				trail = append(trail, trail[len(trail)-1])
				trail[len(trail)-2] = "]"

				if depth >= 1 && depth < 5 {
					eol := strings.TrimSpace(trail[len(trail)-1])
					eol += "\n" + strings.Repeat("  ", depth)
					trail[len(trail)-1] = eol
				}
			}
			return true
		}

		var headTerm Term

		astNodeTermName := func(n ast.Node) string {
			nData := reflect.Indirect(reflect.ValueOf(n))
			nName := nData.Type().Name()
			return strcase.ToSnake(nName)
		}

		switch x := n.(type) {
		case *ast.BinaryExpr:
			headTerm = term.NewAtom(x.Op.String())
		case *ast.UnaryExpr:
			headTerm = term.NewAtom(x.Op.String())
		case *ast.StarExpr:
			headTerm = term.NewAtom("*")
		case *ast.AssignStmt:
			headTerm = Fntr_(astNodeTermName(x), x.Tok.String())
		case *ast.BasicLit:
			headTerm = Fntr_(
				strcase.ToSnake(x.Kind.String()),
				strings.Trim(x.Value, `"' `))
		case *ast.Ident:
			if x.Obj == nil {
				//headTerm = Fntr_("id0", x.Name)
				headTerm = term.NewAtom(x.Name)
			} else if x.Obj.Name != x.Name {
				headTerm = Fntr_("idobj", x.Name, x.Obj.Name)
			} else {
				headTerm = Fntr_("id", x.Obj.Name)
			}
		default:
			headTerm = term.NewAtom(astNodeTermName(n))
		}

		depth++
		trail = append(trail, "[", headTerm.String(), ", ")
		ctr++
		return true
	})
	/*if len(trail) > 0 && strings.TrimSpace(trail[len(trail)-1]) == "," {
		trail[len(trail)-1] = "."
	}*/

	//fmt.Printf("\nD: %d ;; CTR: %d\n\n", depth, ctr)
	return trail
}

func ParseFileToTerms(file string) (termsTrail []string) {
	fset := token.NewFileSet()
	f, _ := parser.ParseFile(fset, file, nil, 0 /*parser.Trace | parser.ParseComments*/)

	return FileToTerms(f, fset)
}

func ProgramToTerms(mainPkg string, outDir string) (fileTermsTrails map[string][]string) {
	fileTermsTrails = map[string][]string{}

	prg := util.LoadProgram(mainPkg)
	fset := prg.Fset

	for _, pkgInfo := range prg.AllPackages {
		////////////////
		if pkgPath := pkgInfo.Pkg.Path(); strings.Contains(pkgPath, "vendor") {
			continue
		}
		////////////////
		for _, f := range pkgInfo.Files {
			var name string
			if f.Name != nil {
				name = f.Name.Name
			} else {
				// failsafe naming :P
				name = fmt.Sprintf("%#v", f.Pos())
			}
			if _, ok := fileTermsTrails[name]; ok {
				name += "_" + time.Now().String()
			}
			fileTermsTrails[name] = FileToTerms(f, fset)
		}
	}

	for fname, trail := range fileTermsTrails {
		mname := term.NewAtom(fname).String()
		preamble := []string{
			`:- encoding(utf8).`,
			`:- module(` + mname + `, [program/2]).`,
		}

		ioutil.WriteFile(
			strings.TrimRight(outDir, "/")+"/"+fname+".pl",
			[]byte(
				strings.Join(preamble, "\n")+"\n"+
					"program("+mname+",\n"+
					strings.Join(trail, "")+
					"\n)."),
			os.ModePerm)
	}
	return fileTermsTrails
}
