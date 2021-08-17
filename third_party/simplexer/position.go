package simplexer

import (
	"fmt"
	"strings"
)

// Position in the file.
type Position struct {
	Line   int
	Column int
}

// Convert to string.
func (p Position) String() string {
	return fmt.Sprintf("[line:%d, column:%d]", p.Line, p.Column)
}

// Position.Before will check p is before than x.
func (p Position) Before(x Position) bool {
	return p.Line < x.Line || (p.Line == x.Line && p.Column < x.Column)
}

// Position.After will check p is after than x.
func (p Position) After(x Position) bool {
	return p.Line > x.Line || (p.Line == x.Line && p.Column > x.Column)
}

func shiftPos(p Position, s string) Position {
	lines := strings.Split(s, "\n")
	lineShift := len(lines) - 1

	if lineShift == 0 {
		p.Column += len(lines[0])
	} else {
		p.Column = len(lines[len(lines)-1])
	}
	p.Line += lineShift

	return p
}
