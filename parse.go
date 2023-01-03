package main

type Function struct {
	next   *Function
	name   string
	params *VarList

	node      *Node
	locals    *VarList
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
	ND_IF
	ND_WHILE
	ND_FOR
	ND_BLOCK
	ND_FUNCALL
	ND_EXPR_STMT
	ND_VAR
	ND_NUM
)

type Node struct {
	kind NodeKind
	next *Node
	lhs  *Node
	rhs  *Node

	cond *Node
	then *Node
	els  *Node
	init *Node
	inc  *Node

	body *Node

	funcname string
	args     *Node

	variable *Variable
	val      int
}

type VarList struct {
	next     *VarList
	variable *Variable
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

func (p *Parser) NewLVar(name string) *Variable {
	v := &Variable{
		name: name,
	}

	vl := &VarList{
		variable: v,
		next:     p.locals,
	}
	p.locals = vl
	return v
}

type Parser struct {
	token  *Token
	locals *VarList
}

func NewParser(token *Token) *Parser {
	return &Parser{
		token: token,
	}
}

func (p *Parser) Program() *Function {
	head := Function{}
	cur := &head

	for !p.token.AtEOF() {
		cur.next = p.function()
		cur = cur.next
	}
	return head.next
}

func (p *Parser) readFuncParams() *VarList {
	if p.consume(")") {
		return nil
	}

	head := &VarList{
		variable: p.NewLVar(p.expectIdent()),
	}
	cur := head

	for !p.consume(")") {
		p.expect(",")
		cur.next = &VarList{}
		cur.next.variable = p.NewLVar(p.expectIdent())
		cur = cur.next
	}

	return head
}

func (p *Parser) function() *Function {
	p.locals = nil

	fn := &Function{}
	fn.name = p.expectIdent()
	p.expect("(")
	fn.params = p.readFuncParams()
	p.expect("{")

	head := &Node{}
	cur := head
	for !p.consume("}") {
		cur.next = p.stmt()
		cur = cur.next
	}

	fn.node = head.next
	fn.locals = p.locals
	return fn
}

func (p *Parser) readExprStmt() *Node {
	return NewUnary(ND_EXPR_STMT, p.expr())
}

func (p *Parser) stmt() *Node {
	if p.consume("return") {
		node := NewUnary(ND_RETURN, p.expr())
		p.expect(";")
		return node
	}

	if p.consume("if") {
		node := NewNode(ND_IF)
		p.expect("(")
		node.cond = p.expr()
		p.expect(")")
		node.then = p.stmt()
		if p.consume("else") {
			node.els = p.stmt()
		}
		return node
	}

	if p.consume("while") {
		node := NewNode(ND_WHILE)
		p.expect("(")
		node.cond = p.expr()
		p.expect(")")
		node.then = p.stmt()
		return node
	}

	if p.consume("for") {
		node := NewNode(ND_FOR)
		p.expect("(")
		if !p.consume(";") {
			node.init = p.readExprStmt()
			p.expect(";")
		}
		if !p.consume(";") {
			node.cond = p.expr()
			p.expect(";")
		}
		if !p.consume(")") {
			node.inc = p.readExprStmt()
			p.expect(")")
		}
		node.then = p.stmt()
		return node
	}

	if p.consume("{") {
		head := Node{}
		cur := &head

		for !p.consume("}") {
			cur.next = p.stmt()
			cur = cur.next
		}

		node := NewNode(ND_BLOCK)
		node.body = head.next
		return node
	}

	node := p.readExprStmt()
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

func (p *Parser) funcArgs() *Node {
	if p.consume(")") {
		return nil
	}

	head := p.assign()
	cur := head
	for p.consume(",") {
		cur.next = p.assign()
		cur = cur.next
	}
	p.expect(")")
	return head
}

func (p *Parser) primary() *Node {
	if p.consume("(") {
		node := p.expr()
		p.expect(")")
		return node
	}

	if token := p.consumeIdent(); token != nil {
		if p.consume("(") {
			node := NewNode(ND_FUNCALL)
			node.funcname = token.str
			node.args = p.funcArgs()
			return node
		}
		v := p.findVariable(token)
		if v == nil {
			v = p.NewLVar(token.str)
		}
		return NewVarNode(v)
	}

	return NewNum(p.expectNumber())
}

func (p *Parser) findVariable(token *Token) *Variable {
	for vl := p.locals; vl != nil; vl = vl.next {
		v := vl.variable
		if token.str == v.name {
			return v
		}
	}
	return nil
}
