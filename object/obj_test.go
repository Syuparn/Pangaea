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
				(NewPanStr("a")).SymHash(): Pair{NewPanStr("a"), &PanInt{1}},
			}),
			`{"a": 1}`,
		},
		{
			PanObjInstance(&map[SymHash]Pair{
				(NewPanStr("a")).SymHash():  Pair{NewPanStr("a"), NewPanStr("A")},
				(NewPanStr("_b")).SymHash(): Pair{NewPanStr("_b"), NewPanStr("B")},
			}),
			`{"_b": "B", "a": "A"}`,
		},
		{
			PanObjInstance(&map[SymHash]Pair{
				(NewPanStr("foo?")).SymHash(): Pair{NewPanStr("foo?"), &PanBool{true}},
				(NewPanStr("b")).SymHash():    Pair{NewPanStr("b"), NewPanStr("B")},
			}),
			`{"b": "B", "foo?": true}`,
		},
		{
			PanObjInstance(&map[SymHash]Pair{
				(NewPanStr("foo?")).SymHash(): Pair{NewPanStr("foo?"), &PanBool{true}},
				(NewPanStr("b")).SymHash(): Pair{
					NewPanStr("b"),
					// NOTE: `&(NewPanObjInstance(...))` is syntax error
					PanObjInstancePtr(&map[SymHash]Pair{
						(NewPanStr("c")).SymHash(): Pair{NewPanStr("c"), NewPanStr("C")},
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
		{PanObjInstancePtr(&map[SymHash]Pair{}), BuiltInObjObj, "BuiltInObjObj"},
		{BuiltInIntObj, BuiltInNumObj, "BuiltInNumObj"},
		{BuiltInFloatObj, BuiltInNumObj, "BuiltInNumObj"},
		{BuiltInObjObj, BuiltInBaseObj, "BuiltInBaseObj"},
	}

	for _, tt := range tests {
		if tt.obj.Proto() != tt.expected {
			t.Fatalf("Proto is not %s. got=%T (%+v)",
				tt.expectedName, tt.obj.Proto(), tt.obj.Proto())
		}
	}
}

func TestObjKeys(t *testing.T) {
	tests := []struct {
		obj      PanObject
		expected []SymHash
	}{
		{
			PanObjInstancePtr(&map[SymHash]Pair{}),
			[]SymHash{},
		},
		{
			PanObjInstancePtr(&map[SymHash]Pair{
				GetSymHash("a"): Pair{Key: NewPanStr("a"), Value: NewPanInt(1)},
			}),
			[]SymHash{
				GetSymHash("a"),
			},
		},
		// keys are ordered alphabetically
		{
			PanObjInstancePtr(&map[SymHash]Pair{
				GetSymHash("a"): Pair{Key: NewPanStr("a"), Value: NewPanInt(1)},
				GetSymHash("b"): Pair{Key: NewPanStr("b"), Value: NewPanInt(2)},
			}),
			[]SymHash{
				GetSymHash("a"),
				GetSymHash("b"),
			},
		},
		{
			PanObjInstancePtr(&map[SymHash]Pair{
				GetSymHash("c"): Pair{Key: NewPanStr("c"), Value: NewPanInt(1)},
				GetSymHash("b"): Pair{Key: NewPanStr("b"), Value: NewPanInt(1)},
			}),
			[]SymHash{
				GetSymHash("b"),
				GetSymHash("c"),
			},
		},
		// keys only contain public str
		{
			PanObjInstancePtr(&map[SymHash]Pair{
				GetSymHash("a"):  Pair{Key: NewPanStr("a"), Value: NewPanInt(1)},
				GetSymHash("=="): Pair{Key: NewPanStr("=="), Value: NewPanInt(2)},
				GetSymHash("_b"): Pair{Key: NewPanStr("_b"), Value: NewPanInt(3)},
				GetSymHash("c "): Pair{Key: NewPanStr("c "), Value: NewPanInt(4)},
				GetSymHash("d"):  Pair{Key: NewPanStr("d"), Value: NewPanInt(5)},
			}),
			[]SymHash{
				GetSymHash("a"),
				GetSymHash("d"),
			},
		},
	}

	for _, tt := range tests {
		obj, ok := tt.obj.(*PanObj)
		if !ok {
			t.Fatalf("obj is not PanObj. got=%T(%s)",
				tt.obj, tt.obj.Type())
		}

		if len(*obj.Keys) != len(tt.expected) {
			t.Fatalf("wrong keys length: expected=%d, got=%d",
				len(tt.expected), len(*obj.Keys))
		}

		for i, key := range *obj.Keys {
			if key != tt.expected[i] {
				t.Errorf("keys[%d] in %s is wrong. expected=%d, got=%d",
					i, tt.obj.Inspect(), tt.expected[i], key)
			}
		}
	}
}

func TestObjPrivateKeys(t *testing.T) {
	tests := []struct {
		obj      PanObject
		expected []SymHash
	}{
		{
			PanObjInstancePtr(&map[SymHash]Pair{}),
			[]SymHash{},
		},
		{
			PanObjInstancePtr(&map[SymHash]Pair{
				GetSymHash("a"): Pair{Key: NewPanStr("a"), Value: NewPanInt(1)},
			}),
			[]SymHash{},
		},
		// keys are ordered alphabetically
		{
			PanObjInstancePtr(&map[SymHash]Pair{
				GetSymHash("_a"): Pair{Key: NewPanStr("_a"), Value: NewPanInt(1)},
				GetSymHash("_b"): Pair{Key: NewPanStr("_b"), Value: NewPanInt(2)},
			}),
			[]SymHash{
				GetSymHash("_a"),
				GetSymHash("_b"),
			},
		},
		{
			PanObjInstancePtr(&map[SymHash]Pair{
				GetSymHash("_c"): Pair{Key: NewPanStr("_c"), Value: NewPanInt(1)},
				GetSymHash("_b"): Pair{Key: NewPanStr("_b"), Value: NewPanInt(1)},
			}),
			[]SymHash{
				GetSymHash("_b"),
				GetSymHash("_c"),
			},
		},
		// private keys only contain private str
		{
			PanObjInstancePtr(&map[SymHash]Pair{
				GetSymHash("a"):  Pair{Key: NewPanStr("a"), Value: NewPanInt(1)},
				GetSymHash("=="): Pair{Key: NewPanStr("=="), Value: NewPanInt(2)},
				GetSymHash("_b"): Pair{Key: NewPanStr("_b"), Value: NewPanInt(3)},
				GetSymHash("c "): Pair{Key: NewPanStr("c "), Value: NewPanInt(4)},
				GetSymHash("d"):  Pair{Key: NewPanStr("d"), Value: NewPanInt(5)},
			}),
			[]SymHash{
				GetSymHash("=="),
				GetSymHash("_b"),
				GetSymHash("c "),
			},
		},
	}

	for _, tt := range tests {
		obj, ok := tt.obj.(*PanObj)
		if !ok {
			t.Fatalf("obj is not PanObj. got=%T(%s)",
				tt.obj, tt.obj.Type())
		}

		if len(*obj.PrivateKeys) != len(tt.expected) {
			t.Fatalf("wrong keys length: expected=%d, got=%d",
				len(tt.expected), len(*obj.PrivateKeys))
		}

		for i, key := range *obj.PrivateKeys {
			if key != tt.expected[i] {
				t.Errorf("keys[%d] in %s is wrong. expected=%d, got=%d",
					i, tt.obj.Inspect(), tt.expected[i], key)
			}
		}
	}
}

// checked by compiler (this function works nothing)
func testObjIsPanObject() {
	var _ PanObject = &PanObj{}
}
