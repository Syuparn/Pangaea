package evaluator

import (
	"github.com/Syuparn/pangaea/ast"
	"github.com/Syuparn/pangaea/object"
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
	case ast.RaiseJump:
		return evalJumpIfRaise(node, env, cond)
	default:
		err := object.NewNotImplementedErr("the stmt is not implemented yet")
		return appendStackTrace(err, node.Source())
	}
}

func evalJumpIfYield(
	node *ast.JumpIfStmt,
	env *object.Env,
	cond object.PanObject,
) object.PanObject {
	if !isTruthy(cond, env) {
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
	if !isTruthy(cond, env) {
		// do nothing and keep func evaluating
		return object.BuiltInNil
	}

	// return if cond is truthy
	// NOTE: wrap by ReturnObj to tell evalProgram to stop evaluation
	ret := Eval(node.JumpStmt, env)
	if err, ok := ret.(*object.PanErr); ok {
		return appendStackTrace(err, node.Source())
	}
	return ret
}

func evalJumpIfDefer(
	node *ast.JumpIfStmt,
	env *object.Env,
	cond object.PanObject,
) object.PanObject {
	if !isTruthy(cond, env) {
		// do nothing and keep func evaluating
		return object.BuiltInNil
	}
	return &object.DeferObj{Node: node.JumpStmt.Val}
}

func isTruthy(obj object.PanObject, env *object.Env) bool {
	if b, ok := obj.(*object.PanBool); ok {
		return b == object.BuiltInTrue
	}

	// use (obj).B to check truthy/falsy
	bSym := object.NewPanStr("B")
	cond := builtInCallProp(env, object.EmptyPanObjPtr(),
		object.EmptyPanObjPtr(), obj, bSym)
	return cond == object.BuiltInTrue
}

func evalJumpIfRaise(
	node *ast.JumpIfStmt,
	env *object.Env,
	cond object.PanObject,
) object.PanObject {
	if !isTruthy(cond, env) {
		// do nothing and keep func evaluating
		return object.BuiltInNil
	}

	// return if cond is truthy
	// NOTE: wrap by ReturnObj to tell evalProgram to stop evaluation
	ret := Eval(node.JumpStmt, env)
	if err, ok := ret.(*object.PanErr); ok {
		return appendStackTrace(err, node.Source())
	}

	// unwrap errWrapper to re-raise error
	if w, ok := ret.(*object.PanErrWrapper); ok {
		err := w.PanErr
		return appendStackTrace(&err, node.Source())
	}

	return ret
}
