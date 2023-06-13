package builtin

import (
	"testing"

	"github.com/Syuparn/pangaea/object"
)

func TestNewHandler(t *testing.T) {
	tests := []struct {
		name     string
		args     []object.PanObject
		expected *panHandler
	}{
		{
			"new handler",
			[]object.PanObject{
				object.NewPanStr("GET"),
				object.NewPanStr("/foo"),
				dummyCallback("{|res| {body: \"ok\"}}"),
			},
			&panHandler{
				"GET",
				"/foo",
				toHandler(object.NewEnv(), dummyCallback("{|res| {body: \"ok\"}}")),
			},
		},
	}

	for _, tt := range tests {
		ret := newHandler(object.NewEnv(), object.EmptyPanObjPtr(), tt.args...)

		actual, ok := ret.(*panHandler)
		if !ok {
			t.Errorf("ret must be *panHandler: got=%T (%v)", ret, ret)
		}

		if actual.method != tt.expected.method {
			t.Errorf("wrong method: expected=%s, got=%s",
				tt.expected.method, actual.method)
		}
		if actual.path != tt.expected.path {
			t.Errorf("wrong path: expected=%s, got=%s",
				tt.expected.path, actual.path)
		}
		// TODO: compare handlers
	}
}

func TestNewHandlerErr(t *testing.T) {
	tests := []struct {
		name     string
		args     []object.PanObject
		expected *object.PanErr
	}{
		{
			"insufficient args",
			[]object.PanObject{
				object.NewPanStr("GET"),
				object.NewPanStr("/foo"),
			},
			object.NewTypeErr("newHandler requires at least 3 args"),
		},
		{
			"args[0] is not str",
			[]object.PanObject{
				object.NewPanInt(1),
				object.NewPanStr("/foo"),
				dummyCallback("{|res| {body: \"ok\"}}"),
			},
			object.NewTypeErr("`1` cannot be treated as str"),
		},
		{
			"args[1] is not str",
			[]object.PanObject{
				object.NewPanStr("GET"),
				object.NewPanInt(1),
				dummyCallback("{|res| {body: \"ok\"}}"),
			},
			object.NewTypeErr("`1` cannot be treated as str"),
		},
		{
			"args[2] is not func",
			[]object.PanObject{
				object.NewPanStr("GET"),
				object.NewPanStr("/foo"),
				object.NewPanInt(1),
			},
			object.NewTypeErr("`1` cannot be treated as func"),
		},
	}

	for _, tt := range tests {
		ret := newHandler(object.NewEnv(), object.EmptyPanObjPtr(), tt.args...)

		if ret.Type() != object.ErrType {
			t.Fatalf("error must be raised: %s", ret.Inspect())
		}
		if ret.Inspect() != tt.expected.Inspect() {
			t.Errorf("wrong value. expected=%v, got=%v", tt.expected.Inspect(), ret.Inspect())
		}
	}
}
