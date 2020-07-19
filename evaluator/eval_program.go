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

	for _, stmt := range stmts {
		val = Eval(stmt, env)

		if err, ok := val.(*object.PanErr); ok {
			return appendStackTrace(err, stmt.Source())
		}
	}

	return val
}
