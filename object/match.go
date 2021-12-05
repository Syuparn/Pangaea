package object

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

// Repr returns pritty-printed string of this object.
func (m *PanMatch) Repr() string {
	return m.Inspect()
}

// Proto returns proto of this object.
func (m *PanMatch) Proto() PanObject {
	return BuiltInMatchObj
}

// Zero returns zero value of this object.
func (m *PanMatch) Zero() PanObject {
	// TODO: implement zero value
	return m
}

// MatchWrapper is a wrapper for match literal ast node.
// NOTE: keep loose coupling to ast.MatchLiteral and PanMatch
// ast.MatchLiteral implements MatchWrapper
type MatchWrapper interface {
	String() string
}
