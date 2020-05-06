package object

import (
	"testing"
)

func TestObjType(t *testing.T) {
	obj := PanObj{}
	if obj.Type() != OBJ_TYPE {
		t.Fatalf("wrong type: expected=%s, got=%s", OBJ_TYPE, obj.Type())
	}
}

func TestObjInspect(t *testing.T) {
	tests := []struct {
		obj      PanObj
		expected string
	}{
		// keys are sorted so that Inspect() always returns same output
		{
			PanObjInstance(&map[SymHash]Pair{}),
			`{}`,
		},
		{
			PanObjInstance(&map[SymHash]Pair{
				(&PanStr{"a"}).SymHash(): Pair{&PanStr{"a"}, &PanInt{1}},
			}),
			`{"a": 1}`,
		},
		{
			PanObjInstance(&map[SymHash]Pair{
				(&PanStr{"a"}).SymHash():  Pair{&PanStr{"a"}, &PanStr{"A"}},
				(&PanStr{"_b"}).SymHash(): Pair{&PanStr{"_b"}, &PanStr{"B"}},
			}),
			`{"_b": "B", "a": "A"}`,
		},
		{
			PanObjInstance(&map[SymHash]Pair{
				(&PanStr{"foo?"}).SymHash(): Pair{&PanStr{"foo?"}, &PanBool{true}},
				(&PanStr{"b"}).SymHash():    Pair{&PanStr{"b"}, &PanStr{"B"}},
			}),
			`{"b": "B", "foo?": true}`,
		},
		{
			PanObjInstance(&map[SymHash]Pair{
				(&PanStr{"foo?"}).SymHash(): Pair{&PanStr{"foo?"}, &PanBool{true}},
				(&PanStr{"b"}).SymHash(): Pair{
					&PanStr{"b"},
					// NOTE: `&(NewPanObjInstance(...))` is syntax error
					PanObjInstancePtr(&map[SymHash]Pair{
						(&PanStr{"c"}).SymHash(): Pair{&PanStr{"c"}, &PanStr{"C"}},
					}),
				},
			}),
			`{"b": {"c": "C"}, "foo?": true}`,
		},
	}

	for _, tt := range tests {
		if tt.obj.Inspect() != tt.expected {
			t.Errorf("wrong output: expected=%s, got=%s",
				tt.expected, tt.obj.Inspect())
		}
	}
}

func TestObjProto(t *testing.T) {
	tests := []struct {
		obj          PanObject
		expected     PanObject
		expectedName string
	}{
		{PanObjInstancePtr(&map[SymHash]Pair{}), builtInObjObj, "builtInObjObj"},
		{builtInIntObj, builtInNumObj, "builtInNumObj"},
		{builtInFloatObj, builtInNumObj, "builtInNumObj"},
		{builtInObjObj, builtInBaseObj, "builtInBaseObj"},
	}

	for _, tt := range tests {
		if tt.obj.Proto() != tt.expected {
			t.Fatalf("Proto is not %s. got=%T (%+v)",
				tt.expectedName, tt.obj.Proto(), tt.obj.Proto())
		}
	}

}

// checked by compiler (this function works nothing)
func testObjIsPanObject() {
	var _ PanObject = &PanObj{}
}
