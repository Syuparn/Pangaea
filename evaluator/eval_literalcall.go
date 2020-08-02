package evaluator

import (
	"../ast"
	"../object"
)

func evalLiteralCall(node *ast.LiteralCallExpr, env *object.Env) object.PanObject {
	recv := Eval(node.Receiver, env)
	if err, ok := recv.(*object.PanErr); ok {
		return appendStackTrace(err, node.Source())
	}

	fObj := Eval(node.Func, env)
	if err, ok := fObj.(*object.PanErr); ok {
		return appendStackTrace(err, node.Source())
	}

	// TODO: handle ancestors of func
	f, ok := fObj.(*object.PanFunc)
	if !ok {
		err := object.NewTypeErr("literal call must be func")
		return appendStackTrace(err, node.Source())
	}

	var chainArg object.PanObject = object.BuiltInNil
	if node.Chain.Arg != nil {
		chainArg = Eval(node.Chain.Arg, env)
	}
	if err, ok := chainArg.(*object.PanErr); ok {
		return appendStackTrace(err, node.Source())
	}

	// TODO: handle args/kwargs

	switch node.Chain.Main {
	case ast.Scalar:
		return evalScalarLiteralCall(node, env, f, recv)
	case ast.List:
		return evalListLiteralCall(node, env, f, recv)
	case ast.Reduce:
		return evalReduceLiteralCall(node, env, f, recv, chainArg)
	default:
		return nil
	}
}

func evalScalarLiteralCall(
	node *ast.LiteralCallExpr,
	env *object.Env,
	f *object.PanFunc,
	recv object.PanObject,
) object.PanObject {
	args := literalCallArgs(recv, f)
	// prepend f itself to args
	args = append([]object.PanObject{f}, args...)
	return _evalLiteralCall(node, env, f, args)
}

func evalListLiteralCall(
	node *ast.LiteralCallExpr,
	env *object.Env,
	f *object.PanFunc,
	recv object.PanObject,
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

		args := literalCallArgs(nextRet, f)
		// prepend f itself to args
		args = append([]object.PanObject{f}, args...)

		// treat next value as recv
		evaluatedElem := _evalLiteralCall(node, env, f, args)

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

func evalReduceLiteralCall(
	node *ast.LiteralCallExpr,
	env *object.Env,
	f *object.PanFunc,
	recv object.PanObject,
	chainArg object.PanObject,
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

		// prepend iteration value to args
		args := []object.PanObject{f, acc, nextRet}
		// reduce each iteration values
		ret := _evalLiteralCall(node, env, f, args)
		if err, ok := ret.(*object.PanErr); ok {
			return appendStackTrace(err, node.Source())
		}
		acc = ret
	}

	return acc
}

func _evalLiteralCall(
	node *ast.LiteralCallExpr,
	env *object.Env,
	f *object.PanFunc,
	args []object.PanObject,
) object.PanObject {
	kwargs := object.EmptyPanObjPtr()
	ret := evalFuncCall(env, kwargs, args...)

	if err, ok := ret.(*object.PanErr); ok {
		return appendStackTrace(err, node.Source())
	}
	return ret
}

func literalCallArgs(recv object.PanObject, f *object.PanFunc) []object.PanObject {
	// NOTE: extract elems
	// TODO: handle ancestors of arr
	if len(f.Args().Elems) > 1 && recv.Type() == object.ARR_TYPE {
		return recv.(*object.PanArr).Elems
	}

	return []object.PanObject{recv}
}
