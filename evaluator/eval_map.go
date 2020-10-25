package evaluator

import (
	"github.com/Syuparn/pangaea/ast"
	"github.com/Syuparn/pangaea/object"
	"fmt"
)

func evalMap(node *ast.MapLiteral, env *object.Env) object.PanObject {
	pairMap := map[object.HashKey]object.Pair{}
	// non-scalar objects are stored in array instead of map
	nonHashablePairs := []object.Pair{}

	for _, pairNode := range node.Pairs {
		pair, err := evalMapPair(pairNode, env)

		if err != nil {
			return appendStackTrace(err, node.Source())
		}

		if s, ok := pair.Key.(object.PanScalar); ok {
			// NOTE: ignore duplicated keys (`%{'a: 1, 'a: 2}` is same as `%{'a: 1}`)
			if _, exists := pairMap[s.Hash()]; !exists {
				pairMap[s.Hash()] = pair
			}
		} else {
			if !existsNonHashableKey(env, nonHashablePairs, pair) {
				nonHashablePairs = append(nonHashablePairs, pair)
			}
		}
	}

	// unpack objExpansion elements (like `**a`)
	panErr := appendEmbeddedElems(node, env, pairMap, nonHashablePairs)

	if panErr != nil {
		return appendStackTrace(panErr, node.Source())
	}

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
) *object.PanErr {
	for _, expElem := range node.EmbeddedExprs {
		evaluated := Eval(expElem, env)

		if err, ok := evaluated.(*object.PanErr); ok {
			return appendStackTrace(err, expElem.Source())
		}

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
			err := object.NewTypeErr(
				fmt.Sprintf("cannot use `**` unpacking for `%s`",
					evaluated.Inspect()))
			return appendStackTrace(err, expElem.Source())
		}
	}

	return nil
}

func existsNonHashableKey(
	env *object.Env,
	nonHashablePairs []object.Pair,
	newPair object.Pair,
) bool {
	eqSym := object.NewPanStr("==")

	for _, pair := range nonHashablePairs {
		// use == method for each pair comparison
		ret := builtInCallProp(env, object.EmptyPanObjPtr(),
			object.EmptyPanObjPtr(), newPair.Key, eqSym, pair.Key)
		if ret == object.BuiltInTrue {
			return true
		}
	}
	return false
}
