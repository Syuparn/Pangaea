package props

import (
	"fmt"
	"io"

	"github.com/Syuparn/pangaea/object"
	"github.com/tanaton/dtoa"
)

// ObjProps provides built-in props for Obj.
// NOTE: Some Obj props are defind by native code (not by this function).
func ObjProps(propContainer map[string]object.PanObject) map[string]object.PanObject {
	// NOTE: inject some built-in functions which relate to parser or evaluator
	return map[string]object.PanObject{
		"!": f(
			func(
				env *object.Env, kwargs *object.PanObj, args ...object.PanObject,
			) object.PanObject {
				if len(args) < 1 {
					return object.NewTypeErr("! requires at least 1 arg")
				}

				// get args[0].B
				bSym := object.NewPanStr("B")
				objBool := propContainer["Obj_callProp"].(*object.PanBuiltIn).Fn(
					env, object.EmptyPanObjPtr(),
					object.EmptyPanObjPtr(), args[0], bSym,
				)

				if objBool == object.BuiltInTrue {
					return object.BuiltInFalse
				}
				return object.BuiltInTrue
			},
		),
		"_iter": f(
			func(
				env *object.Env, kwargs *object.PanObj, args ...object.PanObject,
			) object.PanObject {
				if len(args) < 1 {
					return object.NewTypeErr("Obj#_iter requires at least 1 arg")
				}

				self, ok := object.TraceProtoOfObj(args[0])
				if !ok {
					return object.NewTypeErr(`Obj#_iter cannot be applied to \1`)
				}

				return object.NewPanBuiltInIter(objIter(self), env)
			},
		),
		"_name": object.NewPanStr("Obj"),
		"B": f(
			func(
				env *object.Env, kwargs *object.PanObj, args ...object.PanObject,
			) object.PanObject {
				if len(args) < 1 {
					return object.NewTypeErr("Obj#B requires at least 1 arg")
				}
				self, ok := object.TraceProtoOfObj(args[0])
				if !ok {
					return object.NewTypeErr(`\1 must be obj`)
				}

				if len(*self.Pairs) == 0 {
					return object.BuiltInFalse
				}
				return object.BuiltInTrue
			},
		),
		"callProp": propContainer["Obj_callProp"],
		"items": f(
			func(
				env *object.Env, kwargs *object.PanObj, args ...object.PanObject,
			) object.PanObject {
				if len(args) < 1 {
					return object.NewTypeErr("Obj#values requires at least 1 arg")
				}

				withPrivate := false
				if pair, ok := propIn(kwargs, "private?"); ok {
					withPrivate = (pair.Value == object.BuiltInTrue)
				}

				self, ok := object.TraceProtoOfObj(args[0])
				if !ok {
					return object.NewPanArr()
				}

				items := []object.PanObject{}
				for _, keyHash := range *self.Keys {
					pair, _ := (*self.Pairs)[keyHash]
					items = append(items, object.NewPanArr(pair.Key, pair.Value))
				}

				if withPrivate {
					for _, keyHash := range *self.PrivateKeys {
						pair, _ := (*self.Pairs)[keyHash]
						items = append(items, object.NewPanArr(pair.Key, pair.Value))
					}
				}

				return object.NewPanArr(items...)
			},
		),
		"keys": f(
			func(
				env *object.Env, kwargs *object.PanObj, args ...object.PanObject,
			) object.PanObject {
				if len(args) < 1 {
					return object.NewTypeErr("Obj#keys requires at least 1 arg")
				}

				withPrivate := false
				if pair, ok := propIn(kwargs, "private?"); ok {
					withPrivate = (pair.Value == object.BuiltInTrue)
				}

				self, ok := object.TraceProtoOfObj(args[0])
				if !ok {
					return object.NewPanArr()
				}

				keys := []object.PanObject{}
				for _, keyHash := range *self.Keys {
					keys = append(keys, (*self.Pairs)[keyHash].Key)
				}

				if withPrivate {
					for _, keyHash := range *self.PrivateKeys {
						keys = append(keys, (*self.Pairs)[keyHash].Key)
					}
				}

				return object.NewPanArr(keys...)
			},
		),
		"p": f(
			func(
				env *object.Env, kwargs *object.PanObj, args ...object.PanObject,
			) object.PanObject {
				if len(args) < 1 {
					return object.NewTypeErr("Obj#p requires at least 1 arg")
				}

				// TODO: pass io by kwarg
				// find IO object
				ioObj, ok := env.Get(object.GetSymHash("IO"))
				if !ok {
					return object.NewNameErr("name `IO` is not defined.")
				}
				panIO, ok := ioObj.(*object.PanIO)
				if !ok {
					return object.NewTypeErr("name `IO` is not io object")
				}

				// get args[0].S
				sSym := object.NewPanStr("S")
				sRet := propContainer["Obj_callProp"].(*object.PanBuiltIn).Fn(
					env, object.EmptyPanObjPtr(),
					object.EmptyPanObjPtr(), args[0], sSym,
				)

				str, ok := object.TraceProtoOfStr(sRet)
				if !ok {
					return object.NewTypeErr(`\1.S must be str`)
				}

				// get kwarg end (default: breakline)
				endPair, ok := (*kwargs.Pairs)[object.GetSymHash("end")]
				if !ok {
					// print
					io.WriteString(panIO.Out, str.Value)
					// breakline
					io.WriteString(panIO.Out, "\n")
					return object.BuiltInNil
				}

				endStr, ok := object.TraceProtoOfStr(endPair.Value)
				if !ok {
					return object.NewTypeErr("end must be str")
				}

				// print
				io.WriteString(panIO.Out, str.Value)
				io.WriteString(panIO.Out, endStr.Value)
				return object.BuiltInNil
			},
		),
		"new": f(
			func(
				env *object.Env, kwargs *object.PanObj, args ...object.PanObject,
			) object.PanObject {
				if len(args) < 2 {
					// NOTE: not error because insufficient args are filled with nil!
					return object.EmptyPanObjPtr()
				}
				if o, ok := args[1].(*object.PanObj); ok {
					// NOTE: proto info is ignored (the returned value is a child of Obj)
					return object.PanObjInstancePtr(o.Pairs)
				}

				// if \2 is nil, return empty obj
				if _, ok := object.TraceProtoOfNil(args[1]); ok {
					return object.EmptyPanObjPtr()
				}

				// wrap by PanObj (return {_value: \2})
				return object.PanObjInstancePtr(&map[object.SymHash]object.Pair{
					object.GetSymHash("_value"): {Key: valueSym, Value: args[1]},
				})
			},
		),
		"repr": f(
			func(
				env *object.Env, kwargs *object.PanObj, args ...object.PanObject,
			) object.PanObject {
				if len(args) < 1 {
					return object.NewTypeErr("Obj#repr requires at least 1 arg")
				}

				return object.NewPanStr(args[0].Repr())
			},
		),
		"S": f(
			func(
				env *object.Env, kwargs *object.PanObj, args ...object.PanObject,
			) object.PanObject {
				if len(args) < 1 {
					return object.NewTypeErr("Obj#S requires at least 1 arg")
				}
				return formattedStr(args[0], kwargs)
			},
		),
		"traverse": f(
			func(
				env *object.Env, kwargs *object.PanObj, args ...object.PanObject,
			) object.PanObject {
				if len(args) < 1 {
					return object.NewTypeErr("Obj#traverse requires at least 1 arg")
				}

				if k, ok := (*kwargs.Pairs)[object.GetSymHash("key")]; ok {
					key, ok := object.TraceProtoOfStr(k.Value)
					if !ok {
						return object.NewTypeErr(
							fmt.Sprintf("%s cannot be treated as str", k.Value.Repr()))
					}

					return traverseSelectedObject(args[0], key)
				}

				return traverseObject(args[0])
			},
		),
		"try": f(
			func(
				env *object.Env, kwargs *object.PanObj, args ...object.PanObject,
			) object.PanObject {
				if len(args) < 1 {
					return object.NewTypeErr("Obj#try requires at least 1 arg")
				}
				return toEitherVal(args[0])
			},
		),
		"values": f(
			func(
				env *object.Env, kwargs *object.PanObj, args ...object.PanObject,
			) object.PanObject {
				if len(args) < 1 {
					return object.NewTypeErr("Obj#values requires at least 1 arg")
				}

				withPrivate := false
				if pair, ok := propIn(kwargs, "private?"); ok {
					withPrivate = (pair.Value == object.BuiltInTrue)
				}

				self, ok := object.TraceProtoOfObj(args[0])
				if !ok {
					return object.NewPanArr()
				}

				values := []object.PanObject{}
				for _, keyHash := range *self.Keys {
					values = append(values, (*self.Pairs)[keyHash].Value)
				}

				if withPrivate {
					for _, keyHash := range *self.PrivateKeys {
						values = append(values, (*self.Pairs)[keyHash].Value)
					}
				}

				return object.NewPanArr(values...)
			},
		),
		"which": f(
			func(
				env *object.Env, kwargs *object.PanObj, args ...object.PanObject,
			) object.PanObject {
				if len(args) < 2 {
					return object.NewTypeErr("Obj#which requires at least 2 args")
				}

				propName, ok := object.TraceProtoOfStr(args[1])
				if !ok {
					return object.NewTypeErr(
						fmt.Sprintf("%s cannot be treated as str", args[1].Repr()))
				}

				if owner, ok := object.FindPropOwner(args[0], propName.SymHash()); ok {
					return owner
				}
				return object.BuiltInNil
			},
		),
	}
}

