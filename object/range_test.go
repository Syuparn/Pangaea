package object

import (
	"testing"
)

func TestRangeType(t *testing.T) {
	obj := PanRange{}
	if obj.Type() != RangeType {
		t.Fatalf("wrong type: expected=%s, got=%s", RangeType, obj.Type())
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
		{
			PanRange{&PanInt{1}, &PanInt{2}, &PanInt{3}},
			"(1:2:3)",
		},
		{
			PanRange{&PanNil{}, &PanInt{20}, &PanInt{-1}},
			"(nil:20:-1)",
		},
		{
			PanRange{NewPanStr("a"), NewPanStr("z"), &PanInt{-1}},
			`("a":"z":-1)`,
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
	if r.Proto() != BuiltInRangeObj {
		t.Fatalf("Proto is not BuiltInRangeObj. got=%T (%+v)",
			r.Proto(), r.Proto())
	}
}

// checked by compiler (this function works nothing)
func testRangeIsPanObject() {
	var _ PanObject = &PanRange{&PanInt{1}, &PanInt{2}, &PanInt{3}}
}
