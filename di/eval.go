package di

import (
	"github.com/Syuparn/pangaea/evaluator"
	"github.com/Syuparn/pangaea/object"
	"github.com/Syuparn/pangaea/parser"
)

// eval is a wrapper of evaluator.Eval.
func eval(src *parser.Reader, env *object.Env) object.PanObject {
	node, err := parser.Parse(src)
	if err != nil {
		e := object.NewSyntaxErr("failed to parse")
		e.StackTrace = err.Error()
		return e
	}

	result := evaluator.Eval(node, env)
	return result
}
