package evaluator

import (
	"../ast"
	"../object"
	"fmt"
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

	val, ok := env.Get(object.GetSymHash(ident.Value))

	if !ok {
		err := object.NewNameErr(
			fmt.Sprintf("name `%s` is not defined.", ident.String()))
		return appendStackTrace(err, ident.Source())
	}

	return val
}
