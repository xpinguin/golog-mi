// -*-syntax: go -*-
//(just for completeness)
package main

import (
	"fmt"
	_ "runtime" //XXX: neccessary for the ssa/interp.Interpret
)

/**************************************

(*ast.FuncDecl)(0xc0000eea50)({
 Doc: (*ast.CommentGroup)(<nil>),
 Recv: (*ast.FieldList)(<nil>),
 Name: (*ast.Ident)(0xc000004540)(Square),
 Type: (*ast.FuncType)(0xc0000048e0)({
  Func: (token.Pos) 150,
  Params: (*ast.FieldList)(0xc0000ee840)({
   Opening: (token.Pos) 161,
   List: ([]*ast.Field) (len=1 cap=1) {
    (*ast.Field)(0xc000038540)({
     Doc: (*ast.CommentGroup)(<nil>),
     Names: ([]*ast.Ident) (len=1 cap=1) {
      (*ast.Ident)(0xc000004560)(A)
     },
     Type: (*ast.Ident)(0xc000004580)(int),
     Tag: (*ast.BasicLit)(<nil>),
     Comment: (*ast.CommentGroup)(<nil>)
    })
   },
   Closing: (token.Pos) 167
  }),
  Results: (*ast.FieldList)(0xc0000ee870)({
   Opening: (token.Pos) 169,
   List: ([]*ast.Field) (len=1 cap=1) {
    (*ast.Field)(0xc0000385c0)({
     Doc: (*ast.CommentGroup)(<nil>),
     Names: ([]*ast.Ident) (len=1 cap=1) {
      (*ast.Ident)(0xc0000045a0)(R)
     },
     Type: (*ast.Ident)(0xc0000045c0)(int),
     Tag: (*ast.BasicLit)(<nil>),
     Comment: (*ast.CommentGroup)(<nil>)
    })
   },
   Closing: (token.Pos) 175
  })
 }),
 Body: (*ast.BlockStmt)(0xc0000eea20)({
  Lbrace: (token.Pos) 177,
  List: ([]ast.Stmt) (len=4 cap=4) {
   (*ast.AssignStmt)(0xc000038740)({
    Lhs: ([]ast.Expr) (len=1 cap=1) {
     (*ast.Ident)(0xc000004600)(factor)
    },
    TokPos: (token.Pos) 188,
    Tok: (token.Token) :=,
    Rhs: ([]ast.Expr) (len=1 cap=1) {
     (*ast.BasicLit)(0xc000004620)({
      ValuePos: (token.Pos) 191,
      Kind: (token.Token) INT,
      Value: (string) (len=1) "1"
     })
    }
   }),
   (*ast.AssignStmt)(0xc000038780)({
    Lhs: ([]ast.Expr) (len=1 cap=1) {
     (*ast.Ident)(0xc000004640)(ind)
    },
    TokPos: (token.Pos) 199,
    Tok: (token.Token) :=,
    Rhs: ([]ast.Expr) (len=1 cap=1) {
     (*ast.Ident)(0xc000004660)(A)
    }
   }),
   (*ast.ForStmt)(0xc000038a40)({
    For: (token.Pos) 206,
    Init: (ast.Stmt) <nil>,
    Cond: (*ast.BinaryExpr)(0xc0000ee900)({
     X: (*ast.Ident)(0xc0000046a0)(ind),
     OpPos: (token.Pos) 214,
     Op: (token.Token) !=,
     Y: (*ast.BasicLit)(0xc0000046c0)({
      ValuePos: (token.Pos) 217,
      Kind: (token.Token) INT,
      Value: (string) (len=1) "0"
     })
    }),
    Post: (ast.Stmt) <nil>,
    Body: (*ast.BlockStmt)(0xc0000ee9f0)({
     Lbrace: (token.Pos) 219,
     List: ([]ast.Stmt) (len=3 cap=4) {
      (*ast.AssignStmt)(0xc000038940)({
       Lhs: ([]ast.Expr) (len=1 cap=1) {
        (*ast.Ident)(0xc0000046e0)(R)
       },
       TokPos: (token.Pos) 226,
       Tok: (token.Token) =,
       Rhs: ([]ast.Expr) (len=1 cap=1) {
        (*ast.BinaryExpr)(0xc0000ee960)({
         X: (*ast.Ident)(0xc000004700)(R),
         OpPos: (token.Pos) 230,
         Op: (token.Token) +,
         Y: (*ast.Ident)(0xc000004720)(factor)
        })
       }
      }),
      (*ast.AssignStmt)(0xc000038980)({
       Lhs: ([]ast.Expr) (len=1 cap=1) {
        (*ast.Ident)(0xc000004740)(factor)
       },
       TokPos: (token.Pos) 249,
       Tok: (token.Token) =,
       Rhs: ([]ast.Expr) (len=1 cap=1) {
        (*ast.BinaryExpr)(0xc0000ee990)({
         X: (*ast.Ident)(0xc000004760)(factor),
         OpPos: (token.Pos) 258,
         Op: (token.Token) +,
         Y: (*ast.BasicLit)(0xc000004780)({
          ValuePos: (token.Pos) 260,
          Kind: (token.Token) INT,
          Value: (string) (len=1) "2"
         })
        })
       }
      }),
      (*ast.AssignStmt)(0xc0000389c0)({
       Lhs: ([]ast.Expr) (len=1 cap=1) {
        (*ast.Ident)(0xc0000047c0)(ind)
       },
       TokPos: (token.Pos) 269,
       Tok: (token.Token) =,
       Rhs: ([]ast.Expr) (len=1 cap=1) {
        (*ast.BinaryExpr)(0xc0000ee9c0)({
         X: (*ast.Ident)(0xc0000047e0)(ind),
         OpPos: (token.Pos) 275,
         Op: (token.Token) -,
         Y: (*ast.BasicLit)(0xc000004800)({
          ValuePos: (token.Pos) 277,
          Kind: (token.Token) INT,
          Value: (string) (len=1) "1"
         })
        })
       }
      })
     },
     Rbrace: (token.Pos) 281
    })
   }),
   (*ast.ReturnStmt)(0xc000004820)({
    Return: (token.Pos) 285,
    Results: ([]ast.Expr) <nil>
   })
  },
  Rbrace: (token.Pos) 293
 })
})

***************************************/
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

