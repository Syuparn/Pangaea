package props

import (
	"github.com/Syuparn/pangaea/object"
)

var eqSym = object.NewPanStr("==")
var sSym = object.NewPanStr("S")
var valueSym = object.NewPanStr("_value")
var prefixMinusSym = object.NewPanStr("-%")
var divSym = object.NewPanStr("/")
var floorDivSym = object.NewPanStr("//")

func propIn(obj *object.PanObj, propName string) (object.Pair, bool) {
	propSym := object.GetSymHash(propName)
	pair, ok := (*obj.Pairs)[propSym]
	return pair, ok
}
