package evaluator

import (
	"github.com/Syuparn/pangaea/object"
)

func callProp(recv object.PanObject, propHash object.SymHash) (object.PanObject, bool) {
	// trace prototype chains
	for obj := recv; obj != nil; obj = obj.Proto() {
		prop, ok := findProp(obj, propHash)

		if ok {
			return prop, true
		}
	}
	return nil, false
}

func findProp(o object.PanObject, propHash object.SymHash) (object.PanObject, bool) {
	obj, ok := o.(*object.PanObj)
	if !ok {
		return nil, false
	}

	elem, ok := (*obj.Pairs)[propHash]

	if !ok {
		return nil, false
	}

	return elem.Value, true
}
