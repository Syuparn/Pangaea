package evaluator

import (
	"github.com/Syuparn/pangaea/ast"
	"github.com/Syuparn/pangaea/object"
)

func evalArgs(argNodes []ast.Expr, env *object.Env) ([]object.PanObject, *object.PanErr) {
	args := []object.PanObject{}
	for _, argNode := range argNodes {
		arg := Eval(argNode, env)

		if err, ok := arg.(*object.PanErr); ok {
			appendStackTrace(err, argNode.Source())
			return []object.PanObject{}, err
		}

		args = append(args, arg)
	}

	return args, nil
}
