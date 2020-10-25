// Source code in this file is inherited and modified from
// "Writing an Interpreter in Go" https://interpreterbook.com/
// MIT License | Copyright (c) 2016-2017 Thorsten Ball
// see https://opensource.org/licenses/MIT

package object

import (
	"hash/fnv"
)

// NOTE: hash key map works 2~8 times fast as string key map
// (SHA-1 key is rejected because it is slower than string key)
var symHashTable = make(map[string]SymHash)

// SymHash is symbol hash to refer str literal efficiently.
// This is lighter than HashKey and only for str.
type SymHash = uint64

// StrTable is a table to store PanStr. Str literal can be obtained by symhash.
// NOTE: This table is shared globally.
var StrTable = make(map[SymHash]*PanStr)

// HashKey is a hash for obj/map indexing.
// All Scalar objects have its own hash.
type HashKey struct {
	// to distinguish different type values with same hash
	Type  PanObjType
	Value uint64
}

// GetSymHash gets symbol hash of the string.
func GetSymHash(str string) SymHash {
	if symHash, ok := symHashTable[str]; ok {
		return symHash
	}
	h := fnv.New64a()
	h.Write([]byte(str))
	symHash := h.Sum64()

	symHashTable[str] = symHash

	// set PanStr object corresponding to created hash
	// to generate PanStr from SymHash
	StrTable[symHash] = NewPanStr(str)

	return symHash
}

// SymHash2Str gets str literal from symbol hash.
func SymHash2Str(h SymHash) (PanObject, bool) {
	strObj, ok := StrTable[h]
	return strObj, ok
}
