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

// Node is a base interface of ast node.
type Node interface {
	TokenLiteral() string
	String() string
	Source() *Source
}

// Expr is an interface of ast expression node.
type Expr interface {
	Node
	isExpr() // dummy method not to be taken for Stmt
}

// Stmt is an interface of ast statement node.
type Stmt interface {
	Node
	isStmt() // dummy method not to be taken for Expr
}

// Program is an ast node of program, corresponding to the whole source code.
type Program struct {
	Stmts []Stmt
	Src   *Source
}

// TokenLiteral returns token given by lexer.
func (p *Program) TokenLiteral() string { return "" }
func (p *Program) String() string {
	stmts := []string{}
	for _, s := range p.Stmts {
		stmts = append(stmts, s.String())
	}
	return strings.Join(stmts, "\n")
}

// Source returns stacktrace infomation used for error massages.
func (p *Program) Source() *Source { return p.Stmts[0].Source() }

// ExprStmt is an ast node of statement which consists of just one expression.
type ExprStmt struct {
	Token string
	Expr  Expr
	Src   *Source
}

func (es *ExprStmt) isStmt() {}

// TokenLiteral returns token given by lexer.
func (es *ExprStmt) TokenLiteral() string { return es.Token }

// Source returns stacktrace infomation used for error massages.
func (es *ExprStmt) Source() *Source { return es.Src }

func (es *ExprStmt) String() string {
	if es.Expr != nil {
		return es.Expr.String()
	}
	return ""
}

// JumpType is an Enum for jump keywords, used for JumpStmt.
type JumpType int

const (
	// ReturnJump is JumpType for `return`.
	ReturnJump JumpType = iota
	// RaiseJump is JumpType for `raise`.
	RaiseJump
	// YieldJump is JumpType for `yield`.
	YieldJump
	// DeferJump is JumpType for `defer`.
	// NOTE: treat `defer` as JumpType because syntax is equivalent to the other JumpTypes
	DeferJump
)

func jumpString(j JumpType) string {
	return map[JumpType]string{
		ReturnJump: "return",
		RaiseJump:  "raise",
		YieldJump:  "yield",
		DeferJump:  "defer",
	}[j]
}

// JumpStmt is an ast node of statement such as `return` or `yield`.
type JumpStmt struct {
	Token    string
	Val      Expr
	JumpType JumpType
	Src      *Source
}

func (js *JumpStmt) isStmt() {}

// TokenLiteral returns token given by lexer.
func (js *JumpStmt) TokenLiteral() string { return js.Token }

// Source returns stacktrace infomation used for error massages.
func (js *JumpStmt) Source() *Source { return js.Src }
func (js *JumpStmt) String() string {
	var out bytes.Buffer

	out.WriteString(jumpString(js.JumpType) + " ")

	if js.Val != nil {
		out.WriteString(js.Val.String())
	}

	return out.String()
}

// JumpIfStmt is an ast node of JumpStmt with `if` clause.
type JumpIfStmt struct {
	Token    string
	JumpStmt *JumpStmt
	Cond     Expr
	Src      *Source
}

func (js *JumpIfStmt) isStmt() {}

// TokenLiteral returns token given by lexer.
func (js *JumpIfStmt) TokenLiteral() string { return js.Token }

// Source returns stacktrace infomation used for error massages.
func (js *JumpIfStmt) Source() *Source { return js.Src }
func (js *JumpIfStmt) String() string {
	var out bytes.Buffer
	out.WriteString(js.JumpStmt.String())
	out.WriteString(" if ")
	out.WriteString(js.Cond.String())
	return out.String()
}

// IdentAttr is enum for identifier kinds.
type IdentAttr int

const (
	// NormalIdent is a Type for ordinally ident like `a`.
	NormalIdent IdentAttr = iota
	// ArgIdent is a Type for arg ident like `\1`.
	ArgIdent
	// KwargIdent is a Type for kwarg ident like `\a`.
	KwargIdent
)