/*
%%% === RUN   TestTrivial
**\/
# Name: sff.Square
# Package: sff
# Location: sff.go:9:6
**\/-*- syntax: prolog -*-/
func_info('Square', [blocks(4),name('sff.Square'),pkg(sff),src('sff.go:9:6')], []).

**\/-*- syntax: go-ssa -*-/
func Square(A int) (R int):
0:                                                                entry P:0 S:1
        jump 3
1:                                                             for.body P:1 S:1
        t0 = t3 + t4                                                        int
        t1 = t4 + 2:int                                                     int
        t2 = t5 - 1:int                                                     int
        jump 3
2:                                                             for.done P:1 S:0
        return t3
3:                                                             for.loop P:2 S:2
        t3 = phi [0: 0:int, 1: t0] #R                                       int
        t4 = phi [0: 1:int, 1: t1] #factor                                  int
        t5 = phi [0: A, 1: t2] #ind                                         int
        t6 = t5 != 0:int                                                   bool
        if t6 goto 1 else 2
**\/

%% -*- syntax: prolog -*-
%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%
%%% =====================================
'000': jmp('001') ,
'001': =:(phi([block(0, []),block(2, load(t0))]), t3) ,
=:(phi([block(0, []),block(2, load(t1))]), t4) ,
if(=:('!='(=:(phi([block(0, load('A')),block(2, load(t2))]), t5), []), t6), jmp('002'), jmp('003')) ,
'002': =:(+(=:(phi([block(0, []),block(2, load(t0))]), t3), =:(phi([block(0, []),block(2, load(t1))]), t4)), t0) ,
=:(+(=:(phi([block(0, []),block(2, load(t1))]), t4), []), t1) ,
=:(-(=:(phi([block(0, load('A')),block(2, load(t2))]), t5), []), t2) ,
jmp('001') ,
'003': return([=:(phi([block(0, []),block(2, load(t0))]), t3)]) ,
%%% =====================================.

%%% =====================================
%%% === RUN   TestSSAInterp
%: Square( 10 ) = 100
%: from within
%%% --- PASS: TestSSAInterp (0.02s)
%%% PASS
*/

//XXX: ssa/interp.Interpret panics upon an empty `init`'s body
func init() {
	const Accu = 10
	fmt.Println("Square(", Accu, ") =", Square(Accu))
}

