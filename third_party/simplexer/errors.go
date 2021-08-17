package simplexer

import "fmt"

// The error that returns when found an unknown token.
type UnknownTokenError struct {
	Literal  string
	Position Position
}

// Get error message as string.
func (se UnknownTokenError) Error() string {
	return fmt.Sprintf("%d:%d:UnknownTokenError: %#v", se.Position.Line+1, se.Position.Column+1, se.Literal)
}
