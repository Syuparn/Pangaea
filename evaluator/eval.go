package evaluator

import (
	"../ast"
	"../object"
)

func Eval(node ast.Node, env *object.Env) object.PanObject {
	switch node := node.(type) {
	// Program
	case *ast.Program:
		return evalProgram(node, env)
	// Stmt
	case *ast.ExprStmt:
		return Eval(node.Expr, env)
	// Expr
	case *ast.IntLiteral:
		return &object.PanInt{Value: node.Value}
	case *ast.FloatLiteral:
		return &object.PanFloat{Value: node.Value}
	case *ast.StrLiteral:
		return &object.PanStr{Value: node.Value}
	}

	return nil
}
