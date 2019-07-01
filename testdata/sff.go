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

//XXX: ssa/interp.Interpret panics upon empty `init`'s body
func init() {
	fmt.Println("from init()")
}

func main() {
	fmt.Println("from within")
}
