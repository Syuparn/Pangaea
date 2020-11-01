package di

import (
	"strings"

	"github.com/Syuparn/pangaea/evaluator"
	"github.com/Syuparn/pangaea/object"
	"github.com/Syuparn/pangaea/parser"
)

func strEvalEnv(
	_ *object.Env,
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

	env := object.NewEnv()

	result := evaluator.Eval(node, env)
	if result.Type() == object.ErrType {
		return result
	}

	return env.Items()
}
