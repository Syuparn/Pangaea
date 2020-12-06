package evaluator

import (
	"github.com/Syuparn/pangaea/ast"
	"github.com/Syuparn/pangaea/object"
)

func findProxyLiteral(recv object.PanObject) (object.PanObject, bool) {
	proxy, isMissing := evalProp("_literalProxy", recv)
	if isMissing {
		return nil, false
	}
	if _, ok := proxy.(*object.PanErr); ok {
		return nil, false
	}
	return proxy, true
}

func _evalProxyLiteral(
	env *object.Env,
	recv object.PanObject,
	fObj object.PanObject,
	proxy object.PanObject,
) object.PanObject {
	// TODO: duck typing (allow all objs with `call` prop)
	f, ok := object.TraceProtoOfFunc(fObj)
	if !ok {
		return object.NewTypeErr("literal call must be func")
	}

	args := []object.PanObject{f}
	kwargs := object.EmptyPanObjPtr()
	ret := evalCall(env, recv, proxy, args, kwargs)
	if err, ok := ret.(*object.PanErr); ok {
		return err
	}

	return ret
}

func evalProxyLiteral(
	node ast.Node,
	env *object.Env,
	recv object.PanObject,
	// TODO: duck typing (allow all objs with `call` prop)
	f *object.PanFunc,
	proxy object.PanObject,
) object.PanObject {
	args := []object.PanObject{f}
	kwargs := object.EmptyPanObjPtr()
	ret := evalCall(env, recv, proxy, args, kwargs)
	if err, ok := ret.(*object.PanErr); ok {
		return appendStackTrace(err, node.Source())
	}

	return ret
}
