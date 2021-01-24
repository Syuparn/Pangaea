package evaluator

import (
	"github.com/Syuparn/pangaea/ast"
	"github.com/Syuparn/pangaea/object"
)

func evalInt(node *ast.IntLiteral, env *object.Env) object.PanObject {
	return object.NewPanInt(node.Value)
}
