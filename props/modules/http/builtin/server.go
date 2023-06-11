package builtin

import (
	"context"
	"fmt"

	"github.com/Syuparn/pangaea/object"
)

func serve(env *object.Env, kwargs *object.PanObj, args ...object.PanObject) object.PanObject {
	if len(args) < 2 {
		return object.NewTypeErr("serve requires at least 2 args")
	}

	srv, ok := args[0].(*panServer)
	if !ok {
		return object.NewTypeErr(fmt.Sprintf("`%s` cannot be treated as server", args[0].Inspect()))
	}

	urlStr, ok := object.TraceProtoOfStr(args[1])
	if !ok {
		return object.NewTypeErr(fmt.Sprintf("`%s` cannot be treated as str", args[1].Inspect()))
	}

	err := srv.server.Start(urlStr.Value)
	return object.NewPanErr(err.Error())
}

func serveBackground(env *object.Env, kwargs *object.PanObj, args ...object.PanObject) object.PanObject {
	if len(args) < 2 {
		return object.NewTypeErr("serveBackground requires at least 2 args")
	}

	srv, ok := args[0].(*panServer)
	if !ok {
		return object.NewTypeErr(fmt.Sprintf("`%s` cannot be treated as server", args[0].Inspect()))
	}

	urlStr, ok := object.TraceProtoOfStr(args[1])
	if !ok {
		return object.NewTypeErr(fmt.Sprintf("`%s` cannot be treated as str", args[1].Inspect()))
	}

	go func() {
		srv.server.Start(urlStr.Value)
	}()
	return object.BuiltInNil
}

func stop(env *object.Env, kwargs *object.PanObj, args ...object.PanObject) object.PanObject {
	if len(args) < 1 {
		return object.NewTypeErr("stop requires at least 1 arg")
	}

	srv, ok := args[0].(*panServer)
	if !ok {
		return object.NewTypeErr(fmt.Sprintf("`%s` cannot be treated as server", args[0].Inspect()))
	}

	err := srv.server.Shutdown(context.Background())
	if err != nil {
		return object.NewPanErr(err.Error())
	}
	return object.BuiltInNil
}

func newServer(env *object.Env, kwargs *object.PanObj, args ...object.PanObject) object.PanObject {
	srv := NewPanServer()

	for _, arg := range args {
		handler, ok := arg.(*panHandler)
		if !ok {
			return object.NewValueErr(fmt.Sprintf("`%s` cannot be treated as handler", arg.Inspect()))
		}
		if errObj := handler.register(srv.server); errObj != nil {
			return errObj
		}
	}

	return srv
}
