%{

package parser

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"strconv"

	"github.com/macrat/simplexer"

	"../ast"
)

%}

%union{
    token *simplexer.Token
    chain *ast.Chain
	ident *ast.Ident
	expr  ast.Expr
	argList *ast.ArgList
	paramList *ast.ParamList
	exprList []ast.Expr
	pair *ast.Pair
	stmt  ast.Stmt
	stmts []ast.Stmt
	program *ast.Program
}

%type<program> program
%type<stmts> stmts
%type<stmt> stmt exprStmt
%type<expr> expr literal infixExpr prefixExpr callExpr funcLiteral arrLiteral
%type<pair> pair
%type<argList> argList callArgs
%type<paramList> paramList funcParams
%type<exprList> exprList
%type<ident> ident
%type<chain> chain
%type<token> opMethod breakLine
%type<token> lBrace lParen lBracket comma

%token<token> INT
%token<token> DOUBLE_STAR PLUS MINUS STAR SLASH BANG
%token<token> ADD_CHAIN MAIN_CHAIN
%token<token> IDENT PRIVATE_IDENT
%token<token> LPAREN RPAREN COMMA COLON LBRACE RBRACE VERT LBRACKET RBRACKET
%token<token> RET SEMICOLON
%left PLUS MINUS
%left STAR SLASH
%left DOUBLE_STAR
%left ADD_CHAIN MAIN_CHAIN
%left UNARY_OP

%% 

program
	: stmts
	{
		$$ = &ast.Program{Stmts: $1}
		yylex.(*Lexer).program = $$
		yylex.(*Lexer).curRule = "program -> stmts"
	}
	| RET stmts
	{
		$$ = &ast.Program{Stmts: $2}
		yylex.(*Lexer).program = $$
		yylex.(*Lexer).curRule = "program -> RET stmts"
	}

stmts
	: stmt
	{
		$$ = []ast.Stmt{$1}
		yylex.(*Lexer).curRule = "stmts -> stmt"
	}
	| stmts breakLine stmt
	{
		$$ = append($1, $3)
		yylex.(*Lexer).curRule = "stmts -> stmts breakLine stmt"
	}
	| stmts breakLine
	{
		$$ = $1
		yylex.(*Lexer).curRule = "stmts -> stmts breakLine"
	}

stmt
	: exprStmt
	{
		$$ = $1
		yylex.(*Lexer).curRule = "stmt -> exprStmt"
	}

exprStmt
	: expr
	{
		$$ = &ast.ExprStmt{
			Token: "(exprStmt)",
			Expr: $1,
		}
		yylex.(*Lexer).curRule = "exprStmt -> expr"
	}

expr
	: literal
	{
		$$ = $1
		yylex.(*Lexer).curRule = "expr -> literal"
	}
	| infixExpr
	{
		$$ = $1
		yylex.(*Lexer).curRule = "expr -> infixExpr"
	}
	| prefixExpr
	{
		$$ = $1
		yylex.(*Lexer).curRule = "expr -> prefixExpr"
	}
	| callExpr
	{
		$$ = $1
		yylex.(*Lexer).curRule = "expr -> callExpr"
	}
	| ident
	{
		$$ = $1
		yylex.(*Lexer).curRule = "expr -> ident"
	}

ident
	: IDENT
	{
		$$ = &ast.Ident{
			Token: $1.Literal,
			Value: $1.Literal,
			Src: yylex.(*Lexer).Source,
			IsPrivate: false,
		}
		yylex.(*Lexer).curRule = "ident -> IDENT"
	}
	| PRIVATE_IDENT
	{
		$$ = &ast.Ident{
			Token: $1.Literal,
			Value: $1.Literal,
			Src: yylex.(*Lexer).Source,
			IsPrivate: true,
		}
		yylex.(*Lexer).curRule = "ident -> PRIVATE_IDENT"
	}

literal
	: INT
	{
		n, _ := strconv.ParseInt($1.Literal, 10, 64)
		$$ = &ast.IntLiteral{
			Token: $1.Literal,
			Value: n,
			Src: yylex.(*Lexer).Source,
		}
		yylex.(*Lexer).curRule = "literal -> INT"
	}
	| funcLiteral
	{
		$$ = $1
		yylex.(*Lexer).curRule = "literal -> funcLiteral"
	}
	| arrLiteral
	{
		$$ = $1
		yylex.(*Lexer).curRule = "literal -> arrLiteral"
	}

