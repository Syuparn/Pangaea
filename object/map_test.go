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
			PanMap{&map[HashKey]Pair{}, &[]Pair{}},
			`%{}`,
		},
		{
			PanMap{&map[HashKey]Pair{
				(&PanStr{"a"}).Hash(): Pair{&PanStr{"a"}, &PanInt{1}},
			}, &[]Pair{}},
			`%{"a": 1}`,
		},
		{
			PanMap{&map[HashKey]Pair{
				(&PanStr{"a"}).Hash():  Pair{&PanStr{"a"}, &PanStr{"A"}},
				(&PanStr{"_b"}).Hash(): Pair{&PanStr{"_b"}, &PanStr{"B"}},
			}, &[]Pair{}},
			`%{"_b": "B", "a": "A"}`,
		},
		{
			PanMap{&map[HashKey]Pair{
				(&PanStr{"foo?"}).Hash(): Pair{&PanStr{"foo?"}, &PanBool{true}},
				(&PanStr{"b"}).Hash():    Pair{&PanStr{"b"}, &PanStr{"B"}},
			}, &[]Pair{}},
			`%{"b": "B", "foo?": true}`,
		},
		{
			PanMap{&map[HashKey]Pair{
				(&PanInt{1}).Hash():     Pair{&PanInt{1}, &PanStr{"a"}},
				(&PanBool{true}).Hash(): Pair{&PanBool{true}, &PanStr{"B"}},
			}, &[]Pair{}},
			`%{1: "a", true: "B"}`,
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
			}, &[]Pair{}},
			`%{"b": {"c": "C"}, "foo?": true}`,
		},
		{
			PanMap{&map[HashKey]Pair{
				(&PanStr{"foo?"}).Hash(): Pair{&PanStr{"foo?"}, &PanBool{true}},
				(&PanStr{"b"}).Hash(): Pair{
					&PanStr{"b"},
					&PanMap{&map[HashKey]Pair{
						(&PanStr{"c"}).Hash(): Pair{&PanStr{"c"}, &PanStr{"C"}},
					}, &[]Pair{}},
				},
			}, &[]Pair{}},
			`%{"b": %{"c": "C"}, "foo?": true}`,
		},
		// Map can use non-hashable object as key
		// (indexing non-hashable is implemented by
		// one-by-one `==` method comparizon)
		// order of non-hashable pairs is same as struct initialization
		{
			PanMap{
				&map[HashKey]Pair{},
				&[]Pair{
					Pair{&PanArr{[]PanObject{&PanInt{1}}}, &PanInt{1}},
				},
			},
			"%{[1]: 1}",
		},
		{
			PanMap{
				&map[HashKey]Pair{},
				&[]Pair{
					Pair{&PanArr{[]PanObject{&PanInt{1}, &PanInt{2}}}, &PanInt{1}},
					Pair{
						&PanMap{
							&map[HashKey]Pair{
								(&PanStr{"a"}).Hash(): Pair{&PanStr{"a"}, &PanStr{"b"}},
							},
							&[]Pair{},
						},
						&PanBool{false},
					},
				},
			},
			`%{[1, 2]: 1, %{"a": "b"}: false}`,
		},
		// order: hashable (sorted by key Inspect), non-hashable
		{
			PanMap{
				&map[HashKey]Pair{
					(&PanInt{-2}).Hash():  Pair{&PanInt{-2}, &PanStr{"minus two"}},
					(&PanStr{"a"}).Hash(): Pair{&PanStr{"a"}, &PanStr{"A"}},
					(&PanStr{"z"}).Hash(): Pair{&PanStr{"z"}, &PanStr{"Z"}},
				},
				&[]Pair{
					Pair{
						PanObjInstancePtr(&map[SymHash]Pair{
							(&PanStr{"foo"}).SymHash(): Pair{&PanStr{"foo"}, &PanInt{1}},
						}),
						&PanNil{},
					},
				},
			},
			`%{"a": "A", "z": "Z", -2: "minus two", {"foo": 1}: nil}`,
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
