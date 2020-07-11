package evaluator

import (
	"../ast"
	"../object"
)

type FuncWrapperImpl struct {
	codeStr string
	args    *object.PanArr
	kwargs  *object.PanObj
	body    *[]ast.Stmt
}

func (fw *FuncWrapperImpl) String() string {
	return fw.codeStr
}

func (fw *FuncWrapperImpl) Args() *object.PanArr {
	return fw.args
}

func (fw *FuncWrapperImpl) Kwargs() *object.PanObj {
	return fw.kwargs
}

func (fw *FuncWrapperImpl) Body() *[]ast.Stmt {
	return fw.body
}
