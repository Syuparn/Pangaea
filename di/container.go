package di

import (
	"github.com/Syuparn/pangaea/evaluator"
	"github.com/Syuparn/pangaea/object"
)

// InjectBuiltInProps injects and sets up builtin object props,
// which are defined in "props" package.
func InjectBuiltInProps() {
	propContainer := mergePropContainers(
		NewPropContainer(),
		evaluator.NewPropContainer(),
	)
	evaluator.InjectBuiltInProps(propContainer)
}

// NewPropContainer returns container of built-in props
// which are injected into built-in object.
func NewPropContainer() map[string]object.PanObject {
	return map[string]object.PanObject{
		// name format: "ObjName_propName"
		"Str_eval": &object.PanBuiltIn{Fn: strEval},
	}
}

func mergePropContainers(
	containers ...map[string]object.PanObject,
) map[string]object.PanObject {
	mergedCtn := map[string]object.PanObject{}
	for _, ctn := range containers {
		for k, v := range ctn {
			mergedCtn[k] = v
		}
	}

	return mergedCtn
}
