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

func TestMapRepr(t *testing.T) {
	tests := []struct {
		obj      *PanMap
		expected string
	}{
		// keys are sorted so that Repr() always returns same output
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
		// if key or value has prop _name, show it instead
		{
			NewPanMap(
				Pair{
					PanObjInstancePtr(&map[SymHash]Pair{
						GetSymHash("_name"): {NewPanStr("_name"), NewPanStr("foo")},
					}),
					NewPanInt(1),
				},
			),
			`%{foo: 1}`,
		},
		{
			NewPanMap(
				Pair{
					NewPanInt(2),
					PanObjInstancePtr(&map[SymHash]Pair{
						GetSymHash("_name"): {NewPanStr("_name"), NewPanStr("bar")},
					}),
				},
			),
			`%{2: bar}`,
		},
	}

	for _, tt := range tests {
		actual := tt.obj.Repr()
		if actual != tt.expected {
			t.Errorf("wrong output: expected=%s, got=%s",
				tt.expected, actual)
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

func TestInheritedMapProto(t *testing.T) {
	mapChild := ChildPanObjPtr(BuiltInMapObj, EmptyPanObjPtr())
	a := NewInheritedMap(mapChild)
	if a.Proto() != mapChild {
		t.Fatalf("Proto is not mapChild. got=%T (%s)",
			a.Proto(), a.Proto().Inspect())
	}
}

func TestMapZero(t *testing.T) {
	tests := []struct {
		name string
		obj  *PanMap
	}{
		{"%{'a: 1}", NewPanMap(Pair{NewPanStr("a"), NewPanInt(1)})},
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
		{[]Pair{
			{NewPanStr("foo"), NewPanInt(2)},
			{NewPanStr("bar"), BuiltInNil},
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

		if len(*actual.HashKeys) != len(tt.pairs) {
			t.Fatalf("wrong HashKeys length. expected=%d, got=%d",
				len(tt.pairs), len(*actual.HashKeys))
		}

		for i, pair := range tt.pairs {
			h, _ := pair.Key.(PanScalar)

			// hashkey order check
			if actualHash := (*actual.HashKeys)[i]; actualHash != h.Hash() {
				t.Errorf("wrong HashKeys[%d]. expected=%v, got=%v",
					i, h.Hash(), actualHash)
			}

			// pair existence check
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

		if len(*actual.HashKeys) != 0 {
			t.Fatalf("wrong HashKeys length. expected=%d, got=%d",
				0, len(*actual.HashKeys))
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

	if len(*actual.HashKeys) != 0 {
		t.Fatalf("wrong HashKeys length. expected=%d, got=%d",
			0, len(*actual.HashKeys))
	}
}

func TestNewPanMapWithDuplicatedKeys(t *testing.T) {
	tests := []struct {
		pairs                    []Pair
		expectedHashKeys         []HashKey
		expectedPairs            map[HashKey]Pair
		expectedNonHashablePairs []Pair
	}{
		// if same keys are passed, set only first one
		{
			[]Pair{
				{NewPanStr("a"), NewPanInt(1)},
				{NewPanStr("a"), NewPanInt(2)},
			},
			[]HashKey{
				NewPanStr("a").Hash(),
			},
			map[HashKey]Pair{
				NewPanStr("a").Hash(): {NewPanStr("a"), NewPanInt(1)},
			},
			[]Pair{},
		},
		// NOTE: duplication of non-hashable keys is not checked!
		// (because '== method comparison is required)
	}

	for _, tt := range tests {
		actual := NewPanMap(tt.pairs...)

		if len(*actual.HashKeys) != len(tt.expectedHashKeys) {
			t.Fatalf("wrong hashKey length. expected=%d, got=%d",
				len(tt.expectedHashKeys), len(*actual.HashKeys))
		}

		if len(*actual.Pairs) != len(tt.expectedPairs) {
			t.Fatalf("wrong pair length. expected=%d, got=%d",
				len(tt.expectedPairs), len(*actual.Pairs))
		}

		if len(*actual.NonHashablePairs) != len(tt.expectedNonHashablePairs) {
			t.Fatalf("wrong NonHashablePair length. expected=%d, got=%d",
				len(tt.expectedNonHashablePairs), len(*actual.NonHashablePairs))
		}

		for i, h := range tt.expectedHashKeys {
			actualHash := (*actual.HashKeys)[i]

			if actualHash != h {
				t.Errorf("wrong hashkeys[%d]. expected=%v, got=%v",
					i, h, actualHash)
			}
		}

		for hash, pair := range tt.expectedPairs {
			actualPair, ok := (*actual.Pairs)[hash]

			if !ok {
				t.Fatalf("key %s is not found in pair", pair.Key.Inspect())
			}

			if actualPair.Key.(PanScalar).Hash() != pair.Key.(PanScalar).Hash() {
				t.Errorf("wrong key. expected=%s, got=%s",
					pair.Key.Inspect(), actualPair.Key.Inspect())
			}

			if actualPair.Value != pair.Value {
				t.Errorf("wrong value. expected=%s, got=%s",
					pair.Value.Inspect(), actualPair.Value.Inspect())
			}
		}

		for i, pair := range tt.expectedNonHashablePairs {
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

func TestNewInheritedMap(t *testing.T) {
	// child of Map
	proto := ChildPanObjPtr(BuiltInMapObj, EmptyPanObjPtr())

	tests := []struct {
		pairs []Pair
	}{
		{[]Pair{}},
		{[]Pair{
			{NewPanStr("foo"), NewPanInt(2)},
		}},
		{[]Pair{
			{NewPanStr("foo"), NewPanInt(2)},
			{NewPanStr("bar"), BuiltInNil},
		}},
	}

	for _, tt := range tests {
		actual := NewInheritedMap(proto, tt.pairs...)

		if len(*actual.Pairs) != len(tt.pairs) {
			t.Fatalf("wrong pair length. expected=%d, got=%d",
				len(tt.pairs), len(*actual.Pairs))
		}

		if len(*actual.NonHashablePairs) != 0 {
			t.Fatalf("wrong NonHashablePair length. expected=%d, got=%d",
				0, len(*actual.NonHashablePairs))
		}

		if len(*actual.HashKeys) != len(tt.pairs) {
			t.Fatalf("wrong HashKeys length. expected=%d, got=%d",
				len(tt.pairs), len(*actual.HashKeys))
		}

		for i, pair := range tt.pairs {
			h, _ := pair.Key.(PanScalar)

			// hashkey order check
			if actualHash := (*actual.HashKeys)[i]; actualHash != h.Hash() {
				t.Errorf("wrong HashKeys[%d]. expected=%v, got=%v",
					i, h.Hash(), actualHash)
			}

			// pair existence check
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
