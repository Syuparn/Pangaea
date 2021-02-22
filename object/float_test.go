package object

import (
	"math"
	"testing"
)

func TestFloatType(t *testing.T) {
	floatObj := PanFloat{1.5}
	if floatObj.Type() != FloatType {
		t.Fatalf("wrong type: expected=%s, got=%s", FloatType, floatObj.Type())
	}
}

func TestFloatInspect(t *testing.T) {
	tests := []struct {
		obj      PanFloat
		expected string
	}{
		{PanFloat{1.5}, "1.500000"},
		{PanFloat{0.2}, "0.200000"},
		{PanFloat{-4.33}, "-4.330000"},
		{PanFloat{123.45}, "123.450000"},
	}

	for _, tt := range tests {
		if tt.obj.Inspect() != tt.expected {
			t.Errorf("wrong output: expected=%s, got=%s",
				tt.expected, tt.obj.Inspect())
		}
	}
}

func TestFloatProto(t *testing.T) {
	i := PanFloat{1.4}
	if i.Proto() != BuiltInFloatObj {
		t.Fatalf("Proto of float is not BuiltInFloatObj. got=%T (%+v)",
			i.Proto(), i.Proto())
	}
}

func TestFloatHash(t *testing.T) {
	tests := []struct {
		obj      PanFloat
		expected uint64
	}{
		// Float64bits convert float64 to uint64 with same bit pattern
		{PanFloat{12.3}, math.Float64bits(12.3)},
		{PanFloat{-2.6}, math.Float64bits(-2.6)},
		{PanFloat{1234567890123.45}, math.Float64bits(1234567890123.45)},
		{PanFloat{0.0}, math.Float64bits(0.0)},
	}

	for _, tt := range tests {
		h := tt.obj.Hash()

		if h.Type != FloatType {
			t.Fatalf("hash type must be FloatType. got=%s", h.Type)
		}

		if h.Value != tt.expected {
			t.Errorf("wrong hash key: got=%d, expected=%d",
				h.Value, tt.expected)
		}
	}
}

// checked by compiler (this function works nothing)
func testFloatIsPanObject() {
	var _ PanObject = &PanFloat{1.5}
}

func testFloatIsPanScalar() {
	var _ PanScalar = &PanFloat{1.5}
}

func TestNewPanFloat(t *testing.T) {
	tests := []struct {
		f float64
	}{
		{1.5},
		{5.0},
		{-1.2},
	}

	for _, tt := range tests {
		actual := NewPanFloat(tt.f)
		if actual.Value != tt.f {
			t.Errorf("wrong value. expected=%f, got=%f", tt.f, actual.Value)
		}
	}
}
