package evaluator

import (
	"github.com/Syuparn/pangaea/ast"
	"github.com/Syuparn/pangaea/object"
)

func evalRange(node *ast.RangeLiteral, env *object.Env) object.PanObject {
	// NOTE: *ast.RangeLiteral has nil (not NilLiteral) if nothing is set
	return object.NewPanRange(
		evalOrNil(node.Start, env),
		evalOrNil(node.Stop, env),
		evalOrNil(node.Step, env),
	)
}

func evalOrNil(node ast.Node, env *object.Env) object.PanObject {
	if node == nil {
		return object.BuiltInNil
	}
	return Eval(node, env)
}
