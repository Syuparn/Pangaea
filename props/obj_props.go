package props

import (
	"../object"
	"github.com/tanaton/dtoa"
	"io"
)

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

				self, ok := traceProtoOf(args[0], isObj)
				if !ok {
					return object.NewTypeErr(`Obj#_iter cannot be applied to \1`)
				}

				return &object.PanBuiltInIter{
					Fn:  objIter(self.(*object.PanObj)),
					Env: env, // not used
				}
			},
		),
		"B": f(
			func(
				env *object.Env, kwargs *object.PanObj, args ...object.PanObject,
			) object.PanObject {
				if len(args) < 1 {
					return object.NewTypeErr("Obj#B requires at least 1 arg")
				}
				self, ok := traceProtoOf(args[0], isObj)
				if !ok {
					return object.NewTypeErr(`\1 must be obj`)
				}

				if len(*self.(*object.PanObj).Pairs) == 0 {
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

				self, ok := traceProtoOf(args[0], isObj)
				if !ok {
					return &object.PanArr{Elems: []object.PanObject{}}
				}
				obj, _ := self.(*object.PanObj)

				items := []object.PanObject{}
				for _, keyHash := range *obj.Keys {
					pair, _ := (*obj.Pairs)[keyHash]
					items = append(items, &object.PanArr{Elems: []object.PanObject{
						pair.Key,
						pair.Value,
					}})
				}

				if withPrivate {
					for _, keyHash := range *obj.PrivateKeys {
						pair, _ := (*obj.Pairs)[keyHash]
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

				self, ok := traceProtoOf(args[0], isObj)
				if !ok {
					return &object.PanArr{Elems: []object.PanObject{}}
				}
				obj, _ := self.(*object.PanObj)

				keys := []object.PanObject{}
				for _, keyHash := range *obj.Keys {
					keys = append(keys, (*obj.Pairs)[keyHash].Key)
				}

				if withPrivate {
					for _, keyHash := range *obj.PrivateKeys {
						keys = append(keys, (*obj.Pairs)[keyHash].Key)
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

				str, ok := traceProtoOf(sRet, isStr)
				if !ok {
					return object.NewTypeErr(`\1.S must be str`)
				}

				// print
				io.WriteString(panIO.Out, str.(*object.PanStr).Value)
				// breakline
				io.WriteString(panIO.Out, "\n")

				return object.BuiltInNil
			},
		),
		"repr": f(
			func(
				env *object.Env, kwargs *object.PanObj, args ...object.PanObject,
			) object.PanObject {
				if len(args) < 1 {
					return object.NewTypeErr("Obj#repr requires at least 1 arg")
				}
				return object.NewPanStr(args[0].Inspect())
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

				self, ok := traceProtoOf(args[0], isObj)
				if !ok {
					return &object.PanArr{Elems: []object.PanObject{}}
				}
				obj, _ := self.(*object.PanObj)

				values := []object.PanObject{}
				for _, keyHash := range *obj.Keys {
					values = append(values, (*obj.Pairs)[keyHash].Value)
				}

				if withPrivate {
					for _, keyHash := range *obj.PrivateKeys {
						values = append(values, (*obj.Pairs)[keyHash].Value)
					}
				}

				return &object.PanArr{Elems: values}
			},
		),
	}
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
		yieldIdx += 1
		return yielded
	}
}
