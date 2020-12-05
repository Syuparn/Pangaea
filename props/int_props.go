package props

import (
	"fmt"
	"math/big"

	"github.com/Syuparn/pangaea/object"
)

// IntProps provides built-in props for Int.
// NOTE: Some Int props are defind by native code (not by this function).
func IntProps(propContainer map[string]object.PanObject) map[string]object.PanObject {
	// NOTE: inject some built-in functions which relate to parser or evaluator
	return map[string]object.PanObject{
		"<=>": f(
			func(
				env *object.Env, kwargs *object.PanObj, args ...object.PanObject,
			) object.PanObject {
				self, other, err := checkIntInfixArgs(args, "<=>", object.NewPanInt(0))
				if err != nil {
					return err
				}

				selfVal := self.Value
				otherVal := other.Value
				var res int64

				if selfVal > otherVal {
					res = 1
				} else if selfVal == otherVal {
					res = 0
				} else {
					res = -1
				}

				return object.NewPanInt(res)
			},
		),
		// NOTE: this cannot be removed (Comparable uses Int#== internally)
		"==": f(
			func(
				env *object.Env, kwargs *object.PanObj, args ...object.PanObject,
			) object.PanObject {
				if len(args) < 2 {
					return object.NewTypeErr("== requires at least 2 args")
				}

				// necessary for Int itself! (guarantee `Int == Int`)
				if args[0] == object.BuiltInIntObj && args[1] == object.BuiltInIntObj {
					return object.BuiltInTrue
				}

				self, ok := object.TraceProtoOfInt(args[0])
				if !ok {
					return object.BuiltInFalse
				}
				other, ok := object.TraceProtoOfInt(args[1])
				if !ok {
					return object.BuiltInFalse
				}

				if self.Value == other.Value {
					return object.BuiltInTrue
				}
				return object.BuiltInFalse
			},
		),
		// NOTE: this cannot be removed (Comparable uses Int#!= internally)
		"!=": f(
			func(
				env *object.Env, kwargs *object.PanObj, args ...object.PanObject,
			) object.PanObject {
				if len(args) < 2 {
					return object.NewTypeErr("== requires at least 2 args")
				}

				// necessary for Int itself! (guarantee `Int == Int`)
				if args[0] == object.BuiltInIntObj && args[1] == object.BuiltInIntObj {
					return object.BuiltInTrue
				}

				self, ok := object.TraceProtoOfInt(args[0])
				if !ok {
					return object.BuiltInFalse
				}
				other, ok := object.TraceProtoOfInt(args[1])
				if !ok {
					return object.BuiltInFalse
				}

				if self.Value != other.Value {
					return object.BuiltInTrue
				}
				return object.BuiltInFalse
			},
		),
		"-%": f(
			func(
				env *object.Env, kwargs *object.PanObj, args ...object.PanObject,
			) object.PanObject {
				if len(args) < 1 {
					return object.NewTypeErr("\\- requires at least 1 arg")
				}

				self, ok := object.TraceProtoOfInt(args[0])
				if !ok {
					return object.NewTypeErr("\\1 must be int")
				}

				res := -self.Value
				return object.NewPanInt(res)
			},
		),
		"/~": f(
			func(
				env *object.Env, kwargs *object.PanObj, args ...object.PanObject,
			) object.PanObject {
				if len(args) < 1 {
					return object.NewTypeErr("/~ requires at least 1 arg")
				}

				self, ok := object.TraceProtoOfInt(args[0])
				if !ok {
					return object.NewTypeErr("\\1 must be int")
				}

				res := ^self.Value
				return object.NewPanInt(res)
			},
		),
		"+": f(
			func(
				env *object.Env, kwargs *object.PanObj, args ...object.PanObject,
			) object.PanObject {
				self, other, err := checkIntInfixArgs(args, "+", object.NewPanInt(0))
				if err != nil {
					return err
				}

				res := self.Value + other.Value
				return object.NewPanInt(res)
			},
		),
		"-": f(
			func(
				env *object.Env, kwargs *object.PanObj, args ...object.PanObject,
			) object.PanObject {
				self, other, err := checkIntInfixArgs(args, "-", object.NewPanInt(0))
				if err != nil {
					return err
				}

				res := self.Value - other.Value
				return object.NewPanInt(res)
			},
		),
		"*": f(
			func(
				env *object.Env, kwargs *object.PanObj, args ...object.PanObject,
			) object.PanObject {
				self, other, err := checkIntInfixArgs(args, "*", object.NewPanInt(1))
				if err != nil {
					return err
				}

				res := self.Value * other.Value
				return object.NewPanInt(res)
			},
		),
		"/": f(
			func(
				env *object.Env, kwargs *object.PanObj, args ...object.PanObject,
			) object.PanObject {
				self, other, err := checkIntInfixArgs(args, "/", object.NewPanInt(1))
				if err != nil {
					return err
				}

				if other.Value == 0 {
					return object.NewZeroDivisionErr("cannot be divided by 0")
				}

				// truediv
				res := float64(self.Value) / float64(other.Value)
				return &object.PanFloat{Value: res}
			},
		),
		"//": f(
			func(
				env *object.Env, kwargs *object.PanObj, args ...object.PanObject,
			) object.PanObject {
				self, other, err := checkIntInfixArgs(args, "//", object.NewPanInt(1))
				if err != nil {
					return err
				}

				if other.Value == 0 {
					return object.NewZeroDivisionErr("cannot be divided by 0")
				}

				// floordiv
				res := self.Value / other.Value

				// HACK: convert round to floor
				if res < 0 && self.Value%other.Value != 0 {
					return object.NewPanInt(res - 1)
				}

				return object.NewPanInt(res)
			},
		),
		"%": f(
			func(
				env *object.Env, kwargs *object.PanObj, args ...object.PanObject,
			) object.PanObject {
				self, other, err := checkIntInfixArgs(args, "%", object.NewPanInt(0))
				if err != nil {
					return err
				}

				if other.Value == 0 {
					return object.NewZeroDivisionErr("cannot be divided by 0")
				}

				res := self.Value % other.Value
				return object.NewPanInt(res)
			},
		),
		"_incBy": f(
			func(
				env *object.Env, kwargs *object.PanObj, args ...object.PanObject,
			) object.PanObject {
				self, other, err := checkIntInfixArgs(args, "_incBy", object.NewPanInt(0))
				if err != nil {
					return err
				}

				res := self.Value + other.Value
				return object.NewPanInt(res)
			},
		),
		"_iter": f(
			func(
				env *object.Env, kwargs *object.PanObj, args ...object.PanObject,
			) object.PanObject {
				if len(args) < 1 {
					return object.NewTypeErr("Int#_iter requires at least 1 arg")
				}

				self, ok := object.TraceProtoOfInt(args[0])
				if !ok {
					return object.NewTypeErr("\\1 must be int")
				}

				return &object.PanBuiltInIter{
					Fn:  intIter(self),
					Env: env, // not used
				}
			},
		),
		"at": propContainer["Int_at"],
		"B": f(
			func(
				env *object.Env, kwargs *object.PanObj, args ...object.PanObject,
			) object.PanObject {
				if len(args) < 1 {
					return object.NewTypeErr("Int#B requires at least 1 arg")
				}
				self, ok := object.TraceProtoOfInt(args[0])
				if !ok {
					return object.NewTypeErr(`\1 must be int`)
				}

				if self.Value == 0 {
					return object.BuiltInFalse
				}
				return object.BuiltInTrue
			},
		),
		"prime?": f(
			func(
				env *object.Env, kwargs *object.PanObj, args ...object.PanObject,
			) object.PanObject {
				if len(args) < 1 {
					return object.NewTypeErr("Int#prime? requires at least 1 arg")
				}
				self, ok := object.TraceProtoOfInt(args[0])
				if !ok {
					return object.NewTypeErr(`\1 must be int`)
				}

				if self.Value <= 1 {
					return object.BuiltInFalse
				}

				// NOTE: ProbablyPrime is 100% accurate if n < 2^64
				var n big.Int
				// self.Value must be positive
				n.SetUint64(uint64(self.Value))
				if ok := n.ProbablyPrime(0); ok {
					return object.BuiltInTrue
				}

				return object.BuiltInFalse
			},
		),
	}
}

