package main

type TypeKind int

const (
	TY_INT TypeKind = iota
	TY_PTR
)

type Type struct {
	kind TypeKind
	base *Type
}

var intType *Type = &Type{
	kind: TY_INT,
}

func pointerTo(base *Type) *Type {
	return &Type{
		kind: TY_PTR,
		base: base,
	}
}
