package main

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
	ND_RETURN
	ND_NUM
)

type Node struct {
	kind NodeKind
	next *Node
	lhs  *Node
	rhs  *Node
	val  int
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

type Parser struct {
	token *Token
}

func NewParser(token *Token) *Parser {
	return &Parser{
		token: token,
	}
}

func (p *Parser) Program() *Node {
	head := Node{}
	cur := &head

	for !p.token.AtEOF() {
		cur.next = p.stmt()
		cur = cur.next
	}

	return head.next
}

func (p *Parser) stmt() *Node {
	if p.consume("return") {
		node := NewUnary(ND_RETURN, p.expr())
		p.expect(";")
		return node
	}

	node := p.expr()
	p.expect(";")
	return node
}

func (p *Parser) expr() *Node {
	return p.equality()
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

	return NewNum(p.expectNumber())
}
