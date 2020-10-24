package props

import (
	"../object"
)

var eqSym = object.NewPanStr("==")

func propIn(obj *object.PanObj, propName string) (object.Pair, bool) {
	propSym := object.GetSymHash(propName)
	pair, ok := (*obj.Pairs)[propSym]
	return pair, ok
}
