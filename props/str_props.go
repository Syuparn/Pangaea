package props

import (
	"../object"
	"fmt"
	"strings"
)

func StrProps(propContainer map[string]object.PanObject) map[string]object.PanObject {
	// NOTE: inject some built-in functions which relate to parser or evaluator
	return map[string]object.PanObject{
		"<=>": f(
			func(
				env *object.Env, kwargs *object.PanObj, args ...object.PanObject,
			) object.PanObject {
				self, other, err := checkStrInfixArgs(args, "<=>")
				if err != nil {
					return err
				}

				selfVal := self.(*object.PanStr).Value
				otherVal := other.(*object.PanStr).Value

				return object.NewPanInt(int64(strings.Compare(selfVal, otherVal)))
			},
		),
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

				self, ok := object.TraceProtoOfStr(args[0])
				if !ok {
					return object.NewTypeErr("\\1 must be str")
				}

				strBytes := []byte(self.Value)
				negBytes := []byte{}
				for _, b := range strBytes {
					negBytes = append(negBytes, ^b)
				}

				return object.NewPanStr(string(negBytes))
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
				return object.NewPanStr(res)
			},
		),
		"_incBy": f(
			func(
				env *object.Env, kwargs *object.PanObj, args ...object.PanObject,
			) object.PanObject {
				if len(args) < 2 {
					return object.NewTypeErr("Str#_incBy requires at least 2 args")
				}

				self, ok := object.TraceProtoOfStr(args[0])
				if !ok {
					return object.NewTypeErr("\\1 must be str")
				}

				nInt, ok := object.TraceProtoOfInt(args[1])
				if !ok {
					return object.NewTypeErr("\\2 must be int")
				}
				n := nInt.Value

				runes := []rune(self.Value)
				increasedRune := runes[len(runes)-1] + rune(n)
				newRunes := append(runes[0:len(runes)-1], increasedRune)

				return object.NewPanStr(string(newRunes))
			},
		),
		"_iter": f(
			func(
				env *object.Env, kwargs *object.PanObj, args ...object.PanObject,
			) object.PanObject {
				if len(args) < 1 {
					return object.NewTypeErr("Str#_iter requires at least 1 arg")
				}

				self, ok := object.TraceProtoOfStr(args[0])
				if !ok {
					return object.NewTypeErr("\\1 must be int")
				}

				return &object.PanBuiltInIter{
					Fn:  strIter(self),
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
				self, ok := object.TraceProtoOfStr(args[0])
				if !ok {
					return object.NewTypeErr(`\1 must be str`)
				}

				if self.Value == "" {
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
		yielded := object.NewPanStr(string(runes[yieldIdx]))
		yieldIdx += 1
		return yielded
	}
}
