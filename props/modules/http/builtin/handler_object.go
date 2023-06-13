package builtin

import (
	"bytes"
	"fmt"
	"os"

	"github.com/labstack/echo/v4"

	"github.com/Syuparn/pangaea/evaluator"
	"github.com/Syuparn/pangaea/object"
)

// handlerType is a type of panHandler.
const handlerType = "HandlerType"

// panHandler is object of arr literal.
type panHandler struct {
	method  string
	path    string
	handler echo.HandlerFunc
}

// Type returns type of this PanObject.
func (h *panHandler) Type() object.PanObjType {
	return handlerType
}

// Inspect returns formatted source code of this object.
func (h *panHandler) Inspect() string {
	return "[handler]"
}

// Repr returns pritty-printed string of this object.
func (h *panHandler) Repr() string {
	return "[handler]"
}

// Proto returns proto of this object.
func (h *panHandler) Proto() object.PanObject {
	return object.BuiltInObjObj
}

// Zero returns zero value of this object.
func (h *panHandler) Zero() object.PanObject {
	return h
}

func (h *panHandler) register(e *echo.Echo) *object.PanErr {
	switch h.method {
	case "GET":
		e.GET(h.path, h.handler)
		return nil
	case "POST":
		e.POST(h.path, h.handler)
		return nil
	case "PUT":
		e.PUT(h.path, h.handler)
		return nil
	case "DELETE":
		e.DELETE(h.path, h.handler)
		return nil
	case "PATCH":
		e.PATCH(h.path, h.handler)
		return nil
	}

	return object.NewValueErr(fmt.Sprintf("method `%s` cannot be used", h.method))
}

// newPanHandler returns new handler object.
func newPanHandler(env *object.Env, method string, path string, callback *object.PanFunc) *panHandler {
	return &panHandler{
		method:  method,
		path:    path,
		handler: toHandler(env, callback),
	}
}

func toHandler(env *object.Env, callback *object.PanFunc) echo.HandlerFunc {
	call := evaluator.NewPropContainer()["Func_call"].(*object.PanBuiltIn)

	handler := func(c echo.Context) error {
		reqObj := requestToObj(c)
		res := call.Fn(env, object.EmptyPanObjPtr(), callback, reqObj)
		if res.Type() == object.ErrType {
			fmt.Fprintln(os.Stderr, res.Inspect())
			return fmt.Errorf(res.Inspect())
		}
		// NOTE: all objects inherit Obj
		resObj, _ := object.TraceProtoOfObj(res)

		status := 200

		if headersPair, ok := (*resObj.Pairs)[object.GetSymHash("headers")]; ok {
			// NOTE: all objects inherit Obj
			headersObj, _ := object.TraceProtoOfObj(headersPair.Value)
			for k, v := range *headersObj.Pairs {
				keyObj, _ := object.SymHash2Str(k)
				key := keyObj.(*object.PanStr).Value

				valueStr, ok := object.TraceProtoOfStr(v.Value)
				if !ok {
					errObj := object.NewTypeErr(fmt.Sprintf("`%s` cannot be treated as str", v.Value.Inspect()))
					fmt.Fprintln(os.Stderr, errObj.Inspect())
					return fmt.Errorf(errObj.Inspect())
				}

				c.Response().Header().Add(key, valueStr.Value)
			}
		}

		if statusPair, ok := (*resObj.Pairs)[object.GetSymHash("status")]; ok {
			if statusPair.Value.Type() == object.IntType {
				status = int(statusPair.Value.(*object.PanInt).Value)
			}
		}

		bodyPair, ok := (*resObj.Pairs)[object.GetSymHash("body")]
		if !ok {
			return c.String(status, "")
		}
		bodyStr, ok := object.TraceProtoOfStr(bodyPair.Value)
		if !ok {
			errObj := object.NewTypeErr(fmt.Sprintf("`%s` cannot be treated as str", bodyPair.Value.Inspect()))
			fmt.Fprintln(os.Stderr, errObj.Inspect())
			return fmt.Errorf(errObj.Inspect())
		}

		if isJSONPair, ok := (*resObj.Pairs)[object.GetSymHash("_isJSON")]; ok {
			if isJSONPair.Value == object.BuiltInTrue {
				return c.JSONBlob(status, []byte(bodyStr.Value))
			}
		}

		return c.String(status, bodyStr.Value)
	}

	return handler
}

func requestToObj(c echo.Context) object.PanObject {
	req := c.Request()

	var b bytes.Buffer
	if req.Body != nil {
		b.ReadFrom(req.Body)
		defer req.Body.Close()
	}

	headers := map[string]object.PanObject{}
	for k, v := range req.Header {
		elems := []object.PanObject{}
		for _, h := range v {
			elems = append(elems, object.NewPanStr(h))
		}

		headers[k] = object.NewPanArr(elems...)
	}

	queries := map[string]object.PanObject{}
	for k, v := range c.QueryParams() {
		elems := []object.PanObject{}
		for _, h := range v {
			elems = append(elems, object.NewPanStr(h))
		}

		queries[k] = object.NewPanArr(elems...)
	}

	params := map[string]object.PanObject{}
	for _, name := range c.ParamNames() {
		params[name] = object.NewPanStr(c.Param(name))
	}

	return mapToObj(map[string]object.PanObject{
		"method":  object.NewPanStr(req.Method),
		"host":    object.NewPanStr(req.Host),
		"url":     object.NewPanStr(req.URL.String()),
		"body":    object.NewPanStr(b.String()),
		"headers": mapToObj(headers),
		"queries": mapToObj(queries),
		"params":  mapToObj(params),
	})
}