func checkIntInfixArgs(
	args []object.PanObject,
	propName string,
	nilAs *object.PanInt,
) (*object.PanInt, *object.PanInt, *object.PanErr) {
	if len(args) < 2 {
		return nil, nil, object.NewTypeErr(propName + " requires at least 2 args")
	}

	self, ok := object.TraceProtoOfInt(args[0])
	if !ok {
		return nil, nil, object.NewTypeErr(
			fmt.Sprintf("`%s` cannot be treated as int", args[0].Inspect()))
	}
	other, ok := object.TraceProtoOfInt(args[1])
	if !ok {
		// NOTE: nil is treated as nilAs (0 in `+` and 1 in `*` for example)
		_, ok := object.TraceProtoOfNil(args[1])
		if ok {
			return self, nilAs, nil
		}

		return nil, nil, object.NewTypeErr(
			fmt.Sprintf("`%s` cannot be treated as int", args[1].Inspect()))
	}

	return self, other, nil
}

func intIter(i *object.PanInt) object.BuiltInFunc {
	yieldNum := int64(1)

	return func(
		env *object.Env, kwargs *object.PanObj, args ...object.PanObject,
	) object.PanObject {
		if yieldNum > i.Value {
			return object.NewStopIterErr("iter stopped")
		}
		yielded := object.NewPanInt(yieldNum)
		yieldNum++
		return yielded
	}
}
