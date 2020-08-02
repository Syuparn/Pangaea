package evaluator

import (
	"../ast"
	"../object"
)

func evalIf(node *ast.IfExpr, env *object.Env) object.PanObject {
	cond := Eval(node.Cond, env)
	if err, ok := cond.(*object.PanErr); ok {
		return appendStackTrace(err, node.Source())
	}

	if isTruthy(cond, env) {
		then := Eval(node.Then, env)
		if err, ok := then.(*object.PanErr); ok {
			return appendStackTrace(err, node.Source())
		}
		return then
	}

	if node.Else == nil {
		return object.BuiltInNil
	}

	else_ := Eval(node.Else, env)
	if err, ok := else_.(*object.PanErr); ok {
		return appendStackTrace(err, node.Source())
	}
	return else_
}
