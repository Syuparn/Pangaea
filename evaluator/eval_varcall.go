package evaluator

import (
	"github.com/Syuparn/pangaea/ast"
	"github.com/Syuparn/pangaea/object"
)

func evalVarCall(node *ast.VarCallExpr, env *object.Env) object.PanObject {
	recv, err := extractRecv(node.Receiver, env)
	if err != nil {
		return appendStackTrace(err, node.Source())
	}

	f := Eval(node.Var, env)
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
