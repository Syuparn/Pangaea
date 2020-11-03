package props

import (
	"fmt"

	"github.com/Syuparn/pangaea/object"
)

// KernelProps provides built-in props for Kernel.
// NOTE: Some Kernel props are defind by native code (not by this function).
func KernelProps(propContainer map[string]object.PanObject) map[string]object.PanObject {
	// NOTE: inject some built-in functions which relate to parser or evaluator
	return map[string]object.PanObject{
		"assert": f(
			func(
				env *object.Env, kwargs *object.PanObj, args ...object.PanObject,
			) object.PanObject {
				if len(args) < 1 {
					return object.NewTypeErr("assert requires at least 1 arg")
				}

				// convert args[0] to bool
				bSym := object.NewPanStr("B")
				objBool := propContainer["Obj_callProp"].(*object.PanBuiltIn).Fn(
					env, object.EmptyPanObjPtr(),
					object.EmptyPanObjPtr(), args[0], bSym,
				)

				if objBool == object.BuiltInTrue {
					return object.BuiltInNil
				}

				return object.NewAssertionErr(fmt.Sprintf("%s is not truty.",
					args[0].Inspect()))
			},
		),
	}
}
