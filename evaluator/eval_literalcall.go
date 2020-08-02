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

	args := literalCallArgs(recv, f)
	// prepend f itself to args
	args = append([]object.PanObject{f}, args...)

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
		return evalScalarLiteralCall(node, env, f, args)
	default:
		return nil
	}
}

func evalScalarLiteralCall(
	node *ast.LiteralCallExpr,
	env *object.Env,
	f object.PanObject,
	args []object.PanObject,
) object.PanObject {
	kwargs := object.EmptyPanObjPtr()
	ret := evalFuncCall(env, kwargs, args...)

	if err, ok := f.(*object.PanErr); ok {
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
