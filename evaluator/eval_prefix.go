package evaluator

import (
	"../ast"
	"../object"
)

func evalPrefix(node *ast.PrefixExpr, env *object.Env) object.PanObject {
	// `*` expansion out of arr is invalid
	if node.Operator == `*` {
		e := object.NewSyntaxErr("cannot use `*` unpacking outside of Arr.")
		return appendStackTrace(e, node.Source())
	}

	// TODO: eval prefix
	return nil
}
