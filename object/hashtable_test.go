package object

import (
	"testing"
)

func TestGetSymHash(t *testing.T) {
	hash1 := GetSymHash("foo")
	hash2 := GetSymHash("foo")

	if hash1 != hash2 {
		t.Errorf("hash of same string must be same. expected=%v, got=%v",
			hash1, hash2)
	}
}

func TestSymHashTable(t *testing.T) {
	str := "foo"
	hash1 := GetSymHash(str)

	found, ok := symHashTable[str]

	if !ok {
		t.Fatalf("hash not found in SymHashTable(%+v)",
			symHashTable)
	}

	if found != hash1 {
		t.Errorf("wrong output. expected=%v, got=%v",
			str, found)
	}
}

func TestSymHash2Str(t *testing.T) {
	str := "foo"
	hash := GetSymHash(str)

	strObj1 := SymHash2Str(hash)
	strObj2 := SymHash2Str(hash)

	if strObj1 != strObj2 {
		t.Errorf("strObj1 and strObj2 must be same. strObj1=%v, strObj2=%v",
			strObj1, strObj2)
	}

	if strObj1.Type() != STR_TYPE {
		t.Fatalf("type of strObj1 must be STR_TYPE. got=%s",
			strObj1.Type())
	}

	if strObj1.Inspect() != str {
		t.Errorf("wrong output. expected=%v, got=%v",
			str, strObj1.Inspect())
	}
}

func TestStrTable(t *testing.T) {
	str := "foo"
	hash := GetSymHash(str)

	found, ok := StrTable[hash]

	if !ok {
		t.Fatalf("str not found in StrTable(%+v)",
			StrTable)
	}

	if found.Type() != STR_TYPE {
		t.Fatalf("StrTable must return STR_TYPE PanObject. got=%s",
			found.Type())
	}

	if found.Inspect() != str {
		t.Errorf("wrong output. expected=%v, got=%v",
			str, found.Inspect())
	}
}