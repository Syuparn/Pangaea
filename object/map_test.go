package object

import (
	"testing"
)

func TestMapType(t *testing.T) {
	obj := PanMap{}
	if obj.Type() != MAP_TYPE {
		t.Fatalf("wrong type: expected=%s, got=%s", MAP_TYPE, obj.Type())
	}
}

func TestMapInspect(t *testing.T) {
	tests := []struct {
		obj      PanMap
		expected string
	}{
		// keys are sorted so that Inspect() always returns same output
		{
			PanMap{&map[HashKey]Pair{}},
			`%{}`,
		},
		{
			PanMap{&map[HashKey]Pair{
				(&PanStr{"a"}).Hash(): Pair{&PanStr{"a"}, &PanInt{1}},
			}},
			`%{"a": 1}`,
		},
		{
			PanMap{&map[HashKey]Pair{
				(&PanStr{"a"}).Hash():  Pair{&PanStr{"a"}, &PanStr{"A"}},
				(&PanStr{"_b"}).Hash(): Pair{&PanStr{"_b"}, &PanStr{"B"}},
			}},
			`%{"_b": "B", "a": "A"}`,
		},
		{
			PanMap{&map[HashKey]Pair{
				(&PanStr{"foo?"}).Hash(): Pair{&PanStr{"foo?"}, &PanBool{true}},
				(&PanStr{"b"}).Hash():    Pair{&PanStr{"b"}, &PanStr{"B"}},
			}},
			`%{"b": "B", "foo?": true}`,
		},
		{
			PanMap{&map[HashKey]Pair{
				(&PanStr{"foo?"}).Hash(): Pair{&PanStr{"foo?"}, &PanBool{true}},
				(&PanStr{"b"}).Hash(): Pair{
					&PanStr{"b"},
					// NOTE: `&(NewPanObjInstance(...))` is syntax error
					PanObjInstancePtr(&map[SymHash]Pair{
						(&PanStr{"c"}).SymHash(): Pair{&PanStr{"c"}, &PanStr{"C"}},
					}),
				},
			}},
			`%{"b": {"c": "C"}, "foo?": true}`,
		},
		{
			PanMap{&map[HashKey]Pair{
				(&PanStr{"foo?"}).Hash(): Pair{&PanStr{"foo?"}, &PanBool{true}},
				(&PanStr{"b"}).Hash(): Pair{
					&PanStr{"b"},
					&PanMap{&map[HashKey]Pair{
						(&PanStr{"c"}).Hash(): Pair{&PanStr{"c"}, &PanStr{"C"}},
					}},
				},
			}},
			`%{"b": %{"c": "C"}, "foo?": true}`,
		},
	}

	for _, tt := range tests {
		if tt.obj.Inspect() != tt.expected {
			t.Errorf("wrong output: expected=%s, got=%s",
				tt.expected, tt.obj.Inspect())
		}
	}
}

func TestMapProto(t *testing.T) {
	m := PanMap{}
	if m.Proto() != BuiltInMapObj {
		t.Fatalf("Proto is not BuiltInMapObj. got=%T (%+v)",
			m.Proto(), m.Proto())
	}
}

// checked by compiler (this function works nothing)
func testMapIsPanObject() {
	var _ PanObject = &PanMap{}
}
