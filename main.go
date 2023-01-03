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
	token := Tokenize(input)
	parser := NewParser(token)
	node := parser.Program()

	Codegen(node)
}
