package evaluator

import (
	"../ast"
	"../object"
)

func evalJumpIfStmt(node *ast.JumpIfStmt, env *object.Env) object.PanObject {
	cond := Eval(node.Cond, env)
	if err, ok := cond.(*object.PanErr); ok {
		appendStackTrace(err, node.Source())
		return err
	}

	// if cond is truthy, behave equivalent to jumpstmt
	// TODO: use cond.B to handle truthy/falsy
	if cond == object.BuiltInTrue {
		return Eval(node.JumpStmt, env)
	}

	switch node.JumpStmt.JumpType {
	case ast.ReturnJump:
		// do nothing and keep func evaluating
		return object.BuiltInNil
	case ast.YieldJump:
		// raise err
		err := object.NewStopIterErr("iter stopped")
		return appendStackTrace(err, node.Source())
	default:
		// TODO: handle raise
		err := object.NewNotImplementedErr("the stmt is not implemented yet")
		return appendStackTrace(err, node.Source())
	}
}
