package evaluator

import (
	"../ast"
	"../object"
)

func evalVarCall(node *ast.VarCallExpr, env *object.Env) object.PanObject {
	recv, err := extractRecv(node.Receiver, env)
	if err != nil {
		return appendStackTrace(err, node.Source())
	}

	fObj := Eval(node.Var, env)
	if err, ok := fObj.(*object.PanErr); ok {
		return appendStackTrace(err, node.Source())
	}

	// TODO: duck typing (allow all objs with `call` prop)
	f, ok := object.TraceProtoOfFunc(fObj)
	if !ok {
		err := object.NewTypeErr("var call must be func")
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
