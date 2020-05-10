package object

import (
	"bytes"
)

const FUNC_TYPE = "FUNC_TYPE"

type PanFunc struct {
	FuncWrapper
	FuncType FuncType
	Env      *Env
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
	return BuiltInFuncObj
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
// ast.FuncComponent implements FuncWrapper
type FuncWrapper interface {
	String() string
}
