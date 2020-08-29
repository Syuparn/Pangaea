package evaluator

import (
	"../ast"
	"../object"
)

func evalInfix(node *ast.InfixExpr, env *object.Env) object.PanObject {
	left := Eval(node.Left, env)
	if err, ok := left.(*object.PanErr); ok {
		return appendStackTrace(err, node.Source())
	}

	right := Eval(node.Right, env)
	if err, ok := right.(*object.PanErr); ok {
		return appendStackTrace(err, node.Source())
	}

	propSym := object.NewPanStr(node.Operator)

	// same as `Obj.callProp(left, propSym, right)`, which is evaluated to
	// `left.^propSym(right)`
	ret := builtInCallProp(env, object.EmptyPanObjPtr(),
		object.EmptyPanObjPtr(), left, propSym, right)

	if err, ok := ret.(*object.PanErr); ok {
		return appendStackTrace(err, node.Source())
	}

	return ret
}
