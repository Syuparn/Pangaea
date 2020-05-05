package object

import (
	"testing"
)

func TestIntType(t *testing.T) {
	intObj := PanInt{10}
	if intObj.Type() != INT_TYPE {
		t.Fatalf("wrong type: expected=%s, got=%s", INT_TYPE, intObj.Type())
	}
}

func TestIntInspect(t *testing.T) {
	tests := []struct {
		obj      PanInt
		expected string
	}{
		{PanInt{10}, "10"},
		{PanInt{1}, "1"},
		{PanInt{-4}, "-4"},
		{PanInt{12345}, "12345"},
	}

	for _, tt := range tests {
		if tt.obj.Inspect() != tt.expected {
			t.Errorf("wrong output: expected=%s, got=%s",
				tt.expected, tt.obj.Inspect())
		}
	}
}
