package di

import (
	"strings"

	"github.com/Syuparn/pangaea/object"
)

func strEvalEnv(
	_ *object.Env,
	kwargs *object.PanObj,
	args ...object.PanObject,
) object.PanObject {
	if len(args) < 1 {
		return object.NewTypeErr("evalEnv requires at least 1 arg")
	}

	self, ok := args[0].(*object.PanStr)
	if !ok {
		return object.NewTypeErr("\\1 must be str")
	}

	env := object.NewEnv()
	result := eval(strings.NewReader(self.Value), env)
	if result.Type() == object.ErrType {
		return result
	}

	return env.Items()
}
