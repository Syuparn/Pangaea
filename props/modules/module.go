package modules

import (
	"github.com/Syuparn/pangaea/object"
	"github.com/Syuparn/pangaea/props/modules/dummy"
)

type ModuleFactory = func() map[string]object.PanObject

var Modules = map[string]ModuleFactory{
	"dummy": dummy.New,
}
