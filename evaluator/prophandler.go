package evaluator

import (
	"../object"
)

func builtInCallProp(
	env *object.Env,
	kwargs *object.PanObj,
	args ...object.PanObject,
) object.PanObject {
	if len(args) < 3 {
		return object.BuiltInNil
	}

	obj := args[1]
	propName, ok := args[2].(*object.PanStr)
	if !ok {
		return object.BuiltInNil
	}

	ret, ok := callProp(obj, object.GetSymHash(propName.Value))
	if !ok {
		return object.BuiltInNil
	}

	// (recv, args_for_call...)
	argsToPass := append([]object.PanObject{obj}, args[3:]...)

	switch f := ret.(type) {
	case *object.PanFunc:
		return evalPanFuncCall(f, env, kwargs, argsToPass...)
	case *object.PanBuiltIn:
		return f.Fn(env, kwargs, argsToPass...)
	default:
		// not callable
		return ret
	}
}
