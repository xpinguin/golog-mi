package mi

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"strings"
	"testing"

	"github.com/xpinguin/golog-mi/ext/go/ast2"
)

func TestInspectAST(t *testing.T) {
	fset := token.NewFileSet()
	f, _ := parser.ParseFile(fset, "testdata/squares.go", nil, 0 /*parser.Trace | parser.ParseComments*/)

	//inspect := inspector.New([]*ast.File{f})
	ast.Inspect(f, func(n ast.Node) bool {
		if n == nil {
			return false
		}
		switch x := n.(type) {
		case *ast.File:
			return true

		case *ast.TypeSpec:
			//spew.Dump(x)
		////////////////
		case ast.Spec:

		case *ast.FuncDecl:
			if id := x.Name; id != nil {
				if strings.HasPrefix(id.Name, "Square") {
					ast2.Print(fset, n)
					//spew.Dump(x)
				}
			}
		////////////////
		case ast.Decl:
			//spew.Dump(x)

		case ast.Stmt:
		case ast.Expr:
		default:
			fmt.Println("\n==========\n==========")
			//ast2.Print(fset, n)
			fmt.Println("-------------")
			//spew.Dump(x)
			fmt.Print("=====================.\n\n")
		}
		return false // XXX: temporary search-space limiter
	})
}
