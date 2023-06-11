package builtin

import (
	"context"
	"net/http"
	"testing"
	"time"

	"github.com/Syuparn/pangaea/object"
	"github.com/labstack/echo/v4"
)

func TestNewServer(t *testing.T) {
	tests := []struct {
		name    string
		req     *http.Request
		handler *panHandler
	}{
		{
			"method: GET",
			must(http.NewRequest("GET", "http://localhost:50000", nil)),
			newPanHandler(object.NewEnv(), "GET", "/", dummyCallback(`{|res| "ok"}`)),
		},
		{
			"method: POST",
			must(http.NewRequest("POST", "http://localhost:50000", nil)),
			newPanHandler(object.NewEnv(), "POST", "/", dummyCallback(`{|res| "ok"}`)),
		},
		{
			"method: PUT",
			must(http.NewRequest("PUT", "http://localhost:50000", nil)),
			newPanHandler(object.NewEnv(), "PUT", "/", dummyCallback(`{|res| "ok"}`)),
		},
		{
			"method: DELETE",
			must(http.NewRequest("DELETE", "http://localhost:50000", nil)),
			newPanHandler(object.NewEnv(), "DELETE", "/", dummyCallback(`{|res| "ok"}`)),
		},
		{
			"method: PATCH",
			must(http.NewRequest("PATCH", "http://localhost:50000", nil)),
			newPanHandler(object.NewEnv(), "PATCH", "/", dummyCallback(`{|res| "ok"}`)),
		},
		{
			"path matching",
			must(http.NewRequest("GET", "http://localhost:50000/foo/bar", nil)),
			newPanHandler(object.NewEnv(), "GET", "/foo/bar", dummyCallback(`{|res| "ok"}`)),
		},
	}

	for _, tt := range tests {
		tt := tt // pin

		t.Run(tt.name, func(t *testing.T) {
			env := object.NewEnv()
			ret := newServer(env, object.EmptyPanObjPtr(), tt.handler)

			if ret.Type() == object.ErrType {
				t.Fatalf("error raised: %s", ret.Inspect())
			}
			srv, ok := ret.(*panServer)
			if !ok {
				t.Fatalf("srv is not *panServer: got=%T", ret)
			}

			defer srv.server.Shutdown(context.TODO())
			go func() {
				err := srv.server.Start(":50000")
				t.Log(err.Error())
			}()
			time.Sleep(100 * time.Millisecond)

			client := &http.Client{}
			res, err := client.Do(tt.req)
			if err != nil {
				t.Fatalf("error raised: %s", err)
			}

			if res.StatusCode != 200 {
				t.Errorf("wrong status. expected=%v, got=%v", 200, res.StatusCode)
			}
		})
	}
}

func TestNewServerNotMatched(t *testing.T) {
	tests := []struct {
		name     string
		req      *http.Request
		handler  *panHandler
		expected int
	}{
		{
			"method is not matched",
			must(http.NewRequest("POST", "http://localhost:50000", nil)),
			newPanHandler(object.NewEnv(), "GET", "/", dummyCallback(`{|res| "ok"}`)),
			http.StatusMethodNotAllowed,
		},
		{
			"path is not matched",
			must(http.NewRequest("GET", "http://localhost:50000/foo", nil)),
			newPanHandler(object.NewEnv(), "GET", "/bar", dummyCallback(`{|res| "ok"}`)),
			http.StatusNotFound,
		},
	}

	for _, tt := range tests {
		tt := tt // pin

		t.Run(tt.name, func(t *testing.T) {
			env := object.NewEnv()
			ret := newServer(env, object.EmptyPanObjPtr(), tt.handler)

			if ret.Type() == object.ErrType {
				t.Fatalf("error raised: %s", ret.Inspect())
			}
			srv, ok := ret.(*panServer)
			if !ok {
				t.Fatalf("srv is not *panServer: got=%T", ret)
			}

			defer srv.server.Shutdown(context.TODO())
			go func() {
				err := srv.server.Start(":50000")
				t.Log(err.Error())
			}()
			time.Sleep(100 * time.Millisecond)

			client := &http.Client{}
			res, err := client.Do(tt.req)
			if err != nil {
				t.Fatalf("error raised: %s", err)
			}

			if res.StatusCode != tt.expected {
				t.Errorf("wrong status. expected=%v, got=%v", tt.expected, res.StatusCode)
			}
		})
	}
}

