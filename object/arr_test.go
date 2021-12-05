package object

import (
	"testing"
)

func TestArrType(t *testing.T) {
	obj := NewPanArr()
	if obj.Type() != ArrType {
		t.Fatalf("wrong type: expected=%s, got=%s", ArrType, obj.Type())
	}
}

func TestArrInspect(t *testing.T) {
	tests := []struct {
		obj      *PanArr
		expected string
	}{
		{
			NewPanArr(),
			`[]`,
		},
		{
			NewPanArr(NewPanInt(1)),
			`[1]`,
		},
		{
			NewPanArr(NewPanStr("foo")),
			`["foo"]`,
		},
		{
			NewPanArr(NewPanFloat(1.0)),
			`[1.000000]`,
		},
		{
			NewPanArr(NewPanInt(1), NewPanInt(-10)),
			`[1, -10]`,
		},
		{
			NewPanArr(NewPanInt(1), NewPanStr("foo"), BuiltInFalse),
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

func TestArrRepr(t *testing.T) {
	tests := []struct {
		obj      *PanArr
		expected string
	}{
		{
			NewPanArr(),
			`[]`,
		},
		{
			NewPanArr(NewPanInt(1)),
			`[1]`,
		},
		{
			NewPanArr(NewPanStr("foo")),
			`["foo"]`,
		},
		{
			NewPanArr(NewPanFloat(1.0)),
			`[1.000000]`,
		},
		{
			NewPanArr(NewPanInt(1), NewPanInt(-10)),
			`[1, -10]`,
		},
		{
			NewPanArr(NewPanInt(1), NewPanStr("foo"), BuiltInFalse),
			`[1, "foo", false]`,
		},
		// elem's '_name can be used
		{
			NewPanArr(PanObjInstancePtr(&map[SymHash]Pair{
				GetSymHash("_name"): {NewPanStr("_name"), NewPanStr("foo")},
			})),
			`[foo]`,
		},
		{
			NewPanArr(
				NewPanArr(PanObjInstancePtr(&map[SymHash]Pair{
					GetSymHash("_name"): {NewPanStr("_name"), NewPanStr("foo")},
				})),
			),
			`[[foo]]`,
		},
	}

	for _, tt := range tests {
		if tt.obj.Repr() != tt.expected {
			t.Errorf("wrong output: expected=%s, got=%s",
				tt.expected, tt.obj.Inspect())
		}
	}
}

func TestArrProto(t *testing.T) {
	a := NewPanArr()
	if a.Proto() != BuiltInArrObj {
		t.Fatalf("Proto is not BuiltInArrObj. got=%T (%+v)",
			a.Proto(), a.Proto())
	}
}

func TestInheritedArrProto(t *testing.T) {
	arrChild := ChildPanObjPtr(NewPanArr(), EmptyPanObjPtr())
	a := NewInheritedArr(arrChild)
	if a.Proto() != arrChild {
		t.Fatalf("Proto is not arrChild. got=%T (%s)",
			a.Proto(), a.Proto().Inspect())
	}
}

func TestArrZero(t *testing.T) {
	// Arr.bear
	arrChild := ChildPanObjPtr(BuiltInArrObj, EmptyPanObjPtr())

	tests := []struct {
		name string
		obj  *PanArr
	}{
		{"[1]", NewPanArr(NewPanInt(1))},
		{"Child(1)", NewInheritedArr(arrChild, NewPanInt(1))},
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

// checked by compiler (this function works nothing)
func testArrIsPanObject() {
	var _ PanObject = NewPanArr()
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

func TestNewInheritedArr(t *testing.T) {
	// child of Arr
	proto := ChildPanObjPtr(NewPanArr(), EmptyPanObjPtr())

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
		actual := NewInheritedArr(proto, tt.elems...)

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
