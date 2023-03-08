package di

import (
	"strings"

	"github.com/Syuparn/pangaea/object"
)

func strEvalEnv(
	env *object.Env,
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

	// NOTE: object.NewEnv cannot be used because an empty env does not have built-in objects
	newEnv := object.NewEnclosedEnv(env.Global())
	result := eval(strings.NewReader(self.Value), newEnv)
	if result.Type() == object.ErrType {
		return result
	}

	return newEnv.Items()
}