infixExpr
	: expr STAR expr
	{
		$$ = &ast.InfixExpr{
			Token: $2.Literal,
			Left: $1,
			Operator: $2.Literal,
			Right: $3,
			Src: yylex.(*Lexer).Source,
		}
		yylex.(*Lexer).curRule = "infixExpr -> expr STAR expr"
	}
	| expr SLASH expr
	{
		$$ = &ast.InfixExpr{
			Token: $2.Literal,
			Left: $1,
			Operator: $2.Literal,
			Right: $3,
			Src: yylex.(*Lexer).Source,
		}
		yylex.(*Lexer).curRule = "infixExpr -> expr SLASH expr"
	}
	| expr PLUS expr
	{
		$$ = &ast.InfixExpr{
			Token: $2.Literal,
			Left: $1,
			Operator: $2.Literal,
			Right: $3,
			Src: yylex.(*Lexer).Source,
		}
		yylex.(*Lexer).curRule = "infixExpr -> expr PLUS expr"
	}
	| expr MINUS expr
	{
		$$ = &ast.InfixExpr{
			Token: $2.Literal,
			Left: $1,
			Operator: $2.Literal,
			Right: $3,
			Src: yylex.(*Lexer).Source,
		}
		yylex.(*Lexer).curRule = "infixExpr -> expr MINUS expr"
	}

prefixExpr
	: PLUS expr %prec UNARY_OP
	{
		$$ = &ast.PrefixExpr{
			Token: $1.Literal,
			Operator: $1.Literal,
			Right: $2,
			Src: yylex.(*Lexer).Source,
		}
		yylex.(*Lexer).curRule = "prefixExpr -> PLUS expr"
	}
	| MINUS expr %prec UNARY_OP
	{
		$$ = &ast.PrefixExpr{
			Token: $1.Literal,
			Operator: $1.Literal,
			Right: $2,
			Src: yylex.(*Lexer).Source,
		}
		yylex.(*Lexer).curRule = "prefixExpr -> MINUS expr"
	}
	| STAR expr %prec UNARY_OP
	{
		$$ = &ast.PrefixExpr{
			Token: $1.Literal,
			Operator: $1.Literal,
			Right: $2,
			Src: yylex.(*Lexer).Source,
		}
		yylex.(*Lexer).curRule = "prefixExpr -> STAR expr"
	}
	| DOUBLE_STAR expr %prec UNARY_OP
	{
		$$ = &ast.PrefixExpr{
			Token: $1.Literal,
			Operator: $1.Literal,
			Right: $2,
			Src: yylex.(*Lexer).Source,
		}
		yylex.(*Lexer).curRule = "prefixExpr -> DOUBLE_STAR expr"
	}
	| BANG expr %prec UNARY_OP
	{
		$$ = &ast.PrefixExpr{
			Token: $1.Literal,
			Operator: $1.Literal,
			Right: $2,
			Src: yylex.(*Lexer).Source,
		}
		yylex.(*Lexer).curRule = "prefixExpr -> BANG expr"
	}

arrLiteral
	: lBracket RBRACKET
	{
		$$ = &ast.ArrLiteral{
			Token: $1.Literal,
			Elems: []ast.Expr{},
			Src: yylex.(*Lexer).Source,

		}
		yylex.(*Lexer).curRule = "arrLiteral -> lBracket RBRACKET"
	}
	| lBracket exprList RBRACKET
	{
		$$ = &ast.ArrLiteral{
			Token: $1.Literal,
			Elems: $2,
			Src: yylex.(*Lexer).Source,

		}
		yylex.(*Lexer).curRule = "arrLiteral -> lBracket exprList RBRACKET"
	}

funcLiteral
	: lBrace funcParams RBRACE
	{
		$$ = &ast.FuncLiteral{
			Token: $1.Literal,
			Args: $2.Args,
			Kwargs: $2.Kwargs,
			Body: []ast.Stmt{},
			Src: yylex.(*Lexer).Source,
		}
		yylex.(*Lexer).curRule = "funcLiteral -> lBrace funcParams RBRACE"
	}
	| lBrace stmts RBRACE
	{
		$$ = &ast.FuncLiteral{
			Token: $1.Literal,
			Args: []*ast.Ident{},
			Kwargs: map[*ast.Ident]ast.Expr{},
			Body: $2,
			Src: yylex.(*Lexer).Source,
		}
		yylex.(*Lexer).curRule = "funcLiteral -> lBrace stmts RBRACE"
	}
	| lBrace funcParams stmts RBRACE
	{
		$$ = &ast.FuncLiteral{
			Token: $1.Literal,
			Args: $2.Args,
			Kwargs: $2.Kwargs,
			Body: $3,
			Src: yylex.(*Lexer).Source,
		}
		yylex.(*Lexer).curRule = "funcLiteral -> lBrace funcParams stmts RBRACE"
	}

