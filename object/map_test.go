package object

import (
	"testing"
)

func TestMapType(t *testing.T) {
	obj := PanMap{}
	if obj.Type() != MapType {
		t.Fatalf("wrong type: expected=%s, got=%s", MapType, obj.Type())
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
				(NewPanStr("a")).Hash(): Pair{NewPanStr("a"), &PanInt{1}},
			}, &[]Pair{}},
			`%{"a": 1}`,
		},
		{
			PanMap{&map[HashKey]Pair{
				(NewPanStr("a")).Hash():  Pair{NewPanStr("a"), NewPanStr("A")},
				(NewPanStr("_b")).Hash(): Pair{NewPanStr("_b"), NewPanStr("B")},
			}, &[]Pair{}},
			`%{"_b": "B", "a": "A"}`,
		},
		{
			PanMap{&map[HashKey]Pair{
				(NewPanStr("foo?")).Hash(): Pair{NewPanStr("foo?"), &PanBool{true}},
				(NewPanStr("b")).Hash():    Pair{NewPanStr("b"), NewPanStr("B")},
			}, &[]Pair{}},
			`%{"b": "B", "foo?": true}`,
		},
		{
			PanMap{&map[HashKey]Pair{
				(&PanInt{1}).Hash():     Pair{&PanInt{1}, NewPanStr("a")},
				(&PanBool{true}).Hash(): Pair{&PanBool{true}, NewPanStr("B")},
			}, &[]Pair{}},
			`%{1: "a", true: "B"}`,
		},
		{
			PanMap{&map[HashKey]Pair{
				(NewPanStr("foo?")).Hash(): Pair{NewPanStr("foo?"), &PanBool{true}},
				(NewPanStr("b")).Hash(): Pair{
					NewPanStr("b"),
					// NOTE: `&(NewPanObjInstance(...))` is syntax error
					PanObjInstancePtr(&map[SymHash]Pair{
						(NewPanStr("c")).SymHash(): Pair{NewPanStr("c"), NewPanStr("C")},
					}),
				},
			}, &[]Pair{}},
			`%{"b": {"c": "C"}, "foo?": true}`,
		},
		{
			PanMap{&map[HashKey]Pair{
				(NewPanStr("foo?")).Hash(): Pair{NewPanStr("foo?"), &PanBool{true}},
				(NewPanStr("b")).Hash(): Pair{
					NewPanStr("b"),
					&PanMap{&map[HashKey]Pair{
						(NewPanStr("c")).Hash(): Pair{NewPanStr("c"), NewPanStr("C")},
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
								(NewPanStr("a")).Hash(): Pair{NewPanStr("a"), NewPanStr("b")},
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
					(&PanInt{-2}).Hash():    Pair{&PanInt{-2}, NewPanStr("minus two")},
					(NewPanStr("a")).Hash(): Pair{NewPanStr("a"), NewPanStr("A")},
					(NewPanStr("z")).Hash(): Pair{NewPanStr("z"), NewPanStr("Z")},
				},
				&[]Pair{
					Pair{
						PanObjInstancePtr(&map[SymHash]Pair{
							(NewPanStr("foo")).SymHash(): Pair{NewPanStr("foo"), &PanInt{1}},
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
