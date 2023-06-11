package builtin

import "github.com/Syuparn/pangaea/object"

func New() map[string]object.PanObject {
	return map[string]object.PanObject{
		"newHandler":      object.NewPanBuiltInFunc(newHandler),
		"newServer":       object.NewPanBuiltInFunc(newServer),
		"request":         object.NewPanBuiltInFunc(request),
		"serve":           object.NewPanBuiltInFunc(serve),
		"serveBackground": object.NewPanBuiltInFunc(serveBackground),
		"stop":            object.NewPanBuiltInFunc(stop),
	}
}
