package object

import (
	"testing"
)

func TestMapType(t *testing.T) {
	obj := NewEmptyPanMap()
	if obj.Type() != MapType {
		t.Fatalf("wrong type: expected=%s, got=%s", MapType, obj.Type())
	}
}

func TestMapInspect(t *testing.T) {
	tests := []struct {
		obj      *PanMap
		expected string
	}{
		// keys are sorted so that Inspect() always returns same output
		{
			NewEmptyPanMap(),
			`%{}`,
		},
		{
			NewPanMap(
				Pair{NewPanStr("a"), NewPanInt(1)},
			),
			`%{"a": 1}`,
		},
		{
			NewPanMap(
				Pair{NewPanStr("a"), NewPanStr("A")},
				Pair{NewPanStr("_b"), NewPanStr("B")},
			),
			`%{"_b": "B", "a": "A"}`,
		},
		{
			NewPanMap(
				Pair{NewPanStr("foo?"), BuiltInTrue},
				Pair{NewPanStr("b"), NewPanStr("B")},
			),
			`%{"b": "B", "foo?": true}`,
		},
		{
			NewPanMap(
				Pair{NewPanInt(1), NewPanStr("a")},
				Pair{BuiltInTrue, NewPanStr("B")},
			),
			`%{1: "a", true: "B"}`,
		},
		{
			NewPanMap(
				Pair{NewPanStr("foo?"), &PanBool{true}},
				Pair{
					NewPanStr("b"),
					NewPanMap(
						Pair{NewPanStr("c"), NewPanStr("C")},
					),
				},
			),
			`%{"b": %{"c": "C"}, "foo?": true}`,
		},
		// Map can use non-hashable object as key
		// (indexing non-hashable is implemented by
		// one-by-one `==` method comparizon)
		// order of non-hashable pairs is same as struct initialization
		{
			NewPanMap(
				Pair{NewPanArr(NewPanInt(1)), NewPanInt(1)},
			),
			"%{[1]: 1}",
		},
		{
			NewPanMap(
				Pair{NewPanArr(NewPanInt(1), NewPanInt(2)), NewPanInt(1)},
				Pair{
					NewPanMap(
						Pair{NewPanStr("a"), NewPanStr("b")},
					),
					BuiltInFalse,
				},
			),
			`%{[1, 2]: 1, %{"a": "b"}: false}`,
		},
		// order: hashable (sorted by key Inspect), non-hashable
		{
			NewPanMap(
				Pair{NewPanInt(-2), NewPanStr("minus two")},
				Pair{NewPanStr("a"), NewPanStr("A")},
				Pair{NewPanStr("z"), NewPanStr("Z")},
				Pair{
					PanObjInstancePtr(&map[SymHash]Pair{
						(NewPanStr("foo")).SymHash(): {NewPanStr("foo"), NewPanInt(1)},
					}),
					NewPanNil(),
				},
			),
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
	m := NewEmptyPanMap()
	if m.Proto() != BuiltInMapObj {
		t.Fatalf("Proto is not BuiltInMapObj. got=%T (%+v)",
			m.Proto(), m.Proto())
	}
}

// checked by compiler (this function works nothing)
func testMapIsPanObject() {
	var _ PanObject = NewEmptyPanMap()
}

func TestNewPanMapWithScalarKeys(t *testing.T) {
	tests := []struct {
		pairs []Pair
	}{
		{[]Pair{}},
		{[]Pair{
			{NewPanStr("foo"), NewPanInt(2)},
		}},
	}

	for _, tt := range tests {
		actual := NewPanMap(tt.pairs...)

		if len(*actual.Pairs) != len(tt.pairs) {
			t.Fatalf("wrong pair length. expected=%d, got=%d",
				len(tt.pairs), len(*actual.Pairs))
		}

		if len(*actual.NonHashablePairs) != 0 {
			t.Fatalf("wrong NonHashablePair length. expected=%d, got=%d",
				0, len(*actual.NonHashablePairs))
		}

		for _, pair := range tt.pairs {
			h, _ := pair.Key.(PanScalar)
			actualPair, ok := (*actual.Pairs)[h.Hash()]
			if !ok {
				t.Fatalf("key %s is not found.", pair.Key.Inspect())
			}

			if actualPair.Key != pair.Key {
				t.Errorf("wrong key. expected=%s, got=%s",
					pair.Key.Inspect(), actualPair.Key.Inspect())
			}

			if actualPair.Value != pair.Value {
				t.Errorf("wrong value. expected=%s, got=%s",
					pair.Value.Inspect(), actualPair.Value.Inspect())
			}
		}
	}
}

func TestNewPanMapWithNonScalarKeys(t *testing.T) {
	tests := []struct {
		pairs []Pair
	}{
		{[]Pair{
			{NewPanArr(), NewPanInt(2)},
		}},
	}

	for _, tt := range tests {
		actual := NewPanMap(tt.pairs...)

		if len(*actual.Pairs) != 0 {
			t.Fatalf("wrong pair length. expected=%d, got=%d",
				0, len(*actual.Pairs))
		}

		if len(*actual.NonHashablePairs) != len(tt.pairs) {
			t.Fatalf("wrong NonHashablePair length. expected=%d, got=%d",
				len(tt.pairs), len(*actual.NonHashablePairs))
		}

		for i, pair := range tt.pairs {
			actualPair := (*actual.NonHashablePairs)[i]

			if actualPair.Key != pair.Key {
				t.Errorf("wrong key. expected=%s, got=%s",
					pair.Key.Inspect(), actualPair.Key.Inspect())
			}

			if actualPair.Value != pair.Value {
				t.Errorf("wrong value. expected=%s, got=%s",
					pair.Value.Inspect(), actualPair.Value.Inspect())
			}
		}
	}
}

func TestNewEmptyPanMap(t *testing.T) {
	actual := NewEmptyPanMap()

	if len(*actual.Pairs) != 0 {
		t.Fatalf("wrong pair length. expected=%d, got=%d",
			0, len(*actual.Pairs))
	}

	if len(*actual.NonHashablePairs) != 0 {
		t.Fatalf("wrong NonHashablePair length. expected=%d, got=%d",
			0, len(*actual.NonHashablePairs))
	}
}
