package props

import (
	"../object"
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
	}
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
