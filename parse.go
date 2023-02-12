package main

import (
	"fmt"
)

type Function struct {
	name   string
	params []*Variable

	node      []Node
	locals    []*Variable
	stackSize int
}

func (p *Parser) pushVar(name string, ty Type, isLocal bool) *Variable {
	v := &Variable{
		name:    name,
		ty:      ty,
		isLocal: isLocal,
	}
	if isLocal {
		p.locals = append([]*Variable{v}, p.locals...)
	} else {
		p.globals = append([]*Variable{v}, p.globals...)
	}

	p.scope = append([]*Variable{v}, p.scope...)
	return v
}

var labelCount int

func (p *Parser) newLabel() string {
	label := fmt.Sprintf(".L.data.%d", labelCount)
	labelCount++
	return label
}

type Parser struct {
	token   *Token
	locals  []*Variable
	globals []*Variable
	scope   []*Variable
}

func NewParser(token *Token) *Parser {
	return &Parser{
		token: token,
	}
}

func (p *Parser) isFunction() bool {
	tok := p.token
	p.baseType()
	isFunc := p.consumeIdent() != nil && p.consume("(")
	p.token = tok
	return isFunc
}

func (p *Parser) Program() *Program {
	funcs := []*Function{}

	for !p.token.AtEOF() {
		if p.isFunction() {
			funcs = append(funcs, p.function())
		} else {
			p.globalVar()
		}
	}
	prog := &Program{
		globals: p.globals,
		funcs:   funcs,
	}
	return prog
}

func (p *Parser) baseType() Type {
	var ty Type
	if p.consume("char") {
		ty = charType
	} else if p.consume("int") {
		ty = intType
	} else {
		ty = p.structDecl()
	}
	for p.consume("*") {
		ty = NewPointerType(ty)
	}
	return ty
}

func (p *Parser) readTypeSuffix(base Type) Type {
	if !p.consume("[") {
		return base
	}
	size := p.expectNumber()
	p.expect(("]"))
	base = p.readTypeSuffix(base)
	return NewArrayType(base, size)
}

func (p *Parser) structDecl() Type {
	p.expect("struct")
	p.expect("{")

	members := []*Member{}

	for !p.consume("}") {
		members = append(members, p.structMember())
		//members = append([]*Member{p.structMember()}, members...)
	}

	ty := NewStructType(members)
	offset := 0
	for i := range ty.members {
		ty.members[i].offset = offset
		offset += ty.members[i].ty.size()
	}

	return ty
}

func (p *Parser) structMember() *Member {
	ty := p.baseType()
	name := p.expectIdent()
	m := &Member{
		ty:   ty,
		name: name,
	}
	m.ty = p.readTypeSuffix(m.ty)
	p.expect(";")
	return m
}

func (p *Parser) readFuncParam() *Variable {
	ty := p.baseType()
	name := p.expectIdent()
	ty = p.readTypeSuffix(ty)
	return p.pushVar(name, ty, true)
}

func (p *Parser) readFuncParams() []*Variable {
	if p.consume(")") {
		return nil
	}

	l := []*Variable{p.readFuncParam()}

	for !p.consume(")") {
		p.expect(",")
		l = append(l, p.readFuncParam())
	}

	return l
}

func (p *Parser) function() *Function {
	p.locals = []*Variable{}

	fn := &Function{}
	p.baseType()
	fn.name = p.expectIdent()
	p.expect("(")
	fn.params = p.readFuncParams()
	p.expect("{")

	l := []Node{}
	for !p.consume("}") {
		l = append(l, p.stmt())
	}

	fn.node = l
	fn.locals = p.locals
	return fn
}

func (f *Function) AddType() {
	for i := range f.node {
		f.node[i].AddType()
	}
}

func (p *Parser) globalVar() {
	ty := p.baseType()
	name := p.expectIdent()
	ty = p.readTypeSuffix(ty)
	p.expect(";")
	p.pushVar(name, ty, false)
}

func (p *Parser) declaration() Node {
	ty := p.baseType()
	ident := p.expectIdent()
	ty = p.readTypeSuffix(ty)
	v := p.pushVar(ident, ty, true)
	if p.consume(";") {
		return NewNull()
	}

	p.expect("=")
	lhs := NewVarNode(v)
	rhs := p.expr()
	p.expect(";")
	node := NewAssign(lhs, rhs)
	return NewExpressionStatement(node)
}

func (p *Parser) readExprStmt() Node {
	return NewExpressionStatement(p.expr())
}

func (p *Parser) isTypeName() bool {
	return p.peek("char") || p.peek("int") || p.peek("struct")
}

func (p *Parser) stmt() Node {
	node := p.stmt2()
	node.AddType()
	return node
}

