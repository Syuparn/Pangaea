package evaluator

import (
	"fmt"

	"github.com/Syuparn/pangaea/ast"
	"github.com/Syuparn/pangaea/object"
)

func evalMap(node *ast.MapLiteral, env *object.Env) object.PanObject {
	pairs := []object.Pair{}
	// stored for duplicate key check
	// (NOTE: hashable keys duplication is checked in NewPanMap)
	nonHashablePairs := []object.Pair{}

	for _, pairNode := range node.Pairs {
		pair, err := evalMapPair(pairNode, env)

		if err != nil {
			return appendStackTrace(err, node.Source())
		}

		if _, ok := pair.Key.(object.PanScalar); ok {
			pairs = append(pairs, pair)
		} else {
			if !existsNonHashableKey(env, nonHashablePairs, pair) {
				// stored only if key does not exist
				pairs = append(pairs, pair)
				nonHashablePairs = append(nonHashablePairs, pair)
			}
		}
	}

	// unpack objExpansion elements (like `**a`)
	embeddedPairs, panErr := extractEmbeddedElems(node, env, nonHashablePairs)

	if panErr != nil {
		return appendStackTrace(panErr, node.Source())
	}
	pairs = append(pairs, embeddedPairs...)

	return object.NewPanMap(pairs...)
}

func extractEmbeddedElems(
	node *ast.MapLiteral,
	env *object.Env,
	nonHashablePairs []object.Pair,
) ([]object.Pair, *object.PanErr) {
	pairs := []object.Pair{}

	for _, expElem := range node.EmbeddedExprs {
		evaluated := Eval(expElem, env)

		if err, ok := evaluated.(*object.PanErr); ok {
			return nil, appendStackTrace(err, expElem.Source())
		}

		switch e := evaluated.(type) {
		case *object.PanMap:
			for _, pair := range *e.Pairs {
				pairs = append(pairs, pair)
			}
			for _, nPair := range *e.NonHashablePairs {
				if !existsNonHashableKey(env, nonHashablePairs, nPair) {
					// stored only if key does not exist
					pairs = append(pairs, nPair)
					nonHashablePairs = append(nonHashablePairs, nPair)
				}
			}

		case *object.PanObj:
			for _, pair := range *e.Pairs {
				pairs = append(pairs, pair)
			}

		default:
			err := object.NewTypeErr(
				fmt.Sprintf("cannot use `**` unpacking for `%s`",
					evaluated.Inspect()))
			return nil, appendStackTrace(err, expElem.Source())
		}
	}

	return pairs, nil
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
