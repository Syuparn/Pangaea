package props

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/dlclark/regexp2"
	"github.com/lithammer/dedent"

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

				// NOTE: Str's descendants also call this
				return object.NewInheritedStr(args[0].Proto(), string(negBytes))
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
				// NOTE: Str's descendants also call this
				return object.NewInheritedStr(args[0].Proto(), res)
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
						fmt.Sprintf("%s cannot be treated as str", args[0].Repr()))
				}

				nInt, ok := object.TraceProtoOfInt(args[1])
				if !ok {
					return object.NewTypeErr(
						fmt.Sprintf("%s cannot be treated as int", args[1].Repr()))
				}
				n := nInt.Value

				if n < 0 {
					return object.NewValueErr(
						fmt.Sprintf("%s is not positive", args[1].Repr()))
				}

				res := strings.Repeat(self.Value, int(n))
				// NOTE: Str's descendants also call this
				return object.NewInheritedStr(args[0].Proto(), res)
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
				//                        pattern     RegexOptions
				p, cerr := regexp2.Compile(sep.Value, 0)
				if cerr != nil {
					return object.NewValueErr(
						fmt.Sprintf("%s is invalid regex pattern", sep.Repr()))
				}
				// TODO: replace to Split() function when regexp2 implements it
				splitted, serr := splitBySep(p, self.Value)
				if serr != nil {
					return object.NewPanErr(
						fmt.Sprintf("unexpectedly failed to find match: %s", serr.Error()))
				}

				strs := []object.PanObject{}
				for _, s := range splitted {
					// exclude empty elems
					if s != "" {
						// NOTE: Str's descendants also call this
						strs = append(strs, object.NewInheritedStr(args[0].Proto(), s))
					}
				}

				return object.NewPanArr(strs...)
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

				// NOTE: Str's descendants also call this
				return object.NewInheritedStr(args[0].Proto(), string(newRunes))
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

				// TODO: enable to iterate args[0]'s brothers
				return object.NewPanBuiltInIter(strIter(self), env)
			},
		),
		"_name": object.NewPanStr("Str"),
		"at":    propContainer["Str_at"],
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
		"dedent": f(
			func(
				env *object.Env, kwargs *object.PanObj, args ...object.PanObject,
			) object.PanObject {
				if len(args) < 1 {
					return object.NewTypeErr("Str#dedent requires at least 1 arg")
				}
				self, ok := object.TraceProtoOfStr(args[0])
				if !ok {
					return object.NewTypeErr(
						fmt.Sprintf("%s cannot be treated as str", args[0].Repr()))
				}

				res := dedent.Dedent(self.Value)
				// NOTE: Str's descendants also call this
				return object.NewInheritedStr(args[0].Proto(), res)
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
						fmt.Sprintf("%s cannot be treated as str", args[0].Repr()))
				}

				f, err := strconv.ParseFloat(self.Value, 64)
				if err != nil {
					return object.NewValueErr(
						fmt.Sprintf("%s cannot be converted into float", args[0].Repr()))
				}
				return object.NewPanFloat(f)
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
						fmt.Sprintf("%s cannot be treated as str", args[0].Repr()))
				}

				b, ok := (*kwargs.Pairs)[object.GetSymHash("base")]
				// when kwarg `base` is not specified
				if !ok {
					return parseIntInBase(self, 10)
				}

				base, ok := object.TraceProtoOfInt(b.Value)
				if !ok {
					return object.NewTypeErr(
						fmt.Sprintf("%s cannot be treated as int", b.Value.Repr()))
				}

				return parseIntInBase(self, base.Value)
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

				// NOTE: Str's descendants also call this
				return object.NewInheritedStr(args[0].Proto(), strings.ToLower(self.Value))
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

				//                        pattern     RegexOptions
				p, err := regexp2.Compile(pattern.Value, 0)
				if err != nil {
					return object.NewValueErr(
						fmt.Sprintf("%s is invalid regex pattern", pattern.Repr()))
				}
				match, err := p.FindStringMatch(self.Value)
				if err != nil {
					return object.NewPanErr(
						fmt.Sprintf("unexpectedly failed to find match: %s", err.Error()))
				}

				if match == nil {
					return object.NewPanArr()
				}

				elems := []object.PanObject{}
				for _, group := range match.Groups() {
					// NOTE: Str's descendants also call this
					elems = append(elems, object.NewInheritedStr(args[0].Proto(), group.String()))
				}

				return object.NewPanArr(elems...)
			},
		),
		"new": f(
			func(
				env *object.Env, kwargs *object.PanObj, args ...object.PanObject,
			) object.PanObject {
				if len(args) < 2 {
					return object.NewTypeErr("Str#new requires at least 2 args")
				}
				str, ok := object.TraceProtoOfStr(args[1])
				if ok {
					// NOTE: Str's descendants also call this
					return object.NewInheritedStr(args[0], str.Value)
				}

				// if \2 is not str, return \2.S
				return propContainer["Obj_callProp"].(*object.PanBuiltIn).Fn(
					env, object.EmptyPanObjPtr(),
					object.EmptyPanObjPtr(), args[1], sSym,
				)
			},
		),
		"ord": f(
			func(
				env *object.Env, kwargs *object.PanObj, args ...object.PanObject,
			) object.PanObject {
				if len(args) < 1 {
					return object.NewTypeErr("Str#ord requires at least 1 arg")
				}
				self, ok := object.TraceProtoOfStr(args[0])
				if !ok {
					return object.NewTypeErr(
						fmt.Sprintf("%s cannot be treated as str", args[0].Repr()))
				}

				runes := []rune(self.Value)
				if len(runes) != 1 {
					return object.NewValueErr(
						fmt.Sprintf("length must be 1. got %d (%s)", len(runes), args[0].Repr()))
				}

				return object.NewPanInt(int64(runes[0]))
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

				//                        pattern     RegexOptions
				p, err := regexp2.Compile(pattern.Value, 0)
				if err != nil {
					return object.NewValueErr(
						fmt.Sprintf("%s is invalid regex pattern", pattern.Repr()))
				}

				replaced, err := p.Replace(self.Value, sub.Value, -1, -1)
				if err != nil {
					return object.NewValueErr(
						fmt.Sprintf("%s is invalid regex pattern", sub.Repr()))
				}
				// HACK: handle \U...\E (uppercase) and \L...E (lowercase)
				// TODO: remain (\U|\L)...\E in original string
				ucPattern := regexp2.MustCompile(`\\U(.*?)\\E`, 0)
				ucReplaced, err := ucPattern.ReplaceFunc(
					replaced, func(m regexp2.Match) string {
						return strings.ToUpper(m.Groups()[1].String())
					}, -1, -1)

				if err != nil {
					return object.NewValueErr(
						fmt.Sprintf("%s is invalid regex pattern", sub.Repr()))
				}

				lcPattern := regexp2.MustCompile(`\\L(.*?)\\E`, 0)
				lcReplaced, err := lcPattern.ReplaceFunc(
					ucReplaced, func(m regexp2.Match) string {
						return strings.ToLower(m.Groups()[1].String())
					}, -1, -1)

				if err != nil {
					return object.NewValueErr(
						fmt.Sprintf("%s is invalid regex pattern", sub.Repr()))
				}

				// NOTE: Str's descendants also call this
				return object.NewInheritedStr(args[0].Proto(), lcReplaced)
			},
		),
		"sym?": f(
			func(
				env *object.Env, kwargs *object.PanObj, args ...object.PanObject,
			) object.PanObject {
				if len(args) < 1 {
					return object.NewTypeErr("Str#sym? requires at least 1 arg")
				}
				self, ok := object.TraceProtoOfStr(args[0])
				if !ok {
					return object.NewTypeErr(
						fmt.Sprintf("%s cannot be treated as str",
							args[0].Repr()))
				}

				if self.IsSym {
					return object.BuiltInTrue
				}
				return object.BuiltInFalse
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

				// NOTE: Str's descendants also call this
				return object.NewInheritedStr(args[0].Proto(), strings.ToUpper(self.Value))
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
			fmt.Sprintf("%s cannot be treated as str", args[0].Repr()))
	}
	other, ok := object.TraceProtoOfStr(args[1])
	if !ok {
		// NOTE: nil is treated as ""
		_, ok := object.TraceProtoOfNil(args[1])
		if ok {
			return self, object.NewPanStr(""), nil
		}

		return nil, nil, object.NewTypeErr(
			fmt.Sprintf("%s cannot be treated as str", args[1].Repr()))
	}

	return self, other, nil
}

