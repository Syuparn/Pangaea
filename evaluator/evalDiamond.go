package evaluator

import (
	"github.com/Syuparn/pangaea/ast"
	"github.com/Syuparn/pangaea/object"
)

func evalDiamond(node *ast.DiamondLiteral, env *object.Env) object.PanObject {
	return object.BuiltInDiamondObj
}
