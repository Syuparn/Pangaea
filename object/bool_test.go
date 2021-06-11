package object

import (
	"testing"
)

func TestBoolType(t *testing.T) {
	boolObj := PanBool{true}
	if boolObj.Type() != BoolType {
		t.Fatalf("wrong type: expected=%s, got=%s", BoolType, boolObj.Type())
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

func TestBoolRepr(t *testing.T) {
	tests := []struct {
		obj      PanBool
		expected string
	}{
		{PanBool{true}, "true"},
		{PanBool{false}, "false"},
	}

	for _, tt := range tests {
		if tt.obj.Repr() != tt.expected {
			t.Errorf("wrong output: expected=%s, got=%s",
				tt.expected, tt.obj.Inspect())
		}
	}
}

func TestBoolProto(t *testing.T) {
	tests := []struct {
		obj          PanBool
		expected     PanObject
		expectedName string
	}{
		{PanBool{true}, BuiltInOneInt, "BuiltInOneInt"},
		{PanBool{false}, BuiltInZeroInt, "BuiltInZeroInt"},
	}

	for _, tt := range tests {
		if tt.obj.Proto() != tt.expected {
			t.Fatalf("Proto is not %s. got=%T (%+v)",
				tt.expectedName, tt.obj.Proto(), tt.obj.Proto())
		}
	}

}

func TestBoolHash(t *testing.T) {
	tests := []struct {
		obj      PanBool
		expected int
	}{
		{PanBool{true}, 1},
		{PanBool{false}, 0},
	}

	for _, tt := range tests {
		h := tt.obj.Hash()

		if h.Type != BoolType {
			t.Fatalf("hash type must be BoolType. got=%s", h.Type)
		}

		if h.Value != uint64(tt.expected) {
			t.Errorf("wrong hash key: got=%d, expected=%d",
				h.Value, uint64(tt.expected))
		}
	}
}

func TestBuiltInBools(t *testing.T) {
	tests := []struct {
		obj      *PanBool
		expected bool
	}{
		{BuiltInTrue, true},
		{BuiltInFalse, false},
	}

	for _, tt := range tests {
		if tt.obj.Value != tt.expected {
			t.Errorf("wrong value. got=%t, expected=%t", tt.obj.Value, tt.expected)
		}
	}
}

// checked by compiler (this function works nothing)
func testBoolIsPanObject() {
	var _ PanObject = &PanBool{true}
}

func testBoolIsPanScalar() {
	var _ PanScalar = &PanBool{false}
}
