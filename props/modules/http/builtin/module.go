package builtin

import "github.com/Syuparn/pangaea/object"

func New() map[string]object.PanObject {
	return map[string]object.PanObject{
		"request": object.NewPanBuiltInFunc(request),
	}
}
