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
	pair *ast.Pair
	stmt  ast.Stmt
	stmts []ast.Stmt
	program *ast.Program
}

%type<program> program
%type<stmts> stmts
%type<stmt> stmt exprStmt
%type<expr> expr literal infixExpr callExpr
%type<pair> pair
%type<argList> argList
%type<ident> ident
%type<chain> chain
%type<token> opMethod
%token<token> INT
%token<token> PREC1_OPERATOR PREC2_OPERATOR
%token<token> ADD_CHAIN MAIN_CHAIN
%token<token> IDENT PRIVATE_IDENT
%token<token> LPAREN RPAREN COMMA COLON
%left PREC1_OPERATOR
%left PREC2_OPERATOR
%left ADD_CHAIN MAIN_CHAIN

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
	| callExpr
	{
		$$ = $1
	}
	| ident
	{
		$$ = $1
	}

ident
	: IDENT
	{
		$$ = &ast.Ident{
			Token: $1.Literal,
			Value: $1.Literal,
			Src: yylex.(*Lexer).source,
			IsPrivate: false,
		}
	}
	| PRIVATE_IDENT
	{
		$$ = &ast.Ident{
			Token: $1.Literal,
			Value: $1.Literal,
			Src: yylex.(*Lexer).source,
			IsPrivate: true,
		}
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
			Src: yylex.(*Lexer).source,
		}
	}
	| expr chain ident LPAREN RPAREN
	{
		$$ = &ast.PropCallExpr{
			Token: "(propCall)",
			Chain: $2,
			Receiver: $1,
			Prop: $3,
			Args: nil,
			Kwargs: nil,
			Src: yylex.(*Lexer).source,
		}
	}
	| expr chain ident LPAREN argList RPAREN
	{
		$$ = &ast.PropCallExpr{
			Token: "(propCall)",
			Chain: $2,
			Receiver: $1,
			Prop: $3,
			Args: $5.Args,
			Kwargs: $5.Kwargs,
			Src: yylex.(*Lexer).source,
		}
	}
	| expr chain opMethod
	{
		opIdent := &ast.Ident{
			Token: $3.Literal,
			Value: $3.Literal,
			Src: yylex.(*Lexer).source,
			IsPrivate: true,
		}
		$$ = &ast.PropCallExpr{
			Token: "(propCall)",
			Chain: $2,
			Receiver: $1,
			Prop: opIdent,
			Args: nil,
			Kwargs: nil,
			Src: yylex.(*Lexer).source,
		}
	}
	| expr chain opMethod LPAREN RPAREN
	{
		opIdent := &ast.Ident{
			Token: $3.Literal,
			Value: $3.Literal,
			Src: yylex.(*Lexer).source,
			IsPrivate: true,
		}
		$$ = &ast.PropCallExpr{
			Token: "(propCall)",
			Chain: $2,
			Receiver: $1,
			Prop: opIdent,
			Args: nil,
			Kwargs: nil,
			Src: yylex.(*Lexer).source,
		}
	}
	| expr chain opMethod LPAREN argList RPAREN
	{
		opIdent := &ast.Ident{
			Token: $3.Literal,
			Value: $3.Literal,
			Src: yylex.(*Lexer).source,
			IsPrivate: true,
		}
		$$ = &ast.PropCallExpr{
			Token: "(propCall)",
			Chain: $2,
			Receiver: $1,
			Prop: opIdent,
			Args: $5.Args,
			Kwargs: $5.Kwargs,
			Src: yylex.(*Lexer).source,
		}
	}


opMethod
	: PREC1_OPERATOR
	{
		$$ = $1
	}
	| PREC2_OPERATOR
	{
		$$ = $1
	}

chain
	: ADD_CHAIN MAIN_CHAIN
	{
		$$ = ast.MakeChain($1.Literal, $2.Literal, nil)
	}
	| MAIN_CHAIN
	{
		$$ = ast.MakeChain("", $1.Literal, nil)
	}
	| MAIN_CHAIN LPAREN expr RPAREN
	{
		$$ = ast.MakeChain("", $1.Literal, $3)
	}
	| ADD_CHAIN MAIN_CHAIN LPAREN expr RPAREN
	{
		$$ = ast.MakeChain($1.Literal, $2.Literal, $4)
	}

argList
	: argList COMMA expr
	{
		$$ = $1.AppendArg($3)
	}
	| argList COMMA pair
	{
		$$ = $1.AppendKwarg($3.Key, $3.Val)
	}
	| expr
	{
		$$ = ast.ExprToArgList($1)
	}
	| pair
	{
		$$ = ast.PairToArgList($1)
	}

pair
	: ident COLON expr
	{
		$$ = &ast.Pair{Key: $1, Val: $3}
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
	simplexer.NewRegexpTokenType(LPAREN, `\(`),
	simplexer.NewRegexpTokenType(RPAREN, `\)`),
	simplexer.NewRegexpTokenType(COMMA, `,`),
	simplexer.NewRegexpTokenType(COLON, `:`),
	simplexer.NewRegexpTokenType(PREC1_OPERATOR, `[-+]`),
	simplexer.NewRegexpTokenType(PREC2_OPERATOR, `[*/]`),
	simplexer.NewRegexpTokenType(ADD_CHAIN, `[&~=]`),
	simplexer.NewRegexpTokenType(MAIN_CHAIN, `[\.@$]`),
	simplexer.NewRegexpTokenType(IDENT, `[a-zA-Z][a-zA-Z0-9_]*([!?])?`),
	simplexer.NewRegexpTokenType(PRIVATE_IDENT, `_[a-zA-Z][a-zA-Z0-9_]*([!?])?`),
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
