package evaluator

import (
	"../ast"
	"../object"
)

func evalJumpStmt(node *ast.JumpStmt, env *object.Env) object.PanObject {
	// if stmt is `defer`, not evaluate
	if node.JumpType == ast.DeferJump {
		return &object.DeferObj{Node: node.Val}
	}

	val := Eval(node.Val, env)

	if err, ok := val.(*object.PanErr); ok {
		appendStackTrace(err, node.Source())
		return err
	}

	switch node.JumpType {
	case ast.ReturnJump:
		return &object.ReturnObj{PanObject: val}
	case ast.YieldJump:
		return &object.YieldObj{PanObject: val}
	default:
		// TODO: handle raise
		err := object.NewNotImplementedErr("the stmt is not implemented yet")
		return appendStackTrace(err, node.Source())
	}
}
