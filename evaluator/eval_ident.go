package evaluator

import (
	"fmt"
	"github.com/Syuparn/pangaea/ast"
	"github.com/Syuparn/pangaea/object"
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
