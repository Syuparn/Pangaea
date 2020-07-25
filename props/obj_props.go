package props

import (
	"../object"
	"github.com/tanaton/dtoa"
)

func ObjProps(propContainer map[string]object.PanObject) map[string]object.PanObject {
	// NOTE: inject some built-in functions which relate to parser or evaluator
	return map[string]object.PanObject{
		"callProp": propContainer["Obj_callProp"],
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