funcParams
	: VERT VERT
	{
		$$ = &ast.ParamList{
			Args: []*ast.Ident{},
			Kwargs: map[*ast.Ident]ast.Expr{},
		}
		yylex.(*Lexer).curRule = "funcParams -> VERT VERT"
	}
	| VERT paramList VERT
	{
		$$ = $2
		yylex.(*Lexer).curRule = "funcParams -> VERT paramList VERT"
	}
	| VERT paramList RET VERT
	{
		$$ = $2
		yylex.(*Lexer).curRule = "funcParams -> VERT paramList RET VERT"
	}
	|  VERT VERT RET
	{
		$$ = &ast.ParamList{
			Args: []*ast.Ident{},
			Kwargs: map[*ast.Ident]ast.Expr{},
		}
		yylex.(*Lexer).curRule = "funcParams -> VERT VERT RET"
	}
	| VERT paramList VERT RET
	{
		$$ = $2
		yylex.(*Lexer).curRule = "funcParams -> VERT paramList VERT RET"
	}
	| VERT paramList RET VERT RET
	{
		$$ = $2
		yylex.(*Lexer).curRule = "funcParams -> VERT paramList RET VERT RET"
	}

callExpr
	: expr chain ident
	{
		$$ = &ast.PropCallExpr{
			Token: "(propCall)",
			Chain: $2,
			Receiver: $1,
			Prop: $3,
			Args: nil,
			Kwargs: nil,
			Src: yylex.(*Lexer).Source,
		}
		yylex.(*Lexer).curRule = "callExpr -> expr chain ident"
	}
	| expr chain ident callArgs
	{
		$$ = &ast.PropCallExpr{
			Token: "(propCall)",
			Chain: $2,
			Receiver: $1,
			Prop: $3,
			Args: $4.Args,
			Kwargs: $4.Kwargs,
			Src: yylex.(*Lexer).Source,
		}
		yylex.(*Lexer).curRule = "callExpr -> expr chain ident callArgs"
	}
	| expr chain opMethod
	{
		opIdent := &ast.Ident{
			Token: $3.Literal,
			Value: $3.Literal,
			Src: yylex.(*Lexer).Source,
			IsPrivate: true,
		}
		$$ = &ast.PropCallExpr{
			Token: "(propCall)",
			Chain: $2,
			Receiver: $1,
			Prop: opIdent,
			Args: nil,
			Kwargs: nil,
			Src: yylex.(*Lexer).Source,
		}
		yylex.(*Lexer).curRule = "callExpr -> expr chain opMethod"
	}
	| expr chain opMethod callArgs
	{
		opIdent := &ast.Ident{
			Token: $3.Literal,
			Value: $3.Literal,
			Src: yylex.(*Lexer).Source,
			IsPrivate: true,
		}
		$$ = &ast.PropCallExpr{
			Token: "(propCall)",
			Chain: $2,
			Receiver: $1,
			Prop: opIdent,
			Args: $4.Args,
			Kwargs: $4.Kwargs,
			Src: yylex.(*Lexer).Source,
		}
		yylex.(*Lexer).curRule = "callExpr -> expr chain opMethod callArgs"
	}

callArgs
	: lParen RPAREN
	{
		$$ = &ast.ArgList{
			Args: []ast.Expr{},
			Kwargs: map[*ast.Ident]ast.Expr{},
		}
		yylex.(*Lexer).curRule = "callArgs -> lParen RPAREN"
	}
	| lParen argList RPAREN
	{
		$$ = $2
		yylex.(*Lexer).curRule = "callArgs -> lParen argList RPAREN"
	}
	| lParen argList RET RPAREN
	{
		$$ = $2
		yylex.(*Lexer).curRule = "callArgs -> lParen argList RET RPAREN"
	}

opMethod
	: PLUS
	{
		$$ = $1
		yylex.(*Lexer).curRule = "opMethod -> PLUS"
	}
	| MINUS
	{
		$$ = $1
		yylex.(*Lexer).curRule = "opMethod -> MINUS"
	}
	| STAR
	{
		$$ = $1
		yylex.(*Lexer).curRule = "opMethod -> STAR"
	}
	| SLASH
	{
		$$ = $1
		yylex.(*Lexer).curRule = "opMethod -> SLASH"
	}

