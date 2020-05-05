package object

import (
	"testing"
)

func TestBoolType(t *testing.T) {
	boolObj := PanBool{true}
	if boolObj.Type() != BOOL_TYPE {
		t.Fatalf("wrong type: expected=%s, got=%s", BOOL_TYPE, boolObj.Type())
	}
}

func TestBoolInspect(t *testing.T) {
	tests := []struct {
		obj      PanBool
		expected string
	}{
		{PanBool{true}, "true"},
		{PanBool{false}, "false"},
	}

	for _, tt := range tests {
		if tt.obj.Inspect() != tt.expected {
			t.Errorf("wrong output: expected=%s, got=%s",
				tt.expected, tt.obj.Inspect())
		}
	}
}

func TestBoolProto(t *testing.T) {
	b := PanBool{true}
	if b.Proto() != builtInIntObj {
		t.Fatalf("Proto of int is not BuiltinBoolObj. got=%T (%+v)",
			b.Proto(), b.Proto())
	}
}

// checked by compiler (this function works nothing)
//func testIntIsPanObject() {
//	var _ PanObject = &PanInt{10}
//}
