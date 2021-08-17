package simplexer

import (
	"io"
	"strings"
)

// Defined default values for properties of Lexer as a package value.
var (
	DefaultWhitespace = NewPatternTokenType(-1, []string{" ", "\t", "\r", "\n"})

	DefaultTokenTypes = []TokenType{
		NewRegexpTokenType(IDENT, `[a-zA-Z_][a-zA-Z0-9_]*`),
		NewRegexpTokenType(NUMBER, `[0-9]+(?:\.[0-9]+)?`),
		NewRegexpTokenType(STRING, `\"([^"]*)\"`),
		NewRegexpTokenType(OTHER, `.`),
	}
)

/*
The lexical analyzer.

Whitespace is a TokenType for skipping characters like whitespaces.
The default value is simplexer.DefaultWhitespace.
Won't skip any characters if Whitespace is nil.

TokenTypes is an array of TokenType.
Lexer will sequential check TokenTypes, and return first matched token.
Default is simplexer.DefaultTokenTypes.

Please be careful, Lexer will never use it even if append TokenType after OTHER.
Because OTHER will accept any single character.
*/
type Lexer struct {
	reader     io.Reader
	buf        string
	loadedLine string
	nextPos    Position
	Whitespace TokenType
	TokenTypes []TokenType
}

// Make a new Lexer.
func NewLexer(reader io.Reader) *Lexer {
	l := new(Lexer)
	l.reader = reader

	l.Whitespace = DefaultWhitespace
	l.TokenTypes = DefaultTokenTypes

	return l
}

func (l *Lexer) readBufIfNeed() {
	if len(l.buf) < 1024 {
		buf := make([]byte, 2048)
		l.reader.Read(buf)
		l.buf += strings.TrimRight(string(buf), "\x00")
	}
}

func (l *Lexer) consumeBuffer(t *Token) {
	if t == nil {
		return
	}

	l.buf = l.buf[len(t.Literal):]

	l.nextPos = shiftPos(l.nextPos, t.Literal)

	if idx := strings.LastIndex(t.Literal, "\n"); idx >= 0 {
		l.loadedLine = t.Literal[idx+1:]
	} else {
		l.loadedLine += t.Literal
	}
}

func (l *Lexer) skipWhitespace() {
	if l.Whitespace == nil {
		return
	}

	for true {
		l.readBufIfNeed()

		if t := l.Whitespace.FindToken(l.buf, l.nextPos); t != nil {
			l.consumeBuffer(t)
		} else {
			break
		}
	}
}

func (l *Lexer) makeError() error {
	for shift, _ := range l.buf {
		if l.Whitespace != nil && l.Whitespace.FindToken(l.buf[shift:], l.nextPos) != nil {
			return UnknownTokenError{
				Literal:  l.buf[:shift],
				Position: l.nextPos,
			}
		}

		for _, tokenType := range l.TokenTypes {
			if tokenType.FindToken(l.buf[shift:], l.nextPos) != nil {
				return UnknownTokenError{
					Literal:  l.buf[:shift],
					Position: l.nextPos,
				}
			}
		}
	}

	return UnknownTokenError{
		Literal:  l.buf,
		Position: l.nextPos,
	}
}

/*
Peek the first token in the buffer.

Returns nil as *Token if the buffer is empty.
*/
func (l *Lexer) Peek() (*Token, error) {
	for _, tokenType := range l.TokenTypes {
		l.skipWhitespace()

		l.readBufIfNeed()
		if t := tokenType.FindToken(l.buf, l.nextPos); t != nil {
			return t, nil
		}
	}

	if len(l.buf) > 0 {
		return nil, l.makeError()
	}

	return nil, nil
}

/*
Scan will get the first token in the buffer and remove it from the buffer.

This function using Lexer.Peek. Please read document of Peek.
*/
func (l *Lexer) Scan() (*Token, error) {
	t, e := l.Peek()

	l.consumeBuffer(t)

	return t, e
}

/*
GetCurrentLine returns line of last scanned token.
*/
func (l *Lexer) GetLastLine() string {
	l.readBufIfNeed()

	if idx := strings.Index(l.buf, "\n"); idx >= 0 {
		return l.loadedLine + l.buf[:strings.Index(l.buf, "\n")]
	} else {
		return l.loadedLine + l.buf
	}
}