func main() {
	/*********************************

		 0  *ast.FuncDecl {
	     1  .  Name: *ast.Ident {
	     2  .  .  NamePos: testdata/sff.go:77:6
	     3  .  .  Name: "main"
	     4  .  .  Obj: *ast.Object {
	     5  .  .  .  Kind: func
	     6  .  .  .  Name: "main"
	     7  .  .  .  Decl: *(obj @ 0)
	     8  .  .  }
	     9  .  }
	    10  .  Type: *ast.FuncType {
	    11  .  .  Func: testdata/sff.go:77:1
	    12  .  .  Params: *ast.FieldList {
	    13  .  .  .  Opening: testdata/sff.go:77:10
	    14  .  .  .  Closing: testdata/sff.go:77:11
	    15  .  .  }
	    16  .  }
	    17  .  Body: *ast.BlockStmt {
	    18  .  .  Lbrace: testdata/sff.go:77:13
	    19  .  .  List: []ast.Stmt (len = 1) {
	    20  .  .  .  0: *ast.ExprStmt {
	    21  .  .  .  .  X: *ast.CallExpr {
	    22  .  .  .  .  .  Fun: *ast.SelectorExpr {
	    23  .  .  .  .  .  .  X: *ast.Ident {
	    24  .  .  .  .  .  .  .  NamePos: testdata/sff.go:78:2
	    25  .  .  .  .  .  .  .  Name: "fmt"
	    26  .  .  .  .  .  .  }
	    27  .  .  .  .  .  .  Sel: *ast.Ident {
	    28  .  .  .  .  .  .  .  NamePos: testdata/sff.go:78:6
	    29  .  .  .  .  .  .  .  Name: "Println"
	    30  .  .  .  .  .  .  }
	    31  .  .  .  .  .  }
	    32  .  .  .  .  .  Lparen: testdata/sff.go:78:13
	    33  .  .  .  .  .  Args: []ast.Expr (len = 1) {
	    34  .  .  .  .  .  .  0: *ast.BasicLit {
	    35  .  .  .  .  .  .  .  ValuePos: testdata/sff.go:78:14
	    36  .  .  .  .  .  .  .  Kind: STRING
	    37  .  .  .  .  .  .  .  Value: "\"from within\""
	    38  .  .  .  .  .  .  }
	    39  .  .  .  .  .  }
	    40  .  .  .  .  .  Ellipsis: -
	    41  .  .  .  .  .  Rparen: testdata/sff.go:78:27
	    42  .  .  .  .  }
	    43  .  .  .  }
	    44  .  .  }
	    45  .  .  Rbrace: testdata/sff.go:79:1
	    46  .  }
	    47  }

		**********************************

		(*ast.FuncDecl)(0xc0000fad50)({
			Doc: (*ast.CommentGroup)(<nil>),
			Recv: (*ast.FieldList)(<nil>),
		  	Name: (*ast.Ident)(0xc000106a60)(main),

			Type: (*ast.FuncType)(0xc000106b40)({
			  Func: (token.Pos) 2590,
			  Params: (*ast.FieldList)(0xc0000facc0)({
			    Opening: (token.Pos) 2599,
			    List: ([]*ast.Field) <nil>,
			    Closing: (token.Pos) 2600
			  }),
			  Results: (*ast.FieldList)(<nil>)
			}),

			Body: (*ast.BlockStmt)(0xc0000fad20)({
			  Lbrace: (token.Pos) 2602,
			  List: ([]ast.Stmt) (len=1 cap=1) {
			    (*ast.ExprStmt)(0xc00014d7a0)({
				   	X: (*ast.CallExpr)(0xc000104d00)({
					    Fun: (*ast.SelectorExpr)(0xc000106b00)({
						    X: (*ast.Ident)(0xc000106ac0)(fmt),
						    Sel: (*ast.Ident)(0xc000106ae0)(Println)
					    }),
					    Lparen: (token.Pos) 4493,
					    Args: ([]ast.Expr) (len=1 cap=1) {
					     (*ast.BasicLit)(0xc000106b20)({
						      ValuePos: (token.Pos) 4494,
						      Kind: (token.Token) STRING,
						      Value: (string) (len=13) "\"from within\""
					     })
					    },
				   		Ellipsis: (token.Pos) 0,
					    Rparen: (token.Pos) 4507
				   })
			    })
			  },
			  Rbrace: (token.Pos) 4510
			})
		})

		**********************************/
	fmt.Println("from within")
}
