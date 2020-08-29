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
		"_iter": f(
			func(
				env *object.Env, kwargs *object.PanObj, args ...object.PanObject,
			) object.PanObject {
				if len(args) < 1 {
					return object.NewTypeErr("Range#_iter requires at least 1 arg")
				}

				self, ok := traceProtoOf(args[0], isRange)
				if !ok {
					return object.NewTypeErr("\\1 must be Range")
				}

				range_, _ := self.(*object.PanRange)
				stepInt, err := stepIntOf(range_)
				if err != nil {
					return err
				}

				// call prop `_incBy(range.step)`
				incBySym := object.NewPanStr("_incBy")
				next := func(n object.PanObject) object.PanObject {
					return propContainer["Obj_callProp"].(*object.PanBuiltIn).Fn(
						env, object.EmptyPanObjPtr(),
						object.EmptyPanObjPtr(), n, incBySym, stepInt,
					)
				}

				// call prop `<=>`
				spaceshipSym := object.NewPanStr("<=>")
				reachesStop := func(
					r *object.PanRange,
					o object.PanObject,
				) (bool, *object.PanErr) {
					// o <=> r.Stop
					res := propContainer["Obj_callProp"].(*object.PanBuiltIn).Fn(
						env, object.EmptyPanObjPtr(),
						object.EmptyPanObjPtr(), o, spaceshipSym, r.Stop,
					)
					if err, ok := res.(*object.PanErr); ok {
						return false, err
					}

					resInt, ok := traceProtoOf(res, isInt)
					if !ok {
						return false, object.NewValueErr(`<=> returned non-int value`)
					}

					if stepInt.Value > 0 {
						// o <=> r.Stop is 0 or 1 if o >= r.Stop
						return resInt.(*object.PanInt).Value != -1, nil
					}
					// o <=> r.Stop is -1 or 0 if o <= r.Stop
					return resInt.(*object.PanInt).Value != 1, nil
				}

				return &object.PanBuiltInIter{
					Fn:  rangeIter(range_, next, reachesStop),
					Env: env, // not used
				}
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

func stepIntOf(r *object.PanRange) (*object.PanInt, *object.PanErr) {
	// default step = 1
	if r.Step.Type() == object.NIL_TYPE {
		return object.NewPanInt(1), nil
	}

	// step must be proto of int
	step, ok := traceProtoOf(r.Step, isInt)
	if !ok {
		return nil, object.NewValueErr("step must be Int")
	}

	stepInt, _ := step.(*object.PanInt)
	if stepInt.Value == 0 {
		return nil, object.NewValueErr("cannot use 0 for range step")
	}

	return stepInt, nil
}

func rangeIter(
	r *object.PanRange,
	next func(object.PanObject) object.PanObject,
	reachesStop func(*object.PanRange, object.PanObject) (bool, *object.PanErr),
) object.BuiltInFunc {
	current := r.Start

	return func(
		env *object.Env, kwargs *object.PanObj, args ...object.PanObject,
	) object.PanObject {
		reached, err := reachesStop(r, current)
		if err != nil {
			return err
		}

		if reached {
			return object.NewStopIterErr("iter stopped")
		}
		yielded := current
		current = next(current)
		return yielded
	}
}
