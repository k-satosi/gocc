package main

type TypeKind int

type Type interface {
	size() int
}

type CharType struct{}

func NewCharType() *CharType {
	return &CharType{}
}

func (c *CharType) size() int {
	return 1
}

type IntType struct{}

func NewIntType() *IntType {
	return &IntType{}
}

func (i *IntType) size() int {
	return 8
}

type PointerType struct {
	Type
	base Type
}

func NewPointerType(base Type) *PointerType {
	return &PointerType{
		base: base,
	}
}

func (p *PointerType) size() int {
	return 8
}

var charType Type = NewCharType()
var intType Type = NewIntType()

type ArrayType struct {
	Type
	base Type
	len  int
}

func NewArrayType(base Type, len int) *ArrayType {
	return &ArrayType{
		base: base,
		len:  len,
	}
}

func (a *ArrayType) size() int {
	return a.base.size() * a.len
}
