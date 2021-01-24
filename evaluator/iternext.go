package evaluator

import (
	"fmt"
	"github.com/Syuparn/pangaea/object"
)

// NOTE: Iter#next is immutable! (this overwrites iter's env)
func iterNext(
	env *object.Env,
	kwargs *object.PanObj,
	args ...object.PanObject,
) object.PanObject {
	if len(args) < 1 {
		return object.NewTypeErr("Iter#next requires at least 1 arg")
	}

	// ignore args (args are set in env by Iter#new method or recur in last loop)
	return evalIterCall(args[0])
}

func evalIterCall(self object.PanObject) object.PanObject {
	switch f := self.(type) {
	case *object.PanFunc:
		// inject var `recur`
		f.Env.InjectRecur(recur(f))
		retVal := evalStmts(*f.Body(), f.Env)

		if err, ok := retVal.(*object.PanErr); ok {
			return appendStackTrace(err, (*f.Body())[0].Source())
		}

		return retVal
	case *object.PanBuiltInIter:
		return f.Fn(f.Env, object.EmptyPanObjPtr() /*empty kwargs*/)
	default:
		return object.NewTypeErr(
			fmt.Sprintf("`%s` is not callable.", self.Inspect()))
	}
}

func recur(iter *object.PanFunc) object.BuiltInFunc {
	// func that replaces iter Env to args/kwargs

	return func(
		env *object.Env,
		kwargs *object.PanObj,
		args ...object.PanObject,
	) object.PanObject {
		// replace iter.Env
		// newEnv is in same closure as old env
		newEnv := object.NewEnclosedEnv(iter.Env.Outer())
		assignArgsToEnv(newEnv, iter.Args().Elems, iter.Kwargs(), args, kwargs)
		iter.Env = newEnv

		return object.BuiltInNil
	}
}
