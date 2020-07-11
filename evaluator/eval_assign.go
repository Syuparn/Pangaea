package evaluator

import (
	"../ast"
	"../object"
)

func evalAssign(node *ast.AssignExpr, env *object.Env) object.PanObject {
	val := Eval(node.Right, env)

	symHash := object.GetSymHash(node.Left.Value)
	env.Set(symHash, val)
	return val
}
