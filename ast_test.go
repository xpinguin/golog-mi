package mi

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"reflect"
	"strings"
	"testing"

	"github.com/mndrix/golog/term"

	"github.com/iancoleman/strcase"
)

func TestInspectAST(t *testing.T) {
	fset := token.NewFileSet()
	f, _ := parser.ParseFile(fset, "golog_bridge.go" /*"testdata/squares.go"*/, nil, 0 /*parser.Trace | parser.ParseComments*/)

	/*env := map[ast.Node]*ast.Object{}
	ns0 := map[string]*ast.Object{}
	ns := map[string]ast.Node{}*/

	depth := 0 // increase - [, decrease - ]
	ctr := 0   // preorder counter (nil's are ignored)
	trail := []string{}

	ast.Inspect(f, func(n ast.Node) bool {
		if n == nil {
			depth--
			if strings.TrimSpace(trail[len(trail)-1]) == "," {
				trail = append(trail, trail[len(trail)-1])
				trail[len(trail)-2] = "]"
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
			headTerm = Fntr_("id", x.Name)
			if x.Obj != nil {
				headTerm = Fntr_("ref", x.Obj.Name, headTerm)
			}
		default:
			headTerm = term.NewAtom(astNodeTermName(n))
		}

		depth++
		trail = append(trail, "[", headTerm.String(), ", ")
		ctr++
		return true
	})
	if len(trail) > 0 && strings.TrimSpace(trail[len(trail)-1]) == "," {
		trail[len(trail)-1] = "."
	}

	fmt.Printf("%v\n", strings.Join(trail, ""))
	fmt.Printf("\nD: %d ;; CTR: %d\n\n", depth, ctr)
	//spew.Dump(trace)
}
