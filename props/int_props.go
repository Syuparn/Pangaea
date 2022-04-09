package props

import (
	"fmt"
	"math"
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

				// NOTE: Int's descendants also call this
				return object.NewInheritedInt(args[0].Proto(), res)
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

				self, ok := object.TraceProtoOfInt(args[0])
				if !ok {
					return object.BuiltInFalse
				}
				other, ok := object.TraceProtoOfInt(args[1])
				if !ok {
					return object.BuiltInFalse
				}

				if self.Proto() != other.Proto() {
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
					return object.NewTypeErr("!= requires at least 2 args")
				}

				self, ok := object.TraceProtoOfInt(args[0])
				if !ok {
					// cannot be handled
					return object.BuiltInFalse
				}
				other, ok := object.TraceProtoOfInt(args[1])
				if !ok {
					return object.BuiltInTrue
				}

				if self.Proto() != other.Proto() {
					return object.BuiltInTrue
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
				// NOTE: Int's descendants also call this
				return object.NewInheritedInt(args[0].Proto(), res)
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
				// NOTE: Int's descendants also call this
				return object.NewInheritedInt(args[0].Proto(), res)
			},
		),
		"+": f(
			func(
				env *object.Env, kwargs *object.PanObj, args ...object.PanObject,
			) object.PanObject {
				self, other, err := checkIntInfixArgs(args, "+", object.NewPanInt(0))
				if err == nil {
					// NOTE: Int's descendants also call this
					return object.NewInheritedInt(args[0].Proto(), self.Value+other.Value)
				}

				if fself, fother, ferr := checkFloatInfixArgs(args, "+", object.NewPanFloat(0)); ferr == nil {
					return object.NewPanFloat(fself.Value + fother.Value)
				}

				return err
			},
		),
		"-": f(
			func(
				env *object.Env, kwargs *object.PanObj, args ...object.PanObject,
			) object.PanObject {
				self, other, err := checkIntInfixArgs(args, "-", object.NewPanInt(0))
				if err == nil {
					// NOTE: Int's descendants also call this
					return object.NewInheritedInt(args[0].Proto(), self.Value-other.Value)
				}

				if fself, fother, ferr := checkFloatInfixArgs(args, "-", object.NewPanFloat(0)); ferr == nil {
					return object.NewPanFloat(fself.Value - fother.Value)
				}

				return err
			},
		),
		"*": f(
			func(
				env *object.Env, kwargs *object.PanObj, args ...object.PanObject,
			) object.PanObject {
				self, other, err := checkIntInfixArgs(args, "*", object.NewPanInt(1))
				if err == nil {
					// NOTE: Int's descendants also call this
					return object.NewInheritedInt(args[0].Proto(), self.Value*other.Value)
				}

				if fself, fother, ferr := checkFloatInfixArgs(args, "*", object.NewPanFloat(1)); ferr == nil {
					return object.NewPanFloat(fself.Value * fother.Value)
				}

				return err
			},
		),
		"**": f(
			func(
				env *object.Env, kwargs *object.PanObj, args ...object.PanObject,
			) object.PanObject {
				self, other, err := checkIntInfixArgs(args, "**", object.NewPanInt(1))
				if err == nil {
					res := math.Pow(float64(self.Value), float64(other.Value))
					// check if f is integer
					if math.Floor(res) == res {
						// NOTE: Int's descendants also call this
						return object.NewInheritedInt(args[0].Proto(), int64(res))
					}
					return object.NewPanFloat(res)
				}

				if fself, fother, ferr := checkFloatInfixArgs(args, "**", object.NewPanFloat(1)); ferr == nil {
					return object.NewPanFloat(math.Pow(fself.Value, fother.Value))
				}

				return err
			},
		),
		"/": f(
			func(
				env *object.Env, kwargs *object.PanObj, args ...object.PanObject,
			) object.PanObject {
				self, other, err := checkIntInfixArgs(args, "/", object.NewPanInt(1))
				if err == nil {
					if other.Value == 0 {
						return object.NewZeroDivisionErr("cannot be divided by 0")
					}
					// truediv
					return object.NewPanFloat(float64(self.Value) / float64(other.Value))
				}

				if fself, fother, ferr := checkFloatInfixArgs(args, "/", object.NewPanFloat(1)); ferr == nil {
					if fother.Value == 0 {
						return object.NewZeroDivisionErr("cannot be divided by 0")
					}
					return object.NewPanFloat(fself.Value / fother.Value)
				}

				return err
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
					// NOTE: Int's descendants also call this
					return object.NewInheritedInt(args[0].Proto(), res-1)
				}

				// NOTE: Int's descendants also call this
				return object.NewInheritedInt(args[0].Proto(), res)
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
				// NOTE: Int's descendants also call this
				return object.NewInheritedInt(args[0].Proto(), res)
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
				// NOTE: Int's descendants also call this
				return object.NewInheritedInt(args[0].Proto(), res)
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

				return object.NewPanBuiltInIter(intIter(self), env)
			},
		),
		"_name": object.NewPanStr("Int"),
		"at":    propContainer["Int_at"],
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
		// NOTE: override bear to set arr zero values
		"bear": f(
			func(
				env *object.Env, kwargs *object.PanObj, args ...object.PanObject,
			) object.PanObject {
				if len(args) < 1 {
					return object.NewTypeErr("Int#bear requires at least 1 arg")
				}
				proto := args[0]

				// since true and false are singletons and must not be inherited, bear returns self
				if proto == object.BuiltInTrue || proto == object.BuiltInFalse {
					return proto
				}

				// if soruce is not specified as arg, use {}
				src := object.EmptyPanObjPtr()
				if len(args) >= 2 {
					// `proto.bear(arg)`
					arg, ok := args[1].(*object.PanObj)
					if !ok {
						return object.NewTypeErr("Int#bear requires obj literal src")
					}
					src = arg
				}

				// if self is concrete value, use it as zero value
				if proto.Type() == object.IntType {
					return object.ChildPanObjPtr(proto, src, object.WithZero(proto))
				}

				// otherwise zero value is 0 inherited from the generated obj
				return object.ChildPanObjPtr(proto, src,
					object.WithZeroFromSelf(func(o *object.PanObj) object.PanObject {
						return object.NewInheritedInt(o, 0)
					}))
			},
		),
		"chr": f(
			func(
				env *object.Env, kwargs *object.PanObj, args ...object.PanObject,
			) object.PanObject {
				if len(args) < 1 {
					return object.NewTypeErr("Int#chr requires at least 1 arg")
				}
				self, ok := object.TraceProtoOfInt(args[0])
				if !ok {
					return object.NewTypeErr(
						fmt.Sprintf("%s cannot be treated as int", args[0].Repr()))
				}

				// NOTE: convert int64 to string via rune to tell the conversion is intentional
				// otherwise test warns below:
				//    conversion from int64 to string yields a string of one rune, not a string of digits (did you mean fmt.Sprint(x)?)
				return object.NewPanStr(string(rune(self.Value)))
			},
		),
		"new": f(
			func(
				env *object.Env, kwargs *object.PanObj, args ...object.PanObject,
			) object.PanObject {
				if len(args) < 2 {
					return object.NewTypeErr("Int#new requires at least 2 args")
				}

				proto := args[0]
				// since true and false are singletons and must not be inherited, new inherits Int
				if proto == object.BuiltInTrue || proto == object.BuiltInFalse {
					proto = object.BuiltInIntObj
				}

				i, ok := object.TraceProtoOfInt(args[1])
				if ok {

					// NOTE: Int's descendants also call this
					return object.NewInheritedInt(proto, int64(i.Value))
				}

				// float is rounded down
				f, ok := object.TraceProtoOfFloat(args[1])
				if ok {
					// NOTE: Int's descendants also call this
					return object.NewInheritedInt(proto, int64(f.Value))
				}

				return object.NewTypeErr(
					fmt.Sprintf("%s cannot be treated as int", args[1].Repr()))
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
		"sqrt": f(
			func(
				env *object.Env, kwargs *object.PanObj, args ...object.PanObject,
			) object.PanObject {
				if len(args) < 1 {
					return object.NewTypeErr("Int#sqrt requires at least 1 arg")
				}
				self, ok := object.TraceProtoOfInt(args[0])
				if !ok {
					return object.NewTypeErr(
						fmt.Sprintf("%s cannot be treated as int", args[0].Repr()))
				}

				if self.Value < 0 {
					return object.NewValueErr(
						fmt.Sprintf("sqrt of %s is not a real number", self.Repr()))
				}

				return object.NewPanFloat(math.Sqrt(float64(self.Value)))
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
			fmt.Sprintf("%s cannot be treated as int", args[0].Repr()))
	}
	other, ok := object.TraceProtoOfInt(args[1])
	if !ok {
		// NOTE: nil is treated as nilAs (0 in `+` and 1 in `*` for example)
		if _, ok := object.TraceProtoOfNil(args[1]); ok {
			return self, nilAs, nil
		}

		return nil, nil, object.NewTypeErr(
			fmt.Sprintf("%s cannot be treated as int", args[1].Repr()))
	}

	return self, other, nil
}

func formattedIntStr(self *object.PanInt, kwargs *object.PanObj) object.PanObject {
	// if kwarg `base` is specified
	if b, ok := (*kwargs.Pairs)[object.GetSymHash("base")]; ok {
		base, ok := object.TraceProtoOfInt(b.Value)
		if !ok {
			return object.NewTypeErr(
				fmt.Sprintf("%s cannot be treated as int", b.Value.Repr()))
		}

		return baseString(self, base)
	}

	return object.NewPanStr(self.Inspect())
}

func baseString(self, base *object.PanInt) object.PanObject {
	// NOTE: do not handle base > 36 because Str#I (strconv.Itoa) cannot re-convert it
	if base.Value < 2 || base.Value > 36 {
		return object.NewValueErr(
			fmt.Sprintf("base %s must be within (2:37)", base.Repr()))
	}
	return object.NewPanStr(big.NewInt(self.Value).Text(int(base.Value)))
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
