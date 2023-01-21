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
	for i := range prog.funcs {
		prog.funcs[i].AddType()
	}

	for i := range prog.funcs {
		offset := 0
		for j := range prog.funcs[i].locals {
			offset += prog.funcs[i].locals[j].ty.size()
			prog.funcs[i].locals[j].offset = offset
		}
		prog.funcs[i].stackSize = offset
	}

	prog.Codegen()
}
