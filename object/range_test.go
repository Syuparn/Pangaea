package object

import (
	"testing"
)

func TestRangeType(t *testing.T) {
	obj := PanRange{}
	if obj.Type() != RANGE_TYPE {
		t.Fatalf("wrong type: expected=%s, got=%s", RANGE_TYPE, obj.Type())
	}
}

func TestRangeInspect(t *testing.T) {
	tests := []struct {
		obj      PanRange
		expected string
	}{
		{
			PanRange{&PanNil{}, &PanNil{}, &PanNil{}},
			"(nil:nil:nil)",
		},
	}

	for _, tt := range tests {
		if tt.obj.Inspect() != tt.expected {
			t.Errorf("wrong output: expected=%s, got=%s",
				tt.expected, tt.obj.Inspect())
		}
	}
}

func TestRangeProto(t *testing.T) {
	r := PanRange{}
	if r.Proto() != builtInRangeObj {
		t.Fatalf("Proto is not BuiltinRangeObj. got=%T (%+v)",
			r.Proto(), r.Proto())
	}
}

// checked by compiler (this function works nothing)
func testRangeIsPanObject() {
	var _ PanObject = &PanRange{&PanInt{1}, &PanInt{2}, &PanInt{3}}
}
