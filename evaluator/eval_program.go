package evaluator

import (
	"../ast"
	"../object"
)

func evalProgram(p *ast.Program, env *object.Env) object.PanObject {
	return evalStmts(p.Stmts, env)
}

func evalStmts(stmts []ast.Stmt, env *object.Env) object.PanObject {
	ret, deferObjs := _evalStmts(stmts, env)
	err := evalDefer(deferObjs, env)
	if err != nil {
		return err
	}

	return ret
}

func _evalStmts(
	stmts []ast.Stmt,
	env *object.Env,
) (object.PanObject, []object.DeferObj) {
	// NOTE: if program has no lines, it is evaluated as `nil`
	var val object.PanObject = object.BuiltInNil
	// if `yield` exists in the stmts, this value is set
	var yielded object.PanObject = nil
	// deferred elems are evaluated after return
	deferObjs := []object.DeferObj{}

	for _, stmt := range stmts {
		val = Eval(stmt, env)

		if err, ok := val.(*object.PanErr); ok {
			return appendStackTrace(err, stmt.Source()), deferObjs
		}

		// unwrap ReturnObj
		if ret, ok := val.(*object.ReturnObj); ok {
			return ret.PanObject, deferObjs
		}

		if _defer, ok := val.(*object.DeferObj); ok {
			deferObjs = append(deferObjs, *_defer)
		}

		if val.Type() == object.YieldType {
			y := val.(*object.YieldObj).PanObject

			if err, ok := y.(*object.PanErr); ok {
				return appendStackTrace(err, stmt.Source()), deferObjs
			}

			// NOTE: only first yield is valid
			if yielded == nil {
				yielded = y
			}
		}
	}

	// NOTE: return value precedence:
	// last stmt < yield
	if yielded != nil {
		return yielded, deferObjs
	}

	return val, deferObjs
}

func evalDefer(deferObjs []object.DeferObj, env *object.Env) *object.PanErr {
	for _, o := range deferObjs {
		ret := Eval(o.Node, env)
		if err, ok := ret.(*object.PanErr); ok {
			return err
		}
	}

	return nil
}
