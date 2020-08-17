package props

import (
	"../object"
	"fmt"
)

func StrProps(propContainer map[string]object.PanObject) map[string]object.PanObject {
	// NOTE: inject some built-in functions which relate to parser or evaluator
	return map[string]object.PanObject{
		"==": f(
			func(
				env *object.Env, kwargs *object.PanObj, args ...object.PanObject,
			) object.PanObject {
				if len(args) < 2 {
					return object.NewTypeErr("== requires at least 2 args")
				}

				// necessary for Str itself! (guarantee `Str == Str`)
				if args[0] == object.BuiltInStrObj && args[1] == object.BuiltInStrObj {
					return object.BuiltInTrue
				}

				self, ok := args[0].(*object.PanStr)
				if !ok {
					return object.BuiltInFalse
				}
				other, ok := args[1].(*object.PanStr)
				if !ok {
					return object.BuiltInFalse
				}

				if self.Hash() == other.Hash() {
					return object.BuiltInTrue
				}

				return object.BuiltInFalse
			},
		),
		"/~": f(
			func(
				env *object.Env, kwargs *object.PanObj, args ...object.PanObject,
			) object.PanObject {
				if len(args) < 1 {
					return object.NewTypeErr("/~ requires at least 1 arg")
				}

				self, ok := traceProtoOf(args[0], isStr)
				if !ok {
					return object.NewTypeErr("\\1 must be str")
				}

				strBytes := []byte(self.(*object.PanStr).Value)
				negBytes := []byte{}
				for _, b := range strBytes {
					negBytes = append(negBytes, ^b)
				}

				return &object.PanStr{Value: string(negBytes)}
			},
		),
		"+": f(
			func(
				env *object.Env, kwargs *object.PanObj, args ...object.PanObject,
			) object.PanObject {
				self, other, err := checkStrInfixArgs(args, "+")
				if err != nil {
					return err
				}

				res := self.(*object.PanStr).Value + other.(*object.PanStr).Value
				return &object.PanStr{Value: res}
			},
		),
		"_iter": f(
			func(
				env *object.Env, kwargs *object.PanObj, args ...object.PanObject,
			) object.PanObject {
				if len(args) < 1 {
					return object.NewTypeErr("Str#_iter requires at least 1 arg")
				}

				self, ok := traceProtoOf(args[0], isStr)
				if !ok {
					return object.NewTypeErr("\\1 must be int")
				}

				return &object.PanBuiltInIter{
					Fn:  strIter(self.(*object.PanStr)),
					Env: env, // not used
				}
			},
		),
		"at": propContainer["Str_at"],
		"B": f(
			func(
				env *object.Env, kwargs *object.PanObj, args ...object.PanObject,
			) object.PanObject {
				if len(args) < 1 {
					return object.NewTypeErr("Str#B requires at least 1 arg")
				}
				self, ok := traceProtoOf(args[0], isStr)
				if !ok {
					return object.NewTypeErr(`\1 must be str`)
				}

				if self.(*object.PanStr).Value == "" {
					return object.BuiltInFalse
				}
				return object.BuiltInTrue
			},
		),
	}
}

func checkStrInfixArgs(
	args []object.PanObject,
	propName string,
) (object.PanObject, object.PanObject, *object.PanErr) {
	if len(args) < 2 {
		return nil, nil, object.NewTypeErr("== requires at least 2 args")
	}

	self, ok := args[0].(*object.PanStr)
	if !ok {
		return nil, nil, object.NewTypeErr(
			fmt.Sprintf("`%s` cannot be treated as int", args[0].Inspect()))
	}
	other, ok := args[1].(*object.PanStr)
	if !ok {
		return nil, nil, object.NewTypeErr(
			fmt.Sprintf("`%s` cannot be treated as int", args[0].Inspect()))
	}

	return self, other, nil
}

func strIter(s *object.PanStr) object.BuiltInFunc {
	yieldIdx := 0
	runes := []rune(s.Value)

	return func(
		env *object.Env, kwargs *object.PanObj, args ...object.PanObject,
	) object.PanObject {
		if yieldIdx >= len(runes) {
			return object.NewStopIterErr("iter stopped")
		}
		yielded := &object.PanStr{Value: string(runes[yieldIdx])}
		yieldIdx += 1
		return yielded
	}
}
