package object

import (
	"testing"
)

func TestArrType(t *testing.T) {
	obj := PanArr{}
	if obj.Type() != ArrType {
		t.Fatalf("wrong type: expected=%s, got=%s", ArrType, obj.Type())
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
			PanArr{[]PanObject{NewPanStr("foo")}},
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
			PanArr{[]PanObject{&PanInt{1}, NewPanStr("foo"), &PanBool{false}}},
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

func TestNewPanArr(t *testing.T) {
	tests := []struct {
		elems []PanObject
	}{
		{[]PanObject{}},
		{[]PanObject{
			NewPanInt(2),
			NewPanStr("foo"),
		}},
	}

	for _, tt := range tests {
		actual := NewPanArr(tt.elems...)

		if len(actual.Elems) != len(tt.elems) {
			t.Fatalf("wrong length. expected=%d, got=%d",
				len(actual.Elems), len(tt.elems))
		}

		for i, e := range actual.Elems {
			if e != tt.elems[i] {
				t.Errorf("elems[%d] is wrong. expected=%s, got=%s",
					i, tt.elems[i].Inspect(), e.Inspect())
			}
		}
	}
}
