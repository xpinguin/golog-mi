package main

import (
	"fmt"
	_ "runtime" //XXX: neccessary for the ssa/interp.Interpret
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

/*
%%% === RUN   TestTrivial
**\/
# Name: sff.Square
# Package: sff
# Location: sff.go:9:6
**\/
func_info('Square', [blocks(4),name('sff.Square'),pkg(sff),src('sff.go:9:6')], []).

**\/
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
	fmt.Println("from within")
}
