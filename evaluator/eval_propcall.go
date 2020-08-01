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

	var chainArg object.PanObject = object.BuiltInNil
	if node.Chain.Arg != nil {
		chainArg = Eval(node.Chain.Arg, env)
	}
	if err, ok := chainArg.(*object.PanErr); ok {
		return appendStackTrace(err, node.Source())
	}

	args, kwargs, err := evalCallArgs(node, env)
	if err != nil {
		return appendStackTrace(err, node.Source())
	}

	switch node.Chain.Main {
	case ast.Scalar:
		return evalScalarPropCall(node, env, recv, args, kwargs)
	case ast.List:
		return evalListPropCall(node, env, recv, args, kwargs)
	case ast.Reduce:
		return evalReducePropCall(node, env, recv, chainArg, args, kwargs)
	default:
		return nil
	}
}

func evalListPropCall(
	node *ast.PropCallExpr,
	env *object.Env,
	recv object.PanObject,
	args []object.PanObject,
	kwargs *object.PanObj,
) object.PanObject {
	iter, err := iterOf(env, recv)
	if err != nil {
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
		evaluatedElem := evalScalarPropCall(node, env, nextRet, args, kwargs)

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

func evalReducePropCall(
	node *ast.PropCallExpr,
	env *object.Env,
	recv object.PanObject,
	chainArg object.PanObject,
	args []object.PanObject,
	kwargs *object.PanObj,
) object.PanObject {
	iter, err := iterOf(env, recv)
	if err != nil {
		return appendStackTrace(err, node.Source())
	}

	// call `next` prop until StopIterErr raises
	acc := chainArg
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

		prop := evalProp(node, nextRet)
		if err, ok := prop.(*object.PanErr); ok {
			return appendStackTrace(err, node.Source())
		}

		// prepend iteration value to args
		argsToPass := append([]object.PanObject{nextRet}, args...)
		// reduce each iteration values
		ret := evalCall(env, acc, prop, argsToPass, kwargs)
		if err, ok := ret.(*object.PanErr); ok {
			return appendStackTrace(err, node.Source())
		}
		acc = ret
	}

	return acc
}

func iterOf(
	env *object.Env,
	obj object.PanObject,
) (object.PanObject, *object.PanErr) {
	iterSym := &object.PanStr{Value: "_iter"}
	iter := builtInCallProp(env, object.EmptyPanObjPtr(),
		object.EmptyPanObjPtr(), obj, iterSym)

	if err, ok := iter.(*object.PanErr); ok {
		return nil, err
	}

	if iter == object.BuiltInNil {
		err := object.NewTypeErr("recv must have prop `_iter`")
		return nil, err
	}

	return iter, nil
}

func evalScalarPropCall(
	node *ast.PropCallExpr,
	env *object.Env,
	recv object.PanObject,
	args []object.PanObject,
	kwargs *object.PanObj,
) object.PanObject {
	prop := evalProp(node, recv)
	if err, ok := prop.(*object.PanErr); ok {
		return appendStackTrace(err, node.Source())
	}

	ret := evalCall(env, recv, prop, args, kwargs)
	if err, ok := ret.(*object.PanErr); ok {
		return appendStackTrace(err, node.Source())
	}

	return ret
}

func evalProp(node *ast.PropCallExpr, recv object.PanObject) object.PanObject {
	propStr := node.Prop.Value
	propHash := object.GetSymHash(propStr)

	prop, ok := callProp(recv, propHash)

	if !ok {
		err := object.NewNoPropErr(
			fmt.Sprintf("property `%s` is not defined.", propStr))
		return appendStackTrace(err, node.Source())
	}

	return prop
}

func evalCallArgs(
	node *ast.PropCallExpr,
	env *object.Env,
) ([]object.PanObject, *object.PanObj, *object.PanErr) {
	args, err := evalArgs(node.Args, env)
	if err != nil {
		return nil, nil, err
	}

	kwargs, err := evalKwargs(node.Kwargs, env)
	if err != nil {
		return nil, nil, err
	}

	return args, kwargs, nil
}

func evalCall(
	env *object.Env,
	recv object.PanObject,
	prop object.PanObject,
	args []object.PanObject,
	kwargs *object.PanObj,
) object.PanObject {
	// TODO: handle proto of func
	switch prop := prop.(type) {
	case *object.PanBuiltIn:
		return evalBuiltInFuncMethodCall(env, recv, prop, args, kwargs)
	case *object.PanFunc:
		return evalFuncMethodCall(env, recv, prop, args, kwargs)
	default:
		return prop
	}
}

func evalBuiltInFuncMethodCall(
	env *object.Env,
	recv object.PanObject,
	f *object.PanBuiltIn,
	args []object.PanObject,
	kwargs *object.PanObj,
) object.PanObject {
	// prepend recv to args
	args = append([]object.PanObject{recv}, args...)
	return f.Fn(env, kwargs, args...)
}

func evalFuncMethodCall(
	env *object.Env,
	recv object.PanObject,
	f *object.PanFunc,
	args []object.PanObject,
	kwargs *object.PanObj,
) object.PanObject {
	// prepend recv to args
	args = append([]object.PanObject{f, recv}, args...)
	return evalFuncCall(env, kwargs, args...)
}
