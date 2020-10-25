package evaluator

import (
	"github.com/Syuparn/pangaea/ast"
	"github.com/Syuparn/pangaea/object"
)

func evalPinnedIdent(node *ast.PinnedIdent, env *object.Env) object.PanObject {
	// NOTE: pinned keys are evaluated in evalObj/evalMap
	err := object.NewSyntaxErr("cannot use `^` other than key or var chain.")
	return appendStackTrace(err, node.Source())
}
