package object

import (
	"testing"
)

func TestRangeType(t *testing.T) {
	obj := NewPanRange(NewPanNil(), NewPanNil(), NewPanNil())
	if obj.Type() != RangeType {
		t.Fatalf("wrong type: expected=%s, got=%s", RangeType, obj.Type())
	}
}

func TestRangeInspect(t *testing.T) {
	tests := []struct {
		obj      *PanRange
		expected string
	}{
		{
			NewPanRange(NewPanNil(), NewPanNil(), NewPanNil()),
			"(nil:nil:nil)",
		},
		{
			NewPanRange(NewPanInt(1), NewPanInt(2), NewPanInt(3)),
			"(1:2:3)",
		},
		{
			NewPanRange(NewPanNil(), NewPanInt(20), NewPanInt(-1)),
			"(nil:20:-1)",
		},
		{
			NewPanRange(NewPanStr("a"), NewPanStr("z"), NewPanInt(-1)),
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
	r := NewPanRange(NewPanNil(), NewPanNil(), NewPanNil())
	if r.Proto() != BuiltInRangeObj {
		t.Fatalf("Proto is not BuiltInRangeObj. got=%T (%+v)",
			r.Proto(), r.Proto())
	}
}

// checked by compiler (this function works nothing)
func testRangeIsPanObject() {
	var _ PanObject = NewPanRange(NewPanInt(1), NewPanInt(2), NewPanInt(3))
}

func TestNewPanRange(t *testing.T) {
	tests := []struct {
		start PanObject
		stop  PanObject
		step  PanObject
	}{
		{NewPanInt(1), NewPanStr("a"), NewPanNil()},
	}

	for _, tt := range tests {
		actual := NewPanRange(tt.start, tt.stop, tt.step)

		if actual.Start != tt.start {
			t.Errorf("wrong start. expected=%#v, got=%#v",
				tt.start, actual.Start)
		}

		if actual.Stop != tt.stop {
			t.Errorf("wrong stop. expected=%#v, got=%#v",
				tt.stop, actual.Stop)
		}

		if actual.Step != tt.step {
			t.Errorf("wrong step. expected=%#v, got=%#v",
				tt.step, actual.Step)
		}
	}
}
