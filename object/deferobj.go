package object

import (
	"github.com/Syuparn/pangaea/ast"
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

// Repr returns pritty-printed string of this object.
func (o *DeferObj) Repr() string {
	return o.Inspect()
}

// Proto returns proto of this object.
func (o *DeferObj) Proto() PanObject {
	// never called
	panic("deferObj does not have proto")
}

// Zero returns zero value of this object.
func (e *DeferObj) Zero() PanObject {
	// never called
	panic("deferObj does not have zero")
}
