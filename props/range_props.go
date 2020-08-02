package props

import (
	"../object"
)

func RangeProps(propContainer map[string]object.PanObject) map[string]object.PanObject {
	// NOTE: inject some built-in functions which relate to parser or evaluator
	return map[string]object.PanObject{
		"==": f(
			func(
				env *object.Env, kwargs *object.PanObj, args ...object.PanObject,
			) object.PanObject {
				if len(args) < 2 {
					return object.NewTypeErr("== requires at least 2 args")
				}

				// necessary for Range itself! (guarantee `Range == Range`)
				if args[0] == object.BuiltInRangeObj && args[1] == object.BuiltInRangeObj {
					return object.BuiltInTrue
				}

				self, ok := traceProtoOf(args[0], isRange)
				if !ok {
					return object.BuiltInFalse
				}
				other, ok := traceProtoOf(args[1], isRange)
				if !ok {
					return object.BuiltInFalse
				}

				return compRanges(
					self.(*object.PanRange), other.(*object.PanRange), propContainer, env)
			},
		),
		"B": f(
			func(
				env *object.Env, kwargs *object.PanObj, args ...object.PanObject,
			) object.PanObject {
				if len(args) < 1 {
					return object.NewTypeErr("Range#B requires at least 1 arg")
				}
				_, ok := traceProtoOf(args[0], isRange)
				if !ok {
					return object.NewTypeErr(`\1 must be range`)
				}
				return object.BuiltInTrue
			},
		),
	}
}

func compRanges(
	r1 *object.PanRange,
	r2 *object.PanRange,
	propContainer map[string]object.PanObject,
	env *object.Env,
) object.PanObject {
	vals := []struct {
		v1 object.PanObject
		v2 object.PanObject
	}{
		{r1.Start, r2.Start},
		{r1.Stop, r2.Stop},
		{r1.Step, r2.Step},
	}
	for _, val := range vals {
		// == comparison
		res := propContainer["Obj_callProp"].(*object.PanBuiltIn).Fn(
			env, object.EmptyPanObjPtr(),
			object.EmptyPanObjPtr(), val.v1, eqSym, val.v2,
		)
		if res == object.BuiltInFalse {
			return object.BuiltInFalse
		}
	}
	return object.BuiltInTrue
}
