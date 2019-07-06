package mi

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"strings"
	"testing"

	"github.com/davecgh/go-spew/spew"
)

func TestInspectAST(t *testing.T) {
	fset := token.NewFileSet()
	f, _ := parser.ParseFile(fset, "testdata/squares.go", nil, parser.Trace|parser.ParseComments)

	//inspect := inspector.New([]*ast.File{f})
	ast.Inspect(f, func(n ast.Node) bool {
		if n == nil {
			return false
		}
		switch x := n.(type) {
		case *ast.File:
			return true

		////////////////
		case *ast.TypeSpec:
			//spew.Dump(x)
			return false
		case ast.Spec:

		//////////////
		case *ast.FuncDecl:
			if id := x.Name; id != nil && strings.HasPrefix(id.Name, "Square") {
				spew.Dump(x)
			}
			return false
		case ast.Decl:
			//spew.Dump(x)
			return false

		////////////////
		case ast.Stmt:

		////////////////
		case ast.Expr:

		////////////////
		default:
			fmt.Println("\n==========\n==========")
			//ast.Print(fset, n)
			fmt.Println("-------------")
			//spew.Dump(x)
			fmt.Print("=====================.\n\n")
			return false
		}
		return false // XXX: temporary search-space limiter, to be removed
	})
}
