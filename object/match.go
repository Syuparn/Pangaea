package object

import ()

const MATCH_TYPE = "MATCH_TYPE"

type PanMatch struct {
	MatchWrapper
}

func (m *PanMatch) Type() PanObjType {
	return ""
}

func (m *PanMatch) Inspect() string {
	// delegate to MatchWrapper
	return m.MatchWrapper.String()
}

func (m *PanMatch) Proto() PanObject {
	return m
}

// NOTE: keep loose coupling to ast.MatchLiteral and PanMatch
type MatchWrapper interface {
	String() string
}
