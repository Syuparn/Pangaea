package object

import (
	"../ast"
	"bytes"
)

const FUNC_TYPE = "FUNC_TYPE"

type PanFunc struct {
	FuncWrapper
	FuncType FuncType
}

func (f *PanFunc) Type() PanObjType {
	return FUNC_TYPE
}

func (f *PanFunc) Inspect() string {
	var out bytes.Buffer
	out.WriteString(openParen(f.FuncType))
	// delegate to FuncWrapper
	out.WriteString(f.FuncWrapper.String())
	out.WriteString(closeParen(f.FuncType))

	return out.String()
}

func (f *PanFunc) Proto() PanObject {
	return builtInFuncObj
}

type FuncType int

const (
	FUNC_FUNC FuncType = iota
	ITER_FUNC
)

func openParen(t FuncType) string {
	if t == FUNC_FUNC {
		return "{"
	}
	return "<{"
}

func closeParen(t FuncType) string {
	if t == FUNC_FUNC {
		return "}"
	}
	return "}>"
}

// NOTE: keep loose coupling to ast.FuncComponent and PanFunc
type FuncWrapper interface {
	String() string
}

// implement FuncWrapper (FuncComponent has String() method)
type AstFuncWrapper struct {
	ast.FuncComponent
}