// Ident is an ast node of identifier expression.
type Ident struct {
	Token     string
	Value     string
	Src       *Source
	IsPrivate bool
	IdentAttr IdentAttr
}

func (i *Ident) isExpr() {}

// TokenLiteral returns token given by lexer.
func (i *Ident) TokenLiteral() string { return i.Token }
func (i *Ident) String() string       { return i.Value }

// Source returns stacktrace infomation used for error massages.
func (i *Ident) Source() *Source { return i.Src }

// PinnedIdent is an ast node of identifier with pinned operation `^`.
type PinnedIdent struct {
	Ident
}

func (pi *PinnedIdent) String() string {
	return "^" + pi.Ident.String()
}

// CallExpr is an ast node interface for chain expression.
type CallExpr interface {
	ChainToken() string
	ChainArg() Expr
	Expr
}

// PropCallExpr is an ast node of propcall expression like `"hello".p`.
type PropCallExpr struct {
	Token    string
	Chain    *Chain
	Receiver Expr
	Prop     *Ident
	Args     []Expr
	Kwargs   map[*Ident]Expr
	Src      *Source
}

func (pc *PropCallExpr) isExpr() {}

// TokenLiteral returns token given by lexer.
func (pc *PropCallExpr) TokenLiteral() string { return pc.Token }

// ChainToken returns chain string of this call.
func (pc *PropCallExpr) ChainToken() string { return pc.Chain.Token }

// ChainArg returns chain arg expression of this call.
func (pc *PropCallExpr) ChainArg() Expr { return pc.Chain.Arg }

