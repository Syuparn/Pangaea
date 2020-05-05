package object

import (
	"../ast"
	"bytes"
)

const FUNC_TYPE = "FUNC_TYPE"

type PanFunc struct {
	FuncWrapper
}

func (f *PanFunc) Type() PanObjType {
	return ""
}

func (f *PanFunc) Inspect() string {
	var out bytes.Buffer

	out.WriteString("{")
	// delegate to FuncWrapper
	out.WriteString(f.FuncWrapper.String())
	out.WriteString("}")

	return out.String()
}

func (f *PanFunc) Proto() PanObject {
	return f
}

// NOTE: keep loose coupling to ast.FuncComponent and PanFunc
type FuncWrapper interface {
	String() string
}

// implement FuncWrapper (FuncComponent has String() method)
type AstFuncWrapper struct {
	ast.FuncComponent
}
