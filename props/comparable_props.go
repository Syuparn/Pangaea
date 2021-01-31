package props

import (
	"github.com/Syuparn/pangaea/object"
)

// ComparableProps provides built-in props for Comparable.
// NOTE: Some Comparable props are defind by native code (not by this function).
func ComparableProps(propContainer map[string]object.PanObject) map[string]object.PanObject {
	// NOTE: inject some built-in functions which relate to parser or evaluator
	return map[string]object.PanObject{
		"_name": object.NewPanStr("Comparable"),
	}
}
