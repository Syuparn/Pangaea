package evaluator

import (
	"../object"
)

func callProp(recv object.PanObject, propHash object.SymHash) (object.PanObject, bool) {
	obj, ok := recv.(*object.PanObj)

	if !ok {
		// TODO: implement prototype chain
		return nil, false
	}

	elem, ok := (*obj.Pairs)[propHash]

	if !ok {
		return nil, false
	}

	return elem.Value, true
}
