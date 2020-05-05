package object

import (
	"testing"
)

func TestFloatType(t *testing.T) {
	floatObj := PanFloat{1.5}
	if floatObj.Type() != FLOAT_TYPE {
		t.Fatalf("wrong type: expected=%s, got=%s", FLOAT_TYPE, floatObj.Type())
	}
}

func TestFloatInspect(t *testing.T) {
	tests := []struct {
		obj      PanFloat
		expected string
	}{
		{PanFloat{1.5}, "1.5"},
		{PanFloat{0.2}, "0.2"},
		{PanFloat{-4.33}, "-4.33"},
		{PanFloat{123.45}, "123.45"},
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
	if i.Proto() != builtInFloatObj {
		t.Fatalf("Proto of float is not BuiltinFloatObj. got=%T (%+v)",
			i.Proto(), i.Proto())
	}
}

// checked by compiler (this function works nothing)
//func testFloatIsPanObject() {
//	var _ PanObject = &PanFloat{1.5}
//}
