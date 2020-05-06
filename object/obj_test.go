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
			PanObj{&map[SymHash]Pair{}},
			`{}`,
		},
		{
			PanObj{&map[SymHash]Pair{
				(&PanStr{"a"}).SymHash(): Pair{&PanStr{"a"}, &PanInt{1}},
			}},
			`{"a": 1}`,
		},
		{
			PanObj{&map[SymHash]Pair{
				(&PanStr{"a"}).SymHash():  Pair{&PanStr{"a"}, &PanStr{"A"}},
				(&PanStr{"_b"}).SymHash(): Pair{&PanStr{"_b"}, &PanStr{"B"}},
			}},
			`{"_b": "B", "a": "A"}`,
		},
		{
			PanObj{&map[SymHash]Pair{
				(&PanStr{"foo?"}).SymHash(): Pair{&PanStr{"foo?"}, &PanBool{true}},
				(&PanStr{"b"}).SymHash():    Pair{&PanStr{"b"}, &PanStr{"B"}},
			}},
			`{"b": "B", "foo?": true}`,
		},
		{
			PanObj{&map[SymHash]Pair{
				(&PanStr{"foo?"}).SymHash(): Pair{&PanStr{"foo?"}, &PanBool{true}},
				(&PanStr{"b"}).SymHash(): Pair{
					&PanStr{"b"},
					&PanObj{&map[SymHash]Pair{
						(&PanStr{"c"}).SymHash(): Pair{&PanStr{"c"}, &PanStr{"C"}},
					}},
				},
			}},
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
	m := PanObj{}
	if m.Proto() != builtInObjObj {
		t.Fatalf("Proto is not BuiltinObjObj. got=%T (%+v)",
			m.Proto(), m.Proto())
	}
}

// checked by compiler (this function works nothing)
func testObjIsPanObject() {
	var _ PanObject = &PanObj{}
}