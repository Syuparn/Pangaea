package evaluator

import (
	"../ast"
	"../object"
)

func evalArr(node *ast.ArrLiteral, env *object.Env) object.PanObject {
	elems := []object.PanObject{}
	for _, elemNode := range node.Elems {
		// ok if elem is expansion like `[*a]`
		unpackedElems, ok := unpackArrExpansion(elemNode, env)
		if ok {
			elems = append(elems, unpackedElems...)
		} else {
			elems = append(elems, Eval(elemNode, env))
		}
	}
	return &object.PanArr{Elems: elems}
}
