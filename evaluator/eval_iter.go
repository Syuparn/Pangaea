package evaluator

import (
	"../ast"
	"../object"
)

func evalIter(node *ast.IterLiteral, env *object.Env) object.PanObject {
	return evalCallable(node.FuncComponent, env, object.IterFunc)
}
