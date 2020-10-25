package evaluator

import (
	"../ast"
	"../object"
)

func evalLiteralCall(node *ast.LiteralCallExpr, env *object.Env) object.PanObject {
	recv, err := extractRecv(node.Receiver, env)
	if err != nil {
		return appendStackTrace(err, node.Source())
	}

	fObj := Eval(node.Func, env)
	if err, ok := fObj.(*object.PanErr); ok {
		return appendStackTrace(err, node.Source())
	}

	// TODO: duck typing (allow all objs with `call` prop)
	f, ok := object.TraceProtoOfFunc(fObj)
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

	ignoresNil := node.Chain.Additional == ast.Lonely
	recoversNil := node.Chain.Additional == ast.Thoughtful
	squashesNil := node.Chain.Additional != ast.Strict

	switch node.Chain.Main {
	case ast.Scalar:
		return evalScalarLiteralCall(
			node, env, f, recv, ignoresNil, recoversNil)
	case ast.List:
		return evalListLiteralCall(
			node, env, f, recv, ignoresNil, recoversNil, squashesNil)
	case ast.Reduce:
		return evalReduceLiteralCall(
			node, env, f, recv, chainArg, ignoresNil, recoversNil)
	default:
		return nil
	}
}

func evalScalarLiteralCall(
	node ast.Node,
	env *object.Env,
	f *object.PanFunc,
	recv object.PanObject,
	ignoresNil bool,
	recoversNil bool,
) object.PanObject {
	// lonely chain
	if ignoresNil && recv == object.BuiltInNil {
		return recv
	}

	args := literalCallArgs(recv, f)
	// prepend f itself to args
	args = append([]object.PanObject{f}, args...)

	ret := _evalLiteralCall(node, env, f, args)

	// thoughtful chain
	if recoversNil && shouldRecover(ret) {
		return recv
	}

	return ret
}

func evalListLiteralCall(
	node ast.Node,
	env *object.Env,
	f *object.PanFunc,
	recv object.PanObject,
	ignoresNil bool,
	recoversNil bool,
	squashesNil bool,
) object.PanObject {
	iter, err := iterOf(env, recv)
	if err != nil {
		return appendStackTrace(err, node.Source())
	}

	// call `next` prop until StopIterErr raises
	evaluatedElems := []object.PanObject{}
	nextSym := object.NewPanStr("next")
	for {
		// call `(iter).next`
		nextRet := builtInCallProp(env, object.EmptyPanObjPtr(),
			object.EmptyPanObjPtr(), iter, nextSym)

		if isStopIter(nextRet) {
			break
		}

		if err, ok := nextRet.(*object.PanErr); ok {
			// thoughtful chain
			if recoversNil {
				evaluatedElems = append(evaluatedElems, nextRet)
				continue
			}
			return appendStackTrace(err, node.Source())
		}

		args := literalCallArgs(nextRet, f)
		// prepend f itself to args
		args = append([]object.PanObject{f}, args...)

		// treat next value as recv
		evaluatedElem := _evalLiteralCall(node, env, f, args)

		// thoughtful chain
		if recoversNil && shouldRecover(evaluatedElem) {
			evaluatedElems = append(evaluatedElems, nextRet)
		}

		if err, ok := evaluatedElem.(*object.PanErr); ok {
			return appendStackTrace(err, node.Source())
		}

		// ignore nil
		if squashesNil && evaluatedElem == object.BuiltInNil {
			continue
		}

		evaluatedElems = append(evaluatedElems, evaluatedElem)
	}

	return &object.PanArr{Elems: evaluatedElems}
}

func isStopIter(obj object.PanObject) bool {
	err, ok := obj.(*object.PanErr)
	if !ok {
		return false
	}
	return err.ErrType == object.STOP_ITER_ERR
}

func evalReduceLiteralCall(
	node ast.Node,
	env *object.Env,
	f *object.PanFunc,
	recv object.PanObject,
	chainArg object.PanObject,
	ignoresNil bool, // currently not used
	recoversNil bool,
) object.PanObject {
	iter, err := iterOf(env, recv)
	if err != nil {
		return appendStackTrace(err, node.Source())
	}

	// call `next` prop until StopIterErr raises
	acc := chainArg
	nextSym := object.NewPanStr("next")
	for {
		// call `(iter).next`
		nextRet := builtInCallProp(env, object.EmptyPanObjPtr(),
			object.EmptyPanObjPtr(), iter, nextSym)

		if isStopIter(nextRet) {
			break
		}

		if err, ok := nextRet.(*object.PanErr); ok {
			// thoughtful chain
			if recoversNil {
				continue
			}
			return appendStackTrace(err, node.Source())
		}

		// prepend iteration value to args
		args := []object.PanObject{f, acc, nextRet}
		// reduce each iteration values
		ret := _evalLiteralCall(node, env, f, args)

		// thoughtful chain
		if recoversNil && shouldRecover(ret) {
			continue
		}

		if err, ok := ret.(*object.PanErr); ok {
			return appendStackTrace(err, node.Source())
		}
		acc = ret
	}

	return acc
}

func _evalLiteralCall(
	node ast.Node,
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
	// NOTE: extract elems if there are more than 1 params and recv is arr
	if len(f.Args().Elems) > 1 {
		arr, ok := object.TraceProtoOfArr(recv)
		if ok {
			return arr.Elems
		}
	}

	return []object.PanObject{recv}
}
