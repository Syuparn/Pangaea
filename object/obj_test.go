package object

import (
	"testing"
)

func TestObjType(t *testing.T) {
	obj := PanObj{}
	if obj.Type() != ObjType {
		t.Fatalf("wrong type: expected=%s, got=%s", ObjType, obj.Type())
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
				(NewPanStr("a")).SymHash(): {NewPanStr("a"), &PanInt{1}},
			}),
			`{"a": 1}`,
		},
		{
			PanObjInstance(&map[SymHash]Pair{
				(NewPanStr("a")).SymHash():  {NewPanStr("a"), NewPanStr("A")},
				(NewPanStr("_b")).SymHash(): {NewPanStr("_b"), NewPanStr("B")},
			}),
			`{"_b": "B", "a": "A"}`,
		},
		{
			PanObjInstance(&map[SymHash]Pair{
				(NewPanStr("foo?")).SymHash(): {NewPanStr("foo?"), &PanBool{true}},
				(NewPanStr("b")).SymHash():    {NewPanStr("b"), NewPanStr("B")},
			}),
			`{"b": "B", "foo?": true}`,
		},
		{
			PanObjInstance(&map[SymHash]Pair{
				(NewPanStr("foo?")).SymHash(): {NewPanStr("foo?"), &PanBool{true}},
				(NewPanStr("b")).SymHash(): {
					NewPanStr("b"),
					// NOTE: `&(NewPanObjInstance(...))` is syntax error
					PanObjInstancePtr(&map[SymHash]Pair{
						(NewPanStr("c")).SymHash(): {NewPanStr("c"), NewPanStr("C")},
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

func TestObjAddPairs(t *testing.T) {
	testPairA := Pair{NewPanStr("a"), NewPanInt(1)}
	testPairB := Pair{NewPanStr("b"), NewPanInt(2)}
	testPairC := Pair{NewPanStr("c"), NewPanInt(3)}
	testPairC2 := Pair{NewPanStr("c"), NewPanInt(4)}

	tests := []struct {
		obj      *PanObj
		added    map[SymHash]Pair
		expected *PanObj
	}{
		{
			PanObjInstancePtr(&map[SymHash]Pair{}).(*PanObj),
			map[SymHash]Pair{},
			PanObjInstancePtr(&map[SymHash]Pair{}).(*PanObj),
		},
		{
			PanObjInstancePtr(&map[SymHash]Pair{
				(NewPanStr("a")).SymHash(): testPairA,
				(NewPanStr("c")).SymHash(): testPairC,
			}).(*PanObj),
			map[SymHash]Pair{
				(NewPanStr("b")).SymHash(): testPairB,
				(NewPanStr("c")).SymHash(): testPairC2,
			},
			PanObjInstancePtr(&map[SymHash]Pair{
				(NewPanStr("a")).SymHash(): testPairA,
				(NewPanStr("b")).SymHash(): testPairB,
				(NewPanStr("c")).SymHash(): testPairC,
			}).(*PanObj),
		},
	}

	for _, tt := range tests {
		err := tt.obj.AddPairs(&tt.added)
		if err != nil {
			t.Fatalf("err must be nil. got=%s", err.Error())
		}

		// check keys
		if len(*tt.obj.Keys) != len(*tt.expected.Keys) {
			t.Fatalf("wrong keys length: expected=%d, got=%d",
				len(*tt.expected.Keys), len(*tt.obj.Keys))
		}

		for i, key := range *tt.obj.Keys {
			if key != (*tt.expected.Keys)[i] {
				t.Errorf("keys[%d] in %s is wrong. expected=%d, got=%d",
					i, tt.obj.Inspect(), (*tt.expected.Keys)[i], key)
			}
		}

		// check PrivateKeys
		if len(*tt.obj.PrivateKeys) != len(*tt.expected.PrivateKeys) {
			t.Fatalf("wrong keys length: expected=%d, got=%d",
				len(*tt.expected.PrivateKeys), len(*tt.obj.PrivateKeys))
		}

		for i, key := range *tt.obj.PrivateKeys {
			if key != (*tt.expected.PrivateKeys)[i] {
				t.Errorf("Privatekeys[%d] in %s is wrong. expected=%d, got=%d",
					i, tt.obj.Inspect(), (*tt.expected.PrivateKeys)[i], key)
			}
		}

		// check Pairs
		if len(*tt.obj.Pairs) != len(*tt.expected.Pairs) {
			t.Fatalf("wrong pairs length: expected=%d, got=%d",
				len(*tt.expected.Pairs), len(*tt.obj.Pairs))
		}

		for sym, expected := range *tt.expected.Pairs {
			actual, ok := (*tt.obj.Pairs)[sym]
			if !ok {
				sym, _ := SymHash2Str(sym)
				t.Fatalf("key %s not found.", sym)
			}
			if actual != expected {
				t.Errorf("value must be %s. got=%s",
					expected.Value.Inspect(), actual.Value.Inspect())
			}
		}
	}
}

func TestObjAddPairsErr(t *testing.T) {
	obj := EmptyPanObjPtr()
	err := obj.AddPairs(nil)
	if err == nil {
		t.Fatalf("err must not be nil")
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
				GetSymHash("a"): {Key: NewPanStr("a"), Value: NewPanInt(1)},
			}),
			[]SymHash{
				GetSymHash("a"),
			},
		},
		// keys are ordered alphabetically
		{
			PanObjInstancePtr(&map[SymHash]Pair{
				GetSymHash("a"): {Key: NewPanStr("a"), Value: NewPanInt(1)},
				GetSymHash("b"): {Key: NewPanStr("b"), Value: NewPanInt(2)},
			}),
			[]SymHash{
				GetSymHash("a"),
				GetSymHash("b"),
			},
		},
		{
			PanObjInstancePtr(&map[SymHash]Pair{
				GetSymHash("c"): {Key: NewPanStr("c"), Value: NewPanInt(1)},
				GetSymHash("b"): {Key: NewPanStr("b"), Value: NewPanInt(1)},
			}),
			[]SymHash{
				GetSymHash("b"),
				GetSymHash("c"),
			},
		},
		// keys only contain public str
		{
			PanObjInstancePtr(&map[SymHash]Pair{
				GetSymHash("a"):  {Key: NewPanStr("a"), Value: NewPanInt(1)},
				GetSymHash("=="): {Key: NewPanStr("=="), Value: NewPanInt(2)},
				GetSymHash("_b"): {Key: NewPanStr("_b"), Value: NewPanInt(3)},
				GetSymHash("c "): {Key: NewPanStr("c "), Value: NewPanInt(4)},
				GetSymHash("d"):  {Key: NewPanStr("d"), Value: NewPanInt(5)},
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
				GetSymHash("a"): {Key: NewPanStr("a"), Value: NewPanInt(1)},
			}),
			[]SymHash{},
		},
		// keys are ordered alphabetically
		{
			PanObjInstancePtr(&map[SymHash]Pair{
				GetSymHash("_a"): {Key: NewPanStr("_a"), Value: NewPanInt(1)},
				GetSymHash("_b"): {Key: NewPanStr("_b"), Value: NewPanInt(2)},
			}),
			[]SymHash{
				GetSymHash("_a"),
				GetSymHash("_b"),
			},
		},
		{
			PanObjInstancePtr(&map[SymHash]Pair{
				GetSymHash("_c"): {Key: NewPanStr("_c"), Value: NewPanInt(1)},
				GetSymHash("_b"): {Key: NewPanStr("_b"), Value: NewPanInt(1)},
			}),
			[]SymHash{
				GetSymHash("_b"),
				GetSymHash("_c"),
			},
		},
		// private keys only contain private str
		{
			PanObjInstancePtr(&map[SymHash]Pair{
				GetSymHash("a"):  {Key: NewPanStr("a"), Value: NewPanInt(1)},
				GetSymHash("=="): {Key: NewPanStr("=="), Value: NewPanInt(2)},
				GetSymHash("_b"): {Key: NewPanStr("_b"), Value: NewPanInt(3)},
				GetSymHash("c "): {Key: NewPanStr("c "), Value: NewPanInt(4)},
				GetSymHash("d"):  {Key: NewPanStr("d"), Value: NewPanInt(5)},
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

func TestEmptyPanObjPtr(t *testing.T) {
	actual := EmptyPanObjPtr()

	if actual.Pairs == nil {
		t.Fatal("Pairs must not be nil.")
	}
	if len(*actual.Pairs) != 0 {
		t.Errorf("len(Pairs) must be 0. got=%d", len(*actual.Pairs))
	}

	if actual.Keys == nil {
		t.Fatal("Keys must not be nil.")
	}
	if len(*actual.Keys) != 0 {
		t.Errorf("len(Keys) must be 0. got=%d", len(*actual.Keys))
	}

	if actual.PrivateKeys == nil {
		t.Fatal("PrivateKeys must not be nil.")
	}
	if len(*actual.PrivateKeys) != 0 {
		t.Errorf("len(PrivateKeys) must be 0. got=%d", len(*actual.PrivateKeys))
	}

	if actual.Proto() != BuiltInObjObj {
		t.Errorf("Proto must be BuiltInObjObj. got=%s", actual.Proto().Inspect())
	}
}

func TestChildPanObjPtr(t *testing.T) {
	tests := []struct {
		proto *PanObj
		src   *PanObj
	}{
		{
			PanObjInstancePtr(&map[SymHash]Pair{
				GetSymHash("a"): {Key: NewPanStr("a"), Value: NewPanInt(1)},
			}).(*PanObj),
			PanObjInstancePtr(&map[SymHash]Pair{
				GetSymHash("b"): {Key: NewPanStr("b"), Value: NewPanInt(2)},
			}).(*PanObj),
		},
	}

	for _, tt := range tests {
		actual := ChildPanObjPtr(tt.proto, tt.src)

		if actual.Proto() != tt.proto {
			t.Errorf("Proto must be tt.proto. got=%s", actual.Proto().Inspect())
		}

		if actual.Pairs != tt.src.Pairs {
			t.Errorf("Pairs must be same as src. expected=%s(%+v), got=%s(%+v)",
				tt.src.Inspect(), tt.src.Pairs, actual.Inspect(), actual.Pairs)
		}

		if actual.Keys != tt.src.Keys {
			t.Errorf("Keys must be same as src. expected=%s(%+v), got=%s(%+v)",
				tt.src.Inspect(), tt.src.Keys, actual.Inspect(), actual.Keys)
		}

		if actual.PrivateKeys != tt.src.PrivateKeys {
			t.Errorf("PrivateKeys must be same as src. expected=%s(%+v), got=%s(%+v)",
				tt.src.Inspect(), tt.src.PrivateKeys,
				actual.Inspect(), actual.PrivateKeys)
		}
	}
}

// checked by compiler (this function works nothing)
func testObjIsPanObject() {
	var _ PanObject = &PanObj{}
}
