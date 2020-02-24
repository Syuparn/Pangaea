// Source code in this file is inherited and modified from
// "Writing an Interpreter in Go" https://interpreterbook.com/
// MIT License | Copyright (c) 2016-2017 Thorsten Ball
// see https://opensource.org/licenses/MIT

package ast

import (
	"bytes"
	"fmt"
	"sort"
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
	stmts := []string{}
	for _, s := range p.Stmts {
		stmts = append(stmts, s.String())
	}
	return strings.Join(stmts, "\n")
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

type IdentAttr int

const (
	NormalIdent IdentAttr = iota
	ArgIdent
	KwargIdent
)

type Ident struct {
	Token     string
	Value     string
	Src       *Source
	IsPrivate bool
	IdentAttr IdentAttr
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
	Kwargs   map[*Ident]Expr
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

	args = append(args, sortedPairStrings(pc.Kwargs)...)

	out.WriteString("(" + strings.Join(args, ", ") + ")")

	return out.String()
}

type LiteralCallExpr struct {
	Token    string
	Chain    *Chain
	Receiver Expr
	Func     *FuncLiteral
	Args     []Expr
	Kwargs   map[*Ident]Expr
	Src      *Source
}

func (lc *LiteralCallExpr) isExpr()              {}
func (lc *LiteralCallExpr) TokenLiteral() string { return lc.Token }
func (lc *LiteralCallExpr) ChainToken() string   { return lc.Chain.Token }
func (lc *LiteralCallExpr) ChainArg() Expr       { return lc.Chain.Arg }
func (lc *LiteralCallExpr) Source() *Source      { return lc.Src }
func (lc *LiteralCallExpr) String() string {
	var out bytes.Buffer
	out.WriteString(lc.Receiver.String())
	out.WriteString(lc.Chain.String())
	out.WriteString(lc.Func.String())

	args := []string{}
	for _, a := range lc.Args {
		args = append(args, a.String())
	}

	args = append(args, sortedPairStrings(lc.Kwargs)...)

	out.WriteString("(" + strings.Join(args, ", ") + ")")

	return out.String()
}

type VarCallExpr struct {
	Token    string
	Chain    *Chain
	Receiver Expr
	Var      *Ident
	Args     []Expr
	Kwargs   map[*Ident]Expr
	Src      *Source
}

func (vc *VarCallExpr) isExpr()              {}
func (vc *VarCallExpr) TokenLiteral() string { return vc.Token }
func (vc *VarCallExpr) ChainToken() string   { return vc.Chain.Token }
func (vc *VarCallExpr) ChainArg() Expr       { return vc.Chain.Arg }
func (vc *VarCallExpr) Source() *Source      { return vc.Src }
func (vc *VarCallExpr) String() string {
	var out bytes.Buffer
	out.WriteString(vc.Receiver.String())
	out.WriteString(vc.Chain.String())
	out.WriteString("^" + vc.Var.String())

	args := []string{}
	for _, a := range vc.Args {
		args = append(args, a.String())
	}

	args = append(args, sortedPairStrings(vc.Kwargs)...)

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

type KwargPair struct {
	Key *Ident
	Val Expr
}

func (qp *KwargPair) String() string {
	return qp.Key.String() + ": " + qp.Val.String()
}

type Pair struct {
	Key Expr
	Val Expr
}

func (p *Pair) String() string {
	return p.Key.String() + ": " + p.Val.String()
}

func sortedPairStrings(pairs map[*Ident]Expr) []string {
	// NOTE: sort kwargs by ident name (otherwise order is random!)
	type p struct {
		k string
		v string
	}

	kwargs := []p{}
	for k, arg := range pairs {
		kwargs = append(kwargs, p{k: k.String(), v: arg.String()})
	}
	sort.Slice(kwargs, func(i, j int) bool { return kwargs[i].k < kwargs[j].k })

	sortedStrings := []string{}
	for _, kwarg := range kwargs {
		sortedStrings = append(sortedStrings, kwarg.k+": "+kwarg.v)
	}
	return sortedStrings
}

func SelfIdentParamList(src *Source) *ParamList {
	selfIdent := &Ident{
		Token:     "self",
		Value:     "self",
		Src:       src,
		IsPrivate: false,
	}

	return &ParamList{
		Args:   []*Ident{selfIdent},
		Kwargs: map[*Ident]Expr{},
	}
}

func IdentToParamList(i *Ident) *ParamList {
	return &ParamList{
		Args:   []*Ident{i},
		Kwargs: map[*Ident]Expr{},
	}
}

func KwargPairToParamList(pair *KwargPair) *ParamList {
	return &ParamList{
		Args:   []*Ident{},
		Kwargs: map[*Ident]Expr{pair.Key: pair.Val},
	}
}

type ParamList struct {
	Args   []*Ident
	Kwargs map[*Ident]Expr
}

func (pl *ParamList) PrependSelf(src *Source) *ParamList {
	selfIdent := &Ident{
		Token:     "self",
		Value:     "self",
		Src:       src,
		IsPrivate: false,
	}
	pl.Args = append([]*Ident{selfIdent}, pl.Args...)
	return pl
}

func (pl *ParamList) AppendArg(arg *Ident) *ParamList {
	pl.Args = append(pl.Args, arg)
	return pl
}

func (pl *ParamList) AppendKwarg(key *Ident, arg Expr) *ParamList {
	pl.Kwargs[key] = arg
	return pl
}

func ExprToArgList(e Expr) *ArgList {
	return &ArgList{
		Args:   []Expr{e},
		Kwargs: map[*Ident]Expr{},
	}
}

func KwargPairToArgList(pair *KwargPair) *ArgList {
	return &ArgList{
		Args:   []Expr{},
		Kwargs: map[*Ident]Expr{pair.Key: pair.Val},
	}
}

type ArgList struct {
	Args   []Expr
	Kwargs map[*Ident]Expr
}

func (al *ArgList) String() string {
	args := []string{}
	for _, arg := range al.Args {
		args = append(args, arg.String())
	}

	// NOTE: sort kwargs by ident name (otherwise order is random!)
	args = append(args, sortedPairStrings(al.Kwargs)...)

	return strings.Join(args, ", ")
}

func (al *ArgList) AppendArg(arg Expr) *ArgList {
	al.Args = append(al.Args, arg)
	return al
}

func (al *ArgList) AppendKwarg(key *Ident, arg Expr) *ArgList {
	al.Kwargs[key] = arg
	return al
}

type PrefixExpr struct {
	Token    string
	Operator string
	Right    Expr
	Src      *Source
}

func (pe *PrefixExpr) isExpr()              {}
func (pe *PrefixExpr) TokenLiteral() string { return pe.Token }
func (pe *PrefixExpr) Source() *Source      { return pe.Src }
func (pe *PrefixExpr) String() string {
	var out bytes.Buffer
	out.WriteString("(")
	out.WriteString(pe.Operator + pe.Right.String())
	out.WriteString(")")

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

type EmbeddedStr struct {
	Token  string
	Former *FormerStrPiece
	Latter string
	Src    *Source
}

func (es *EmbeddedStr) isExpr()              {}
func (es *EmbeddedStr) TokenLiteral() string { return es.Token }
func (es *EmbeddedStr) Source() *Source      { return es.Src }
func (es *EmbeddedStr) String() string {
	return `"` + es.Former.String() + es.Latter + `"`
}

type FormerStrPiece struct {
	Token  string
	Former *FormerStrPiece
	Str    string
	Expr   Expr
}

func (fs *FormerStrPiece) String() string {
	var out bytes.Buffer

	if fs.Former != nil {
		out.WriteString(fs.Former.String())
	}

	out.WriteString(fs.Str)
	out.WriteString(fmt.Sprintf("#{ %s }", fs.Expr.String()))
	return out.String()
}

type StrLiteral struct {
	Token string
	Value string
	IsRaw bool
	Src   *Source
}

func (sl *StrLiteral) isExpr()              {}
func (sl *StrLiteral) TokenLiteral() string { return sl.Token }
func (sl *StrLiteral) Source() *Source      { return sl.Src }
func (sl *StrLiteral) String() string {
	if sl.IsRaw {
		return "`" + sl.Value + "`"
	} else {
		return `"` + sl.Value + `"`
	}
}

type SymLiteral struct {
	Token string
	Value string
	Src   *Source
}

func (sl *SymLiteral) isExpr()              {}
func (sl *SymLiteral) TokenLiteral() string { return sl.Token }
func (sl *SymLiteral) Source() *Source      { return sl.Src }
func (sl *SymLiteral) String() string       { return "'" + sl.Value }

type FuncLiteral struct {
	Token  string
	Args   []*Ident
	Kwargs map[*Ident]Expr
	Body   []Stmt
	Src    *Source
}

func (fl *FuncLiteral) isExpr()              {}
func (fl *FuncLiteral) TokenLiteral() string { return fl.Token }
func (fl *FuncLiteral) Source() *Source      { return fl.Src }
func (fl *FuncLiteral) String() string {
	var out bytes.Buffer
	out.WriteString("{")

	args := []string{}
	for _, a := range fl.Args {
		args = append(args, a.String())
	}
	args = append(args, sortedPairStrings(fl.Kwargs)...)

	out.WriteString("|" + strings.Join(args, ", ") + "| ")

	bodies := []string{}
	for _, stmt := range fl.Body {
		bodies = append(bodies, stmt.String())
	}

	switch len(bodies) {
	case 0:
		// nothing
	case 1:
		out.WriteString(bodies[0])
	default:
		out.WriteString("\n" + strings.Join(bodies, "\n") + "\n")
	}

	out.WriteString("}")

	return out.String()
}

type ObjLiteral struct {
	Token         string
	Pairs         []*Pair
	EmbeddedExprs []Expr
	Src           *Source
}

func (ol *ObjLiteral) isExpr()              {}
func (ol *ObjLiteral) TokenLiteral() string { return ol.Token }
func (ol *ObjLiteral) Source() *Source      { return ol.Src }
func (ol *ObjLiteral) String() string {
	var out bytes.Buffer
	out.WriteString("{")

	elems := []string{}
	for _, pair := range ol.Pairs {
		elems = append(elems, pair.String())
	}

	for _, expr := range ol.EmbeddedExprs {
		elems = append(elems, "**"+expr.String())
	}

	out.WriteString(strings.Join(elems, ", "))

	out.WriteString("}")
	return out.String()
}

type MapLiteral struct {
	Token         string
	Pairs         []*Pair
	EmbeddedExprs []Expr
	Src           *Source
}

func (ml *MapLiteral) isExpr()              {}
func (ml *MapLiteral) TokenLiteral() string { return ml.Token }
func (ml *MapLiteral) Source() *Source      { return ml.Src }
func (ml *MapLiteral) String() string {
	var out bytes.Buffer
	out.WriteString("%{")

	elems := []string{}
	for _, pair := range ml.Pairs {
		elems = append(elems, pair.String())
	}

	for _, expr := range ml.EmbeddedExprs {
		elems = append(elems, "**"+expr.String())
	}

	out.WriteString(strings.Join(elems, ", "))

	out.WriteString("}")
	return out.String()
}

type ArrLiteral struct {
	Token string
	Elems []Expr
	Src   *Source
}

func (al *ArrLiteral) isExpr()              {}
func (al *ArrLiteral) TokenLiteral() string { return al.Token }
func (al *ArrLiteral) Source() *Source      { return al.Src }
func (al *ArrLiteral) String() string {
	var out bytes.Buffer
	out.WriteString("[")

	elems := []string{}

	for _, elem := range al.Elems {
		elems = append(elems, elem.String())
	}

	out.WriteString(strings.Join(elems, ", "))
	out.WriteString("]")
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
	Line         string
	Pos          Position
	TokenLiteral string
}

type Position struct {
	Line   int
	Column int
}

func (p *Position) String() string {
	return fmt.Sprintf("line: %d, col: %d", p.Line, p.Column)
}
