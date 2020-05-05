package object

import (
	"testing"
)

func TestNilType(t *testing.T) {
	nilObj := PanNil{}
	if nilObj.Type() != NIL_TYPE {
		t.Fatalf("wrong type: expected=%s, got=%s", NIL_TYPE, nilObj.Type())
	}
}

func TestNilInspect(t *testing.T) {
	tests := []struct {
		obj      PanNil
		expected string
	}{
		{PanNil{}, "nil"},
	}

	for _, tt := range tests {
		if tt.obj.Inspect() != tt.expected {
			t.Errorf("wrong output: expected=%s, got=%s",
				tt.expected, tt.obj.Inspect())
		}
	}
}

func TestNilProto(t *testing.T) {
	n := PanNil{}
	if n.Proto() != builtInNilObj {
		t.Fatalf("Proto is not BuiltinNilObj. got=%T (%+v)",
			n.Proto(), n.Proto())
	}
}

// checked by compiler (this function works nothing)
func testNilIsPanObject() {
	var _ PanObject = &PanNil{}
}
