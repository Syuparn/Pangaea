package props

import (
	"github.com/Syuparn/pangaea/object"
)

// WrappableProps provides built-in props for Wrappable.
// NOTE: Some Wrappable props are defind by native code (not by this function).
func WrappableProps(propContainer map[string]object.PanObject) map[string]object.PanObject {
	// NOTE: inject some built-in functions which relate to parser or evaluator
	return map[string]object.PanObject{
		"_name": object.NewPanStr("Wrappable"),
	}
}
