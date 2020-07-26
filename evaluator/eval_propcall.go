package evaluator

import (
	"../ast"
	"../object"
	"fmt"
)

func evalPropCall(node *ast.PropCallExpr, env *object.Env) object.PanObject {
	recv := Eval(node.Receiver, env)
	if err, ok := recv.(*object.PanErr); ok {
		return appendStackTrace(err, node.Source())
	}

	switch node.Chain.Main {
	case ast.Scalar:
		return _evalPropCall(node, env, recv)
	case ast.List:
		return evalListPropCall(node, env, recv)
	default:
		return nil
	}
}

func evalListPropCall(
	node *ast.PropCallExpr,
	env *object.Env,
	recv object.PanObject,
) object.PanObject {
	// get `(recv)._iter`
	iterSym := &object.PanStr{Value: "_iter"}
	iter := builtInCallProp(env, object.EmptyPanObjPtr(),
		object.EmptyPanObjPtr(), recv, iterSym)

	if err, ok := iter.(*object.PanErr); ok {
		return appendStackTrace(err, node.Source())
	}

	if iter == object.BuiltInNil {
		err := object.NewTypeErr("recv must have prop `_iter`")
		return appendStackTrace(err, node.Source())
	}

	// call `next` prop until StopIterErr raises
	evaluatedElems := []object.PanObject{}
	nextSym := &object.PanStr{Value: "next"}
	for {
		// call `(iter).next`
		nextRet := builtInCallProp(env, object.EmptyPanObjPtr(),
			object.EmptyPanObjPtr(), iter, nextSym)

		if err, ok := nextRet.(*object.PanErr); ok {
			if err.ErrType == object.STOP_ITER_ERR {
				break
			}
			// if err is not StopIterErr, don't catch
			return appendStackTrace(err, node.Source())
		}

		// treat next value as recv
		evaluatedElem := _evalPropCall(node, env, nextRet)

		if err, ok := evaluatedElem.(*object.PanErr); ok {
			return appendStackTrace(err, node.Source())
		}

		// NOTE: nil is ignored
		if evaluatedElem != object.BuiltInNil {
			evaluatedElems = append(evaluatedElems, evaluatedElem)
		}
	}

	return &object.PanArr{Elems: evaluatedElems}
}

func _evalPropCall(
	node *ast.PropCallExpr,
	env *object.Env,
	recv object.PanObject,
) object.PanObject {
	propStr := node.Prop.Value
	propHash := object.GetSymHash(propStr)

	prop, ok := callProp(recv, propHash)

	if !ok {
		err := object.NewNoPropErr(
			fmt.Sprintf("property `%s` is not defined.", propStr))
		return appendStackTrace(err, node.Source())
	}

	switch prop := prop.(type) {
	case *object.PanBuiltIn:
		return evalBuiltInFuncMethodCall(node, recv, prop, env)
	case *object.PanFunc:
		return evalFuncMethodCall(node, recv, prop, env)
	default:
		return prop
	}
}

func evalBuiltInFuncMethodCall(
	node *ast.PropCallExpr,
	recv object.PanObject,
	f *object.PanBuiltIn,
	env *object.Env,
) object.PanObject {
	args, err := evalArgs(node.Args, env)

	if err != nil {
		return appendStackTrace(err, node.Source())
	}

	kwargs, err := evalKwargs(node.Kwargs, env)

	if err != nil {
		return appendStackTrace(err, node.Source())
	}

	// prepend recv to args
	args = append([]object.PanObject{recv}, args...)

	return f.Fn(env, kwargs, args...)
}

func evalFuncMethodCall(
	node *ast.PropCallExpr,
	recv object.PanObject,
	f *object.PanFunc,
	env *object.Env,
) object.PanObject {
	args, err := evalArgs(node.Args, env)

	if err != nil {
		return appendStackTrace(err, node.Source())
	}

	kwargs, err := evalKwargs(node.Kwargs, env)

	if err != nil {
		return appendStackTrace(err, node.Source())
	}

	// prepend recv to args
	args = append([]object.PanObject{f, recv}, args...)

	return evalFuncCall(env, kwargs, args...)
}
