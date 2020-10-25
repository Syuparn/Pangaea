package evaluator

import (
	"github.com/Syuparn/pangaea/object"
)

// generate copied iter with args set to its env
func iterNew(
	env *object.Env,
	kwargs *object.PanObj,
	args ...object.PanObject,
) object.PanObject {
	if len(args) < 1 {
		return object.NewTypeErr("Iter#new requires at least 1 arg")
	}

	// allow descendant of iter
	self, ok := object.TraceProtoOfFunc(args[0])
	if !ok {
		return object.NewTypeErr("\\1 must be iter")
	}

	// locate env in same closure as self.Env
	newEnv := object.NewEnclosedEnv(self.Env.Outer())
	assignArgsToEnv(newEnv, self.Args().Elems, self.Kwargs(), args[1:], kwargs)
	return &object.PanFunc{
		FuncWrapper: self.FuncWrapper,
		FuncKind:    object.IterFunc,
		Env:         newEnv,
	}
}