chain
	: ADD_CHAIN MAIN_CHAIN
	{
		$$ = ast.MakeChain($1.Literal, $2.Literal, nil)
		yylex.(*Lexer).curRule = "chain -> ADD_CHAIN MAIN_CHAIN"
	}
	| MAIN_CHAIN
	{
		$$ = ast.MakeChain("", $1.Literal, nil)
		yylex.(*Lexer).curRule = "chain -> MAIN_CHAIN"
	}
	| MAIN_CHAIN lParen expr RPAREN
	{
		$$ = ast.MakeChain("", $1.Literal, $3)
		yylex.(*Lexer).curRule = "chain -> MAIN_CHAIN lParen expr RPAREN"
	}
	| ADD_CHAIN MAIN_CHAIN lParen expr RPAREN
	{
		$$ = ast.MakeChain($1.Literal, $2.Literal, $4)
		yylex.(*Lexer).curRule = "chain -> ADD_CHAIN MAIN_CHAIN lParen expr RPAREN"
	}

exprList
	: exprList comma expr
	{
		$$ = append($1, $3)
		yylex.(*Lexer).curRule = "exprList -> exprList comma expr"
	}
	| expr
	{
		$$ = []ast.Expr{$1}
		yylex.(*Lexer).curRule = "exprList -> expr"
	}

argList
	: argList comma expr
	{
		$$ = $1.AppendArg($3)
		yylex.(*Lexer).curRule = "argList -> argList comma expr"
	}
	| argList comma pair
	{
		$$ = $1.AppendKwarg($3.Key, $3.Val)
		yylex.(*Lexer).curRule = "argList -> argList comma pair"
	}
	| expr
	{
		$$ = ast.ExprToArgList($1)
		yylex.(*Lexer).curRule = "argList -> expr"
	}
	| pair
	{
		$$ = ast.PairToArgList($1)
		yylex.(*Lexer).curRule = "argList -> pair"
	}

paramList
	: paramList comma ident
	{
		$$ = $1.AppendArg($3)
		yylex.(*Lexer).curRule = "paramList -> paramList comma ident"
	}
	| paramList comma pair
	{
		$$ = $1.AppendKwarg($3.Key, $3.Val)
		yylex.(*Lexer).curRule = "paramList -> paramList comma pair"
	}
	| ident
	{
		$$ = ast.IdentToParamList($1)
		yylex.(*Lexer).curRule = "paramList -> ident"
	}
	| pair
	{
		$$ = ast.PairToParamList($1)
		yylex.(*Lexer).curRule = "paramList -> pair"
	}

pair
	: ident COLON expr
	{
		$$ = &ast.Pair{Key: $1, Val: $3}
		yylex.(*Lexer).curRule = "pair -> ident COLON expr"
	}

lBrace
	: LBRACE
	{
		$$ = $1
		yylex.(*Lexer).curRule = "lBrace -> LBRACE RET"
	}
	| LBRACE RET
	{
		$$ = $1
		yylex.(*Lexer).curRule = "lBrace -> LBRACE RET"
	}

lParen
	: LPAREN
	{
		$$ = $1
		yylex.(*Lexer).curRule = "lParen -> LPAREN"
	}
	| LPAREN RET
	{
		$$ = $1
		yylex.(*Lexer).curRule = "lParen -> LPAREN RET"
	}

lBracket
	: LBRACKET
	{
		$$ = $1
		yylex.(*Lexer).curRule = "lBracket -> LBRACKET"
	}
	| LBRACKET RET
	{
		$$ = $1
		yylex.(*Lexer).curRule = "lBracket -> LBRACKET RET"
	}

breakLine
	: SEMICOLON
	{
		$$ = $1
		yylex.(*Lexer).curRule = "breakLine -> SEMICOLON"
	}
	| RET
	{
		$$ = $1
		yylex.(*Lexer).curRule = "breakLine -> RET"
	}

comma
	: COMMA
	{
		$$ = $1
		yylex.(*Lexer).curRule = "comma -> COMMA"
	}
	| COMMA RET
	{
		$$ = $1
		yylex.(*Lexer).curRule = "comma -> COMMA RET"
	}

%%

func Parse(src io.Reader) (*ast.Program, error) {	
	lexer := NewLexer(src)
	prog, err := tryParse(src, lexer)
	
	if err != nil {
		return nil, err
	}

	if prog == nil {
		return nil, errors.New("failed to parse")
	}

	program, ok := prog.(*ast.Program)
	if !ok {
		msg := fmt.Sprintf("could not parsed to *ast.Program. got=%+v", program)
		return nil, errors.New(msg)
	}

	return program, nil
}

func tryParse(src io.Reader, l *Lexer) (a ast.Node, e error) {
	// HACK: catch yyParse error by recover
	// (because yacc cannot return error)
	defer func(l *Lexer) {
		if err := recover(); err != nil {
			m := "error occured:"
			if l.Source != nil {
				m = fmt.Sprintf("%s\n%s", m,
					l.ErrMsg()) 
			} else {
				m = m + "\nbefore lexing"
			}
			e = errors.New(m)
		}
	}(l)
	
	yyParse(l)
	return l.program, nil
}

