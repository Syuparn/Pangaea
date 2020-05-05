// Source code in this file is inherited and modified from
// "Writing an Interpreter in Go" https://interpreterbook.com/
// MIT License | Copyright (c) 2016-2017 Thorsten Ball
// see https://opensource.org/licenses/MIT

package object

// NOTE: hash key map works 2~8 times fast as string key map
// (SHA-1 key is rejected because it is slower than string key)
var symHashTable = make(map[string]SymHash)

// for only symbols to call props (lighter than HashKey)
type SymHash uint64

type HashKey struct {
	// to distinguish different type values with same hash
	Type  PanObjType
	Value uint64
}
