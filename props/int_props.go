package props

import (
	"../object"
	"fmt"
)

func IntProps(propContainer map[string]object.PanObject) map[string]object.PanObject {
	// NOTE: inject some built-in functions which relate to parser or evaluator
	return map[string]object.PanObject{
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

				self, ok := traceProtoOf(args[0], isInt)
				if !ok {
					return object.BuiltInFalse
				}
				other, ok := traceProtoOf(args[1], isInt)
				if !ok {
					return object.BuiltInFalse
				}

				if self.(*object.PanInt).Value == other.(*object.PanInt).Value {
					return object.BuiltInTrue
				}
				return object.BuiltInFalse
			},
		),
		"+": f(
			func(
				env *object.Env, kwargs *object.PanObj, args ...object.PanObject,
			) object.PanObject {
				self, other, err := checkIntInfixArgs(args)
				if err != nil {
					return err
				}

				res := self.(*object.PanInt).Value + other.(*object.PanInt).Value
				return object.NewPanInt(res)
			},
		),
		"*": f(
			func(
				env *object.Env, kwargs *object.PanObj, args ...object.PanObject,
			) object.PanObject {
				self, other, err := checkIntInfixArgs(args)
				if err != nil {
					return err
				}

				res := self.(*object.PanInt).Value * other.(*object.PanInt).Value
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

				self, ok := traceProtoOf(args[0], isInt)
				if !ok {
					return object.NewTypeErr("\\1 must be int")
				}

				return &object.PanBuiltInIter{
					Fn:  intIter(self.(*object.PanInt)),
					Env: env, // not used
				}
			},
		),
	}
}

func checkIntInfixArgs(
	args []object.PanObject,
) (object.PanObject, object.PanObject, *object.PanErr) {
	if len(args) < 2 {
		return nil, nil, object.NewTypeErr("+ requires at least 2 args")
	}

	self, ok := traceProtoOf(args[0], isInt)
	if !ok {
		return nil, nil, object.NewTypeErr(
			fmt.Sprintf("`%s` cannot be treated as int", args[0].Inspect()))
	}
	other, ok := traceProtoOf(args[1], isInt)
	if !ok {
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
		yieldNum += 1
		return yielded
	}
}
