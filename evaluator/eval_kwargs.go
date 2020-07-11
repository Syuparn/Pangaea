package evaluator

import (
	"../ast"
	"../object"
)

func evalKwargs(kwargs map[*ast.Ident]ast.Expr, env *object.Env) *object.PanObj {
	pairMap := map[object.SymHash]object.Pair{}

	for k, v := range kwargs {
		val := Eval(v, env)

		paramName := k.String()
		param := &object.PanStr{Value: paramName}
		symHash := object.GetSymHash(paramName)

		// NOTE: ignore duplicated params (`|a: 1, a: 2|` is same as `|a: 1|`)
		if _, exists := pairMap[symHash]; !exists {
			pairMap[symHash] = object.Pair{Key: param, Value: val}
		}
	}

	obj, _ := (object.PanObjInstancePtr(&pairMap)).(*object.PanObj)
	return obj
}
