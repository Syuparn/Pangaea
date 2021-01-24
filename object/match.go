package object

import ()

// MatchType is a type of PanMatch.
const MatchType = "MatchType"

// PanMatch is object of match literal.
type PanMatch struct {
	MatchWrapper
}

// Type returns type of this PanObject.
func (m *PanMatch) Type() PanObjType {
	return MatchType
}

// Inspect returns formatted source code of this object.
func (m *PanMatch) Inspect() string {
	// delegate to MatchWrapper
	return m.MatchWrapper.String()
}

// Proto returns proto of this object.
func (m *PanMatch) Proto() PanObject {
	return BuiltInMatchObj
}

// MatchWrapper is a wrapper for match literal ast node.
// NOTE: keep loose coupling to ast.MatchLiteral and PanMatch
// ast.MatchLiteral implements MatchWrapper
type MatchWrapper interface {
	String() string
}
