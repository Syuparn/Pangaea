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
	if n.Proto() != BuiltInNilObj {
		t.Fatalf("Proto is not BuiltInNilObj. got=%T (%+v)",
			n.Proto(), n.Proto())
	}
}

func TestNilHash(t *testing.T) {
	tests := []struct {
		obj      PanNil
		expected int
	}{
		{PanNil{}, 0},
	}

	for _, tt := range tests {
		h := tt.obj.Hash()

		if h.Type != NIL_TYPE {
			t.Fatalf("hash type must be NIL_TYPE. got=%s", h.Type)
		}

		if h.Value != uint64(tt.expected) {
			t.Errorf("wrong hash key: got=%d, expected=%d",
				h.Value, uint64(tt.expected))
		}
	}
}

// checked by compiler (this function works nothing)
func testNilIsPanObject() {
	var _ PanObject = &PanNil{}
}

func testNilIsScalarObject() {
	var _ PanScalar = &PanNil{}
}
