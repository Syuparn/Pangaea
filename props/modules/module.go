package modules

import (
	"github.com/Syuparn/pangaea/object"
	"github.com/Syuparn/pangaea/props/modules/dummy"
	"github.com/Syuparn/pangaea/props/modules/http/builtin"
)

type ModuleFactory = func() map[string]object.PanObject

var Modules = map[string]ModuleFactory{
	"dummy": dummy.New,
	// NOTE: package is renamed because go does not import `internal` package
	"http/internal": builtin.New,
}
