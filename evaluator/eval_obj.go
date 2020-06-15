package evaluator

import (
	"../ast"
	"../object"
)

func evalObj(node *ast.ObjLiteral, env *object.Env) object.PanObject {
	pairMap := map[object.SymHash]object.Pair{}
	for _, pairNode := range node.Pairs {
		pair := evalPair(pairNode, env)

		// TODO: raise error if key is not PanStr
		panStr, _ := pair.Key.(*object.PanStr)
		symHash := object.GetSymHash(panStr.Value)

		// NOTE: ignore duplicated keys (`{a: 1, a: 2}` is same as `{a: 1}`)
		if _, exists := pairMap[symHash]; !exists {
			pairMap[symHash] = pair
		}
	}

	// unpack objExpansion elements (like `**a`)
	for _, expElem := range node.EmbeddedExprs {
		evaluated := Eval(expElem, env)
		obj, ok := evaluated.(*object.PanObj)

		if !ok {
			// TODO: error handling if evaluated is not *PanObj
			return nil
		}

		for symHash, pair := range *obj.Pairs {
			pairMap[symHash] = pair
		}
	}

	return object.PanObjInstancePtr(&pairMap)
}
