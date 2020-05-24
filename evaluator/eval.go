package evaluator

import (
	"../ast"
	"../object"
)

func Eval(node ast.Node, env *object.Env) object.PanObject {
	return object.BuiltInObjObj
}
