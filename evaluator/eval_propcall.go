package evaluator

import (
	"../ast"
	"../object"
)

func evalPropCall(node *ast.PropCallExpr, env *object.Env) object.PanObject {
	recv := Eval(node.Receiver, env)
	if err, ok := recv.(*object.PanErr); ok {
		return appendStackTrace(err, node.Source())
	}

	propStr := node.Prop.Value
	propHash := object.GetSymHash(propStr)

	retVal, ok := callProp(recv, propHash)

	if !ok {
		// TODO: error handling
		return nil
	}

	return retVal
}