func TestInternalServerError(t *testing.T) {
	tests := []struct {
		name           string
		req            *http.Request
		handler        *panHandler
		expectedStatus int
	}{
		{
			"variable is not defined",
			must(http.NewRequest("GET", "http://localhost:50000", nil)),
			newPanHandler(object.NewEnv(), "GET", "/", dummyCallback(`{|req| a}`)),
			http.StatusInternalServerError,
		},
	}

	for _, tt := range tests {
		tt := tt // pin

		t.Run(tt.name, func(t *testing.T) {
			env := object.NewEnv()
			ret := newServer(env, object.EmptyPanObjPtr(), tt.handler)

			if ret.Type() == object.ErrType {
				t.Fatalf("error raised: %s", ret.Inspect())
			}
			srv, ok := ret.(*panServer)
			if !ok {
				t.Fatalf("srv is not *panServer: got=%T", ret)
			}

			defer srv.server.Shutdown(context.TODO())
			go func() {
				err := srv.server.Start(":50000")
				t.Log(err.Error())
			}()
			time.Sleep(100 * time.Millisecond)

			client := &http.Client{}
			res, err := client.Do(tt.req)
			if err != nil {
				t.Fatalf("error raised: %s", err)
			}

			if res.StatusCode != tt.expectedStatus {
				t.Errorf("wrong status. expected=%v, got=%v", tt.expectedStatus, res.StatusCode)
			}
		})
	}
}

func TestNewServerError(t *testing.T) {
	tests := []struct {
		name     string
		args     []object.PanObject
		expected *object.PanErr
	}{
		{
			"method is invalid",
			[]object.PanObject{
				newPanHandler(object.NewEnv(), "INVALID", "/", dummyCallback(`{|res| "ok"}`)),
			},
			object.NewValueErr("method `INVALID` cannot be used"),
		},
		{
			"arg is not a handler",
			[]object.PanObject{
				newPanHandler(object.NewEnv(), "GET", "/", dummyCallback(`{|res| "ok"}`)),
				object.NewPanInt(1),
			},
			object.NewValueErr("`1` cannot be treated as handler"),
		},
	}

	for _, tt := range tests {
		tt := tt // pin

		t.Run(tt.name, func(t *testing.T) {
			env := object.NewEnv()
			ret := newServer(env, object.EmptyPanObjPtr(), tt.args...)

			if ret.Type() != object.ErrType {
				t.Fatalf("error must be raised: %s", ret.Inspect())
			}
			if ret.Inspect() != tt.expected.Inspect() {
				t.Errorf("wrong value. expected=%v, got=%v", tt.expected.Inspect(), ret.Inspect())
			}
		})
	}
}

func TestServe(t *testing.T) {
	tests := []struct {
		name string
		req  *http.Request
	}{
		{
			"server is running",
			must(http.NewRequest("GET", "http://localhost:50000", nil)),
		},
	}

	for _, tt := range tests {
		tt := tt // pin

		t.Run(tt.name, func(t *testing.T) {
			env := object.NewEnv()
			ret := newServer(env, object.EmptyPanObjPtr(), newPanHandler(object.NewEnv(), "GET", "/", dummyCallback(`{|res| "ok"}`)))
			srv, ok := ret.(*panServer)
			if !ok {
				t.Fatalf("srv is not *panServer: got=%T", ret)
			}

			defer srv.server.Shutdown(context.TODO())
			go func() {
				errObj := serve(env, object.EmptyPanObjPtr(), srv, object.NewPanStr(":50000"))
				t.Log(errObj.Inspect())
			}()
			time.Sleep(100 * time.Millisecond)

			client := &http.Client{}
			res, err := client.Do(tt.req)
			if err != nil {
				t.Fatalf("error raised: %s", err)
			}

			if res.StatusCode != 200 {
				t.Errorf("wrong status. expected=%v, got=%v", 200, res.StatusCode)
			}
		})
	}
}

func TestServeErr(t *testing.T) {
	tests := []struct {
		name     string
		args     []object.PanObject
		expected *object.PanErr
	}{
		{
			"args are insufficient",
			[]object.PanObject{
				newServer(object.NewEnv(), object.EmptyPanObjPtr(), newPanHandler(object.NewEnv(), "GET", "/", dummyCallback(`{|res| "ok"}`))),
			},
			object.NewTypeErr("serve requires at least 2 args"),
		},
		{
			"args[0] is not server",
			[]object.PanObject{
				object.NewPanInt(1),
				object.NewPanStr(":50000"),
			},
			object.NewTypeErr("`1` cannot be treated as server"),
		},
		{
			"args[0] is not str",
			[]object.PanObject{
				newServer(object.NewEnv(), object.EmptyPanObjPtr(), newPanHandler(object.NewEnv(), "GET", "/", dummyCallback(`{|res| "ok"}`))),
				object.NewPanInt(1),
			},
			object.NewTypeErr("`1` cannot be treated as str"),
		},
	}

	for _, tt := range tests {
		tt := tt // pin

		t.Run(tt.name, func(t *testing.T) {
			ret := serve(object.NewEnv(), object.EmptyPanObjPtr(), tt.args...)
			if ret.Type() != object.ErrType {
				t.Fatalf("error must be raised: %s", ret.Inspect())
			}
			if ret.Inspect() != tt.expected.Inspect() {
				t.Errorf("wrong value. expected=%v, got=%v", tt.expected.Inspect(), ret.Inspect())
			}
		})
	}
}

