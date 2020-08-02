package object

import (
	"../ast"
)

const DEFER_TYPE = "DEFER_TYPE"

// wrapper for defer expr
type DeferObj struct {
	Node ast.Expr
}

func (o *DeferObj) Type() PanObjType {
	return DEFER_TYPE
}

func (o *DeferObj) Inspect() string {
	return "deferObj"
}

func (o *DeferObj) Proto() PanObject {
	// never called
	panic("deferObj does not have proto")
}
