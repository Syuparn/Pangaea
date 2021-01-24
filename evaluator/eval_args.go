package evaluator

import (
	"fmt"

	"github.com/Syuparn/pangaea/ast"
	"github.com/Syuparn/pangaea/object"
)

func evalArgs(
	argNodes []ast.Expr,
	env *object.Env,
) ([]object.PanObject, *map[object.SymHash]object.Pair, *object.PanErr) {
	args := []object.PanObject{}
	// NOTE: for syntactic reason, kwarg expansion is in Args as `**` prefixExpr
	// (not in Kwargs)
	unpackedKwargs := object.EmptyPanObjPtr()
	for _, argNode := range argNodes {
		// unpack kwarg expansion (like **obj)
		kwargs, err, ok := unpackObjExpansion(argNode, env)
		if ok {
			if err != nil {
				appendStackTrace(err, argNode.Source())
				return []object.PanObject{}, nil, err
			}
			unpackedKwargs.AddPairs(kwargs)
			continue
		}

		// try to unpack arg expansion (like *arr)
		elems, err, ok := unpackArrExpansion(argNode, env)
		if ok {
			if err != nil {
				appendStackTrace(err, argNode.Source())
				return []object.PanObject{}, nil, err
			}
			args = append(args, elems...)
			continue
		}

		arg := Eval(argNode, env)

		if err, ok := arg.(*object.PanErr); ok {
			appendStackTrace(err, argNode.Source())
			return []object.PanObject{}, nil, err
		}

		args = append(args, arg)
	}

	return args, unpackedKwargs.Pairs, nil
}

func unpackObjExpansion(
	node ast.Node,
	env *object.Env,
) (*map[object.SymHash]object.Pair, *object.PanErr, bool) {
	pref, ok := node.(*ast.PrefixExpr)
	if !ok {
		return nil, nil, false
	}
	if pref.Operator != "**" {
		return nil, nil, false
	}

	o := Eval(pref.Right, env)
	if err, ok := o.(*object.PanErr); ok {
		appendStackTrace(err, node.Source())
		return nil, err, true
	}

	obj, ok := o.(*object.PanObj)
	if !ok {
		err := object.NewTypeErr(fmt.Sprintf(
			"cannot use `**` unpacking for `%s`", o.Inspect()))
		appendStackTrace(err, node.Source())
		return nil, err, true
	}

	return obj.Pairs, nil, true
}
