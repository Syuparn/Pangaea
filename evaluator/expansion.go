package evaluator

import (
	"../ast"
	"../object"
)

func unpackArrExpansion(node ast.Node, env *object.Env) ([]object.PanObject, bool) {
	// unpack elem if it is arr expansion (like `*a`)
	pref, ok := node.(*ast.PrefixExpr)

	if !ok {
		return nil, false
	}

	if pref.Operator != "*" {
		return nil, false
	}

	evaluated := Eval(pref.Right, env)
	arr, ok := evaluated.(*object.PanArr)

	// TODO: error handling if evaluated is not *PanArr
	if !ok {
		return nil, false
	}

	return arr.Elems, true
}
