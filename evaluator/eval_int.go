package evaluator

import (
	"../ast"
	"../object"
)

func evalInt(node *ast.IntLiteral, env *object.Env) object.PanObject {
	// return cached object instead
	if node.Value == 0 {
		return object.BuiltInZeroInt
	}
	if node.Value == 1 {
		return object.BuiltInOneInt
	}
	return &object.PanInt{Value: node.Value}
}
