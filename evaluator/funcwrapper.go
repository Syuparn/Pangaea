package evaluator

import (
	"github.com/Syuparn/pangaea/ast"
	"github.com/Syuparn/pangaea/object"
)

// FuncWrapperImpl is wrapper for ast func code.
// This implements object.FuncWrapper.
type FuncWrapperImpl struct {
	codeStr string
	args    *object.PanArr
	kwargs  *object.PanObj
	body    *[]ast.Stmt
}

// String returns func code.
func (fw *FuncWrapperImpl) String() string {
	return fw.codeStr
}

// Args returns potisional params of func.
func (fw *FuncWrapperImpl) Args() *object.PanArr {
	return fw.args
}

// Kwargs returns keyword params of func.
func (fw *FuncWrapperImpl) Kwargs() *object.PanObj {
	return fw.kwargs
}

// Body returns ast of func body.
func (fw *FuncWrapperImpl) Body() *[]ast.Stmt {
	return fw.body
}
