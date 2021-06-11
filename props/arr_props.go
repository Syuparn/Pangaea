package props

import (
	"fmt"
	"strings"

	"github.com/Syuparn/pangaea/object"
)

// ArrProps provides built-in props for Arr.
// NOTE: Some Arr props are defind by native code (not by this function).
func ArrProps(propContainer map[string]object.PanObject) map[string]object.PanObject {
	// NOTE: inject some built-in functions which relate to parser or evaluator
	return map[string]object.PanObject{
		"==": f(
			func(
				env *object.Env, kwargs *object.PanObj, args ...object.PanObject,
			) object.PanObject {
				if len(args) < 2 {
					return object.NewTypeErr("== requires at least 2 args")
				}

				self, ok := object.TraceProtoOfArr(args[0])
				if !ok {
					return object.BuiltInFalse
				}
				other, ok := object.TraceProtoOfArr(args[1])
				if !ok {
					return object.BuiltInFalse
				}

				return compArrs(self, other, propContainer, env)
			},
		),
		"+": f(
			func(
				env *object.Env, kwargs *object.PanObj, args ...object.PanObject,
			) object.PanObject {
				if len(args) < 2 {
					return object.NewTypeErr("+ requires at least 2 args")
				}

				self, ok := object.TraceProtoOfArr(args[0])
				if !ok {
					return object.NewTypeErr(
						fmt.Sprintf("%s cannot be treated as arr", args[0].Repr()))
				}
				other, ok := object.TraceProtoOfArr(args[1])
				if !ok {
					// NOTE: nil is treated as []
					_, ok := object.TraceProtoOfNil(args[1])
					if ok {
						return self
					}

					return object.NewTypeErr(
						fmt.Sprintf("%s cannot be treated as arr",
							args[1].Repr()))
				}

				// NOTE: no need to copy each elem because they are immutable
				elems := append(self.Elems, other.Elems...)
				return object.NewPanArr(elems...)
			},
		),
		"*": f(
			func(
				env *object.Env, kwargs *object.PanObj, args ...object.PanObject,
			) object.PanObject {
				if len(args) < 2 {
					return object.NewTypeErr("* requires at least 2 args")
				}

				self, ok := object.TraceProtoOfArr(args[0])
				if !ok {
					return object.NewTypeErr(
						fmt.Sprintf("%s cannot be treated as arr",
							args[0].Repr()))
				}
				selfElems := self.Elems

				other, ok := object.TraceProtoOfInt(args[1])
				if !ok {
					return object.NewTypeErr(
						fmt.Sprintf("%s cannot be treated as int", args[0].Repr()))
				}

				// NOTE: no need to copy each elem because they are immutable
				elems := []object.PanObject{}
				for i := int64(0); i < other.Value; i++ {
					elems = append(elems, selfElems...)
				}
				return object.NewPanArr(elems...)
			},
		),
		"_iter": f(
			func(
				env *object.Env, kwargs *object.PanObj, args ...object.PanObject,
			) object.PanObject {
				if len(args) < 1 {
					return object.NewTypeErr("Arr#_iter requires at least 1 arg")
				}

				self, ok := object.TraceProtoOfArr(args[0])
				if !ok {
					return object.NewTypeErr("\\1 must be arr")
				}

				return object.NewPanBuiltInIter(arrIter(self), env)
			},
		),
		"_name": object.NewPanStr("Arr"),
		"at":    propContainer["Arr_at"],
		"B": f(
			func(
				env *object.Env, kwargs *object.PanObj, args ...object.PanObject,
			) object.PanObject {
				if len(args) < 1 {
					return object.NewTypeErr("Arr#B requires at least 1 arg")
				}
				self, ok := object.TraceProtoOfArr(args[0])
				if !ok {
					return object.NewTypeErr(`\1 must be arr`)
				}

				if len(self.Elems) == 0 {
					return object.BuiltInFalse
				}
				return object.BuiltInTrue
			},
		),
		"call": f(
			func(
				env *object.Env, kwargs *object.PanObj, args ...object.PanObject,
			) object.PanObject {
				if len(args) < 1 {
					return object.NewTypeErr("Arr#call requires at least 1 arg")
				}

				return object.NewPanArr(args[1:]...)
			},
		),
		"has?": f(
			func(
				env *object.Env, kwargs *object.PanObj, args ...object.PanObject,
			) object.PanObject {
				if len(args) < 2 {
					return object.NewTypeErr("Arr#has? requires at least 2 args")
				}
				self, ok := object.TraceProtoOfArr(args[0])
				if !ok {
					return object.NewTypeErr(`\1 must be arr`)
				}

				for _, elem := range self.Elems {
					// == comparison
					isEq := propContainer["Obj_callProp"].(*object.PanBuiltIn).Fn(
						env, object.EmptyPanObjPtr(),
						object.EmptyPanObjPtr(), elem, eqSym, args[1],
					)
					if isEq == object.BuiltInTrue {
						return object.BuiltInTrue
					}
				}

				return object.BuiltInFalse
			},
		),
		"join": f(
			func(
				env *object.Env, kwargs *object.PanObj, args ...object.PanObject,
			) object.PanObject {
				if len(args) < 2 {
					return object.NewTypeErr("Arr#join requires at least 2 args")
				}
				self, ok := object.TraceProtoOfArr(args[0])
				if !ok {
					return object.NewTypeErr(`\1 must be arr`)
				}
				joiner, ok := object.TraceProtoOfStr(args[1])
				if !ok {
					return object.NewTypeErr(`\2 must be str`)
				}

				elemStrs := []string{}
				for _, elem := range self.Elems {
					// stringify each elem
					s := propContainer["Obj_callProp"].(*object.PanBuiltIn).Fn(
						env, object.EmptyPanObjPtr(),
						object.EmptyPanObjPtr(), elem, sSym, args[1],
					)
					if s.Type() == object.ErrType {
						return s
					}
					str, ok := object.TraceProtoOfStr(s)
					if !ok {
						return object.NewTypeErr(
							fmt.Sprintf("%s.S returned non-str value %s",
								elem.Repr(), s.Repr()))
					}
					elemStrs = append(elemStrs, str.Value)
				}
				res := strings.Join(elemStrs, joiner.Value)

				return object.NewPanStr(res)
			},
		),
		"len": f(
			func(
				env *object.Env, kwargs *object.PanObj, args ...object.PanObject,
			) object.PanObject {
				if len(args) < 1 {
					return object.NewTypeErr("Arr#len requires at least 1 arg")
				}
				self, ok := object.TraceProtoOfArr(args[0])
				if !ok {
					return object.NewTypeErr(`\1 must be arr`)
				}

				return object.NewPanInt(int64(len(self.Elems)))
			},
		),
		"M": f(
			func(
				env *object.Env, kwargs *object.PanObj, args ...object.PanObject,
			) object.PanObject {
				if len(args) < 1 {
					return object.NewTypeErr("Arr#M requires at least 1 arg")
				}
				self, ok := object.TraceProtoOfArr(args[0])
				if !ok {
					return object.NewTypeErr(
						fmt.Sprintf("%s cannot be treated as arr", args[0].Repr()))
				}

				pairs := make([]object.Pair, len(self.Elems))

				for i, e := range self.Elems {
					arr, ok := object.TraceProtoOfArr(e)
					if !ok {
						return object.NewValueErr(
							fmt.Sprintf(`element %s cannot be treated as arr`, e.Repr()))
					}

					if len(arr.Elems) != 2 {
						return object.NewValueErr(
							fmt.Sprintf(`element %s must have two elements`, arr.Repr()))
					}

					pairs[i] = object.Pair{
						Key:   arr.Elems[0],
						Value: arr.Elems[1],
					}
				}

				return object.NewPanMap(pairs...)
			},
		),
		"new": f(
			func(
				env *object.Env, kwargs *object.PanObj, args ...object.PanObject,
			) object.PanObject {
				if len(args) < 2 {
					return object.NewTypeErr("Arr#new requires at least 2 args")
				}
				arr, ok := object.TraceProtoOfArr(args[1])
				if !ok {
					return object.NewTypeErr(
						fmt.Sprintf("%s cannot be treated as arr", args[1].Repr()))
				}

				return arr
			},
		),
		"O": f(
			func(
				env *object.Env, kwargs *object.PanObj, args ...object.PanObject,
			) object.PanObject {
				if len(args) < 1 {
					return object.NewTypeErr("Arr#O requires at least 1 arg")
				}
				self, ok := object.TraceProtoOfArr(args[0])
				if !ok {
					return object.NewTypeErr(
						fmt.Sprintf("%s cannot be treated as arr", args[0].Repr()))
				}

				pairs := map[object.SymHash]object.Pair{}

				for _, e := range self.Elems {
					arr, ok := object.TraceProtoOfArr(e)
					if !ok {
						return object.NewValueErr(
							fmt.Sprintf(`element %s cannot be treated as arr`, e.Repr()))
					}

					if len(arr.Elems) != 2 {
						return object.NewValueErr(
							fmt.Sprintf(`element %s must have two elements`, arr.Repr()))
					}

					k, ok := object.TraceProtoOfStr(arr.Elems[0])
					if !ok {
						return object.NewValueErr(
							fmt.Sprintf(`element key %s cannot be treated as str`, arr.Elems[0].Repr()))
					}

					pairs[object.GetSymHash(k.Value)] = object.Pair{
						Key:   k,
						Value: arr.Elems[1],
					}
				}

				return object.PanObjInstancePtr(&pairs)
			},
		),
	}
}

func compArrs(
	a1 *object.PanArr,
	a2 *object.PanArr,
	propContainer map[string]object.PanObject,
	env *object.Env,
) object.PanObject {
	if len(a1.Elems) != len(a2.Elems) {
		return object.BuiltInFalse
	}

	for i, e := range a1.Elems {
		// == comparison for both elements
		res := propContainer["Obj_callProp"].(*object.PanBuiltIn).Fn(
			env, object.EmptyPanObjPtr(),
			object.EmptyPanObjPtr(), e, eqSym, a2.Elems[i],
		)
		if res == object.BuiltInFalse {
			return object.BuiltInFalse
		}
	}
	return object.BuiltInTrue
}

func arrIter(arr *object.PanArr) object.BuiltInFunc {
	yieldIdx := 0

	return func(
		env *object.Env, kwargs *object.PanObj, args ...object.PanObject,
	) object.PanObject {
		if yieldIdx >= len(arr.Elems) {
			return object.NewStopIterErr("iter stopped")
		}
		yielded := arr.Elems[yieldIdx]
		yieldIdx++
		return yielded
	}
}
