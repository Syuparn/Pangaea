package evaluator

import (
	"fmt"

	"github.com/Syuparn/pangaea/ast"
	"github.com/Syuparn/pangaea/object"
)

func evalPropCall(node *ast.PropCallExpr, env *object.Env) object.PanObject {
	recv, err := extractRecv(node.Receiver, env)
	if err != nil {
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

	ignoresNil := node.Chain.Additional == ast.Lonely
	recoversNil := node.Chain.Additional == ast.Thoughtful
	squashesNil := node.Chain.Additional != ast.Strict

	switch node.Chain.Main {
	case ast.Scalar:
		return evalScalarPropCall(
			node, env, recv, args, kwargs, ignoresNil, recoversNil)
	case ast.List:
		return evalListPropCall(
			node, env, recv, args, kwargs, ignoresNil, recoversNil, squashesNil)
	case ast.Reduce:
		return evalReducePropCall(
			node, env, recv, chainArg, args, kwargs, ignoresNil, recoversNil)
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
				continue
			}
			return appendStackTrace(err, node.Source())
		}

		// treat next value as recv
		evaluatedElem := evalScalarPropCall(
			node, env, nextRet, args, kwargs, ignoresNil, recoversNil)

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

func evalReducePropCall(
	node *ast.PropCallExpr,
	env *object.Env,
	recv object.PanObject,
	chainArg object.PanObject,
	args []object.PanObject,
	kwargs *object.PanObj,
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

		prop, isMissing := evalProp(node.Prop.Value, acc)
		if err, ok := prop.(*object.PanErr); ok {
			// thoughtful chain
			if recoversNil {
				continue
			}
			return appendStackTrace(err, node.Source())
		}

		// prepend iteration value to args
		argsToPass := append([]object.PanObject{nextRet}, args...)
		// prepend prop name to arg if _missing is called
		if isMissing {
			propName := object.NewPanStr(node.Prop.Value)
			argsToPass = append([]object.PanObject{propName}, argsToPass...)
		}
		// reduce each iteration values
		ret := evalCall(env, acc, prop, argsToPass, kwargs)

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

func shouldRecover(ret object.PanObject) bool {
	switch ret.Type() {
	case "ErrType":
		return true
	case "NilType":
		return true
	default:
		return false
	}
}

func iterOf(
	env *object.Env,
	obj object.PanObject,
) (object.PanObject, *object.PanErr) {
	iterSym := object.NewPanStr("_iter")
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
	ignoresNil bool,
	recoversNil bool,
) object.PanObject {
	// lonely chain
	if recv == object.BuiltInNil && ignoresNil {
		return recv
	}

	prop, isMissing := evalProp(node.Prop.Value, recv)
	if err, ok := prop.(*object.PanErr); ok {
		if recoversNil {
			return recv
		}
		return appendStackTrace(err, node.Source())
	}

	// prepend prop name to arg if _missing is called
	if isMissing {
		propName := object.NewPanStr(node.Prop.Value)
		args = append([]object.PanObject{propName}, args...)
	}

	ret := evalCall(env, recv, prop, args, kwargs)

	// thoughtful chain
	if recoversNil && shouldRecover(ret) {
		return recv
	}

	if err, ok := ret.(*object.PanErr); ok {
		return appendStackTrace(err, node.Source())
	}

	return ret
}

func evalProp(
	propStr string,
	recv object.PanObject,
) (o object.PanObject, isMissing bool) {
	propHash := object.GetSymHash(propStr)

	prop, ok := object.FindPropAlongProtos(recv, propHash)

	if ok {
		return prop, false
	}

	// try to find _missing instead
	missingHash := object.GetSymHash("_missing")
	missing, ok := object.FindPropAlongProtos(recv, missingHash)

	if !ok {
		err := object.NewNoPropErr(
			fmt.Sprintf("property `%s` is not defined.", propStr))
		return err, false
	}

	return missing, true
}

func evalCallArgs(
	node *ast.PropCallExpr,
	env *object.Env,
) ([]object.PanObject, *object.PanObj, *object.PanErr) {
	// NOTE: for syntactic reason, kwarg expansion is in Args as `**` prefixExpr
	// (not in Kwargs)
	args, unpackedKwargs, err := evalArgs(node.Args, env)
	if err != nil {
		return nil, nil, err
	}

	kwargs, err := evalKwargs(node.Kwargs, env)
	if err != nil {
		return nil, nil, err
	}
	kwargs.AddPairs(unpackedKwargs)

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

func extractRecv(
	recvNode ast.Node,
	env *object.Env,
) (object.PanObject, *object.PanErr) {
	// handle anonchain
	if recvNode == nil {
		self, err := extractAnonChainRecv(env)
		return self, err
	}

	recv := Eval(recvNode, env)
	if err, ok := recv.(*object.PanErr); ok {
		return nil, err
	}
	return recv, nil
}

func extractAnonChainRecv(env *object.Env) (object.PanObject, *object.PanErr) {
	// recv is 1st arg in current env
	self, ok := env.Get(object.GetSymHash(`\1`))
	if !ok {
		return nil, object.NewNameErr("name `\\1` is not defined")
	}
	return self, nil
}
