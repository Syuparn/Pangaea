package evaluator

import (
	"../ast"
	"../object"
)

// NOTE: this value is used with err
var emptyPair = object.Pair{Key: nil, Value: nil}

func evalObjPair(node *ast.Pair, env *object.Env) (object.Pair, *object.PanErr) {
	v := Eval(node.Val, env)

	if e, ok := v.(*object.PanErr); ok {
		return emptyPair, appendStackTrace(e, node.Val.Source())
	}

	if ident, ok := node.Key.(*ast.Ident); ok {
		// syntax sugar (`{a: 1}` is same as `{'a: 1}`)
		k := object.NewPanStr(ident.String())
		return object.Pair{Key: k, Value: v}, nil
	}

	// pinned ident works same as ordinal ident
	if pinned, ok := node.Key.(*ast.PinnedIdent); ok {
		k, err := searchPinnedKey(pinned, env)

		if err != nil {
			return emptyPair, appendStackTrace(err, node.Val.Source())
		}

		strK, ok := object.TraceProtoOfStr(k)
		if !ok {
			err := object.NewTypeErr("key of obj must be str")
			return emptyPair, appendStackTrace(err, node.Val.Source())
		}

		return object.Pair{Key: strK, Value: v}, nil
	}

	k := Eval(node.Key, env)

	if e, ok := k.(*object.PanErr); ok {
		return emptyPair, appendStackTrace(e, node.Key.Source())
	}

	return object.Pair{Key: k, Value: v}, nil
}

func evalMapPair(node *ast.Pair, env *object.Env) (object.Pair, *object.PanErr) {
	v := Eval(node.Val, env)

	if err, ok := v.(*object.PanErr); ok {
		return emptyPair, appendStackTrace(err, node.Key.Source())
	}

	// pinned ident works same as ordinal ident
	if pinned, ok := node.Key.(*ast.PinnedIdent); ok {
		k, err := searchPinnedKey(pinned, env)

		if err != nil {
			return emptyPair, appendStackTrace(err, node.Val.Source())
		}

		return object.Pair{Key: k, Value: v}, nil
	}

	k := Eval(node.Key, env)

	if err, ok := k.(*object.PanErr); ok {
		return emptyPair, appendStackTrace(err, node.Key.Source())
	}

	return object.Pair{Key: k, Value: v}, nil
}

func searchPinnedKey(
	p *ast.PinnedIdent,
	env *object.Env,
) (object.PanObject, *object.PanErr) {
	k := Eval(&p.Ident, env)
	if err, ok := k.(*object.PanErr); ok {
		return nil, appendStackTrace(err, p.Source())
	}
	return k, nil
}
