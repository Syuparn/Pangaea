package props

import (
	"github.com/Syuparn/pangaea/object"
)

var eqSym = object.NewPanStr("==")
var sSym = object.NewPanStr("S")

func propIn(obj *object.PanObj, propName string) (object.Pair, bool) {
	propSym := object.GetSymHash(propName)
	pair, ok := (*obj.Pairs)[propSym]
	return pair, ok
}
