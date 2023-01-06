package main

type Function struct {
	name   string
	params []*Variable

	node      []Node
	locals    []*Variable
	stackSize int
}

func (p *Parser) NewLVar(name string) *Variable {
	v := &Variable{
		name: name,
	}
	p.locals = append(p.locals, v)
	return v
}

type Parser struct {
	token  *Token
	locals []*Variable
}

func NewParser(token *Token) *Parser {
	return &Parser{
		token: token,
	}
}

func (p *Parser) Program() []*Function {
	funcs := []*Function{}

	for !p.token.AtEOF() {
		funcs = append(funcs, p.function())
	}
	return funcs
}

func (p *Parser) readFuncParams() []*Variable {
	if p.consume(")") {
		return nil
	}

	v := p.NewLVar(p.expectIdent())
	l := []*Variable{v}

	for !p.consume(")") {
		p.expect(",")
		l = append(l, p.NewLVar(p.expectIdent()))
	}

	return l
}

func (p *Parser) function() *Function {
	p.locals = []*Variable{}

	fn := &Function{}
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

func (p *Parser) readExprStmt() Node {
	return NewExpressionStatement(p.expr())
}

func (p *Parser) stmt() Node {
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

		for !p.consume("}") {
			l = append(l, p.stmt())
		}

		node := NewBlock(l)

		return node
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
			node = NewEqual(&node, p.relational())
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
	} else {
		return p.primary()
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

	if token := p.consumeIdent(); token != nil {
		if p.consume("(") {
			name := token.str
			args := p.funcArgs()
			return NewFuncCall(name, args)
		}
		v := p.findVariable(token)
		if v == nil {
			v = p.NewLVar(token.str)
		}
		return NewVarNode(v)
	}

	return NewNumber(p.expectNumber())
}

func (p *Parser) findVariable(token *Token) *Variable {
	for i := range p.locals {
		if token.str == p.locals[i].name {
			return p.locals[i]
		}
	}
	return nil
}
