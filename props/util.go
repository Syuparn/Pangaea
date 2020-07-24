package props

import (
	"../object"
)

func traceProtoOf(
	obj object.PanObject,
	cond func(object.PanObject) bool,
) (object.PanObject, bool) {
	for o := obj; o.Proto() != nil; o = o.Proto() {
		if cond(o) {
			return o, true
		}
	}
	return nil, false
}
