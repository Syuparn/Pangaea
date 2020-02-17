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
    expr  ast.Expr
	stmt  ast.Stmt
	stmts []ast.Stmt
	program *ast.Program
}

%type<program> program
%type<stmts> stmts
%type<stmt> stmt exprStmt
%type<expr> expr literal infixExpr
%token<token> INT PREC1_OPERATOR PREC2_OPERATOR
%left PREC1_OPERATOR
%left PREC2_OPERATOR

%% 

program
	: stmts
	{
		$$ = &ast.Program{Stmts: $1}
		yylex.(*Lexer).program = $$
	}

stmts
	: stmt
	{
		$$ = []ast.Stmt{$1}
	}
	| stmts stmt
	{
		$$ = append($1, $2)
	}

stmt
	: exprStmt
	{
		$$ = $1
	}

exprStmt
	: expr
	{
		$$ = &ast.ExprStmt{
			Token: "(exprStmt)",
			Expr: $1,
		}
	}

expr
	: literal
	{
		$$ = $1
	}
	| infixExpr
	{
		$$ = $1
	}

literal
	: INT
	{
		n, _ := strconv.ParseInt($1.Literal, 10, 64)
		$$ = &ast.IntLiteral{
			Token: $1.Literal,
			Value: n,
			Src: yylex.(*Lexer).source,
		}
	}

infixExpr
	: expr PREC2_OPERATOR expr
	{
		$$ = &ast.InfixExpr{
			Token: $2.Literal,
			Left: $1,
			Operator: $2.Literal,
			Right: $3,
			Src: yylex.(*Lexer).source,
		}
	}
	| expr PREC1_OPERATOR expr
	{
		$$ = &ast.InfixExpr{
			Token: $2.Literal,
			Left: $1,
			Operator: $2.Literal,
			Right: $3,
			Src: yylex.(*Lexer).source,
		}
	}

%%

func Parse(src io.Reader) (*ast.Program, error) {
	lexer := NewLexer(src)
	yyParse(lexer)

	if lexer.program == nil {
		return nil, errors.New("failed to parse")
	}

	program, ok := lexer.program.(*ast.Program)
	if !ok {
		msg := fmt.Sprintf("could not parsed to *ast.Program. got=%+v", program)
		return nil, errors.New(msg)
	}

	return program, nil
}

type Lexer struct {
	lexer        *simplexer.Lexer
	// NOTE: embed final ast in lexer because yyParse cannot return ast
	program      ast.Node
	source		 *ast.Source
}

var tokenTypes = []simplexer.TokenType{
	simplexer.NewRegexpTokenType(INT, `[0-9]+(\.[0-9]+)?`),
	simplexer.NewRegexpTokenType(PREC1_OPERATOR, `[-+]`),
	simplexer.NewRegexpTokenType(PREC2_OPERATOR, `[*/]`),
}

func NewLexer(reader io.Reader) *Lexer {
	l := simplexer.NewLexer(reader)
	l.TokenTypes = tokenTypes
	return &Lexer{ lexer: l }
}

func (l *Lexer) Lex(lval *yySymType) int {
	token, err := l.lexer.Scan()

	if _, ok := err.(*simplexer.UnknownTokenError); ok {
		l.Error(l.unknownTokenErrMsg())
	} else if err != nil {
		l.Error(l.errMsg())
	}

	if token == nil {
		return -1
	}

	lval.token = token
	l.source = l.convertSourceInfo(token)
	return int(token.Type.GetID())
}

func (l *Lexer) unknownTokenErrMsg() string {
	var out bytes.Buffer
	out.WriteString("Lexer Error: unknown token was found\n")
	out.WriteString("after " + l.source.Pos.String() + "\n")
	out.WriteString(l.source.Line + "\n")
	return out.String()
}

func (l *Lexer) errMsg() string {
	var out bytes.Buffer
	out.WriteString("Lexer Error:\n")
	out.WriteString("after " + l.source.Pos.String() + "\n")
	out.WriteString(l.source.Line + "\n")
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
	}
}
