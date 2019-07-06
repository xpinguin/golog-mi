package main

import "fmt"

/*******************************************

(*ast.FuncDecl)(0xc0000f48a0)({
 Doc: (*ast.CommentGroup)(<nil>),
 Recv: (*ast.FieldList)(<nil>),
 Name: (*ast.Ident)(0xc0000f8360)(Square2),
 Type: (*ast.FuncType)(0xc0000f8600)({
  Func: (token.Pos) 29,
  Params: (*ast.FieldList)(0xc0000f4720)({
   Opening: (token.Pos) 41,
   List: ([]*ast.Field) (len=1 cap=1) {
    (*ast.Field)(0xc0000fc540)({
     Doc: (*ast.CommentGroup)(<nil>),
     Names: ([]*ast.Ident) (len=1 cap=1) {
      (*ast.Ident)(0xc0000f8380)(A)
     },
     Type: (*ast.Ident)(0xc0000f83a0)(int),
     Tag: (*ast.BasicLit)(<nil>),
     Comment: (*ast.CommentGroup)(<nil>)
    })
   },
   Closing: (token.Pos) 47
  }),
  Results: (*ast.FieldList)(0xc0000f4750)({
   Opening: (token.Pos) 49,
   List: ([]*ast.Field) (len=1 cap=1) {
    (*ast.Field)(0xc0000fc5c0)({
     Doc: (*ast.CommentGroup)(<nil>),
     Names: ([]*ast.Ident) (len=1 cap=1) {
      (*ast.Ident)(0xc0000f83c0)(R)
     },
     Type: (*ast.Ident)(0xc0000f83e0)(int),
     Tag: (*ast.BasicLit)(<nil>),
     Comment: (*ast.CommentGroup)(<nil>)
    })
   },
   Closing: (token.Pos) 55
  })
 }),
 Body: (*ast.BlockStmt)(0xc0000f4870)({
  Lbrace: (token.Pos) 57,
  List: ([]ast.Stmt) (len=4 cap=4) {
   (*ast.AssignStmt)(0xc0000fc740)({
    Lhs: ([]ast.Expr) (len=1 cap=1) {
     (*ast.Ident)(0xc0000f8420)(factor)
    },
    TokPos: (token.Pos) 67,
    Tok: (token.Token) :=,
    Rhs: ([]ast.Expr) (len=1 cap=1) {
     (*ast.BasicLit)(0xc0000f8440)({
      ValuePos: (token.Pos) 70,
      Kind: (token.Token) INT,
      Value: (string) (len=1) "1"
     })
    }
   }),
   (*ast.AssignStmt)(0xc0000fc780)({
    Lhs: ([]ast.Expr) (len=1 cap=1) {
     (*ast.Ident)(0xc0000f8460)(ind)
    },
    TokPos: (token.Pos) 77,
    Tok: (token.Token) :=,
    Rhs: ([]ast.Expr) (len=1 cap=1) {
     (*ast.Ident)(0xc0000f8480)(A)
    }
   }),
   (*ast.ForStmt)(0xc0000fca00)({
    For: (token.Pos) 83,
    Init: (ast.Stmt) <nil>,
    Cond: (*ast.BinaryExpr)(0xc0000f47e0)({
     X: (*ast.Ident)(0xc0000f84c0)(ind),
     OpPos: (token.Pos) 91,
     Op: (token.Token) !=,
     Y: (*ast.BasicLit)(0xc0000f84e0)({
      ValuePos: (token.Pos) 94,
      Kind: (token.Token) INT,
      Value: (string) (len=1) "0"
     })
    }),
    Post: (ast.Stmt) <nil>,
    Body: (*ast.BlockStmt)(0xc0000f4840)({
     Lbrace: (token.Pos) 96,
     List: ([]ast.Stmt) (len=3 cap=4) {
      (*ast.AssignStmt)(0xc0000fc900)({
       Lhs: ([]ast.Expr) (len=1 cap=1) {
        (*ast.Ident)(0xc0000f8500)(R)
       },
       TokPos: (token.Pos) 102,
       Tok: (token.Token) +=,
       Rhs: ([]ast.Expr) (len=1 cap=1) {
        (*ast.Ident)(0xc0000f8520)(factor)
       }
      }),
      (*ast.AssignStmt)(0xc0000fc940)({
       Lhs: ([]ast.Expr) (len=1 cap=1) {
        (*ast.Ident)(0xc0000f8540)(factor)
       },
       TokPos: (token.Pos) 121,
       Tok: (token.Token) +=,
       Rhs: ([]ast.Expr) (len=1 cap=1) {
        (*ast.BasicLit)(0xc0000f8560)({
         ValuePos: (token.Pos) 124,
         Kind: (token.Token) INT,
         Value: (string) (len=1) "2"
        })
       }
      }),
      (*ast.AssignStmt)(0xc0000fc980)({
       Lhs: ([]ast.Expr) (len=1 cap=1) {
        (*ast.Ident)(0xc0000f85a0)(ind)
       },
       TokPos: (token.Pos) 132,
       Tok: (token.Token) -=,
       Rhs: ([]ast.Expr) (len=1 cap=1) {
        (*ast.BasicLit)(0xc0000f85c0)({
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
   (*ast.ReturnStmt)(0xc0000f85e0)({
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