func TestServeBackground(t *testing.T) {
	tests := []struct {
		name string
		req  *http.Request
	}{
		{
			"server is running",
			must(http.NewRequest("GET", "http://localhost:50000", nil)),
		},
	}

	for _, tt := range tests {
		tt := tt // pin

		t.Run(tt.name, func(t *testing.T) {
			env := object.NewEnv()
			ret := newServer(env, object.EmptyPanObjPtr(), newPanHandler(object.NewEnv(), "GET", "/", dummyCallback(`{|res| "ok"}`)))
			srv, ok := ret.(*panServer)
			if !ok {
				t.Fatalf("srv is not *panServer: got=%T", ret)
			}

			defer srv.server.Shutdown(context.TODO())
			ret = serveBackground(env, object.EmptyPanObjPtr(), srv, object.NewPanStr(":50000"))
			if ret != object.BuiltInNil {
				t.Errorf("%s must be BuiltInNil", ret.Inspect())
			}

			client := &http.Client{}
			res, err := client.Do(tt.req)
			if err != nil {
				t.Fatalf("error raised: %s", err)
			}

			if res.StatusCode != 200 {
				t.Errorf("wrong status. expected=%v, got=%v", 200, res.StatusCode)
			}
		})
	}
}

func TestServeBackgroundErr(t *testing.T) {
	tests := []struct {
		name     string
		args     []object.PanObject
		expected *object.PanErr
	}{
		{
			"args are insufficient",
			[]object.PanObject{
				newServer(object.NewEnv(), object.EmptyPanObjPtr(), newPanHandler(object.NewEnv(), "GET", "/", dummyCallback(`{|res| "ok"}`))),
			},
			object.NewTypeErr("serveBackground requires at least 2 args"),
		},
		{
			"args[0] is not server",
			[]object.PanObject{
				object.NewPanInt(1),
				object.NewPanStr(":50000"),
			},
			object.NewTypeErr("`1` cannot be treated as server"),
		},
		{
			"args[0] is not str",
			[]object.PanObject{
				newServer(object.NewEnv(), object.EmptyPanObjPtr(), newPanHandler(object.NewEnv(), "GET", "/", dummyCallback(`{|res| "ok"}`))),
				object.NewPanInt(1),
			},
			object.NewTypeErr("`1` cannot be treated as str"),
		},
	}

	for _, tt := range tests {
		tt := tt // pin

		t.Run(tt.name, func(t *testing.T) {
			ret := serveBackground(object.NewEnv(), object.EmptyPanObjPtr(), tt.args...)
			if ret.Type() != object.ErrType {
				t.Fatalf("error must be raised: %s", ret.Inspect())
			}
			if ret.Inspect() != tt.expected.Inspect() {
				t.Errorf("wrong value. expected=%v, got=%v", tt.expected.Inspect(), ret.Inspect())
			}
		})
	}
}

func TestStop(t *testing.T) {
	env := object.NewEnv()
	e := echo.New()
	srv := &panServer{server: e}
	go func() {
		time.Sleep(100 * time.Millisecond)
		stop(env, object.EmptyPanObjPtr(), srv)
	}()

	err := e.Start(":50000")
	if err.Error() != "http: Server closed" {
		t.Errorf("other error occurred: %s", err.Error())
	}
}

func TestStopError(t *testing.T) {
	tests := []struct {
		name     string
		args     []object.PanObject
		expected *object.PanErr
	}{
		{
			"args are insufficient",
			[]object.PanObject{},
			object.NewTypeErr("stop requires at least 1 arg"),
		},
		{
			"args[0] is not server",
			[]object.PanObject{
				object.NewPanInt(1),
			},
			object.NewTypeErr("`1` cannot be treated as server"),
		},
	}

	for _, tt := range tests {
		tt := tt // pin

		t.Run(tt.name, func(t *testing.T) {
			ret := stop(object.NewEnv(), object.EmptyPanObjPtr(), tt.args...)
			if ret.Type() != object.ErrType {
				t.Fatalf("error must be raised: %s", ret.Inspect())
			}
			if ret.Inspect() != tt.expected.Inspect() {
				t.Errorf("wrong value. expected=%v, got=%v", tt.expected.Inspect(), ret.Inspect())
			}
		})
	}
}
