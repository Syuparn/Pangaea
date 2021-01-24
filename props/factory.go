package props

import (
	"github.com/Syuparn/pangaea/object"
)

// PanBuiltIn factory
func f(fn object.BuiltInFunc) object.PanObject {
	return &object.PanBuiltIn{Fn: fn}
}
