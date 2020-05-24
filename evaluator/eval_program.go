package evaluator

import (
	"../ast"
	"../object"
)

func evalProgram(p *ast.Program, env *object.Env) object.PanObject {
	// FIXME: eval multiple stmt lines!
	return Eval(p.Stmts[0], env)
}
