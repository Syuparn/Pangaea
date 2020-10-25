package evaluator

import (
	"github.com/Syuparn/pangaea/ast"
	"github.com/Syuparn/pangaea/object"
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

	_else := Eval(node.Else, env)
	if err, ok := _else.(*object.PanErr); ok {
		return appendStackTrace(err, node.Source())
	}
	return _else
}
