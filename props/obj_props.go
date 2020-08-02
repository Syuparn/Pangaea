package props

import (
	"../object"
	"github.com/tanaton/dtoa"
	"io"
)

func ObjProps(propContainer map[string]object.PanObject) map[string]object.PanObject {
	// NOTE: inject some built-in functions which relate to parser or evaluator
	return map[string]object.PanObject{
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
				sSym := &object.PanStr{Value: "S"}
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
				return &object.PanStr{Value: args[0].Inspect()}
			},
		),
		"S": f(
			func(
				env *object.Env, kwargs *object.PanObj, args ...object.PanObject,
			) object.PanObject {
				if len(args) < 1 {
					return object.NewTypeErr("Obj#repr requires at least 1 arg")
				}
				return &object.PanStr{Value: formattedStr(args[0])}
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
