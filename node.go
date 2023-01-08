package main

type Node interface {
	Gen()
	AddType()
	Type() *Type
}

type Unary interface {
	Node
}

type BinaryNode interface {
	Node
	Lhs() Node
	Rhs() Node
	ApplyOperation()
}

type Binary struct {
	BinaryNode
	lhs Node
	rhs Node
	ty  *Type
}

func (b *Binary) Lhs() Node {
	return b.lhs
}

func (b *Binary) Rhs() Node {
	return b.rhs
}

func (b *Binary) AddType() {
	b.lhs.AddType()
	b.rhs.AddType()
	b.ty = intType
}

// func (b *Binary) ApplyOperation() {}

type Add struct {
	*Binary
}

func NewAdd(lhs Node, rhs Node) *Add {
	return &Add{
		&Binary{
			lhs: lhs,
			rhs: rhs,
		},
	}
}

type PointerAdd struct {
	*Binary
}

func NewPointerAdd(lhs Node, rhs Node) *PointerAdd {
	return &PointerAdd{
		&Binary{
			lhs: lhs,
			rhs: rhs,
		},
	}
}

type Sub struct {
	*Binary
}

func NewSub(lhs Node, rhs Node) *Sub {
	return &Sub{
		&Binary{
			lhs: lhs,
			rhs: rhs,
		},
	}
}

type PointerSub struct {
	*Binary
}

func NewPointterSub(lhs Node, rhs Node) *PointerSub {
	return &PointerSub{
		&Binary{
			lhs: lhs,
			rhs: rhs,
		},
	}
}

type PointerDiff struct {
	*Binary
}

func NewPointterDiff(lhs Node, rhs Node) *PointerDiff {
	return &PointerDiff{
		&Binary{
			lhs: lhs,
			rhs: rhs,
		},
	}
}

type Mul struct {
	*Binary
}

func NewMul(lhs Node, rhs Node) *Mul {
	return &Mul{
		&Binary{
			lhs: lhs,
			rhs: rhs,
		},
	}
}

type Div struct {
	*Binary
}

func NewDiv(lhs Node, rhs Node) *Div {
	return &Div{
		&Binary{
			lhs: lhs,
			rhs: rhs,
		},
	}
}

type Equal struct {
	*Binary
}

func NewEqual(lhs Node, rhs Node) *Equal {
	return &Equal{
		&Binary{
			lhs: lhs,
			rhs: rhs,
		},
	}
}

type NotEqual struct {
	*Binary
}

func NewNotEqual(lhs Node, rhs Node) *NotEqual {
	return &NotEqual{
		&Binary{
			lhs: lhs,
			rhs: rhs,
		},
	}
}

type LessThan struct {
	*Binary
}

func NewLessThan(lhs Node, rhs Node) *LessThan {
	return &LessThan{
		&Binary{
			lhs: lhs,
			rhs: rhs,
		},
	}
}

type LessEqual struct {
	*Binary
}

func NewLessEqual(lhs Node, rhs Node) *LessEqual {
	return &LessEqual{
		&Binary{
			lhs: lhs,
			rhs: rhs,
		},
	}
}

type Assign struct {
	lhs *VarNode
	rhs Node
	ty  *Type
}

func NewAssign(lhs *VarNode, rhs Node) *Assign {
	return &Assign{
		lhs: lhs,
		rhs: rhs,
	}
}

func (a *Assign) AddType() {
	a.ty = a.lhs.ty
}

func (a *Assign) Type() *Type {
	return a.ty
}

type Address struct {
	Unary
	expr Node
	ty   *Type
}

func NewAddress(expr Node) *Address {
	return &Address{
		expr: expr,
	}
}

func (a *Address) AddType() {
	a.ty = pointerTo(a.expr.Type())
}

type Dereference struct {
	Unary
	expr Node
	ty   *Type
}

func NewDereference(expr Node) *Dereference {
	return &Dereference{
		expr: expr,
	}
}

func (d *Dereference) AddType() {
	if d.expr.Type().kind != TY_PTR {
		errorToken(nil, "invalid pointer dereference")
	} else {
		d.ty = d.expr.Type().base
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

func (r *Return) AddType() {}

type If struct {
	Node
	cond Node
	then Node
	els  Node
}

func NewIf(cond Node, then Node, els Node) *If {
	return &If{
		cond: cond,
		then: then,
		els:  els,
	}
}

func (f *If) AddType() {
	f.cond.AddType()
	f.then.AddType()
	if f.els != nil {
		f.els.AddType()
	}
}

type While struct {
	Node
	cond Node
	then Node
}

func NewWhile(cond Node, then Node) *While {
	return &While{
		cond: cond,
		then: then,
	}
}

func (w *While) AddType() {
	w.cond.AddType()
	w.then.AddType()
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

func (f *For) AddType() {
	if f.init != nil {
		f.init.AddType()
	}
	if f.cond != nil {
		f.cond.AddType()
	}
	if f.inc != nil {
		f.inc.AddType()
	}
	f.block.AddType()
}

func (f *For) Type() *Type {
	return nil
}

type Block struct {
	body []Node
}

func NewBlock(block []Node) *Block {
	return &Block{
		body: block,
	}
}

func (b *Block) AddType() {
	for i := range b.body {
		b.body[i].AddType()
	}
}

func (b *Block) Type() *Type {
	return nil
}

type FuncCall struct {
	name string
	args []Node
	ty   *Type
}

func NewFuncCall(name string, args []Node) *FuncCall {
	return &FuncCall{
		name: name,
		args: args,
	}
}

func (f *FuncCall) AddType() {
	for i := range f.args {
		f.args[i].AddType()
	}
	f.ty = intType
}

func (f *FuncCall) Type() *Type {
	return f.ty
}

type ExpressionStatement struct {
	statement Node
}

func NewExpressionStatement(expr Node) *ExpressionStatement {
	return &ExpressionStatement{
		statement: expr,
	}
}

func (e *ExpressionStatement) AddType() {}

func (e *ExpressionStatement) Type() *Type {
	return intType
}

type Variable struct {
	name   string
	ty     *Type
	offset int
}

type VarNode struct {
	Node
	variable *Variable
	ty       *Type
}

func NewVarNode(v *Variable) *VarNode {
	return &VarNode{
		variable: v,
	}
}

func (v *VarNode) AddType() {
	v.ty = v.variable.ty
}

type Number struct {
	Node
	val int
	ty  *Type
}

func NewNumber(val int) *Number {
	return &Number{
		val: val,
	}
}

func (n *Number) AddType() {
	n.ty = intType
}

type Null struct {
	Node
}

func NewNull() *Null {
	return &Null{}
}
