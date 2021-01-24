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

	chainMiddleware := newChainMiddleware(*node.Chain)
	ret := _evalPropCall(env, recv, chainArg, node.Prop.Value, args, kwargs,
		chainMiddleware)
	if err, ok := ret.(*object.PanErr); ok {
		return appendStackTrace(err, node.Source())
	}

	return ret
}

func propCallHandler(
	env *object.Env,
	recv object.PanObject,
	propName string,
	prop object.PanObject,
	chainArg object.PanObject,
	args []object.PanObject,
	kwargs *object.PanObj,
) object.PanObject {
	ret := evalCall(env, recv, prop, args, kwargs)
	if err, ok := ret.(*object.PanErr); ok {
		return err
	}

	return ret
}

func _evalPropCall(
	env *object.Env,
	recv object.PanObject,
	chainArg object.PanObject,
	propName string,
	args []object.PanObject,
	kwargs *object.PanObj,
	middlewares ..._PropCallMiddleware,
) object.PanObject {
	handle := mergePropCallMiddlewares(middlewares...)(propCallHandler)
	return handle(env, recv, propName, nil, chainArg, args, kwargs)
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
