package object

import (
	"hash/fnv"
)

const STR_TYPE = "STR_TYPE"

type PanStr struct {
	Value string
}

func (s *PanStr) Type() PanObjType {
	return STR_TYPE
}

func (s *PanStr) Inspect() string {
	return s.Value
}

func (s *PanStr) Proto() PanObject {
	return builtInStrObj
}

func (s *PanStr) Hash() HashKey {
	return HashKey{STR_TYPE, s.SymHash()}
}

func (s *PanStr) SymHash() SymHash {
	if symHash, ok := symHashTable[s.Value]; ok {
		return symHash
	}
	h := fnv.New64a()
	h.Write([]byte(s.Value))
	return SymHash(h.Sum64())
}
