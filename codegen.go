package main

import (
	"fmt"
)

func genAddr(node *Node) {
	if node.kind == ND_VAR {
		fmt.Printf("  lea rax, [rbp-%d]\n", node.variable.offset)
		fmt.Printf("  push rax\n")
		return
	}

	fmt.Printf("not an lvalue")
}

func load() {
	fmt.Printf("  pop rax\n")
	fmt.Printf("  mov rax, [rax]\n")
	fmt.Printf("  push rax\n")
}

func store() {
	fmt.Printf("  pop rdi\n")
	fmt.Printf("  pop rax\n")
	fmt.Printf("  mov [rax], rdi\n")
	fmt.Printf("  push rdi\n")
}

func gen(node *Node) {
	switch node.kind {
	case ND_NUM:
		fmt.Printf("  push %d\n", node.val)
		return
	case ND_EXPR_STMT:
		gen(node.lhs)
		fmt.Printf("  add rsp, 8\n")
		return
	case ND_VAR:
		genAddr(node)
		load()
		return
	case ND_ASSIGN:
		genAddr(node.lhs)
		gen(node.rhs)
		store()
		return
	case ND_RETURN:
		gen(node.lhs)
		fmt.Printf("  pop rax\n")
		fmt.Printf("  jmp .L.return\n")
		return
	}

	gen(node.lhs)
	gen(node.rhs)

	fmt.Printf("  pop rdi\n")
	fmt.Printf("  pop rax\n")

	switch node.kind {
	case ND_ADD:
		fmt.Printf("  add rax, rdi\n")
	case ND_SUB:
		fmt.Printf("  sub rax, rdi\n")
	case ND_MUL:
		fmt.Printf("  imul rax, rdi\n")
	case ND_DIV:
		fmt.Printf("  cqo\n")
		fmt.Printf("  idiv rdi\n")
	case ND_EQ:
		fmt.Printf("  cmp rax, rdi\n")
		fmt.Printf("  sete al\n")
		fmt.Printf("  movzb rax, al\n")
	case ND_NE:
		fmt.Printf("  cmp rax, rdi\n")
		fmt.Printf("  setne al\n")
		fmt.Printf("  movzb rax, al\n")
	case ND_LT:
		fmt.Printf("  cmp rax, rdi\n")
		fmt.Printf("  setl al\n")
		fmt.Printf("  movzb rax, al\n")
	case ND_LE:
		fmt.Printf("  cmp rax, rdi\n")
		fmt.Printf("  setle al\n")
		fmt.Printf("  movzb rax, al\n")
	}

	fmt.Printf("  push rax\n")
}

func Codegen(prog *Function) {
	fmt.Printf(".intel_syntax noprefix\n")
	fmt.Printf(".global main\n")
	fmt.Printf("main:\n")

	fmt.Printf("  push rbp\n")
	fmt.Printf("  mov rbp, rsp\n")
	fmt.Printf("  sub rsp, %d\n", prog.stackSize)

	for n := prog.node; n != nil; n = n.next {
		gen(n)
	}

	fmt.Printf(".L.return:\n")
	fmt.Printf("  mov rsp, rbp\n")
	fmt.Printf("  pop rbp\n")
	fmt.Printf("  ret\n")
}