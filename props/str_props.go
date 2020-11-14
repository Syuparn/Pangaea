package props

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"

	"github.com/Syuparn/pangaea/object"
)

// StrProps provides built-in props for Str.
// NOTE: Some Str props are defind by native code (not by this function).
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

				return object.NewPanInt(int64(strings.Compare(self.Value, other.Value)))
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

				res := self.Value + other.Value
				return object.NewPanStr(res)
			},
		),
		"*": f(
			func(
				env *object.Env, kwargs *object.PanObj, args ...object.PanObject,
			) object.PanObject {
				if len(args) < 2 {
					return object.NewTypeErr("* requires at least 2 args")
				}

				self, ok := object.TraceProtoOfStr(args[0])
				if !ok {
					return object.NewTypeErr(
						fmt.Sprintf("`%s` cannot be treated as str", args[0].Inspect()))
				}

				nInt, ok := object.TraceProtoOfInt(args[1])
				if !ok {
					return object.NewTypeErr(
						fmt.Sprintf("`%s` cannot be treated as int", args[1].Inspect()))
				}
				n := nInt.Value

				if n < 0 {
					return object.NewValueErr(
						fmt.Sprintf("`%s` is not positive", args[1].Inspect()))
				}

				return object.NewPanStr(strings.Repeat(self.Value, int(n)))
			},
		),
		"/": f(
			func(
				env *object.Env, kwargs *object.PanObj, args ...object.PanObject,
			) object.PanObject {
				self, sep, err := checkStrInfixArgs(args, "/")
				if err != nil {
					return err
				}
				splitted := regexp.MustCompile(sep.Value).Split(self.Value, -1)
				strs := []object.PanObject{}
				for _, s := range splitted {
					// exclude empty elems
					if s != "" {
						strs = append(strs, object.NewPanStr(s))
					}
				}

				return &object.PanArr{Elems: strs}
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
		"eval":    propContainer["Str_eval"],
		"evalEnv": propContainer["Str_evalEnv"],
		"F": f(
			func(
				env *object.Env, kwargs *object.PanObj, args ...object.PanObject,
			) object.PanObject {
				if len(args) < 1 {
					return object.NewTypeErr("Str#F requires at least 1 arg")
				}
				self, ok := object.TraceProtoOfStr(args[0])
				if !ok {
					return object.NewTypeErr(
						fmt.Sprintf("`%s` cannot be treated as str", args[0].Inspect()))
				}

				f, err := strconv.ParseFloat(self.Value, 64)
				if err != nil {
					return object.NewValueErr(
						fmt.Sprintf("`%s` cannot be converted into float", args[0].Inspect()))
				}
				return &object.PanFloat{Value: f}
			},
		),
		"I": f(
			func(
				env *object.Env, kwargs *object.PanObj, args ...object.PanObject,
			) object.PanObject {
				if len(args) < 1 {
					return object.NewTypeErr("Str#I requires at least 1 arg")
				}
				self, ok := object.TraceProtoOfStr(args[0])
				if !ok {
					return object.NewTypeErr(
						fmt.Sprintf("`%s` cannot be treated as str", args[0].Inspect()))
				}

				i, err := strconv.ParseInt(self.Value, 10, 64)
				if err != nil {
					return object.NewValueErr(
						fmt.Sprintf("`%s` cannot be converted into int", args[0].Inspect()))
				}
				return object.NewPanInt(i)
			},
		),
		"lc": f(
			func(
				env *object.Env, kwargs *object.PanObj, args ...object.PanObject,
			) object.PanObject {
				if len(args) < 1 {
					return object.NewTypeErr("Str#lc requires at least 1 arg")
				}
				self, ok := object.TraceProtoOfStr(args[0])
				if !ok {
					return object.NewTypeErr(`\1 must be str`)
				}

				return object.NewPanStr(strings.ToLower(self.Value))
			},
		),
		"len": f(
			func(
				env *object.Env, kwargs *object.PanObj, args ...object.PanObject,
			) object.PanObject {
				if len(args) < 1 {
					return object.NewTypeErr("Str#len requires at least 1 arg")
				}
				self, ok := object.TraceProtoOfStr(args[0])
				if !ok {
					return object.NewTypeErr(`\1 must be str`)
				}

				return object.NewPanInt(int64(len([]rune(self.Value))))
			},
		),
		"match": f(
			func(
				env *object.Env, kwargs *object.PanObj, args ...object.PanObject,
			) object.PanObject {
				if len(args) < 2 {
					return object.NewTypeErr("Str#match requires at least 2 args")
				}
				self, ok := object.TraceProtoOfStr(args[0])
				if !ok {
					return object.NewTypeErr(`\1 must be str`)
				}
				pattern, ok := object.TraceProtoOfStr(args[1])
				if !ok {
					return object.NewTypeErr(`\2 must be str`)
				}

				found := regexp.MustCompile(pattern.Value).FindStringSubmatch(self.Value)
				elems := []object.PanObject{}
				for _, s := range found {
					elems = append(elems, object.NewPanStr(s))
				}

				return &object.PanArr{Elems: elems}
			},
		),
		"sub": f(
			func(
				env *object.Env, kwargs *object.PanObj, args ...object.PanObject,
			) object.PanObject {
				if len(args) < 3 {
					return object.NewTypeErr("Str#sub requires at least 3 args")
				}
				self, ok := object.TraceProtoOfStr(args[0])
				if !ok {
					return object.NewTypeErr(`\1 must be str`)
				}
				pattern, ok := object.TraceProtoOfStr(args[1])
				if !ok {
					return object.NewTypeErr(`\2 must be str`)
				}
				sub, ok := object.TraceProtoOfStr(args[2])
				if !ok {
					return object.NewTypeErr(`\3 must be str`)
				}

				replaced := regexp.MustCompile(pattern.Value).
					ReplaceAllString(self.Value, sub.Value)

				return object.NewPanStr(replaced)
			},
		),
		"uc": f(
			func(
				env *object.Env, kwargs *object.PanObj, args ...object.PanObject,
			) object.PanObject {
				if len(args) < 1 {
					return object.NewTypeErr("Str#uc requires at least 1 arg")
				}
				self, ok := object.TraceProtoOfStr(args[0])
				if !ok {
					return object.NewTypeErr(`\1 must be str`)
				}

				return object.NewPanStr(strings.ToUpper(self.Value))
			},
		),
	}
}

func checkStrInfixArgs(
	args []object.PanObject,
	propName string,
) (*object.PanStr, *object.PanStr, *object.PanErr) {
	if len(args) < 2 {
		return nil, nil, object.NewTypeErr(propName + " requires at least 2 args")
	}

	self, ok := object.TraceProtoOfStr(args[0])
	if !ok {
		return nil, nil, object.NewTypeErr(
			fmt.Sprintf("`%s` cannot be treated as str", args[0].Inspect()))
	}
	other, ok := object.TraceProtoOfStr(args[1])
	if !ok {
		// NOTE: nil is treated as ""
		_, ok := object.TraceProtoOfNil(args[1])
		if ok {
			return self, object.NewPanStr(""), nil
		}

		return nil, nil, object.NewTypeErr(
			fmt.Sprintf("`%s` cannot be treated as str", args[1].Inspect()))
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
		yieldIdx++
		return yielded
	}
}
