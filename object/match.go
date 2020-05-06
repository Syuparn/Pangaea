package object

import ()

const MATCH_TYPE = "MATCH_TYPE"

type PanMatch struct {
	MatchWrapper
}

func (m *PanMatch) Type() PanObjType {
	return MATCH_TYPE
}

func (m *PanMatch) Inspect() string {
	// delegate to MatchWrapper
	return m.MatchWrapper.String()
}

func (m *PanMatch) Proto() PanObject {
	return BuiltInMatchObj
}

// NOTE: keep loose coupling to ast.MatchLiteral and PanMatch
// ast.MatchLiteral implements MatchWrapper
type MatchWrapper interface {
	String() string
}
