package props

import (
	"fmt"

	"github.com/Syuparn/pangaea/object"
)

// NumProps provides built-in props for Num.
// NOTE: Some Num props are defind by native code (not by this function).
func NumProps(propContainer map[string]object.PanObject) map[string]object.PanObject {
	// NOTE: inject some built-in functions which relate to parser or evaluator
	return map[string]object.PanObject{
		"_name": object.NewPanStr("Num"),
		"F": f(
			func(
				env *object.Env, kwargs *object.PanObj, args ...object.PanObject,
			) object.PanObject {
				if len(args) < 1 {
					return object.NewTypeErr("Num#F requires at least 1 arg")
				}

				if f, ok := object.TraceProtoOfFloat(args[0]); ok {
					return f
				}

				if i, ok := object.TraceProtoOfInt(args[0]); ok {
					return &object.PanFloat{Value: float64(i.Value)}
				}

				return object.NewTypeErr(fmt.Sprintf("`%s` cannot be treated as num",
					args[0].Inspect()))
			},
		),
	}
}
