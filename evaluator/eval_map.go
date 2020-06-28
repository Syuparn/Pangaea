package evaluator

import (
	"../ast"
	"../object"
)

func evalMap(node *ast.MapLiteral, env *object.Env) object.PanObject {
	pairMap := map[object.HashKey]object.Pair{}
	// non-scalar objects are stored in array instead of map
	nonHashablePairs := []object.Pair{}

	for _, pairNode := range node.Pairs {
		pair := evalMapPair(pairNode, env)

		if s, ok := pair.Key.(object.PanScalar); ok {
			// NOTE: ignore duplicated keys (`%{'a: 1, 'a: 2}` is same as `%{'a: 1}`)
			if _, exists := pairMap[s.Hash()]; !exists {
				pairMap[s.Hash()] = pair
			}
		} else {
			nonHashablePairs = append(nonHashablePairs, pair)
		}
	}

	// unpack objExpansion elements (like `**a`)
	appendEmbeddedElems(node, env, pairMap, nonHashablePairs)

	return &object.PanMap{
		Pairs:            &pairMap,
		NonHashablePairs: &nonHashablePairs,
	}
}

func appendEmbeddedElems(
	node *ast.MapLiteral,
	env *object.Env,
	pairMap map[object.HashKey]object.Pair,
	nonHashablePairs []object.Pair,
) {
	for _, expElem := range node.EmbeddedExprs {
		evaluated := Eval(expElem, env)

		switch e := evaluated.(type) {
		case *object.PanMap:
			for hash, pair := range *e.Pairs {
				pairMap[hash] = pair
			}
			for _, nPair := range *e.NonHashablePairs {
				nonHashablePairs = append(nonHashablePairs, nPair)
			}

		case *object.PanObj:
			for _, pair := range *e.Pairs {
				s, _ := pair.Key.(object.PanScalar)
				pairMap[s.Hash()] = pair
			}

		default:
			// TODO: error handling if evaluated is not *PanObj or *PanMap
			return
		}
	}
}
