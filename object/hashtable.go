// Source code in this file is inherited and modified from
// "Writing an Interpreter in Go" https://interpreterbook.com/
// MIT License | Copyright (c) 2016-2017 Thorsten Ball
// see https://opensource.org/licenses/MIT

package object

import (
	"hash/fnv"
	"sync"
)

// NOTE: hash key map works 2~8 times fast as string key map
// (SHA-1 key is rejected because it is slower than string key)
var symHashTable = make(map[string]SymHash)

// NOTE: lock is necessary to make access to symHashTable and strTable goroutine-safe
var lock sync.RWMutex

// SymHash is symbol hash to refer str literal efficiently.
// This is lighter than HashKey and only for str.
type SymHash = uint64

// strTable is a table to store PanStr. Str literal can be obtained by symhash.
var strTable = make(map[SymHash]*PanStr)

// HashKey is a hash for obj/map indexing.
// All Scalar objects have its own hash.
type HashKey struct {
	// to distinguish different type values with same hash
	Type  PanObjType
	Value uint64
}

// GetSymHash gets symbol hash of the string.
func GetSymHash(str string) SymHash {
	if symHash, ok := readSymHash(str); ok {
		return symHash
	}

	h := fnv.New64a()
	h.Write([]byte(str))
	symHash := h.Sum64()
	writeSymHash(symHash, str)
	return symHash
}

func readSymHash(str string) (SymHash, bool) {
	// make table access goroutine-safe (RLock allow other goroutines to read)
	lock.RLock()
	defer lock.RUnlock()
	symHash, ok := symHashTable[str]
	return symHash, ok
}

func writeSymHash(symHash SymHash, str string) {
	// make table access goroutine-safe
	lock.Lock()
	defer lock.Unlock()
	symHashTable[str] = symHash
	// set PanStr object corresponding to created hash
	// to generate PanStr from SymHash
	strTable[symHash] = NewPanStr(str)
}

// SymHash2Str gets str literal from symbol hash.
func SymHash2Str(h SymHash) (PanObject, bool) {
	strObj, ok := strTable[h]
	return strObj, ok
}
