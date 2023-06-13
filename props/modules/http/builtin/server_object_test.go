package builtin

import (
	"testing"

	"github.com/Syuparn/pangaea/object"
)

func TestServerType(t *testing.T) {
	obj := NewPanServer()
	if obj.Type() != serverType {
		t.Fatalf("wrong type: expected=%s, got=%s", serverType, obj.Type())
	}
}

func TestServerInspect(t *testing.T) {
	tests := []struct {
		obj      *panServer
		expected string
	}{
		{
			NewPanServer(),
			`[server]`,
		},
	}

	for _, tt := range tests {
		if tt.obj.Inspect() != tt.expected {
			t.Errorf("wrong output: expected=%s, got=%s",
				tt.expected, tt.obj.Inspect())
		}
	}
}

func TestServerRepr(t *testing.T) {
	tests := []struct {
		obj      *panServer
		expected string
	}{
		{
			NewPanServer(),
			`[server]`,
		},
	}

	for _, tt := range tests {
		if tt.obj.Repr() != tt.expected {
			t.Errorf("wrong output: expected=%s, got=%s",
				tt.expected, tt.obj.Inspect())
		}
	}
}

func TestServerProto(t *testing.T) {
	a := NewPanServer()
	if a.Proto() != object.BuiltInObjObj {
		t.Fatalf("Proto is not object.BuiltInObjObj. got=%T (%+v)",
			a.Proto(), a.Proto())
	}
}

func TestServerZero(t *testing.T) {
	tests := []struct {
		obj *panServer
	}{
		{NewPanServer()},
	}

	for _, tt := range tests {
		tt := tt // pin

		actual := tt.obj.Zero()

		if actual != tt.obj {
			t.Errorf("zero must be itself (%#v). got=%s (%#v)",
				tt.obj, actual.Repr(), actual)
		}
	}
}

// checked by compiler (this function works nothing)
func testServerIsPanObject() {
	var _ object.PanObject = NewPanServer()
}
