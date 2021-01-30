package evaluator

import (
	"github.com/Syuparn/pangaea/ast"
	"github.com/Syuparn/pangaea/object"
)

func evalAssign(node *ast.AssignExpr, env *object.Env) object.PanObject {
	val := Eval(node.Right, env)

	if err, ok := val.(*object.PanErr); ok {
		return appendStackTrace(err, node.Source())
	}

	symHash := object.GetSymHash(node.Left.Value)
	env.Set(symHash, val)
	return val
}