func parseIntInBase(s *object.PanStr, base int64) object.PanObject {
	if base < 2 || base > 36 {
		return object.NewValueErr(
			fmt.Sprintf("base %s must be within (2:37)", object.NewPanInt(base).Repr()))
	}

	i, err := strconv.ParseInt(s.Value, int(base), 64)
	if err != nil {
		return object.NewValueErr(
			fmt.Sprintf("%s cannot be converted into int", s.Repr()))
	}
	return object.NewPanInt(i)
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

// TODO: delete it when Split() is implemented in regexp2
func splitBySep(re *regexp2.Regexp, str string) ([]string, error) {
	runes := []rune(str)
	splitted := []string{}
	m, err := re.FindRunesMatch(runes)
	if err != nil {
		return nil, err
	}
	// if nothing is matched, just return original str
	if m == nil {
		return []string{str}, nil
	}

	sepFrom, sepTo := m.Index, (m.Index + m.Length)
	splitted = append(splitted, string(runes[0:sepFrom]))
	for m, err := re.FindNextMatch(m); m != nil; m, err = re.FindNextMatch(m) {
		if err != nil {
			return nil, err
		}
		sepFrom = m.Index
		splitted = append(splitted, string(runes[sepTo:sepFrom]))
		sepTo = m.Index + m.Length
	}
	splitted = append(splitted, string(runes[sepTo:]))

	return splitted, nil
}
