package evaluator

import (
	"github.com/Syuparn/pangaea/ast"
	"github.com/Syuparn/pangaea/object"
)

func evalLiteralCall(node *ast.LiteralCallExpr, env *object.Env) object.PanObject {
	recv, err := extractRecv(node.Receiver, env)
	if err != nil {
		return appendStackTrace(err, node.Source())
	}

	f := Eval(node.Func, env)
	if err, ok := f.(*object.PanErr); ok {
		return appendStackTrace(err, node.Source())
	}

	var chainArg object.PanObject = object.BuiltInNil
	if node.Chain.Arg != nil {
		chainArg = Eval(node.Chain.Arg, env)
	}
	if err, ok := chainArg.(*object.PanErr); ok {
		return appendStackTrace(err, node.Source())
	}

	chainMiddleware := newLiteralCallChainMiddleware(*node.Chain)
	ret := _evalLiteralCall(env, recv, chainArg, f,
		chainMiddleware, literalProxyMiddleware)
	if err, ok := ret.(*object.PanErr); ok {
		return appendStackTrace(err, node.Source())
	}

	return ret
}

func literalProxyMiddleware(next _LiteralCallMiddlewareHandler) _LiteralCallMiddlewareHandler {
	return func(
		env *object.Env,
		recv object.PanObject,
		chainArg object.PanObject,
		args []object.PanObject,
		kwargs *object.PanObj,
	) object.PanObject {
		// if recv has _proxyLiteral, call it instead
		if proxy, ok := findProxyLiteral(recv); ok {
			return _evalProxyLiteral(env, recv, args[0], proxy)
		}
		return next(env, recv, chainArg, args, kwargs)
	}
}

func literalCallHandler(
	env *object.Env,
	recv object.PanObject,
	chainArg object.PanObject,
	args []object.PanObject,
	kwargs *object.PanObj,
) object.PanObject {
	// TODO: duck typing (allow all objs with `call` prop)
	f, ok := object.TraceProtoOfFunc(args[0])
	if !ok {
		return object.NewTypeErr("literal call must be func")
	}

	// unpack recv for params
	argsToPass := append([]object.PanObject{f}, literalCallArgs(recv, f)...)

	return evalFuncCall(env, kwargs, argsToPass...)
}

func _evalLiteralCall(
	env *object.Env,
	recv object.PanObject,
	chainArg object.PanObject,
	f object.PanObject,
	middlewares ..._LiteralCallMiddleware,
) object.PanObject {
	args := []object.PanObject{f}
	kwargs := object.EmptyPanObjPtr()
	handle := mergeLiteralCallMiddlewares(middlewares...)(literalCallHandler)
	return handle(env, recv, chainArg, args, kwargs)
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
