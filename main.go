package main

import (
	"fmt"
	"os"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Printf("usage: %v \"<program>\"\n", os.Args[0])
		return
	}
	input := os.Args[1]
	compile(input)
}

func compile(input string) {
	token := Tokenize(input)
	parser := NewParser(token)
	prog := parser.Program()

	for fn := prog; fn != nil; fn = fn.next {
		offset := 0
		for vl := fn.locals; vl != nil; vl = vl.next {
			offset += 8
			vl.variable.offset = offset
		}
		prog.stackSize = offset
	}

	Codegen(prog)
}
