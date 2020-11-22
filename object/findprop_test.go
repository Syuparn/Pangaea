package object

import (
	"testing"
)

func TestFindProp(t *testing.T) {
	obj1 := PanObjInstancePtr(&map[SymHash]Pair{
		GetSymHash("a"): Pair{Key: NewPanStr("a"), Value: BuiltInTrue},
	})

	tests := []struct {
		obj          PanObject
		propName     string
		expectedProp PanObject
		expectedBool bool
	}{
		{
			obj1,
			"a",
			BuiltInTrue,
			true,
		},
		{
			obj1,
			"b",
			nil,
			false,
		},
		{
			NewPanInt(3),
			"b",
			nil,
			false,
		},
	}

	for _, tt := range tests {
		actual, ok := findProp(tt.obj, GetSymHash(tt.propName))

		if ok != tt.expectedBool {
			t.Errorf("ok in %s is wrong: expected=%T, got=%T",
				tt.obj.Inspect(), tt.expectedBool, ok)
		}

		if actual != tt.expectedProp {
			t.Errorf("found prop is wrong: expected=%s, got=%s",
				actual.Inspect(), tt.obj.Inspect())
		}
	}
}

func TestFindPropAlongProtos(t *testing.T) {
	obj1 := PanObjInstancePtr(&map[SymHash]Pair{
		GetSymHash("a"): Pair{Key: NewPanStr("a"), Value: BuiltInTrue},
	}).(*PanObj)
	obj2 := PanObjInstancePtr(&map[SymHash]Pair{
		GetSymHash("b"): Pair{Key: NewPanStr("b"), Value: BuiltInFalse},
	}).(*PanObj)

	tests := []struct {
		obj          PanObject
		propName     string
		expectedProp PanObject
		expectedBool bool
	}{
		{
			obj1,
			"a",
			BuiltInTrue,
			true,
		},
		{
			obj1,
			"b",
			nil,
			false,
		},
		{
			NewPanInt(3),
			"b",
			nil,
			false,
		},
		// find proto's prop
		{
			ChildPanObjPtr(obj1, obj2),
			"a",
			BuiltInTrue,
			true,
		},
	}

	for _, tt := range tests {
		actual, ok := FindPropAlongProtos(tt.obj, GetSymHash(tt.propName))

		if ok != tt.expectedBool {
			t.Errorf("ok in %s is wrong: expected=%T, got=%T",
				tt.obj.Inspect(), tt.expectedBool, ok)
		}

		if actual != tt.expectedProp {
			t.Errorf("found prop is wrong: expected=%s, got=%s",
				actual.Inspect(), tt.obj.Inspect())
		}
	}
}
