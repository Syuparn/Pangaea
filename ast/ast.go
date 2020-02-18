// Source code in this file is inherited and modified from
// "Writing an Interpreter in Go" https://interpreterbook.com/
// MIT License | Copyright (c) 2016-2017 Thorsten Ball
// see https://opensource.org/licenses/MIT

package ast

import (
	"bytes"
	"fmt"
	"strings"
)

type Node interface {
	TokenLiteral() string
	String() string
	Source() *Source
}

type Expr interface {
	Node
	isExpr() // dammy method not to be taken for Stmt
}

type Stmt interface {
	Node
	isStmt() // dammy method not to be taken for Expr
}

type Program struct {
	Stmts []Stmt
	Src   *Source
}

func (p *Program) TokenLiteral() string { return "" }
func (p *Program) String() string {
	var out bytes.Buffer

	for _, s := range p.Stmts {
		out.WriteString(s.String())
	}

	return out.String()
}
func (p *Program) Source() *Source { return p.Stmts[0].Source() }

type ExprStmt struct {
	Token string
	Expr  Expr
	Src   *Source
}

func (es *ExprStmt) isStmt()              {}
func (es *ExprStmt) TokenLiteral() string { return es.Token }
func (es *ExprStmt) Source() *Source      { return es.Src }

func (es *ExprStmt) String() string {
	if es.Expr != nil {
		return es.Expr.String()
	}
	return ""
}

type Ident struct {
	Token     string
	Value     string
	Src       *Source
	IsPrivate bool
}

func (i *Ident) isExpr()              {}
func (i *Ident) TokenLiteral() string { return i.Token }
func (i *Ident) String() string       { return i.Value }
func (i *Ident) Source() *Source      { return i.Src }

type CallExpr interface {
	ChainToken() string
	ChainArg() Expr
	Expr
}

type PropCallExpr struct {
	Token    string
	Chain    *Chain
	Receiver Expr
	Prop     *Ident
	Args     []Expr
	Kwargs   map[Ident]Expr
	Src      *Source
}

func (pc *PropCallExpr) isExpr()              {}
func (pc *PropCallExpr) TokenLiteral() string { return pc.Token }
func (pc *PropCallExpr) ChainToken() string   { return pc.Chain.Token }
func (pc *PropCallExpr) ChainArg() Expr       { return pc.Chain.Arg }
func (pc *PropCallExpr) Source() *Source      { return pc.Src }
func (pc *PropCallExpr) String() string {
	var out bytes.Buffer
	out.WriteString(pc.Receiver.String())
	out.WriteString(pc.Chain.String())
	out.WriteString(pc.Prop.String())

	args := []string{}
	for _, a := range pc.Args {
		args = append(args, a.String())
	}

	out.WriteString("(" + strings.Join(args, ", ") + ")")

	return out.String()
}

type MainChain int

const (
	Scalar MainChain = iota
	List
	Reduce
)

type AdditionalChain int

const (
	Vanilla AdditionalChain = iota
	Lonely
	Thoughtful
	Strict
)

func MakeChain(addChain string, mainChain string, chainArg Expr) *Chain {
	var addChainMap = map[string]AdditionalChain{
		"":  Vanilla,
		"&": Lonely,
		"~": Thoughtful,
		"=": Strict,
	}

	var mainChainMap = map[string]MainChain{
		".": Scalar,
		"@": List,
		"$": Reduce,
	}

	return &Chain{
		Token:      addChain + mainChain,
		Additional: addChainMap[addChain],
		Main:       mainChainMap[mainChain],
		Arg:        chainArg,
	}
}

type Chain struct {
	Token      string
	Additional AdditionalChain
	Main       MainChain
	Arg        Expr
}

func (c *Chain) String() string {
	var out bytes.Buffer
	out.WriteString(c.Token)
	if c.Arg != nil {
		out.WriteString("(" + c.Arg.String() + ")")
	}
	return out.String()
}

type InfixExpr struct {
	Token    string // i.e.: "+"
	Left     Expr
	Operator string
	Right    Expr
	Src      *Source
}

func (ie *InfixExpr) isExpr()              {}
func (ie *InfixExpr) TokenLiteral() string { return ie.Token }
func (ie *InfixExpr) Source() *Source      { return ie.Src }
func (ie *InfixExpr) String() string {
	var out bytes.Buffer
	out.WriteString("(")
	out.WriteString(ie.Left.String())
	out.WriteString(" " + ie.Operator + " ")
	out.WriteString(ie.Right.String())
	out.WriteString(")")

	return out.String()
}

type IntLiteral struct {
	Token string
	Value int64
	Src   *Source
}

func (il *IntLiteral) isExpr()              {}
func (il *IntLiteral) TokenLiteral() string { return il.Token }
func (il *IntLiteral) Source() *Source      { return il.Src }
func (il *IntLiteral) String() string       { return il.Token }

// for error message
type Source struct {
	Line string
	Pos  Position
}

type Position struct {
	Line   int
	Column int
}

func (p *Position) String() string {
	return fmt.Sprintf("line: %d, col: %d", p.Line, p.Column)
}
