package object

import (
	"testing"
)

func TestNilType(t *testing.T) {
	nilObj := NewPanNil()
	if nilObj.Type() != NilType {
		t.Fatalf("wrong type: expected=%s, got=%s", NilType, nilObj.Type())
	}
}

func TestNilInspect(t *testing.T) {
	tests := []struct {
		obj      *PanNil
		expected string
	}{
		{NewPanNil(), "nil"},
	}

	for _, tt := range tests {
		if tt.obj.Inspect() != tt.expected {
			t.Errorf("wrong output: expected=%s, got=%s",
				tt.expected, tt.obj.Inspect())
		}
	}
}

func TestNilRepr(t *testing.T) {
	tests := []struct {
		obj      *PanNil
		expected string
	}{
		{NewPanNil(), "nil"},
	}

	for _, tt := range tests {
		if tt.obj.Repr() != tt.expected {
			t.Errorf("wrong output: expected=%s, got=%s",
				tt.expected, tt.obj.Inspect())
		}
	}
}

func TestNilProto(t *testing.T) {
	n := NewPanNil()
	if n.Proto() != BuiltInNilObj {
		t.Fatalf("Proto is not BuiltInNilObj. got=%T (%+v)",
			n.Proto(), n.Proto())
	}
}

func TestNilZero(t *testing.T) {
	tests := []struct {
		name     string
		obj      *PanNil
		expected PanObject
	}{
		{"nil", NewPanNil(), BuiltInNil},
	}

	for _, tt := range tests {
		tt := tt // pin

		t.Run(tt.name, func(t *testing.T) {
			actual := tt.obj.Zero()
			if actual != tt.expected {
				t.Errorf("zero value must be %s. got=%s",
					tt.expected.Repr(), actual.Repr())
			}
		})
	}
}

func TestNilHash(t *testing.T) {
	tests := []struct {
		obj      *PanNil
		expected int
	}{
		{NewPanNil(), 0},
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
	var _ PanObject = NewPanNil()
}

func testNilIsScalarObject() {
	var _ PanScalar = NewPanNil()
}

func TestNewPanNil(t *testing.T) {
	actual := NewPanNil()
	if actual != BuiltInNil {
		// NOTE: actual is *PanNil but not (golang) nil!
		t.Errorf("actual must be BuiltInNil. got %#v", actual)
	}
}
