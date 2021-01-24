package evaluator

import (
	"fmt"
	"github.com/Syuparn/pangaea/ast"
	"github.com/Syuparn/pangaea/object"
)

func unpackArrExpansion(
	node ast.Node,
	env *object.Env,
) ([]object.PanObject, *object.PanErr, bool) {
	// unpack elem if it is arr expansion (like `*a`)
	pref, ok := node.(*ast.PrefixExpr)

	if !ok {
		return nil, nil, false
	}

	if pref.Operator != "*" {
		return nil, nil, false
	}

	evaluated := Eval(pref.Right, env)

	if err, ok := evaluated.(*object.PanErr); ok {
		return nil, appendStackTrace(err, pref.Source()), true
	}

	arr, ok := object.TraceProtoOfArr(evaluated)

	if !ok {
		err := object.NewTypeErr(
			fmt.Sprintf("cannot use `*` unpacking for `%s`", evaluated.Inspect()))
		return nil, appendStackTrace(err, pref.Source()), true
	}

	return arr.Elems, nil, true
}
