package props

import (
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

				return &object.PanBuiltInIter{
					Fn:  objIter(self),
					Env: env, // not used
				}
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
					return &object.PanArr{Elems: []object.PanObject{}}
				}

				items := []object.PanObject{}
				for _, keyHash := range *self.Keys {
					pair, _ := (*self.Pairs)[keyHash]
					items = append(items, &object.PanArr{Elems: []object.PanObject{
						pair.Key,
						pair.Value,
					}})
				}

				if withPrivate {
					for _, keyHash := range *self.PrivateKeys {
						pair, _ := (*self.Pairs)[keyHash]
						items = append(items, &object.PanArr{Elems: []object.PanObject{
							pair.Key,
							pair.Value,
						}})
					}
				}

				return &object.PanArr{Elems: items}
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
					return &object.PanArr{Elems: []object.PanObject{}}
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

				return &object.PanArr{Elems: keys}
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

				return reprStr(args[0])
			},
		),
		"S": f(
			func(
				env *object.Env, kwargs *object.PanObj, args ...object.PanObject,
			) object.PanObject {
				if len(args) < 1 {
					return object.NewTypeErr("Obj#repr requires at least 1 arg")
				}
				return object.NewPanStr(formattedStr(args[0]))
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
					return &object.PanArr{Elems: []object.PanObject{}}
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

				return &object.PanArr{Elems: values}
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

func formattedStr(o object.PanObject) string {
	switch o := o.(type) {
	case *object.PanStr:
		// not quoted
		return o.Value
	case *object.PanFloat:
		// round values for readability
		return dtoaWrapper(o.Value)
	default:
		// same as repr()
		return o.Inspect()
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

		yielded := &object.PanArr{Elems: []object.PanObject{
			pair.Key,
			pair.Value,
		}}
		yieldIdx++
		return yielded
	}
}

func reprStr(obj object.PanObject) *object.PanStr {
	// if self has '_name and it is str, use it
	// NOTE: only refer _name in self (NOT PROTO)!,
	// otherwise all children are shown as _name
	if name, ok := extractName(obj); ok {
		return name
	}

	return object.NewPanStr(obj.Inspect())
}

func extractName(obj object.PanObject) (*object.PanStr, bool) {
	o, ok := obj.(*object.PanObj)
	if !ok {
		return nil, false
	}
	pair, ok := (*o.Pairs)[object.GetSymHash("_name")]
	if !ok {
		return nil, false
	}

	name, ok := object.TraceProtoOfStr(pair.Value)
	if !ok {
		return nil, false
	}

	return name, true
}
