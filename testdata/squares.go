package main

import "fmt"

/*******************************************

*ast.FuncDecl({
 Doc: *ast.CommentGroup(<nil>),
 Recv: *ast.FieldList(<nil>),
 Name: *ast.Ident(Square2),
 Type: *ast.FuncType({
  Func: (token.Pos) 29,
  Params: *ast.FieldList({
   Opening: (token.Pos) 41,
   List: []*ast.Field (len=1 cap=1) {
    *ast.Field({
     Doc: *ast.CommentGroup(<nil>),
     Names: []*ast.Ident (len=1 cap=1) {
      *ast.Ident(A)
     },
     Type: *ast.Ident(int),
     Tag: *ast.BasicLit(<nil>),
     Comment: *ast.CommentGroup(<nil>)
    })
   },
   Closing: (token.Pos) 47
  }),
  Results: *ast.FieldList({
   Opening: (token.Pos) 49,
   List: []*ast.Field (len=1 cap=1) {
    *ast.Field({
     Doc: *ast.CommentGroup(<nil>),
     Names: []*ast.Ident (len=1 cap=1) {
      *ast.Ident(R)
     },
     Type: *ast.Ident(int),
     Tag: *ast.BasicLit(<nil>),
     Comment: *ast.CommentGroup(<nil>)
    })
   },
   Closing: (token.Pos) 55
  })
 }),
 Body: *ast.BlockStmt({
  Lbrace: (token.Pos) 57,
  List: ([]ast.Stmt) (len=4 cap=4) {
   *ast.AssignStmt({
    Lhs: ([]ast.Expr) (len=1 cap=1) {
     *ast.Ident(factor)
    },
    TokPos: (token.Pos) 67,
    Tok: (token.Token) :=,
    Rhs: ([]ast.Expr) (len=1 cap=1) {
     *ast.BasicLit({
      ValuePos: (token.Pos) 70,
      Kind: (token.Token) INT,
      Value: (string) (len=1) "1"
     })
    }
   }),
   *ast.AssignStmt({
    Lhs: ([]ast.Expr) (len=1 cap=1) {
     *ast.Ident(ind)
    },
    TokPos: (token.Pos) 77,
    Tok: (token.Token) :=,
    Rhs: ([]ast.Expr) (len=1 cap=1) {
     *ast.Ident(A)
    }
   }),
   *ast.ForStmt({
    For: (token.Pos) 83,
    Init: (ast.Stmt) <nil>,
    Cond: *ast.BinaryExpr({
     X: *ast.Ident(ind),
     OpPos: (token.Pos) 91,
     Op: (token.Token) !=,
     Y: *ast.BasicLit({
      ValuePos: (token.Pos) 94,
      Kind: (token.Token) INT,
      Value: (string) (len=1) "0"
     })
    }),
    Post: (ast.Stmt) <nil>,
    Body: *ast.BlockStmt({
     Lbrace: (token.Pos) 96,
     List: ([]ast.Stmt) (len=3 cap=4) {
      *ast.AssignStmt({
       Lhs: ([]ast.Expr) (len=1 cap=1) {
        *ast.Ident(R)
       },
       TokPos: (token.Pos) 102,
       Tok: (token.Token) +=,
       Rhs: ([]ast.Expr) (len=1 cap=1) {
        *ast.Ident(factor)
       }
      }),
      *ast.AssignStmt({
       Lhs: ([]ast.Expr) (len=1 cap=1) {
        *ast.Ident(factor)
       },
       TokPos: (token.Pos) 121,
       Tok: (token.Token) +=,
       Rhs: ([]ast.Expr) (len=1 cap=1) {
        *ast.BasicLit({
         ValuePos: (token.Pos) 124,
         Kind: (token.Token) INT,
         Value: (string) (len=1) "2"
        })
       }
      }),
      *ast.AssignStmt({
       Lhs: ([]ast.Expr) (len=1 cap=1) {
        *ast.Ident(ind)
       },
       TokPos: (token.Pos) 132,
       Tok: (token.Token) -=,
       Rhs: ([]ast.Expr) (len=1 cap=1) {
        *ast.BasicLit({
         ValuePos: (token.Pos) 135,
         Kind: (token.Token) INT,
         Value: (string) (len=1) "1"
        })
       }
      })
     },
     Rbrace: (token.Pos) 138
    })
   }),
   *ast.ReturnStmt({
    Return: (token.Pos) 141,
    Results: ([]ast.Expr) <nil>
   })
  },
  Rbrace: (token.Pos) 148
 })
})

********************************************/

func Square2(A int) (R int) {
	factor := 1
	ind := A
	for ind != 0 {
		R += factor
		factor += 2
		ind -= 1
	}
	return
}

func init() {
	const Accu = 10
	fmt.Println("Square(", Accu, ") =", Square2(Accu))
}

func main() {
}
