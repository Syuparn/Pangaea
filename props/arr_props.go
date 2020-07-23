package props

import (
	"../object"
)

func ArrProps(propContainer map[string]object.PanObject) map[string]object.PanObject {
	// NOTE: inject some built-in functions which relate to parser or evaluator
	return map[string]object.PanObject{
		"at": propContainer["Arr_at"],
	}
}
