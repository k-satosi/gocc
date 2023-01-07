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

func (n *Number) Gen() {
	fmt.Printf("  push %d\n", n.val)
}

func (e *ExpressionStatement) Gen() {
	gen(e.statement)
	fmt.Printf("  add rsp, 8\n")
}

func (v *VarNode) Gen() {
	genAddr(v)
	load()
}

func (a *Assign) Gen() {
	genAddr(a.lhs)
	gen(a.rhs)
	store()
}

func (i *If) Gen() {
	labelseq++
	seq := labelseq
	if i.els != nil {
		gen(i.cond)
		fmt.Printf("  pop rax\n")
		fmt.Printf("  cmp rax, 0\n")
		fmt.Printf("  je .L.else.%d\n", seq)
		gen(i.then)
		fmt.Printf("  jmp .L.end.%d\n", seq)
		fmt.Printf(".L.else.%d\n", seq)
		gen(i.els)
		fmt.Printf(".L.end.%d:\n", seq)
	} else {
		gen(i.cond)
		fmt.Printf("  pop rax\n")
		fmt.Printf("  cmp rax, 0\n")
		fmt.Printf("  je .L.end.%d\n", seq)
		gen(i.then)
		fmt.Printf(".L.end.%d:\n", seq)
	}
}

func (w *While) Gen() {
	labelseq++
	seq := labelseq
	fmt.Printf(".L.begin.%d:\n", seq)
	gen(w.cond)
	fmt.Printf("  pop rax\n")
	fmt.Printf("  cmp rax, 0\n")
	fmt.Printf("  je .L.end.%d\n", seq)
	gen(w.then)
	fmt.Printf("  jmp .L.begin.%d\n", seq)
	fmt.Printf(".L.end.%d:\n", seq)
}

func (f *For) Gen() {
	labelseq++
	seq := labelseq
	if f.init != nil {
		gen(f.init)
	}
	fmt.Printf(".L.begin.%d:\n", seq)
	if f.cond != nil {
		gen(f.cond)
		fmt.Printf("  pop rax\n")
		fmt.Printf("  cmp rax, 0\n")
		fmt.Printf("  je .L.end.%d\n", seq)
	}
	gen(f.block)
	if f.inc != nil {
		gen(f.inc)
	}
	fmt.Printf("  jmp .L.begin.%d\n", seq)
	fmt.Printf(".L.end.%d:\n", seq)
}

func (b *Block) Gen() {
	for _, n := range b.body {
		gen(n)
	}
}

func (f *FuncCall) Gen() {
	nargs := 0
	for _, arg := range f.args {
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
	fmt.Printf("  call %s\n", f.name)
	fmt.Printf("  jmp .L.end.%d\n", seq)
	fmt.Printf(".L.call.%d:\n", seq)
	fmt.Printf("  sub rsp, 8\n")
	fmt.Printf("  mov rax, 0\n")
	fmt.Printf("  call %s\n", f.name)
	fmt.Printf("  add rsp, 8\n")
	fmt.Printf(".L.end.%d:\n", seq)
	fmt.Printf("  push rax\n")
}

func (r *Return) Gen() {
	gen(r.expr)
	fmt.Printf("  pop rax\n")
	fmt.Printf("  jmp .L.return.%s\n", funcname)
}

func (b *Binary) Gen() {
	gen(b.Lhs())
	gen(b.Rhs())
	fmt.Printf("  pop rdi\n")
	fmt.Printf("  pop rax\n")
}

func (a *Add) Gen() {
	a.Binary.Gen()
	fmt.Printf("  add rax, rdi\n")
	fmt.Printf("  push rax\n")
}

func (s *Sub) Gen() {
	s.Binary.Gen()
	fmt.Printf("  sub rax, rdi\n")
	fmt.Printf("  push rax\n")
}

func (m *Mul) Gen() {
	m.Binary.Gen()
	fmt.Printf("  imul rax, rdi\n")
	fmt.Printf("  push rax\n")
}

func (d *Div) Gen() {
	d.Binary.Gen()
	fmt.Printf("  cqo\n")
	fmt.Printf("  idiv rdi\n")
	fmt.Printf("  push rax\n")
}

func (e *Equal) Gen() {
	e.Binary.Gen()
	fmt.Printf("  cmp rax, rdi\n")
	fmt.Printf("  sete al\n")
	fmt.Printf("  movzb rax, al\n")
	fmt.Printf("  push rax\n")
}

func (n *NotEqual) Gen() {
	n.Binary.Gen()
	fmt.Printf("  cmp rax, rdi\n")
	fmt.Printf("  setne al\n")
	fmt.Printf("  movzb rax, al\n")
	fmt.Printf("  push rax\n")
}

func (l *LessThan) Gen() {
	l.Binary.Gen()
	fmt.Printf("  cmp rax, rdi\n")
	fmt.Printf("  setl al\n")
	fmt.Printf("  movzb rax, al\n")
	fmt.Printf("  push rax\n")
}

func (l *LessEqual) Gen() {
	l.Binary.Gen()
	fmt.Printf("  cmp rax, rdi\n")
	fmt.Printf("  setle al\n")
	fmt.Printf("  movzb rax, al\n")
	fmt.Printf("  push rax\n")
}

func gen(node Node) {
	node.Gen()
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
