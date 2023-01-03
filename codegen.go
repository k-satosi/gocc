package main

import (
	"fmt"
)

var labelseq int

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
	case ND_IF:
		labelseq++
		seq := labelseq
		if node.els != nil {
			gen(node.cond)
			fmt.Printf("  pop rax\n")
			fmt.Printf("  cmp rax, 0\n")
			fmt.Printf("  je .L.else.%d\n", seq)
			gen(node.then)
			fmt.Printf("  jmp .L.end.%d\n", seq)
			fmt.Printf(".L.else.%d\n", seq)
			gen(node.els)
			fmt.Printf(".L.end.%d:\n", seq)
		} else {
			gen(node.cond)
			fmt.Printf("  pop rax\n")
			fmt.Printf("  cmp rax, 0\n")
			fmt.Printf("  je .L.end.%d\n", seq)
			gen(node.then)
			fmt.Printf(".L.end.%d:\n", seq)
		}
		return
	case ND_WHILE:
		labelseq++
		seq := labelseq
		fmt.Printf(".L.begin.%d:\n", seq)
		gen(node.cond)
		fmt.Printf("  pop rax\n")
		fmt.Printf("  cmp rax, 0\n")
		fmt.Printf("  je .L.end.%d\n", seq)
		gen(node.then)
		fmt.Printf("  jmp .L.begin.%d\n", seq)
		fmt.Printf(".L.end.%d:\n", seq)
		return
	case ND_FOR:
		labelseq++
		seq := labelseq
		if node.init != nil {
			gen(node.init)
		}
		fmt.Printf(".L.begin.%d:\n", seq)
		if node.cond != nil {
			gen(node.cond)
			fmt.Printf("  pop rax\n")
			fmt.Printf("  cmp rax, 0\n")
			fmt.Printf("  je .L.end.%d\n", seq)
		}
		gen(node.then)
		if node.inc != nil {
			gen(node.inc)
		}
		fmt.Printf("  jmp .L.begin.%d\n", seq)
		fmt.Printf(".L.end.%d:\n", seq)
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
