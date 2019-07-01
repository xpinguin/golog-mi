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

//XXX: ssa/interp.Interpret panics upon an empty `init`'s body
func init() {
	const Accu = 10
	fmt.Println("Square(", Accu, ") =", Square(Accu))
}

func main() {
	fmt.Println("from within")
}
