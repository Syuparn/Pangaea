package evaluator

import (
	"fmt"
	"github.com/Syuparn/pangaea/ast"
	"github.com/Syuparn/pangaea/object"
)

func evalObj(node *ast.ObjLiteral, env *object.Env) object.PanObject {
	pairMap := map[object.SymHash]object.Pair{}
	for _, pairNode := range node.Pairs {
		pair, err := evalObjPair(pairNode, env)

		if err != nil {
			return appendStackTrace(err, node.Source())
		}

		// NOTE: key must be str. if not, str proto is used instead
		panStr, ok := object.TraceProtoOfStr(pair.Key)

		if !ok {
			err := object.NewTypeErr(
				fmt.Sprintf("cannot use `%s` as Obj key.", pair.Key.Inspect()))
			return appendStackTrace(err, node.Source())
		}

		// replace key with its str proto if it is not str
		pair.Key = panStr
		symHash := object.GetSymHash(panStr.Value)

		// NOTE: ignore duplicated keys (`{a: 1, a: 2}` is same as `{a: 1}`)
		if _, exists := pairMap[symHash]; !exists {
			pairMap[symHash] = pair
		}
	}

	// unpack objExpansion elements (like `**a`)
	for _, expElem := range node.EmbeddedExprs {
		evaluated := Eval(expElem, env)

		if err, ok := evaluated.(*object.PanErr); ok {
			return appendStackTrace(err, expElem.Source())
		}

		obj, ok := evaluated.(*object.PanObj)

		if !ok {
			e := object.NewTypeErr(
				fmt.Sprintf("cannot use `**` unpacking for `%s`",
					evaluated.Inspect()))
			return appendStackTrace(e, expElem.Source())
		}

		for symHash, pair := range *obj.Pairs {
			pairMap[symHash] = pair
		}
	}

	return object.PanObjInstancePtr(&pairMap)
}
