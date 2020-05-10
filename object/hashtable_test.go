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
