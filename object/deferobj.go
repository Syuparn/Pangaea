package object

import (
	"../ast"
)

// DeferType is a type of DeferObj.
const DeferType = "DeferType"

// DeferObj is a wrapper for deferred expr.
type DeferObj struct {
	Node ast.Expr
}

// Type returns type of this PanObject.
func (o *DeferObj) Type() PanObjType {
	return DeferType
}

// Inspect returns formatted source code of this object.
func (o *DeferObj) Inspect() string {
	return "deferObj"
}

// Proto returns proto of this object.
func (o *DeferObj) Proto() PanObject {
	// never called
	panic("deferObj does not have proto")
}
