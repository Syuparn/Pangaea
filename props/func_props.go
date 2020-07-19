package props

import (
	"../object"
)

func FuncProps(propContainer map[string]object.PanObject) map[string]object.PanObject {
	// NOTE: inject some built-in functions which relate to parser or evaluator
	return map[string]object.PanObject{
		"call": propContainer["Func_call"],
	}
}
