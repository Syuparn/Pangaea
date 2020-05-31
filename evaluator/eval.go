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
		return evalInt(node, env)
	case *ast.FloatLiteral:
		return &object.PanFloat{Value: node.Value}
	case *ast.StrLiteral:
		return &object.PanStr{Value: node.Value}
	case *ast.SymLiteral:
		return &object.PanStr{Value: node.Value}
	case *ast.RangeLiteral:
		return evalRange(node, env)
	case *ast.Ident:
		return evalIdent(node, env)
	}

	return nil
}
