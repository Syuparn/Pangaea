package object

import (
	"regexp"
)

const STR_TYPE = "STR_TYPE"

type PanStr struct {
	Value    string
	IsPublic bool
}

func (s *PanStr) Type() PanObjType {
	return STR_TYPE
}

func (s *PanStr) Inspect() string {
	return `"` + s.Value + `"`
}

func (s *PanStr) Proto() PanObject {
	return BuiltInStrObj
}

func (s *PanStr) Hash() HashKey {
	return HashKey{STR_TYPE, s.SymHash()}
}

func (s *PanStr) SymHash() SymHash {
	return GetSymHash(s.Value)
}

func NewPanStr(s string) *PanStr {
	return &PanStr{Value: s, IsPublic: isPublic(s)}
}

var publicPattern = regexp.MustCompile(`^[a-zA-Z][a-zA-Z0-9_]*[!?]?$`)

func isPublic(s string) bool {
	return publicPattern.MatchString(s)
}
