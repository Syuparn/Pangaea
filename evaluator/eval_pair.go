package evaluator

import (
	"../ast"
	"../object"
)

func evalObjPair(node *ast.Pair, env *object.Env) (object.Pair, *object.PanErr) {
	v := Eval(node.Val, env)

	if e, ok := v.(*object.PanErr); ok {
		emptyPair := object.Pair{Key: nil, Value: nil}
		return emptyPair, appendStackTrace(e, node.Val.Source())
	}

	if ident, ok := node.Key.(*ast.Ident); ok {
		// syntax sugar (`{a: 1}` is same as `{'a: 1}`)
		k := object.NewPanStr(ident.String())
		return object.Pair{Key: k, Value: v}, nil
	}

	k := Eval(node.Key, env)

	if e, ok := k.(*object.PanErr); ok {
		emptyPair := object.Pair{Key: nil, Value: nil}
		return emptyPair, appendStackTrace(e, node.Key.Source())
	}

	return object.Pair{Key: k, Value: v}, nil
}

func evalMapPair(node *ast.Pair, env *object.Env) (object.Pair, *object.PanErr) {
	k := Eval(node.Key, env)

	if err, ok := k.(*object.PanErr); ok {
		emptyPair := object.Pair{Key: nil, Value: nil}
		return emptyPair, appendStackTrace(err, node.Key.Source())
	}

	v := Eval(node.Val, env)

	if err, ok := v.(*object.PanErr); ok {
		emptyPair := object.Pair{Key: nil, Value: nil}
		return emptyPair, appendStackTrace(err, node.Key.Source())
	}

	return object.Pair{Key: k, Value: v}, nil
}
