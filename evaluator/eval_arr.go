package evaluator

import (
	"github.com/Syuparn/pangaea/ast"
	"github.com/Syuparn/pangaea/object"
)

func evalArr(node *ast.ArrLiteral, env *object.Env) object.PanObject {
	elems := []object.PanObject{}
	for _, elemNode := range node.Elems {
		// ok if elem is expansion like `[*a]`
		unpackedElems, err, ok := unpackArrExpansion(elemNode, env)
		if ok {
			if err != nil {
				return appendStackTrace(err, elemNode.Source())
			}

			elems = append(elems, unpackedElems...)
		} else {
			elem := Eval(elemNode, env)

			if err, ok := elem.(*object.PanErr); ok {
				return appendStackTrace(err, elemNode.Source())
			}

			elems = append(elems, elem)
		}
	}
	return &object.PanArr{Elems: elems}
}
