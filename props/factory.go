package props

import (
	"../object"
)

// PanBuiltIn factory
func f(fn object.BuiltInFunc) object.PanObject {
	return &object.PanBuiltIn{Fn: fn}
}
