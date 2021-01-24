package evaluator

import (
	"github.com/Syuparn/pangaea/ast"
	"github.com/Syuparn/pangaea/object"
)

func evalPrefix(node *ast.PrefixExpr, env *object.Env) object.PanObject {
	// `*` expansion out of arr is invalid
	if node.Operator == `*` {
		e := object.NewSyntaxErr("cannot use `*` unpacking outside of Arr.")
		return appendStackTrace(e, node.Source())
	}

	right := Eval(node.Right, env)
	if err, ok := right.(*object.PanErr); ok {
		return appendStackTrace(err, node.Source())
	}

	propSym := object.NewPanStr(prefixOpMethodName(node.Operator))

	// same as `Obj.callProp(right, propSym)`, which is evaluated to
	// `right.^propSym`
	ret := builtInCallProp(env, object.EmptyPanObjPtr(),
		object.EmptyPanObjPtr(), right, propSym)

	if err, ok := ret.(*object.PanErr); ok {
		return appendStackTrace(err, node.Source())
	}

	return ret
}

func prefixOpMethodName(op string) string {
	// handle irregular mapping ("+" -> "+%", "-" -> "-%")
	// ("+", "-" are for binary operator)
	switch op {
	case "+":
		return "+%"
	case "-":
		return "-%"
	default:
		return op
	}
}
