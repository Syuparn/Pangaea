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

	switch node.JumpStmt.JumpType {
	case ast.ReturnJump:
		return evalJumpIfReturn(node, env, cond)
	case ast.YieldJump:
		return evalJumpIfYield(node, env, cond)
	case ast.DeferJump:
		return evalJumpIfDefer(node, env, cond)
	default:
		// TODO: handle raise
		err := object.NewNotImplementedErr("the stmt is not implemented yet")
		return appendStackTrace(err, node.Source())
	}
}

func evalJumpIfYield(
	node *ast.JumpIfStmt,
	env *object.Env,
	cond object.PanObject,
) object.PanObject {
	if !isTruthy(cond) {
		// stop iteration
		err := object.NewStopIterErr("iter stopped")
		return appendStackTrace(err, node.Source())
	}

	ret := Eval(node.JumpStmt, env)
	if err, ok := ret.(*object.PanErr); ok {
		return appendStackTrace(err, node.Source())
	}
	return ret
}

func evalJumpIfReturn(
	node *ast.JumpIfStmt,
	env *object.Env,
	cond object.PanObject,
) object.PanObject {
	if !isTruthy(cond) {
		// do nothing and keep func evaluating
		return object.BuiltInNil
	}

	// return if cond is truthy
	// NOTE: wrap by ReturnObj to tell evalProgram to stop evaluation
	ret := Eval(node.JumpStmt, env)
	if err, ok := ret.(*object.PanErr); ok {
		return appendStackTrace(err, node.Source())
	}
	return &object.ReturnObj{PanObject: ret}
}

func evalJumpIfDefer(
	node *ast.JumpIfStmt,
	env *object.Env,
	cond object.PanObject,
) object.PanObject {
	if !isTruthy(cond) {
		// do nothing and keep func evaluating
		return object.BuiltInNil
	}
	return &object.DeferObj{Node: node.JumpStmt.Val}
}

func isTruthy(obj object.PanObject) bool {
	// TODO: use cond.B to check truthy/falsy
	return obj == object.BuiltInTrue
}