type Lexer struct {
	lexer        *simplexer.Lexer
	// NOTE: embed final ast in lexer because yyParse cannot return ast
	program      ast.Node
	Source		 *ast.Source
	curRule		 string
}

var tokenTypes = []simplexer.TokenType{
	simplexer.NewRegexpTokenType(INT, `[0-9]+(\.[0-9]+)?`),
	simplexer.NewRegexpTokenType(RET, `(\r|\n|\r\n)+`),
	simplexer.NewRegexpTokenType(LPAREN, `\(`),
	simplexer.NewRegexpTokenType(RPAREN, `\)`),
	simplexer.NewRegexpTokenType(VERT, `\|`),
	simplexer.NewRegexpTokenType(LBRACE, `\{`),
	simplexer.NewRegexpTokenType(RBRACE, `\}`),
	simplexer.NewRegexpTokenType(LBRACKET, `\[`),
	simplexer.NewRegexpTokenType(RBRACKET, `\]`),
	simplexer.NewRegexpTokenType(COMMA, `,`),
	simplexer.NewRegexpTokenType(COLON, `:`),
	simplexer.NewRegexpTokenType(SEMICOLON, `;`),
	simplexer.NewRegexpTokenType(BANG, `!`),
	simplexer.NewRegexpTokenType(DOUBLE_STAR, `\*\*`),
	simplexer.NewRegexpTokenType(PLUS, `\+`),
	simplexer.NewRegexpTokenType(MINUS, `\-`),
	simplexer.NewRegexpTokenType(STAR, `\*`),
	simplexer.NewRegexpTokenType(SLASH, `/`),
	simplexer.NewRegexpTokenType(ADD_CHAIN, `[&~=]`),
	simplexer.NewRegexpTokenType(MAIN_CHAIN, `[\.@$]`),
	simplexer.NewRegexpTokenType(IDENT, `[a-zA-Z][a-zA-Z0-9_]*([!?])?`),
	simplexer.NewRegexpTokenType(PRIVATE_IDENT, `_[a-zA-Z][a-zA-Z0-9_]*([!?])?`),
}

func NewLexer(reader io.Reader) *Lexer {
	l := simplexer.NewLexer(reader)
	l.TokenTypes = tokenTypes
	// NOTE: remove "\n" from whitespace list
	// to use it stmts separator
	l.Whitespace = simplexer.NewPatternTokenType(
		-1, []string{" ", "\t"})
	return &Lexer{ lexer: l }
}

func (l *Lexer) Lex(lval *yySymType) int {
	token, err := l.lexer.Scan()

	if _, ok := err.(*simplexer.UnknownTokenError); ok {
		l.Error(l.unknownTokenErrMsg())
	} else if err != nil {
		l.Error(l.ErrMsg())
	}

	if token == nil {
		return -1
	}

	lval.token = token
	l.Source = l.convertSourceInfo(token)
	return int(token.Type.GetID())
}

func (l *Lexer) unknownTokenErrMsg() string {
	var out bytes.Buffer
	tok := l.Source.TokenLiteral
	if tok == "\n" {
		tok = "\\n" // for readability
	}
	
	out.WriteString(fmt.Sprintf("Lexer Error: unknown token '%s'was found\n",
		tok))
	out.WriteString("after " + l.Source.Pos.String() + "\n")
	out.WriteString(l.Source.Line + "\n")
	out.WriteString(fmt.Sprintf("in rule: %s\n", l.curRule))
	return out.String()
}

func (l *Lexer) ErrMsg() string {
	var out bytes.Buffer
	tok := l.Source.TokenLiteral
	if tok == "\n" {
		tok = "\\n" // for readability
	}

	out.WriteString(fmt.Sprintf("Lexer Error in token '%s':\n", tok))
	out.WriteString("after " + l.Source.Pos.String() + "\n")
	out.WriteString(l.Source.Line + "\n")
	out.WriteString(fmt.Sprintf("in rule: %s\n", l.curRule))
	return out.String()
}

func (l *Lexer) Error(e string) {
	// NOTE: yacc(yyParse) cannot return err object...
	panic(e)
}

func (l *Lexer) convertSourceInfo(token *simplexer.Token) *ast.Source {
	pos := ast.Position{
		Line: token.Position.Line,
		Column: token.Position.Column,
	}
	return &ast.Source{
		Line: l.lexer.GetLastLine(),
		Pos: pos,
		TokenLiteral: token.Literal,
	}
}
