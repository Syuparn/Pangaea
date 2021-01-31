package props

import (
	"fmt"

	"github.com/Syuparn/pangaea/object"
)

// FuncProps provides built-in props for Func.
// NOTE: Some Func props are defind by native code (not by this function).
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
				fn, ok := object.TraceProtoOfFunc(args[0])
				if ok {
					other, ok := object.TraceProtoOfFunc(args[1])
					if !ok {
						return object.BuiltInFalse
					}
					return compFuncs(fn, other)
				}

				// BuiltInFunc comparison
				builtIn, ok := object.TraceProtoOfBuiltInFunc(args[0])
				if !ok {
					return object.BuiltInFalse
				}
				other, ok := object.TraceProtoOfBuiltInFunc(args[1])
				if !ok {
					return object.BuiltInFalse
				}

				return compBuiltInFuncs(builtIn, other)
			},
		),
		"_name": object.NewPanStr("Func"),
		"B": f(
			func(
				env *object.Env, kwargs *object.PanObj, args ...object.PanObject,
			) object.PanObject {
				if len(args) < 1 {
					return object.NewTypeErr("Func#B requires at least 1 arg")
				}
				_, ok := object.TraceProtoOfFunc(args[0])
				if !ok {
					return object.NewTypeErr(`\1 must be func`)
				}
				return object.BuiltInTrue
			},
		),
		"call": propContainer["Func_call"],
		"new": f(
			func(
				env *object.Env, kwargs *object.PanObj, args ...object.PanObject,
			) object.PanObject {
				if len(args) < 2 {
					return object.NewTypeErr("Func#new requires at least 2 args")
				}
				f, ok := object.TraceProtoOfFunc(args[1])
				if !ok {
					return object.NewTypeErr(
						fmt.Sprintf("%s cannot be treated as func", args[1].Inspect()))
				}

				return f
			},
		),
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
