package evaluator

import (
	"../ast"
	"../object"
)

func evalPair(node *ast.Pair, env *object.Env) object.Pair {
	v := Eval(node.Val, env)

	if ident, ok := node.Key.(*ast.Ident); ok {
		// HACK: syntax sugar (`{a: 1}` is same as `{'a: 1}`)
		k := &object.PanStr{Value: ident.String()}
		return object.Pair{Key: k, Value: v}
	}

	k := Eval(node.Key, env)
	return object.Pair{Key: k, Value: v}
}