// Source returns stacktrace infomation used for error massages.
func (pc *PropCallExpr) Source() *Source { return pc.Src }
func (pc *PropCallExpr) String() string {
	var out bytes.Buffer

	if pc.Receiver != nil {
		out.WriteString(pc.Receiver.String())
	}
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

// LiteralCallExpr is an ast node of literalcall expression like `1.{|i| i * 2}`.
type LiteralCallExpr struct {
	Token    string
	Chain    *Chain
	Receiver Expr
	Func     *FuncLiteral
	Args     []Expr
	Kwargs   map[*Ident]Expr
	Src      *Source
}

func (lc *LiteralCallExpr) isExpr() {}

// TokenLiteral returns token given by lexer.
func (lc *LiteralCallExpr) TokenLiteral() string { return lc.Token }

// ChainToken returns chain string of this call.
func (lc *LiteralCallExpr) ChainToken() string { return lc.Chain.Token }

// ChainArg returns chain arg expression of this call.
func (lc *LiteralCallExpr) ChainArg() Expr { return lc.Chain.Arg }

// Source returns stacktrace infomation used for error massages.
func (lc *LiteralCallExpr) Source() *Source { return lc.Src }
func (lc *LiteralCallExpr) String() string {
	var out bytes.Buffer
	if lc.Receiver != nil {
		out.WriteString(lc.Receiver.String())
	}
	out.WriteString(lc.Chain.String())
	out.WriteString(lc.Func.String())

	return out.String()
}

// VarCallExpr is an ast node of varcall expression like `foo.^bar`.
type VarCallExpr struct {
	Token    string
	Chain    *Chain
	Receiver Expr
	Var      *Ident
	Args     []Expr
	Kwargs   map[*Ident]Expr
	Src      *Source
}

func (vc *VarCallExpr) isExpr() {}

// TokenLiteral returns token given by lexer.
func (vc *VarCallExpr) TokenLiteral() string { return vc.Token }

// ChainToken returns chain string of this call.
func (vc *VarCallExpr) ChainToken() string { return vc.Chain.Token }

// ChainArg returns chain arg expression of this call.
func (vc *VarCallExpr) ChainArg() Expr { return vc.Chain.Arg }

// Source returns stacktrace infomation used for error massages.
func (vc *VarCallExpr) Source() *Source { return vc.Src }
func (vc *VarCallExpr) String() string {
	var out bytes.Buffer
	if vc.Receiver != nil {
		out.WriteString(vc.Receiver.String())
	}
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

// RecvAndChain is a container for receiver and chain of CallExpr.
// This type is only used to make parser rules simple.
type RecvAndChain struct {
	Recv  Expr
	Chain *Chain
}

// MainChain is an Enum for main chain context.
type MainChain int

const (
	// Scalar is a Type for `.` chain.
	Scalar MainChain = iota
	// List is a Type for `@` chain.
	List
	// Reduce is a Type for `$` chain.
	Reduce
)

// AdditionalChain is an Enum for additional chain context.
type AdditionalChain int

const (
	// Vanilla is a Type showing no addtional chains.
	Vanilla AdditionalChain = iota
	// Lonely is a Type for `&` chain.
	Lonely
	// Thoughtful is a Type for `~` chain.
	Thoughtful
	// Strict is a Type for `=` chain.
	Strict
)

// MakeChain makes new Chain from main chain literal and addtional chain literal.
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

// Chain is an ast element for chain like `.`.
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

// KwargPair is an ast element for keyword argument and its value.
type KwargPair struct {
	Key *Ident
	Val Expr
}

func (qp *KwargPair) String() string {
	return qp.Key.String() + ": " + qp.Val.String()
}

// Pair is an ast element for key-value pair.
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

// SelfIdentArgList returns ArgList only with identifier `self`.
func SelfIdentArgList(src *Source) *ArgList {
	return &ArgList{
		Args:   []Expr{selfIdent(src)},
		Kwargs: map[*Ident]Expr{},
	}
}

// ExprToArgList makes ArgList with the passed expression.
func ExprToArgList(e Expr) *ArgList {
	return &ArgList{
		Args:   []Expr{e},
		Kwargs: map[*Ident]Expr{},
	}
}

// KwargPairToArgList makes ArgList with the passed keyword argument.
func KwargPairToArgList(pair *KwargPair) *ArgList {
	return &ArgList{
		Args:   []Expr{},
		Kwargs: map[*Ident]Expr{pair.Key: pair.Val},
	}
}

// ArgList is an ast element for argument list.
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

// AppendArg appends a positional argument, then returns self.
func (al *ArgList) AppendArg(arg Expr) *ArgList {
	al.Args = append(al.Args, arg)
	return al
}

// AppendKwarg appends a keyword argument and its value, then returns self.
func (al *ArgList) AppendKwarg(key *Ident, arg Expr) *ArgList {
	al.Kwargs[key] = arg
	return al
}

// PrependSelf prepends `self` ident to its arguments, then returns self.
func (al *ArgList) PrependSelf(src *Source) *ArgList {
	al.Args = append([]Expr{selfIdent(src)}, al.Args...)
	return al
}

func selfIdent(src *Source) *Ident {
	return &Ident{
		Token:     "self",
		Value:     "self",
		Src:       src,
		IsPrivate: false,
	}
}

// PrefixExpr is an ast node of prefix expression.
type PrefixExpr struct {
	Token    string
	Operator string
	Right    Expr
	Src      *Source
}

func (pe *PrefixExpr) isExpr() {}

// TokenLiteral returns token given by lexer.
func (pe *PrefixExpr) TokenLiteral() string { return pe.Token }

// Source returns stacktrace infomation used for error massages.
func (pe *PrefixExpr) Source() *Source { return pe.Src }
func (pe *PrefixExpr) String() string {
	var out bytes.Buffer
	out.WriteString("(")
	out.WriteString(pe.Operator + pe.Right.String())
	out.WriteString(")")

	return out.String()
}

// InfixExpr is an ast node of infix expression.
type InfixExpr struct {
	Token    string // i.e.: "+"
	Left     Expr
	Operator string
	Right    Expr
	Src      *Source
}

func (ie *InfixExpr) isExpr() {}

// TokenLiteral returns token given by lexer.
func (ie *InfixExpr) TokenLiteral() string { return ie.Token }

// Source returns stacktrace infomation used for error massages.
func (ie *InfixExpr) Source() *Source { return ie.Src }
func (ie *InfixExpr) String() string {
	var out bytes.Buffer
	out.WriteString("(")
	out.WriteString(ie.Left.String())
	out.WriteString(" " + ie.Operator + " ")
	out.WriteString(ie.Right.String())
	out.WriteString(")")

	return out.String()
}

// AssignExpr is an ast node of assign expression like `a := 2`.
type AssignExpr struct {
	Token string // ":="
	Left  *Ident
	Right Expr
	Src   *Source
}

func (ae *AssignExpr) isExpr() {}

// TokenLiteral returns token given by lexer.
func (ae *AssignExpr) TokenLiteral() string { return ae.Token }

// Source returns stacktrace infomation used for error massages.
func (ae *AssignExpr) Source() *Source { return ae.Src }
func (ae *AssignExpr) String() string {
	var out bytes.Buffer
	out.WriteString("(")
	out.WriteString(ae.Left.String())
	out.WriteString(" " + ae.Token + " ")
	out.WriteString(ae.Right.String())
	out.WriteString(")")

	return out.String()
}

// EmbeddedStr is an ast node of embedded string like `"value=#{1 + 1}"`.
type EmbeddedStr struct {
	Token  string
	Former *FormerStrPiece
	Latter string
	Src    *Source
}

func (es *EmbeddedStr) isExpr() {}

// TokenLiteral returns token given by lexer.
func (es *EmbeddedStr) TokenLiteral() string { return es.Token }

// Source returns stacktrace infomation used for error massages.
func (es *EmbeddedStr) Source() *Source { return es.Src }
func (es *EmbeddedStr) String() string {
	return `"` + es.Former.String() + es.Latter + `"`
}

// FormerStrPiece is a former part of EmbeddedStr expression.
// NOTE: EmbeddedStr consists of recursive FormerStrPieces and
// each FormerStrPiece corresponds to string literal and expression.
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

// StrLiteral is an ast node of str literal like `"a"`.
type StrLiteral struct {
	Token string
	Value string
	IsRaw bool
	Src   *Source
}

func (sl *StrLiteral) isExpr() {}

// TokenLiteral returns token given by lexer.
func (sl *StrLiteral) TokenLiteral() string { return sl.Token }

// Source returns stacktrace infomation used for error massages.
func (sl *StrLiteral) Source() *Source { return sl.Src }
func (sl *StrLiteral) String() string {
	if sl.IsRaw {
		return "`" + sl.Value + "`"
	}
	return `"` + sl.Value + `"`
}

// SymLiteral is an ast node of sym literal like `'a`.
type SymLiteral struct {
	Token string
	Value string
	Src   *Source
}

func (sl *SymLiteral) isExpr() {}

// TokenLiteral returns token given by lexer.
func (sl *SymLiteral) TokenLiteral() string { return sl.Token }

// Source returns stacktrace infomation used for error massages.
func (sl *SymLiteral) Source() *Source { return sl.Src }
func (sl *SymLiteral) String() string  { return "'" + sl.Value }

// RangeLiteral is an ast node of range literal like `(1:10)`.
type RangeLiteral struct {
	Token string
	Start Expr
	Stop  Expr
	Step  Expr
	Src   *Source
}

func (rl *RangeLiteral) isExpr() {}

// TokenLiteral returns token given by lexer.
func (rl *RangeLiteral) TokenLiteral() string { return rl.Token }

// Source returns stacktrace infomation used for error massages.
func (rl *RangeLiteral) Source() *Source { return rl.Src }
func (rl *RangeLiteral) String() string {
	var out bytes.Buffer
	out.WriteString("(")

	if rl.Start != nil {
		out.WriteString(rl.Start.String())
	}
	out.WriteString(":")

	if rl.Stop != nil {
		out.WriteString(rl.Stop.String())
	}
	out.WriteString(":")

	if rl.Step != nil {
		out.WriteString(rl.Step.String())
	}

	out.WriteString(")")

	return out.String()
}

// IfExpr is an ast node of if expression.
type IfExpr struct {
	Token string
	Cond  Expr
	Then  Expr
	Else  Expr
	Src   *Source
}

func (ie *IfExpr) isExpr() {}

// TokenLiteral returns token given by lexer.
func (ie *IfExpr) TokenLiteral() string { return ie.Token }

// Source returns stacktrace infomation used for error massages.
func (ie *IfExpr) Source() *Source { return ie.Src }
func (ie *IfExpr) String() string {
	var out bytes.Buffer
	out.WriteString("(")
	out.WriteString(ie.Then.String() + " if " + ie.Cond.String())
	if ie.Else != nil {
		out.WriteString(" else " + ie.Else.String())
	}
	out.WriteString(")")

	return out.String()
}

// FuncLiteral is an ast node of func literal like `{|i| i * 2}`.
type FuncLiteral struct {
	FuncComponent
	Token string
	Src   *Source
}

func (fl *FuncLiteral) isExpr() {}

// TokenLiteral returns token given by lexer.
func (fl *FuncLiteral) TokenLiteral() string { return fl.Token }

// Source returns stacktrace infomation used for error massages.
func (fl *FuncLiteral) Source() *Source { return fl.Src }
func (fl *FuncLiteral) String() string {
	var out bytes.Buffer
	out.WriteString("{")
	out.WriteString(fl.FuncComponent.String())
	out.WriteString("}")

	return out.String()
}

// IterLiteral is an ast node of func literal like `<{|i| yield i if i; recur(i-1)}>`.
type IterLiteral struct {
	FuncComponent
	Token string
	Src   *Source
}

func (il *IterLiteral) isExpr() {}

// TokenLiteral returns token given by lexer.
func (il *IterLiteral) TokenLiteral() string { return il.Token }

// Source returns stacktrace infomation used for error massages.
func (il *IterLiteral) Source() *Source { return il.Src }
func (il *IterLiteral) String() string {
	var out bytes.Buffer
	out.WriteString("<{")
	out.WriteString(il.FuncComponent.String())
	out.WriteString("}>")

	return out.String()
}

// MatchLiteral is an ast node of match literal like `%{|1| 2, |3| 4}`.
type MatchLiteral struct {
	Token    string
	Patterns []*FuncComponent
	Src      *Source
}

func (ml *MatchLiteral) isExpr() {}

// TokenLiteral returns token given by lexer.
func (ml *MatchLiteral) TokenLiteral() string { return ml.Token }

// Source returns stacktrace infomation used for error massages.
func (ml *MatchLiteral) Source() *Source { return ml.Src }
func (ml *MatchLiteral) String() string {
	var out bytes.Buffer
	out.WriteString("%{\n")

	patterns := []string{}
	for _, pat := range ml.Patterns {
		patterns = append(patterns, pat.String())
	}
	out.WriteString(strings.Join(patterns, ",\n"))

	out.WriteString("}")
	return out.String()
}

// FuncComponent is an ast container for elements of func/iter/match literal.
type FuncComponent struct {
	Args   []Expr
	Kwargs map[*Ident]Expr
	Body   []Stmt
	Src    *Source
}

// PrependSelf prepends `self` parameter to its paramters.
// This is helper function to parse methods like `m{}`
// (, which is a syntax sugar of `{|self| }`)
func (fc *FuncComponent) PrependSelf(src *Source) *FuncComponent {
	fc.Args = append([]Expr{selfIdent(src)}, fc.Args...)
	return fc
}

func (fc *FuncComponent) String() string {
	var out bytes.Buffer

	args := []string{}
	for _, a := range fc.Args {
		args = append(args, a.String())
	}
	args = append(args, sortedPairStrings(fc.Kwargs)...)

	out.WriteString("|" + strings.Join(args, ", ") + "|")

	bodies := []string{}
	for _, stmt := range fc.Body {
		bodies = append(bodies, stmt.String())
	}

	switch len(bodies) {
	case 0:
		// prepend space to break args and body(empty)
		out.WriteString(" ")
	case 1:
		// prepend space to break args and body
		out.WriteString(" " + bodies[0])
	default:
		out.WriteString("\n" + strings.Join(bodies, "\n") + "\n")
	}
	return out.String()
}

// DiamondLiteral is an ast node for diamond literal `<>`.
type DiamondLiteral struct {
	Token string
	Src   *Source
}

func (dl *DiamondLiteral) isExpr() {}

// TokenLiteral returns token given by lexer.
func (dl *DiamondLiteral) TokenLiteral() string { return dl.Token }

// Source returns stacktrace infomation used for error massages.
func (dl *DiamondLiteral) Source() *Source { return dl.Src }
func (dl *DiamondLiteral) String() string  { return "<>" }

// ObjLiteral is an ast node for obj literal like `{a: 1}`.
type ObjLiteral struct {
	Token         string
	Pairs         []*Pair
	EmbeddedExprs []Expr
	Src           *Source
}

func (ol *ObjLiteral) isExpr() {}

// TokenLiteral returns token given by lexer.
func (ol *ObjLiteral) TokenLiteral() string { return ol.Token }

// Source returns stacktrace infomation used for error massages.
func (ol *ObjLiteral) Source() *Source { return ol.Src }
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

// MapLiteral is an ast node for map literal like `%{"a": 1}`.
type MapLiteral struct {
	Token         string
	Pairs         []*Pair
	EmbeddedExprs []Expr
	Src           *Source
}

func (ml *MapLiteral) isExpr() {}

// TokenLiteral returns token given by lexer.
func (ml *MapLiteral) TokenLiteral() string { return ml.Token }

// Source returns stacktrace infomation used for error massages.
func (ml *MapLiteral) Source() *Source { return ml.Src }
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

// ArrLiteral is an ast node for arr literal like `[1, 2]`.
type ArrLiteral struct {
	Token string
	Elems []Expr
	Src   *Source
}

func (al *ArrLiteral) isExpr() {}

// TokenLiteral returns token given by lexer.
func (al *ArrLiteral) TokenLiteral() string { return al.Token }

// Source returns stacktrace infomation used for error massages.
func (al *ArrLiteral) Source() *Source { return al.Src }
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

// IntLiteral is an ast node for int literal like `1`.
type IntLiteral struct {
	Token string
	Value int64
	Src   *Source
}

func (il *IntLiteral) isExpr() {}

// TokenLiteral returns token given by lexer.
func (il *IntLiteral) TokenLiteral() string { return il.Token }

// Source returns stacktrace infomation used for error massages.
func (il *IntLiteral) Source() *Source { return il.Src }
func (il *IntLiteral) String() string  { return il.Token }

// FloatLiteral is an ast node for float literal like `1.0`.
type FloatLiteral struct {
	Token string
	Value float64
	Src   *Source
}

func (fl *FloatLiteral) isExpr() {}

// TokenLiteral returns token given by lexer.
func (fl *FloatLiteral) TokenLiteral() string { return fl.Token }

// Source returns stacktrace infomation used for error massages.
func (fl *FloatLiteral) Source() *Source { return fl.Src }
func (fl *FloatLiteral) String() string  { return fl.Token }

// Source is a stacktrace infomation used for error massages.
type Source struct {
	Line         string
	Pos          Position
	TokenLiteral string
}

// Position is a location of source code.
type Position struct {
	Line   int
	Column int
}

func (p *Position) String() string {
	// NOTE: add 1 otherwise first element is shown as 0
	return fmt.Sprintf("line: %d, col: %d", p.Line+1, p.Column+1)
}
