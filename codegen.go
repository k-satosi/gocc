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

func (v *VarNode) genAddr() {
	if v.variable.isLocal {
		fmt.Printf("  lea rax, [rbp-%d]\n", v.variable.offset)
		fmt.Printf("  push rax\n")
	} else {
		fmt.Printf("  push offset %s\n", v.variable.name)
	}
}

func (m *Member) genAddr() {
	if v, ok := m.expr.(*VarNode); ok {
		v.genAddr()
	}
	fmt.Printf("  pop rax\n")
	fmt.Printf("  add rax, %d\n", m.offset)
	fmt.Printf("  push rax\n")
}

func (n *Dereference) genAddr() {
	n.expr.Gen()
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
	e.statement.Gen()
	fmt.Printf("  add rsp, 8\n")
}

func (v *VarNode) Gen() {
	v.genAddr()
	if _, ok := v.variable.ty.(*ArrayType); !ok {
		load()
	}
}

func (m *Member) Gen() {
	m.genAddr()
	if _, ok := m.ty.(*ArrayType); !ok {
		load()
	}
}

func (s *Sizeof) Gen() {
	fmt.Printf("  push %d\n", s.v.Type().size())
}

func (a *Assign) Gen() {
	switch v := a.lhs.(type) {
	case *VarNode:
		v.genAddr()
	case *Dereference:
		v.genAddr()
	case *Member:
		v.genAddr()
	}
	a.rhs.Gen()
	store()
}

func (a *Address) Gen() {
	switch v := a.expr.(type) {
	case *VarNode:
		v.genAddr()
	case *Dereference:
		v.genAddr()
	case *Member:
		v.genAddr()
	}
}

func (d *Dereference) Gen() {
	d.expr.Gen()
	if _, ok := d.ty.(*ArrayType); !ok {
		load()
	}
}

func (i *If) Gen() {
	labelseq++
	seq := labelseq
	if i.els != nil {
		i.cond.Gen()
		fmt.Printf("  pop rax\n")
		fmt.Printf("  cmp rax, 0\n")
		fmt.Printf("  je .L.else.%d\n", seq)
		i.then.Gen()
		fmt.Printf("  jmp .L.end.%d\n", seq)
		fmt.Printf(".L.else.%d\n", seq)
		i.els.Gen()
		fmt.Printf(".L.end.%d:\n", seq)
	} else {
		i.cond.Gen()
		fmt.Printf("  pop rax\n")
		fmt.Printf("  cmp rax, 0\n")
		fmt.Printf("  je .L.end.%d\n", seq)
		i.then.Gen()
		fmt.Printf(".L.end.%d:\n", seq)
	}
}

func (w *While) Gen() {
	labelseq++
	seq := labelseq
	fmt.Printf(".L.begin.%d:\n", seq)
	w.cond.Gen()
	fmt.Printf("  pop rax\n")
	fmt.Printf("  cmp rax, 0\n")
	fmt.Printf("  je .L.end.%d\n", seq)
	w.then.Gen()
	fmt.Printf("  jmp .L.begin.%d\n", seq)
	fmt.Printf(".L.end.%d:\n", seq)
}

func (f *For) Gen() {
	labelseq++
	seq := labelseq
	if f.init != nil {
		f.init.Gen()
	}
	fmt.Printf(".L.begin.%d:\n", seq)
	if f.cond != nil {
		f.cond.Gen()
		fmt.Printf("  pop rax\n")
		fmt.Printf("  cmp rax, 0\n")
		fmt.Printf("  je .L.end.%d\n", seq)
	}
	f.block.Gen()
	if f.inc != nil {
		f.inc.Gen()
	}
	fmt.Printf("  jmp .L.begin.%d\n", seq)
	fmt.Printf(".L.end.%d:\n", seq)
}

func (b *Block) Gen() {
	for _, n := range b.body {
		n.Gen()
	}
}

func (f *FuncCall) Gen() {
	nargs := 0
	for _, arg := range f.args {
		arg.Gen()
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
	r.expr.Gen()
	fmt.Printf("  pop rax\n")
	fmt.Printf("  jmp .L.return.%s\n", funcname)
}

func (b *Binary) Gen() {
	b.Lhs().Gen()
	b.Rhs().Gen()
	fmt.Printf("  pop rdi\n")
	fmt.Printf("  pop rax\n")
}

func (a *Add) Gen() {
	a.Binary.Gen()
	switch t := a.Binary.ty.(type) {
	case *PointerType:
		fmt.Printf("  imul rdi, %d\n", t.base.size())
	case *ArrayType:
		fmt.Printf("  imul rdi, %d\n", t.base.size())
	}
	fmt.Printf("  add rax, rdi\n")
	fmt.Printf("  push rax\n")
}

func (s *Sub) Gen() {
	s.Binary.Gen()
	switch t := s.ty.(type) {
	case *PointerType:
		fmt.Printf("  imul rdi, %d\n", t.size())
	case *ArrayType:
		fmt.Printf("  imul rdi, %d\n", t.size())
	}
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

func (n *Null) Gen() {}

type Program struct {
	globals []*Variable
	funcs   []*Function
}

func (p *Program) emitData() {
	fmt.Printf(".data\n")

	for _, v := range p.globals {
		fmt.Printf("%s:\n", v.name)

		if len(v.contents) == 0 {
			fmt.Printf("  .zero %d\n", v.ty.size())
			continue
		}

		for _, r := range v.contents {
			fmt.Printf("  .byte %d\n", r)
		}
	}
}

func (p *Program) emitText() {
	fmt.Printf(".text\n")

	for _, fn := range p.funcs {
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
			n.Gen()
		}

		fmt.Printf(".L.return.%s:\n", funcname)
		fmt.Printf("  mov rsp, rbp\n")
		fmt.Printf("  pop rbp\n")
		fmt.Printf("  ret\n")
	}
}

func (p *Program) Codegen() {
	fmt.Printf(".intel_syntax noprefix\n")

	p.emitData()
	p.emitText()
}
