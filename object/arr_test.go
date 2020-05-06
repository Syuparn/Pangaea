package object

import (
	"testing"
)

func TestArrType(t *testing.T) {
	obj := PanArr{}
	if obj.Type() != ARR_TYPE {
		t.Fatalf("wrong type: expected=%s, got=%s", ARR_TYPE, obj.Type())
	}
}

func TestArrInspect(t *testing.T) {
	tests := []struct {
		obj      PanArr
		expected string
	}{
		{
			PanArr{},
			`[]`,
		},
		{
			PanArr{[]PanObject{&PanInt{1}}},
			`[1]`,
		},
		{
			PanArr{[]PanObject{&PanStr{"foo"}}},
			`["foo"]`,
		},
		{
			PanArr{[]PanObject{&PanFloat{1.0}}},
			`[1.000000]`,
		},
		{
			PanArr{[]PanObject{&PanInt{1}, &PanInt{-10}}},
			`[1, -10]`,
		},
		{
			PanArr{[]PanObject{&PanInt{1}, &PanStr{"foo"}, &PanBool{false}}},
			`[1, "foo", false]`,
		},
	}

	for _, tt := range tests {
		if tt.obj.Inspect() != tt.expected {
			t.Errorf("wrong output: expected=%s, got=%s",
				tt.expected, tt.obj.Inspect())
		}
	}
}

func TestArrProto(t *testing.T) {
	a := PanArr{}
	if a.Proto() != BuiltInArrObj {
		t.Fatalf("Proto is not BuiltInArrObj. got=%T (%+v)",
			a.Proto(), a.Proto())
	}
}

// checked by compiler (this function works nothing)
func testArrIsPanObject() {
	var _ PanObject = &PanArr{[]PanObject{&PanInt{1}}}
}
