package evaluator

import (
	"../ast"
	"../object"
)

func evalInt(node *ast.IntLiteral, env *object.Env) object.PanObject {
	return object.NewPanInt(node.Value)
}
