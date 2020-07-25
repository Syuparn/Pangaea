package evaluator

import (
	"../ast"
	"../object"
)

func evalProgram(p *ast.Program, env *object.Env) object.PanObject {
	return evalStmts(p.Stmts, env)
}

func evalStmts(stmts []ast.Stmt, env *object.Env) object.PanObject {
	// NOTE: if program has no lines, it is evaluated as `nil`
	var val object.PanObject = object.BuiltInNil
	// if `yield` exists in the stmts, this value is set
	var yielded object.PanObject = nil

	for _, stmt := range stmts {
		val = Eval(stmt, env)

		if err, ok := val.(*object.PanErr); ok {
			return appendStackTrace(err, stmt.Source())
		}

		if val.Type() == object.YIELD_TYPE {
			y := val.(*object.YieldObj).PanObject

			if err, ok := y.(*object.PanErr); ok {
				return appendStackTrace(err, stmt.Source())
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
		return yielded
	}

	return val
}
