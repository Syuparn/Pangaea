package object

import (
	"fmt"
	"regexp"
	"strings"
)

// StrType is a type of PanStr.
const StrType = "StrType"

// used as zero value
var zeroStr = NewPanStr("")

// PanStr is object of str literal.
type PanStr struct {
	Value    string
	IsPublic bool
	IsSym    bool
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
	return &PanStr{Value: s, IsPublic: isPublic(s), IsSym: isSym(s)}
}

var publicPattern = regexp.MustCompile(`^[a-zA-Z][a-zA-Z0-9_]*[!?]?$`)
var privatePattern = regexp.MustCompile(`^_[a-zA-Z][a-zA-Z0-9_]*[!?]?$`)
var argVarPattern = regexp.MustCompile(`^\\[0-9]*$`)
var kwargVarPattern = regexp.MustCompile(`^\\[_a-zA-Z][a-zA-Z0-9_]*[!?]?$`)
var opPattern = regexp.MustCompile(fmt.Sprintf(`^(%s)$`, strings.Join([]string{
	`<=>`, `==`, `!=`, `>=`, `<=`, `>`, `<`, `<<`, `>>`,
	`/&`, `/\|`, `/\^`, `/~`, `!`, `\+`, `\-`, `\*`, `\*\*`,
	`/`, "//", `%`, `\-%`, `\+%`,
}, "|")))

func isPublic(s string) bool {
	return publicPattern.MatchString(s)
}

func isSym(s string) bool {
	switch {
	case publicPattern.MatchString(s):
		return true
	case privatePattern.MatchString(s):
		return true
	case argVarPattern.MatchString(s):
		return true
	case kwargVarPattern.MatchString(s):
		return true
	case opPattern.MatchString(s):
		return true
	default:
		return false
	}
}
