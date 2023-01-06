package main

type Node interface {
}

type Unary interface {
	Node
}

type Binary interface {
	Node
}

type Add struct {
	Binary
	lhs Node
	rhs Node
}

func NewAdd(lhs Node, rhs Node) *Add {
	return &Add{
		lhs: lhs,
		rhs: rhs,
	}
}

type Sub struct {
	Binary
	lhs Node
	rhs Node
}

func NewSub(lhs Node, rhs Node) *Sub {
	return &Sub{
		lhs: lhs,
		rhs: rhs,
	}
}

type Mul struct {
	Binary
	lhs Node
	rhs Node
}

func NewMul(lhs Node, rhs Node) *Mul {
	return &Mul{
		lhs: lhs,
		rhs: rhs,
	}
}

type Div struct {
	Binary
	lhs Node
	rhs Node
}

func NewDiv(lhs Node, rhs Node) *Div {
	return &Div{
		lhs: lhs,
		rhs: rhs,
	}
}

type Equal struct {
	Binary
	lhs Node
	rhs Node
}

func NewEqual(lhs Node, rhs Node) *Equal {
	return &Equal{
		lhs: lhs,
		rhs: rhs,
	}
}

type NotEqual struct {
	Binary
	lhs Node
	rhs Node
}

func NewNotEqual(lhs Node, rhs Node) *NotEqual {
	return &NotEqual{
		lhs: lhs,
		rhs: rhs,
	}
}

type LessThan struct {
	Binary
	lhs Node
	rhs Node
}

func NewLessThan(lhs Node, rhs Node) *LessThan {
	return &LessThan{
		lhs: lhs,
		rhs: rhs,
	}
}

type LessEqual struct {
	Binary
	lhs Node
	rhs Node
}

func NewLessEqual(lhs Node, rhs Node) *LessEqual {
	return &LessEqual{
		lhs: lhs,
		rhs: rhs,
	}
}

type Assign struct {
	Binary
	lhs Node
	rhs Node
}

func NewAssign(lhs Node, rhs Node) *Assign {
	return &Assign{
		lhs: lhs,
		rhs: rhs,
	}
}

type Return struct {
	Unary
	expr Node
}

func NewReturn(expr Node) *Return {
	return &Return{
		expr: expr,
	}
}

type If struct {
	cond Node
	then Node
	els  Node
}

func NewIf(cond Node, then Node, els Node) Node {
	return &If{
		cond: cond,
		then: then,
		els:  els,
	}
}

type While struct {
	cond Node
	then Node
}

func NewWhile(cond Node, then Node) *While {
	return &While{
		cond: cond,
		then: then,
	}
}

type For struct {
	init  Node
	cond  Node
	inc   Node
	block Node
}

func NewFor(init Node, cond Node, inc Node, block Node) *For {
	return &For{
		init:  init,
		cond:  cond,
		inc:   inc,
		block: block,
	}
}

type Block struct {
	body []Node
}

func NewBlock(block []Node) *Block {
	return &Block{
		body: block,
	}
}

type FuncCall struct {
	name string
	args []Node
}

func NewFuncCall(name string, args []Node) *FuncCall {
	return &FuncCall{
		name: name,
		args: args,
	}
}

type ExpressionStatement struct {
	statement Node
}

func NewExpressionStatement(expr Node) *ExpressionStatement {
	return &ExpressionStatement{
		statement: expr,
	}
}

type Variable struct {
	name   string
	offset int
}

type VarNode struct {
	Node
	variable *Variable
}

func NewVarNode(v *Variable) *VarNode {
	return &VarNode{
		variable: v,
	}
}

type Number struct {
	Node
	val int
}

func NewNumber(val int) *Number {
	return &Number{
		val: val,
	}
}
