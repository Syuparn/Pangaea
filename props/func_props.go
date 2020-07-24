package props

import (
	"../object"
	"fmt"
)

func FuncProps(propContainer map[string]object.PanObject) map[string]object.PanObject {
	// NOTE: inject some built-in functions which relate to parser or evaluator
	return map[string]object.PanObject{
		"==": f(
			func(
				env *object.Env, kwargs *object.PanObj, args ...object.PanObject,
			) object.PanObject {
				if len(args) < 2 {
					return object.NewTypeErr("== requires at least 2 args")
				}

				// necessary for Func itself! (guarantee `Func == Func`)
				if args[0] == object.BuiltInFuncObj && args[1] == object.BuiltInFuncObj {
					return object.BuiltInTrue
				}

				// func comparison
				fn, ok := traceProtoOf(args[0], isFunc)
				if ok {
					other, ok := traceProtoOf(args[1], isFunc)
					if !ok {
						return object.BuiltInFalse
					}
					return compFuncs(fn.(*object.PanFunc), other.(*object.PanFunc))
				}

				// BuiltInFunc comparison
				builtIn, ok := traceProtoOf(args[0], isBuiltInFunc)
				if !ok {
					return object.BuiltInFalse
				}
				other, ok := traceProtoOf(args[1], isBuiltInFunc)
				if !ok {
					return object.BuiltInFalse
				}

				return compBuiltInFuncs(
					builtIn.(*object.PanBuiltIn), other.(*object.PanBuiltIn))
			},
		),
		"call": propContainer["Func_call"],
	}
}

func compFuncs(f1 *object.PanFunc, f2 *object.PanFunc) object.PanObject {
	// if src is equivalent, return true
	if f1.Inspect() == f2.Inspect() {
		return object.BuiltInTrue
	}
	return object.BuiltInFalse
}

func compBuiltInFuncs(f1 *object.PanBuiltIn, f2 *object.PanBuiltIn) object.PanObject {
	// if pointer is same, return true
	if fmt.Sprintf("%p", f1) == fmt.Sprintf("%p", f2) {
		return object.BuiltInTrue
	}
	return object.BuiltInFalse
}
