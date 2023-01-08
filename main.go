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

	for i := range prog {
		offset := 0
		for j := range prog[i].locals {
			offset += prog[i].locals[j].ty.size()
			prog[i].locals[j].offset = offset
		}
		prog[i].stackSize = offset
	}

	Codegen(prog)
}