func toEitherVal(o object.PanObject) object.PanObject {
	// wrap o by EitherVal
	pairMap := map[object.SymHash]object.Pair{
		object.GetSymHash("_value"): {Key: object.NewPanStr("_value"), Value: o},
	}

	return object.NewPanObj(&pairMap, object.BuiltInEitherValObj)
}

func formattedStr(o object.PanObject, kwargs *object.PanObj) object.PanObject {
	switch o := o.(type) {
	case *object.PanStr:
		// not quoted
		return object.NewPanStr(o.Value)
	case *object.PanFloat:
		// round values for readability
		return object.NewPanStr(dtoaWrapper(o.Value))
	case *object.PanInt:
		// handle `base` kwarg
		return formattedIntStr(o, kwargs)
	default:
		return object.NewPanStr(o.Inspect())
	}
}

func dtoaWrapper(f float64) string {
	//               buf,    val, maxDecimalPlaces
	buf := dtoa.Dtoa([]byte{}, f, 324)
	return string(buf)
}

func objIter(o *object.PanObj) object.BuiltInFunc {
	yieldIdx := 0
	hashes := *o.Keys

	return func(
		env *object.Env, kwargs *object.PanObj, args ...object.PanObject,
	) object.PanObject {
		if yieldIdx >= len(hashes) {
			return object.NewStopIterErr("iter stopped")
		}
		pair, ok := (*o.Pairs)[hashes[yieldIdx]]
		// must be ok
		if !ok {
			return object.NewValueErr("pair data in obj somehow got changed")
		}

		yielded := object.NewPanArr(pair.Key, pair.Value)

		yieldIdx++
		return yielded
	}
}

