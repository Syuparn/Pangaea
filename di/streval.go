package di

import (
	"strings"

	"github.com/Syuparn/pangaea/evaluator"
	"github.com/Syuparn/pangaea/object"
	"github.com/Syuparn/pangaea/parser"
)

func strEval(
	env *object.Env,
	kwargs *object.PanObj,
	args ...object.PanObject,
) object.PanObject {
	if len(args) < 1 {
		return object.NewTypeErr("eval requires at least 1 arg")
	}

	self, ok := args[0].(*object.PanStr)
	if !ok {
		return object.NewTypeErr("\\1 must be str")
	}

	node, err := parser.Parse(strings.NewReader(self.Value))
	if err != nil {
		e := object.NewSyntaxErr("failed to parse")
		e.StackTrace = err.Error()
		return e
	}

	// TODO: enable to choose whether current env is used or not
	result := evaluator.Eval(node, object.NewEnclosedEnv(env))
	return result
}
