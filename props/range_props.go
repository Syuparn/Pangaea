package props

import (
	"fmt"

	"github.com/Syuparn/pangaea/object"
)

// RangeProps provides built-in props for Range.
// NOTE: Some Range props are defind by native code (not by this function).
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

				self, ok := object.TraceProtoOfRange(args[0])
				if !ok {
					return object.BuiltInFalse
				}
				other, ok := object.TraceProtoOfRange(args[1])
				if !ok {
					return object.BuiltInFalse
				}

				return compRanges(self, other, propContainer, env)
			},
		),
		"_iter": f(
			func(
				env *object.Env, kwargs *object.PanObj, args ...object.PanObject,
			) object.PanObject {
				if len(args) < 1 {
					return object.NewTypeErr("Range#_iter requires at least 1 arg")
				}

				self, ok := object.TraceProtoOfRange(args[0])
				if !ok {
					return object.NewTypeErr("\\1 must be Range")
				}

				stepInt, err := stepIntOf(self)
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

					resInt, ok := object.TraceProtoOfInt(res)
					if !ok {
						return false, object.NewValueErr(`<=> returned non-int value`)
					}

					if stepInt.Value > 0 {
						// o <=> r.Stop is 0 or 1 if o >= r.Stop
						return resInt.Value != -1, nil
					}
					// o <=> r.Stop is -1 or 0 if o <= r.Stop
					return resInt.Value != 1, nil
				}

				return &object.PanBuiltInIter{
					Fn:  rangeIter(self, next, reachesStop),
					Env: env, // not used
				}
			},
		),
		"_name": object.NewPanStr("Range"),
		"B": f(
			func(
				env *object.Env, kwargs *object.PanObj, args ...object.PanObject,
			) object.PanObject {
				if len(args) < 1 {
					return object.NewTypeErr("Range#B requires at least 1 arg")
				}
				_, ok := object.TraceProtoOfRange(args[0])
				if !ok {
					return object.NewTypeErr(`\1 must be range`)
				}
				return object.BuiltInTrue
			},
		),
		"new": f(
			func(
				env *object.Env, kwargs *object.PanObj, args ...object.PanObject,
			) object.PanObject {
				if len(args) < 2 {
					return object.NewTypeErr("Range#new requires at least 2 args")
				}

				// insufficient args are filled by nil (Pangaea call spec)
				switch len(args) {
				case 2:
					return &object.PanRange{
						Start: args[1],
						Stop:  object.BuiltInNil,
						Step:  object.BuiltInNil,
					}
				case 3:
					return &object.PanRange{
						Start: args[1],
						Stop:  args[2],
						Step:  object.BuiltInNil,
					}
				default:
					return &object.PanRange{
						Start: args[1],
						Stop:  args[2],
						Step:  args[3],
					}
				}
			},
		),
		"start": f(
			func(
				env *object.Env, kwargs *object.PanObj, args ...object.PanObject,
			) object.PanObject {
				if len(args) < 1 {
					return object.NewTypeErr("Range#start requires at least 1 arg")
				}
				r, ok := object.TraceProtoOfRange(args[0])
				if !ok {
					return object.NewTypeErr(
						fmt.Sprintf("%s cannot be treated as range", ReprStr(args[0])))
				}

				return r.Start
			},
		),
		"stop": f(
			func(
				env *object.Env, kwargs *object.PanObj, args ...object.PanObject,
			) object.PanObject {
				if len(args) < 1 {
					return object.NewTypeErr("Range#stop requires at least 1 arg")
				}
				r, ok := object.TraceProtoOfRange(args[0])
				if !ok {
					return object.NewTypeErr(
						fmt.Sprintf("%s cannot be treated as range", ReprStr(args[0])))
				}

				return r.Stop
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
	if r.Step.Type() == object.NilType {
		return object.NewPanInt(1), nil
	}

	// step must be proto of int
	step, ok := object.TraceProtoOfInt(r.Step)
	if !ok {
		return nil, object.NewValueErr("step must be Int")
	}

	if step.Value == 0 {
		return nil, object.NewValueErr("cannot use 0 for range step")
	}

	return step, nil
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
