package main

import (
	"fmt"
	"os"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Printf("%v: invalid number of arguments", os.Args[0])
		return
	}
	bytes, err := os.ReadFile(os.Args[1])
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	input := string(bytes)
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
