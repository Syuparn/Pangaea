package evaluator

import (
	"../ast"
	"../object"
	"fmt"
)

func evalIdent(ident *ast.Ident, env *object.Env) object.PanObject {
	val, ok := env.Get(object.GetSymHash(ident.Value))

	if !ok {
		err := object.NewNameErr(
			fmt.Sprintf("name `%s` is not defined", ident.String()))
		return appendStackTrace(err, ident.Source())
	}

	return val
}
