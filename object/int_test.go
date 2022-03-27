package object

import (
	"testing"
)

func TestIntType(t *testing.T) {
	intObj := NewPanInt(10)
	if intObj.Type() != IntType {
		t.Fatalf("wrong type: expected=%s, got=%s", IntType, intObj.Type())
	}
}

func TestIntInspect(t *testing.T) {
	tests := []struct {
		obj      *PanInt
		expected string
	}{
		{NewPanInt(10), "10"},
		{NewPanInt(1), "1"},
		{NewPanInt(-4), "-4"},
		{NewPanInt(12345), "12345"},
	}

	for _, tt := range tests {
		if tt.obj.Inspect() != tt.expected {
			t.Errorf("wrong output: expected=%s, got=%s",
				tt.expected, tt.obj.Inspect())
		}
	}
}

func TestIntRepr(t *testing.T) {
	tests := []struct {
		obj      *PanInt
		expected string
	}{
		{NewPanInt(10), "10"},
		{NewPanInt(1), "1"},
		{NewPanInt(-4), "-4"},
		{NewPanInt(12345), "12345"},
	}

	for _, tt := range tests {
		if tt.obj.Repr() != tt.expected {
			t.Errorf("wrong output: expected=%s, got=%s",
				tt.expected, tt.obj.Inspect())
		}
	}
}

func TestIntProto(t *testing.T) {
	i := NewPanInt(10)
	if i.Proto() != BuiltInIntObj {
		t.Fatalf("Proto is not BuiltInIntObj. got=%T (%+v)",
			i.Proto(), i.Proto())
	}
}

func TestInheritedIntProto(t *testing.T) {
	intChild := ChildPanObjPtr(BuiltInIntObj, EmptyPanObjPtr())
	a := NewInheritedInt(intChild, 10)
	if a.Proto() != intChild {
		t.Fatalf("Proto is not intChild. got=%T (%s)",
			a.Proto(), a.Proto().Inspect())
	}
}

func TestIntZero(t *testing.T) {
	tests := []struct {
		name string
		obj  *PanInt
	}{
		{"2", NewPanInt(2)},
	}

	for _, tt := range tests {
		tt := tt // pin

		t.Run(tt.name, func(t *testing.T) {
			actual := tt.obj.Zero()

			if actual != tt.obj {
				t.Errorf("zero must be itself (%#v). got=%s (%#v)",
					tt.obj, actual.Repr(), actual)
			}
		})
	}
}

func TestIntHash(t *testing.T) {
	tests := []struct {
		obj      *PanInt
		expected int
	}{
		{NewPanInt(10), 10},
		{NewPanInt(-2), -2},
		{NewPanInt(12345678901), 12345678901},
	}

	for _, tt := range tests {
		h := tt.obj.Hash()

		if h.Type != IntType {
			t.Fatalf("hash type must be IntType. got=%s", h.Type)
		}

		if h.Value != uint64(tt.expected) {
			t.Errorf("wrong hash key: got=%d, expected=%d",
				h.Value, uint64(tt.expected))
		}
	}
}

// checked by compiler (this function works nothing)
func testIntIsPanObject() {
	var _ PanObject = NewPanInt(10)
}

func testIntIsPanScalar() {
	var _ PanScalar = NewPanInt(10)
}

func TestNewPanInt(t *testing.T) {
	tests := []struct {
		i int64
	}{
		{1},
		{5},
		{-1},
	}

	for _, tt := range tests {
		actual := NewPanInt(tt.i)
		if actual.Value != tt.i {
			t.Errorf("wrong value. expected=%d, got=%d", tt.i, actual.Value)
		}
	}
}

func TestNewInheritedInt(t *testing.T) {
	// child of Int
	proto := ChildPanObjPtr(BuiltInIntObj, EmptyPanObjPtr())

	tests := []struct {
		i int64
	}{
		{10},
	}

	for _, tt := range tests {
		actual := NewInheritedInt(proto, tt.i)

		if actual.Value != tt.i {
			t.Errorf("value is wrong. expected=%d, got=%d",
				tt.i, actual.Value)
		}
	}
}
