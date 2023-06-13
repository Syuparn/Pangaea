package builtin

import (
	"fmt"

	"github.com/Syuparn/pangaea/object"
)

func newHandler(env *object.Env, kwargs *object.PanObj, args ...object.PanObject) object.PanObject {
	if len(args) < 3 {
		return object.NewTypeErr("newHandler requires at least 3 args")
	}

	method, ok := object.TraceProtoOfStr(args[0])
	if !ok {
		return object.NewTypeErr(fmt.Sprintf("`%s` cannot be treated as str", args[0].Inspect()))
	}

	path, ok := object.TraceProtoOfStr(args[1])
	if !ok {
		return object.NewTypeErr(fmt.Sprintf("`%s` cannot be treated as str", args[1].Inspect()))
	}

	callback, ok := object.TraceProtoOfFunc(args[2])
	if !ok {
		return object.NewTypeErr(fmt.Sprintf("`%s` cannot be treated as func", args[2].Inspect()))
	}

	return newPanHandler(env, method.Value, path.Value, callback)
}
