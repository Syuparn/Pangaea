package props

import (
	"github.com/Syuparn/pangaea/object"
)

// MatchProps provides built-in props for Match.
// NOTE: Some Match props are defind by native code (not by this function).
func MatchProps(propContainer map[string]object.PanObject) map[string]object.PanObject {
	// NOTE: inject some built-in functions which relate to parser or evaluator
	return map[string]object.PanObject{
		"_name": object.NewPanStr("Match"),
	}
}
