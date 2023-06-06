package di

import (
	"strings"

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

	// TODO: enable to choose whether current env is used or not
	result := eval(parser.NewReader(strings.NewReader(self.Value), StrFileName), object.NewEnclosedEnv(env))
	return result
}
