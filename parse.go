package main

type Function struct {
	node      *Node
	locals    *Variable
	stackSize int
}

type NodeKind int

const (
	ND_ADD NodeKind = iota
	ND_SUB
	ND_MUL
	ND_DIV
	ND_EQ
	ND_NE
	ND_LT
	ND_LE
	ND_ASSIGN
	ND_RETURN
	ND_EXPR_STMT
	ND_VAR
	ND_NUM
)

type Node struct {
	kind     NodeKind
	next     *Node
	lhs      *Node
	rhs      *Node
	variable *Variable
	val      int
}

func NewNode(kind NodeKind) *Node {
	return &Node{
		kind: kind,
	}
}

func NewBinary(kind NodeKind, lhs *Node, rhs *Node) *Node {
	return &Node{
		kind: kind,
		lhs:  lhs,
		rhs:  rhs,
	}
}

func NewUnary(kind NodeKind, expr *Node) *Node {
	return &Node{
		kind: kind,
		lhs:  expr,
	}
}

func NewNum(val int) *Node {
	node := NewNode(ND_NUM)
	node.val = val
	return node
}

func NewVarNode(v *Variable) *Node {
	node := NewNode(ND_VAR)
	node.variable = v
	return node
}

func (p *Parser) newLvar(name string) *Variable {
	v := &Variable{
		next: p.locals,
		name: name,
	}
	p.locals = v
	return v
}

type Parser struct {
	token  *Token
	locals *Variable
}

func NewParser(token *Token) *Parser {
	return &Parser{
		token: token,
	}
}

func (p *Parser) Program() *Function {
	head := Node{}
	cur := &head

	for !p.token.AtEOF() {
		cur.next = p.stmt()
		cur = cur.next
	}

	return &Function{
		node:   head.next,
		locals: p.locals,
	}
}

func (p *Parser) stmt() *Node {
	if p.consume("return") {
		node := NewUnary(ND_RETURN, p.expr())
		p.expect(";")
		return node
	}

	node := NewUnary(ND_EXPR_STMT, p.expr())
	p.expect(";")
	return node
}

func (p *Parser) expr() *Node {
	return p.assign()
}

func (p *Parser) assign() *Node {
	node := p.equality()
	if p.consume("=") {
		node = NewBinary(ND_ASSIGN, node, p.assign())
	}

	return node
}

func (p *Parser) equality() *Node {
	node := p.relational()

	for {
		if p.consume("==") {
			node = NewBinary(ND_EQ, node, p.relational())
		} else if p.consume("!=") {
			node = NewBinary(ND_NE, node, p.relational())
		} else {
			return node
		}
	}
}

func (p *Parser) relational() *Node {
	node := p.add()

	for {
		if p.consume("<") {
			node = NewBinary(ND_LT, node, p.add())
		} else if p.consume("<=") {
			node = NewBinary(ND_LE, node, p.add())
		} else if p.consume(">") {
			node = NewBinary(ND_LT, p.add(), node)
		} else if p.consume(">=") {
			node = NewBinary(ND_LE, p.add(), node)
		} else {
			return node
		}
	}
}

func (p *Parser) add() *Node {
	node := p.mul()

	for {
		if p.consume("+") {
			node = NewBinary(ND_ADD, node, p.mul())
		} else if p.consume(("-")) {
			node = NewBinary(ND_SUB, node, p.mul())
		} else {
			return node
		}
	}
}

func (p *Parser) mul() *Node {
	node := p.unary()

	for {
		if p.consume("*") {
			node = NewBinary(ND_MUL, node, p.unary())
		} else if p.consume("/") {
			node = NewBinary(ND_DIV, node, p.unary())
		} else {
			return node
		}
	}
}

func (p *Parser) unary() *Node {
	if p.consume("+") {
		return p.unary()
	} else if p.consume("-") {
		return NewBinary(ND_SUB, NewNum(0), p.unary())
	} else {
		return p.primary()
	}
}

func (p *Parser) primary() *Node {
	if p.consume("(") {
		node := p.expr()
		p.expect(")")
		return node
	}

	if token := p.consumeIdent(); token != nil {
		v := p.findVariable(token)
		if v == nil {
			v = p.newLvar(token.str)
		}
		return NewVarNode(v)
	}

	return NewNum(p.expectNumber())
}

func (p *Parser) findVariable(token *Token) *Variable {
	for v := p.locals; v != nil; v = v.next {
		if token.str == v.name {
			return v
		}
	}
	return nil
}
