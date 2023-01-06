package main

import (
	"fmt"
)

var labelseq int
var argreg = []string{
	"rdi",
	"rsi",
	"rdx",
	"rcx",
	"r8",
	"r9",
}
var funcname string

func genAddr(node Node) {
	if v, ok := node.(*VarNode); ok {
		fmt.Printf("  lea rax, [rbp-%d]\n", v.variable.offset)
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

func gen(node Node) {
	switch v := node.(type) {
	case *Number:
		fmt.Printf("  push %d\n", v.val)
		return
	case *ExpressionStatement:
		gen(v.statement)
		fmt.Printf("  add rsp, 8\n")
		return
	case *VarNode:
		genAddr(node)
		load()
		return
	case *Assign:
		genAddr(v.lhs)
		gen(v.rhs)
		store()
		return
	case *If:
		labelseq++
		seq := labelseq
		if v.els != nil {
			gen(v.cond)
			fmt.Printf("  pop rax\n")
			fmt.Printf("  cmp rax, 0\n")
			fmt.Printf("  je .L.else.%d\n", seq)
			gen(v.then)
			fmt.Printf("  jmp .L.end.%d\n", seq)
			fmt.Printf(".L.else.%d\n", seq)
			gen(v.els)
			fmt.Printf(".L.end.%d:\n", seq)
		} else {
			gen(v.cond)
			fmt.Printf("  pop rax\n")
			fmt.Printf("  cmp rax, 0\n")
			fmt.Printf("  je .L.end.%d\n", seq)
			gen(v.then)
			fmt.Printf(".L.end.%d:\n", seq)
		}
		return
	case *While:
		labelseq++
		seq := labelseq
		fmt.Printf(".L.begin.%d:\n", seq)
		gen(v.cond)
		fmt.Printf("  pop rax\n")
		fmt.Printf("  cmp rax, 0\n")
		fmt.Printf("  je .L.end.%d\n", seq)
		gen(v.then)
		fmt.Printf("  jmp .L.begin.%d\n", seq)
		fmt.Printf(".L.end.%d:\n", seq)
		return
	case *For:
		labelseq++
		seq := labelseq
		if v.init != nil {
			gen(v.init)
		}
		fmt.Printf(".L.begin.%d:\n", seq)
		if v.cond != nil {
			gen(v.cond)
			fmt.Printf("  pop rax\n")
			fmt.Printf("  cmp rax, 0\n")
			fmt.Printf("  je .L.end.%d\n", seq)
		}
		gen(v.block)
		if v.inc != nil {
			gen(v.inc)
		}
		fmt.Printf("  jmp .L.begin.%d\n", seq)
		fmt.Printf(".L.end.%d:\n", seq)
		return
	case *Block:
		for _, n := range v.body {
			gen(n)
		}
		return
	case *FuncCall:
		nargs := 0
		for _, arg := range v.args {
			gen(arg)
			nargs++
		}
		for i := nargs - 1; i >= 0; i-- {
			fmt.Printf("  pop %s\n", argreg[i])
		}

		labelseq++
		seq := labelseq
		fmt.Printf("  mov rax, rsp\n")
		fmt.Printf("  and rax, 15\n")
		fmt.Printf("  jnz .L.call.%d\n", seq)
		fmt.Printf("  mov rax, 0\n")
		fmt.Printf("  call %s\n", v.name)
		fmt.Printf("  jmp .L.end.%d\n", seq)
		fmt.Printf(".L.call.%d:\n", seq)
		fmt.Printf("  sub rsp, 8\n")
		fmt.Printf("  mov rax, 0\n")
		fmt.Printf("  call %s\n", v.name)
		fmt.Printf("  add rsp, 8\n")
		fmt.Printf(".L.end.%d:\n", seq)
		fmt.Printf("  push rax\n")
		return
	case *Return:
		gen(v.expr)
		fmt.Printf("  pop rax\n")
		fmt.Printf("  jmp .L.return.%s\n", funcname)
		return
	}

	f := func(lhs Node, rhs Node) {
		gen(lhs)
		gen(rhs)

		fmt.Printf("  pop rdi\n")
		fmt.Printf("  pop rax\n")
	}

	switch v := node.(type) {
	case *Add:
		f(v.lhs, v.rhs)
		fmt.Printf("  add rax, rdi\n")
	case *Sub:
		f(v.lhs, v.rhs)
		fmt.Printf("  sub rax, rdi\n")
	case *Mul:
		f(v.lhs, v.rhs)
		fmt.Printf("  imul rax, rdi\n")
	case *Div:
		f(v.lhs, v.rhs)
		fmt.Printf("  cqo\n")
		fmt.Printf("  idiv rdi\n")
	case *Equal:
		f(v.lhs, v.rhs)
		fmt.Printf("  cmp rax, rdi\n")
		fmt.Printf("  sete al\n")
		fmt.Printf("  movzb rax, al\n")
	case *NotEqual:
		f(v.lhs, v.rhs)
		fmt.Printf("  cmp rax, rdi\n")
		fmt.Printf("  setne al\n")
		fmt.Printf("  movzb rax, al\n")
	case *LessThan:
		f(v.lhs, v.rhs)
		fmt.Printf("  cmp rax, rdi\n")
		fmt.Printf("  setl al\n")
		fmt.Printf("  movzb rax, al\n")
	case *LessEqual:
		f(v.lhs, v.rhs)
		fmt.Printf("  cmp rax, rdi\n")
		fmt.Printf("  setle al\n")
		fmt.Printf("  movzb rax, al\n")
	}

	fmt.Printf("  push rax\n")
}

func Codegen(prog []*Function) {
	fmt.Printf(".intel_syntax noprefix\n")

	for _, fn := range prog {
		fmt.Printf(".global %s\n", fn.name)
		fmt.Printf("%s:\n", fn.name)
		funcname = fn.name

		fmt.Printf("  push rbp\n")
		fmt.Printf("  mov rbp, rsp\n")
		fmt.Printf("  sub rsp, %d\n", fn.stackSize)

		i := 0
		for _, v := range fn.params {
			fmt.Printf("  mov [rbp-%d], %s\n", v.offset, argreg[i])
			i++
		}

		for _, n := range fn.node {
			gen(n)
		}

		fmt.Printf(".L.return.%s:\n", funcname)
		fmt.Printf("  mov rsp, rbp\n")
		fmt.Printf("  pop rbp\n")
		fmt.Printf("  ret\n")
	}
}
