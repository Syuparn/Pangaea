package simplexer

import (
	"regexp"
	"strconv"
	"strings"
)

// TokenID is Identifier for TokenType.
type TokenID int

// Default token IDs.
const (
	OTHER TokenID = -(iota + 1)
	IDENT
	NUMBER
	STRING
)

/*
Convert to readable string.

Be careful, user added token ID's will convert to UNKNOWN.
*/
func (id TokenID) String() string {
	switch id {
	case OTHER:
		return "OTHER"
	case IDENT:
		return "IDENT"
	case NUMBER:
		return "NUMBER"
	case STRING:
		return "STRING"
	default:
		return "UNKNOWN(" + strconv.Itoa(int(id)) + ")"
	}
}

/*
TokenType is a rule for making Token.

GetID returns TokenID of this TokenType.
TokenID can share with another TokenType.

FindToken returns new Token if the head of first argument was matched with the pattern of this TokenType.
The second argument is a position of the token in the buffer. In almost implement, Position will pass into result Token directly.
*/
type TokenType interface {
	GetID() TokenID
	FindToken(string, Position) *Token
}

/*
RegexpTokenType is a TokenType implement with regexp.

ID is TokenID for this token type.

Re is regular expression of token. It have to starts with "^".
*/
type RegexpTokenType struct {
	ID TokenID
	Re *regexp.Regexp
}

/*
Make new RegexpTokenType.

id is a TokenID of new RegexpTokenType.

re is a regular expression of token.
*/
func NewRegexpTokenType(id TokenID, re string) *RegexpTokenType {
	if !strings.HasPrefix(re, "^") {
		re = "^(?:" + re + ")"
	}
	return &RegexpTokenType{
		ID: id,
		Re: regexp.MustCompile(re),
	}
}

// Get readable string of TokenID.
func (rtt *RegexpTokenType) String() string {
	return rtt.ID.String()
}

// GetID returns id of this token type.
func (rtt *RegexpTokenType) GetID() TokenID {
	return rtt.ID
}

// FindToken returns new Token if s starts with this token.
func (rtt *RegexpTokenType) FindToken(s string, p Position) *Token {
	m := rtt.Re.FindStringSubmatch(s)
	if len(m) > 0 {
		return &Token{
			Type:       rtt,
			Literal:    m[0],
			Submatches: m[1:],
			Position:   p,
		}
	}
	return nil
}

/*
PatternTokenType is dictionary token type.

PatternTokenType has some strings and find token that perfect match they.
*/
type PatternTokenType struct {
	ID       TokenID
	Patterns []string
}

/*
Make new PatternTokenType.

id is a TokenID of new PatternTokenType.

patterns is array of patterns.
*/
func NewPatternTokenType(id TokenID, patterns []string) *PatternTokenType {
	return &PatternTokenType{
		ID:       id,
		Patterns: patterns,
	}
}

// Get readable string of TokenID.
func (ptt *PatternTokenType) String() string {
	return ptt.ID.String()
}

// GetID returns id of token type.
func (ptt *PatternTokenType) GetID() TokenID {
	return ptt.ID
}

// FindToken returns new Token if s starts with this token.
func (ptt *PatternTokenType) FindToken(s string, p Position) *Token {
	for _, x := range ptt.Patterns {
		if strings.HasPrefix(s, x) {
			return &Token{
				Type:     ptt,
				Literal:  x,
				Position: p,
			}
		}
	}
	return nil
}

// A data of found Token.
type Token struct {
	Type       TokenType
	Literal    string   // The string of matched.
	Submatches []string // Submatches of regular expression.
	Position   Position // Position of token.
}
