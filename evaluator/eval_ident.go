package evaluator

import (
	"../ast"
	"../object"
)

func evalIdent(ident *ast.Ident, env *object.Env) object.PanObject {
	// check if ident refers keyword
	switch ident.Value {
	case "true":
		return object.BuiltInTrue
	case "false":
		return object.BuiltInFalse
	case "nil":
		return object.BuiltInNil
	}
	// TODO: return error object
	return nil
}
