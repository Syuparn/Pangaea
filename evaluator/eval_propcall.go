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
	args = append([]object.PanObject{recv}, args...)

	return evalFuncCall(env, kwargs, args...)
}