func (p *Parser) stmt2() Node {
	if p.consume("return") {
		node := NewReturn(p.expr())
		p.expect(";")
		return node
	}

	if p.consume("if") {
		var cond Node
		var then Node
		var els Node
		p.expect("(")
		cond = p.expr()
		p.expect(")")
		then = p.stmt()
		if p.consume("else") {
			els = p.stmt()
		}
		return NewIf(cond, then, els)
	}

	if p.consume("while") {
		p.expect("(")
		cond := p.expr()
		p.expect(")")
		then := p.stmt()
		return NewWhile(cond, then)
	}

	if p.consume("for") {
		var init Node
		var cond Node
		var inc Node
		var block Node
		p.expect("(")
		if !p.consume(";") {
			init = p.readExprStmt()
			p.expect(";")
		}
		if !p.consume(";") {
			cond = p.expr()
			p.expect(";")
		}
		if !p.consume(")") {
			inc = p.readExprStmt()
			p.expect(")")
		}
		block = p.stmt()
		return NewFor(init, cond, inc, block)
	}

	if p.consume("{") {
		l := []Node{}

		sc := p.scope
		for !p.consume("}") {
			l = append(l, p.stmt())
		}
		p.scope = sc

		node := NewBlock(l)

		return node
	}

	if p.isTypeName() {
		return p.declaration()
	}

	node := p.readExprStmt()
	p.expect(";")
	return node
}

func (p *Parser) expr() Node {
	return p.assign()
}

func (p *Parser) assign() Node {
	node := p.equality()
	if p.consume("=") {
		node = NewAssign(node, p.assign())
	}

	return node
}

func (p *Parser) equality() Node {
	node := p.relational()

	for {
		if p.consume("==") {
			node = NewEqual(node, p.relational())
		} else if p.consume("!=") {
			node = NewNotEqual(node, p.relational())
		} else {
			return node
		}
	}
}

func (p *Parser) relational() Node {
	node := p.add()

	for {
		if p.consume("<") {
			node = NewLessThan(node, p.add())
		} else if p.consume("<=") {
			node = NewLessEqual(node, p.add())
		} else if p.consume(">") {
			node = NewLessThan(p.add(), node)
		} else if p.consume(">=") {
			node = NewLessEqual(p.add(), node)
		} else {
			return node
		}
	}
}

func (p *Parser) add() Node {
	node := p.mul()

	for {
		if p.consume("+") {
			node = NewAdd(node, p.mul())
		} else if p.consume(("-")) {
			node = NewSub(node, p.mul())
		} else {
			return node
		}
	}
}

func (p *Parser) mul() Node {
	node := p.unary()

	for {
		if p.consume("*") {
			node = NewMul(node, p.unary())
		} else if p.consume("/") {
			node = NewDiv(node, p.unary())
		} else {
			return node
		}
	}
}

func (p *Parser) unary() Node {
	if p.consume("+") {
		return p.unary()
	} else if p.consume("-") {
		return NewSub(NewNumber(0), p.unary())
	} else if p.consume("&") {
		return NewAddress(p.unary())
	} else if p.consume("*") {
		return NewDereference(p.unary())
	} else {
		return p.postFix()
	}
}

func (p *Parser) postFix() Node {
	node := p.primary()

	for {
		if p.consume("[") {
			exp := NewAdd(node, p.expr())
			p.expect("]")
			node = NewDereference(exp)
			continue
		}

		if p.consume(".") {
			name := p.expectIdent()
			node = NewMember(node, name)
			continue
		}
		return node
	}
}

func (p *Parser) funcArgs() []Node {
	if p.consume(")") {
		return nil
	}

	l := []Node{p.assign()}
	for p.consume(",") {
		l = append(l, p.assign())
	}
	p.expect(")")
	return l
}

func (p *Parser) primary() Node {
	if p.consume("(") {
		node := p.expr()
		p.expect(")")
		return node
	}

	if p.consume("sizeof") {
		return NewSizeof(p.unary())
	}

	if token := p.consumeIdent(); token != nil {
		if p.consume("(") {
			name := token.str
			args := p.funcArgs()
			return NewFuncCall(name, args)
		}
		v := p.findVariable(token)
		if v == nil {
			errorToken(token, "undefined variable")
		}
		return NewVarNode(v)
	}

	tok := p.token
	if tok.kind == TK_STRING {
		p.token = p.token.next

		ty := NewArrayType(charType, len(tok.contents))
		v := p.pushVar(p.newLabel(), ty, false)
		v.contents = tok.contents
		return NewVarNode(v)
	}

	return NewNumber(p.expectNumber())
}

func (p *Parser) findVariable(token *Token) *Variable {
	for i := range p.scope {
		if token.str == p.scope[i].name {
			return p.scope[i]
		}
	}
	return nil
}
