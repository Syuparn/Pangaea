package object

import (
	"testing"
)

func TestNilType(t *testing.T) {
	nilObj := PanNil{}
	if nilObj.Type() != NilType {
		t.Fatalf("wrong type: expected=%s, got=%s", NilType, nilObj.Type())
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

		if h.Type != NilType {
			t.Fatalf("hash type must be NilType. got=%s", h.Type)
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

func TestNewPanNil(t *testing.T) {
	actual := NewPanNil()
	if actual != BuiltInNil {
		// NOTE: actual is *PanNil but not (golang) nil!
		t.Errorf("actual must be BuiltInNil. got %#v", actual)
	}
}
