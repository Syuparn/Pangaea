package props

import (
	"github.com/Syuparn/pangaea/object"
)

// IterableProps provides built-in props for Iterable.
// NOTE: Some Iterable props are defind by native code (not by this function).
func IterableProps(propContainer map[string]object.PanObject) map[string]object.PanObject {
	// NOTE: inject some built-in functions which relate to parser or evaluator
	return map[string]object.PanObject{
		"_name": object.NewPanStr("Iterable"),
	}
}
