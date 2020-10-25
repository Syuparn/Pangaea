package object

import (
	"regexp"
)

// StrType is a type of PanStr.
const StrType = "StrType"

// PanStr is object of str literal.
type PanStr struct {
	Value    string
	IsPublic bool
}

// Type returns type of this PanObject.
func (s *PanStr) Type() PanObjType {
	return StrType
}

// Inspect returns formatted source code of this object.
func (s *PanStr) Inspect() string {
	return `"` + s.Value + `"`
}

// Proto returns proto of this object.
func (s *PanStr) Proto() PanObject {
	return BuiltInStrObj
}

// Hash returns hashkey of this object.
func (s *PanStr) Hash() HashKey {
	return HashKey{StrType, s.SymHash()}
}

// SymHash returns symbol hash of this object, which is used for prop reference.
func (s *PanStr) SymHash() SymHash {
	return GetSymHash(s.Value)
}

// NewPanStr makes new str object.
func NewPanStr(s string) *PanStr {
	return &PanStr{Value: s, IsPublic: isPublic(s)}
}

var publicPattern = regexp.MustCompile(`^[a-zA-Z][a-zA-Z0-9_]*[!?]?$`)

func isPublic(s string) bool {
	return publicPattern.MatchString(s)
}