func traverseObject(obj object.PanObject) *object.PanArr {
	traversedPairs := traversePairs(obj)
	elems := make([]object.PanObject, len(traversedPairs))

	for i, tp := range traversedPairs {
		elems[i] = object.NewPanArr(
			object.NewPanArr(tp.Keys...),
			tp.Value,
		)
	}

	return object.NewPanArr(elems...)
}

func traverseSelectedObject(obj object.PanObject, key object.PanScalar) *object.PanArr {
	traversedPairs := traversePairs(obj)
	filteredPairs := []*traversedPair{}
	for _, tp := range traversedPairs {
		if tp.Contains(key) {
			filteredPairs = append(filteredPairs, tp)
		}
	}

	elems := make([]object.PanObject, len(filteredPairs))

	for i, p := range filteredPairs {
		elems[i] = object.NewPanArr(
			object.NewPanArr(p.Keys...),
			p.Value,
		)
	}

	return object.NewPanArr(elems...)
}

func traversePairs(obj object.PanObject) []*traversedPair {
	switch o := obj.(type) {
	case *object.PanObj:
		return traverseObjPairs(o)
	case *object.PanArr:
		return traverseArrPairs(o)
	default:
		return []*traversedPair{}
	}
}

func traverseObjPairs(obj *object.PanObj) []*traversedPair {
	pairs := []*traversedPair{}

	for _, key := range *obj.Keys {
		p := (*obj.Pairs)[key]
		traversedPairs := traversePairs(p.Value)

		if len(traversedPairs) > 0 {
			// elem is collection
			for _, tp := range traversedPairs {
				pairs = append(pairs, tp.PrependKey(p.Key))
			}
		} else {
			// elem is not collection
			pairs = append(pairs, &traversedPair{
				Keys:  []object.PanObject{p.Key},
				Value: p.Value,
			})
		}
	}

	return pairs
}

func traverseArrPairs(arr *object.PanArr) []*traversedPair {
	pairs := []*traversedPair{}

	for i, e := range arr.Elems {
		key := object.NewPanInt(int64(i))
		traversedPairs := traversePairs(e)

		if len(traversedPairs) > 0 {
			// elem is collection
			for _, tp := range traversedPairs {
				pairs = append(pairs, tp.PrependKey(key))
			}
		} else {
			// elem is not collection
			pairs = append(pairs, &traversedPair{
				Keys:  []object.PanObject{key},
				Value: e,
			})
		}
	}

	return pairs
}

type traversedPair struct {
	Keys  []object.PanObject
	Value object.PanObject
}

func (p *traversedPair) PrependKey(key object.PanObject) *traversedPair {
	return &traversedPair{
		Keys:  append([]object.PanObject{key}, p.Keys...),
		Value: p.Value,
	}
}

func (p *traversedPair) Contains(key object.PanScalar) bool {
	for _, k := range p.Keys {
		if s, ok := k.(object.PanScalar); ok {
			if key.Hash() == s.Hash() {
				return true
			}
		}
	}

	return false
}
