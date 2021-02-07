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
	f, ok := object.TraceProtoOfFunc(args[0])
	if ok {
		// unpack recv for params
		argsToPass := append([]object.PanObject{f}, literalCallArgs(recv, f)...)
		return evalFuncCall(env, kwargs, argsToPass...)
	}

	builtIn, ok := object.TraceProtoOfBuiltInFunc(args[0])
	if ok {
		// TODO: extract elems if there are more than 1 params and recv is arr
		return builtIn.Fn(env, kwargs, recv)
	}

	if call, ok := findCall(args[0]); ok {
		return handleCallable(env, recv, kwargs, args[0], call)
	}

	return object.NewTypeErr("literal call must be func")
}

func handleCallable(
	env *object.Env,
	recv object.PanObject,
	kwargs *object.PanObj,
	callable object.PanObject,
	call object.PanObject,
) object.PanObject {
	f, ok := object.TraceProtoOfFunc(call)
	if ok {
		// callable is passed to call as self
		argsToPass := append([]object.PanObject{f, callable}, literalCallArgs(recv, f)...)
		return evalFuncCall(env, kwargs, argsToPass...)
	}

	// NOTE: builtin func is not supported because it cannot recieve addtional self
	// NOTE: 'call prop in call is not searched to prevent infinite loop
	return object.NewTypeErr("prop 'call in varcall must be func")
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

func findCall(funcSubstitute object.PanObject) (object.PanObject, bool) {
	call, isMissing := evalProp("call", funcSubstitute)
	if isMissing {
		return nil, false
	}
	if _, ok := call.(*object.PanErr); ok {
		return nil, false
	}
	return call, true
}
