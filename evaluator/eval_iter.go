package evaluator

import (
	"github.com/Syuparn/pangaea/ast"
	"github.com/Syuparn/pangaea/object"
)

func evalIter(node *ast.IterLiteral, env *object.Env) object.PanObject {
	return evalCallable(node.FuncComponent, env, object.IterFunc)
}
