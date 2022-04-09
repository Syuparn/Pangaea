package object

import (
	"bytes"

	"github.com/Syuparn/pangaea/ast"
)

// FuncType is a type of PanFunc.
const FuncType = "FuncType"

// PanFunc is object of func literal.
type PanFunc struct {
	FuncWrapper
	FuncKind FuncKind
	Env      *Env
	// TODO: add proto field to fix inheritance
}

// Type returns type of this PanObject.
func (f *PanFunc) Type() PanObjType {
	return FuncType
}

// Inspect returns formatted source code of this object.
func (f *PanFunc) Inspect() string {
	var out bytes.Buffer
	out.WriteString(openParen(f.FuncKind))
	// delegate to FuncWrapper
	out.WriteString(f.FuncWrapper.String())
	out.WriteString(closeParen(f.FuncKind))

	return out.String()
}

// Repr returns pritty-printed string of this object.
func (f *PanFunc) Repr() string {
	return f.Inspect()
}

// Proto returns proto of this object.
func (f *PanFunc) Proto() PanObject {
	if f.FuncKind == IterFunc {
		return BuiltInIterObj
	}
	return BuiltInFuncObj
}

// Zero returns zero value of this object.
func (f *PanFunc) Zero() PanObject {
	// TODO: implement zero value
	return f
}

// FuncKind is a type of func-like objects.
// NOTE: The type is used to designate func and iter because their implementation is
// same type.
type FuncKind int

const (
	// FuncFunc shows PanFunc is func literal.
	FuncFunc FuncKind = iota
	// IterFunc shows PanFunc is iter literal.
	IterFunc
)

func openParen(t FuncKind) string {
	if t == FuncFunc {
		return "{"
	}
	return "<{"
}

func closeParen(t FuncKind) string {
	if t == FuncFunc {
		return "}"
	}
	return "}>"
}

// FuncWrapper is a wrapper for func literal ast node.
// NOTE: use interface to keep loose coupling to ast package
type FuncWrapper interface {
	String() string
	Args() *PanArr
	Kwargs() *PanObj
	Body() *[]ast.Stmt
}

// NewPanFunc returns new func object.
func NewPanFunc(f FuncWrapper, env *Env) *PanFunc {
	return &PanFunc{f, FuncFunc, env}
}

// NewPanIter returns new func object.
func NewPanIter(f FuncWrapper, env *Env) *PanFunc {
	return &PanFunc{f, IterFunc, env}
}
