package evaluator

import (
	"github.com/Syuparn/pangaea/ast"
	"github.com/Syuparn/pangaea/object"
)

func evalInfix(node *ast.InfixExpr, env *object.Env) object.PanObject {
	if isShortCutOperator(node.Operator) {
		return evalShortCutInfix(node, env)
	}

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

func isShortCutOperator(op string) bool {
	switch op {
	case "||":
		return true
	case "&&":
		return true
	}
	return false
}

func evalShortCutInfix(node *ast.InfixExpr, env *object.Env) object.PanObject {
	left := Eval(node.Left, env)
	if err, ok := left.(*object.PanErr); ok {
		return appendStackTrace(err, node.Source())
	}

	if canShortCut(left, node.Operator, env) {
		return left
	}

	return Eval(node.Right, env)
}

func canShortCut(left object.PanObject, op string, env *object.Env) bool {
	// able to shortcut if left is truthy
	bSym := object.NewPanStr("B")
	boolified := builtInCallProp(
		env, object.EmptyPanObjPtr(),
		object.EmptyPanObjPtr(), left, bSym,
	)
	isTruthy := boolified == object.BuiltInTrue

	switch op {
	case "||":
		return isTruthy
	case "&&":
		return !isTruthy
	}
	return false
}
