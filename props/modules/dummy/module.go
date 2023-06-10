package dummy

import "github.com/Syuparn/pangaea/object"

func New() map[string]object.PanObject {
	return map[string]object.PanObject{
		"message": object.NewPanStr("This is a dummy module."),
	}
}
