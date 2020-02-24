%{

package parser

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"sort"
	"strconv"
	"strings"

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
	pairList []*ast.Pair
	kwargPair *ast.KwargPair
	stmt  ast.Stmt
	stmts []ast.Stmt
	program *ast.Program
}

%type<program> program
%type<stmts> stmts
%type<stmt> stmt exprStmt
%type<expr> expr infixExpr prefixExpr callExpr
%type<expr> literal funcLiteral arrLiteral objLiteral mapLiteral strLiteral symLiteral
%type<kwargPair> kwargPair
%type<pair> pair
%type<pairList> pairList
%type<argList> argList callArgs
%type<paramList> paramList funcParams
%type<exprList> exprList kwargExpansionList
%type<ident> ident
%type<chain> chain
%type<token> opMethod breakLine
%type<token> lBrace lParen lBracket comma mapLBrace methodLBrace

%token<token> INT SYMBOL CHAR_STR
%token<token> DOUBLE_STAR PLUS MINUS STAR SLASH BANG DOUBLE_SLASH PERCENT
%token<token> SPACESHIP EQ NEQ LT LE GT GE
%token<token> BIT_LSHIFT BIT_RSHIFT BIT_AND BIT_OR BIT_XOR BIT_NOT
%token<token> AND OR IADD ISUB
%token<token> ADD_CHAIN MAIN_CHAIN
%token<token> IDENT PRIVATE_IDENT
%token<token> LPAREN RPAREN COMMA COLON LBRACE RBRACE VERT LBRACKET RBRACKET
%token<token> MAP_LBRACE METHOD_LBRACE
%token<token> RET SEMICOLON

%left OR
%left AND
%left SPACESHIP EQ NEQ LT LE GT GE
%left BIT_OR BIT_XOR
%left BIT_AND
%left BIT_LSHIFT BIT_RSHIFT
%left PLUS MINUS
%left STAR SLASH DOUBLE_SLASH PERCENT
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
	| RET
	{
		$$ = &ast.Program{Stmts: []ast.Stmt{}}
		yylex.(*Lexer).program = $$
		yylex.(*Lexer).curRule = "program -> RET"
	}
	|
	{
		$$ = &ast.Program{Stmts: []ast.Stmt{}}
		yylex.(*Lexer).program = $$
		yylex.(*Lexer).curRule = "program -> (nothing)"
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
	| objLiteral
	{
		$$ = $1
		yylex.(*Lexer).curRule = "literal -> objLiteral"
	}
	| mapLiteral
	{
		$$ = $1
		yylex.(*Lexer).curRule = "literal -> mapLiteral"
	}
	| strLiteral
	{
		$$ = $1
		yylex.(*Lexer).curRule = "literal -> strLiteral"
	}
	| symLiteral
	{
		$$ = $1
		yylex.(*Lexer).curRule = "literal -> symLiteral"
	}

infixExpr
	: expr DOUBLE_STAR expr
	{
		$$ = &ast.InfixExpr{
			Token: $2.Literal,
			Left: $1,
			Operator: $2.Literal,
			Right: $3,
			Src: yylex.(*Lexer).Source,
		}
		yylex.(*Lexer).curRule = "infixExpr -> expr DOUBLE_STAR expr"
	}
	| expr STAR expr
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
	| expr DOUBLE_SLASH expr
	{
		$$ = &ast.InfixExpr{
			Token: $2.Literal,
			Left: $1,
			Operator: $2.Literal,
			Right: $3,
			Src: yylex.(*Lexer).Source,
		}
		yylex.(*Lexer).curRule = "infixExpr -> expr DOUBLE_SLASH expr"
	}
	| expr PERCENT expr
	{
		$$ = &ast.InfixExpr{
			Token: $2.Literal,
			Left: $1,
			Operator: $2.Literal,
			Right: $3,
			Src: yylex.(*Lexer).Source,
		}
		yylex.(*Lexer).curRule = "infixExpr -> expr PERCENT expr"
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
	| expr BIT_LSHIFT expr
	{
		$$ = &ast.InfixExpr{
			Token: $2.Literal,
			Left: $1,
			Operator: $2.Literal,
			Right: $3,
			Src: yylex.(*Lexer).Source,
		}
		yylex.(*Lexer).curRule = "infixExpr -> expr BIT_LSHIFT expr"
	}
	| expr BIT_RSHIFT expr
	{
		$$ = &ast.InfixExpr{
			Token: $2.Literal,
			Left: $1,
			Operator: $2.Literal,
			Right: $3,
			Src: yylex.(*Lexer).Source,
		}
		yylex.(*Lexer).curRule = "infixExpr -> expr BIT_RSHIFT expr"
	}
	| expr BIT_AND expr
	{
		$$ = &ast.InfixExpr{
			Token: $2.Literal,
			Left: $1,
			Operator: $2.Literal,
			Right: $3,
			Src: yylex.(*Lexer).Source,
		}
		yylex.(*Lexer).curRule = "infixExpr -> expr BIT_AND expr"
	}
	| expr BIT_OR expr
	{
		$$ = &ast.InfixExpr{
			Token: $2.Literal,
			Left: $1,
			Operator: $2.Literal,
			Right: $3,
			Src: yylex.(*Lexer).Source,
		}
		yylex.(*Lexer).curRule = "infixExpr -> expr BIT_OR expr"
	}
	| expr BIT_XOR expr
	{
		$$ = &ast.InfixExpr{
			Token: $2.Literal,
			Left: $1,
			Operator: $2.Literal,
			Right: $3,
			Src: yylex.(*Lexer).Source,
		}
		yylex.(*Lexer).curRule = "infixExpr -> expr BIT_XOR expr"
	}
	| expr SPACESHIP expr
	{
		$$ = &ast.InfixExpr{
			Token: $2.Literal,
			Left: $1,
			Operator: $2.Literal,
			Right: $3,
			Src: yylex.(*Lexer).Source,
		}
		yylex.(*Lexer).curRule = "infixExpr -> expr SPACESHIP expr"
	}
	| expr EQ expr
	{
		$$ = &ast.InfixExpr{
			Token: $2.Literal,
			Left: $1,
			Operator: $2.Literal,
			Right: $3,
			Src: yylex.(*Lexer).Source,
		}
		yylex.(*Lexer).curRule = "infixExpr -> expr EQ expr"
	}
	| expr NEQ expr
	{
		$$ = &ast.InfixExpr{
			Token: $2.Literal,
			Left: $1,
			Operator: $2.Literal,
			Right: $3,
			Src: yylex.(*Lexer).Source,
		}
		yylex.(*Lexer).curRule = "infixExpr -> expr NEQ expr"
	}
	| expr LT expr
	{
		$$ = &ast.InfixExpr{
			Token: $2.Literal,
			Left: $1,
			Operator: $2.Literal,
			Right: $3,
			Src: yylex.(*Lexer).Source,
		}
		yylex.(*Lexer).curRule = "infixExpr -> expr LT expr"
	}
	| expr GT expr
	{
		$$ = &ast.InfixExpr{
			Token: $2.Literal,
			Left: $1,
			Operator: $2.Literal,
			Right: $3,
			Src: yylex.(*Lexer).Source,
		}
		yylex.(*Lexer).curRule = "infixExpr -> expr GT expr"
	}
	| expr LE expr
	{
		$$ = &ast.InfixExpr{
			Token: $2.Literal,
			Left: $1,
			Operator: $2.Literal,
			Right: $3,
			Src: yylex.(*Lexer).Source,
		}
		yylex.(*Lexer).curRule = "infixExpr -> expr LE expr"
	}
	| expr GE expr
	{
		$$ = &ast.InfixExpr{
			Token: $2.Literal,
			Left: $1,
			Operator: $2.Literal,
			Right: $3,
			Src: yylex.(*Lexer).Source,
		}
		yylex.(*Lexer).curRule = "infixExpr -> expr GE expr"
	}
	| expr AND expr
	{
		$$ = &ast.InfixExpr{
			Token: $2.Literal,
			Left: $1,
			Operator: $2.Literal,
			Right: $3,
			Src: yylex.(*Lexer).Source,
		}
		yylex.(*Lexer).curRule = "infixExpr -> expr AND expr"
	}
	| expr OR expr
	{
		$$ = &ast.InfixExpr{
			Token: $2.Literal,
			Left: $1,
			Operator: $2.Literal,
			Right: $3,
			Src: yylex.(*Lexer).Source,
		}
		yylex.(*Lexer).curRule = "infixExpr -> expr OR expr"
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
	| BIT_NOT expr %prec UNARY_OP
	{
		$$ = &ast.PrefixExpr{
			Token: $1.Literal,
			Operator: $1.Literal,
			Right: $2,
			Src: yylex.(*Lexer).Source,
		}
		yylex.(*Lexer).curRule = "prefixExpr -> BIT_NOT expr"
	}

objLiteral
	: lBrace RBRACE
	{
		$$ = &ast.ObjLiteral{
			Token: $1.Literal,
			Pairs: []*ast.Pair{},
			EmbeddedExprs: []ast.Expr{},
			Src: yylex.(*Lexer).Source,
		}
	}
	| lBrace pairList RBRACE
	{
		$$ = &ast.ObjLiteral{
			Token: $1.Literal,
			Pairs: $2,
			EmbeddedExprs: []ast.Expr{},
			Src: yylex.(*Lexer).Source,
		}
	}
	| lBrace pairList RET RBRACE
	{
		$$ = &ast.ObjLiteral{
			Token: $1.Literal,
			Pairs: $2,
			EmbeddedExprs: []ast.Expr{},
			Src: yylex.(*Lexer).Source,
		}
	}
	| lBrace pairList comma RBRACE
	{
		$$ = &ast.ObjLiteral{
			Token: $1.Literal,
			Pairs: $2,
			EmbeddedExprs: []ast.Expr{},
			Src: yylex.(*Lexer).Source,
		}
	}
	| lBrace kwargExpansionList RBRACE
	{
		$$ = &ast.ObjLiteral{
			Token: $1.Literal,
			Pairs: []*ast.Pair{},
			EmbeddedExprs: $2,
			Src: yylex.(*Lexer).Source,
		}
	}
	| lBrace kwargExpansionList RET RBRACE
	{
		$$ = &ast.ObjLiteral{
			Token: $1.Literal,
			Pairs: []*ast.Pair{},
			EmbeddedExprs: $2,
			Src: yylex.(*Lexer).Source,
		}
	}
	| lBrace kwargExpansionList comma RBRACE
	{
		$$ = &ast.ObjLiteral{
			Token: $1.Literal,
			Pairs: []*ast.Pair{},
			EmbeddedExprs: $2,
			Src: yylex.(*Lexer).Source,
		}
	}
	| lBrace pairList comma kwargExpansionList RBRACE
	{
		$$ = &ast.ObjLiteral{
			Token: $1.Literal,
			Pairs: $2,
			EmbeddedExprs: $4,
			Src: yylex.(*Lexer).Source,
		}
	}
	| lBrace pairList comma kwargExpansionList RET RBRACE
	{
		$$ = &ast.ObjLiteral{
			Token: $1.Literal,
			Pairs: $2,
			EmbeddedExprs: $4,
			Src: yylex.(*Lexer).Source,
		}
	}
	| lBrace pairList comma kwargExpansionList comma RBRACE
	{
		$$ = &ast.ObjLiteral{
			Token: $1.Literal,
			Pairs: $2,
			EmbeddedExprs: $4,
			Src: yylex.(*Lexer).Source,
		}
	}

mapLiteral
	: mapLBrace RBRACE
	{
		$$ = &ast.MapLiteral{
			Token: $1.Literal,
			Pairs: []*ast.Pair{},
			EmbeddedExprs: []ast.Expr{},
			Src: yylex.(*Lexer).Source,
		}
	}
	| mapLBrace pairList RBRACE
	{
		$$ = &ast.MapLiteral{
			Token: $1.Literal,
			Pairs: $2,
			EmbeddedExprs: []ast.Expr{},
			Src: yylex.(*Lexer).Source,
		}
	}
	| mapLBrace pairList RET RBRACE
	{
		$$ = &ast.MapLiteral{
			Token: $1.Literal,
			Pairs: $2,
			EmbeddedExprs: []ast.Expr{},
			Src: yylex.(*Lexer).Source,
		}
	}
	| mapLBrace pairList comma RBRACE
	{
		$$ = &ast.MapLiteral{
			Token: $1.Literal,
			Pairs: $2,
			EmbeddedExprs: []ast.Expr{},
			Src: yylex.(*Lexer).Source,
		}
	}
	| mapLBrace kwargExpansionList RBRACE
	{
		$$ = &ast.MapLiteral{
			Token: $1.Literal,
			Pairs: []*ast.Pair{},
			EmbeddedExprs: $2,
			Src: yylex.(*Lexer).Source,
		}
	}
	| mapLBrace kwargExpansionList RET RBRACE
	{
		$$ = &ast.MapLiteral{
			Token: $1.Literal,
			Pairs: []*ast.Pair{},
			EmbeddedExprs: $2,
			Src: yylex.(*Lexer).Source,
		}
	}
	| mapLBrace kwargExpansionList comma RBRACE
	{
		$$ = &ast.MapLiteral{
			Token: $1.Literal,
			Pairs: []*ast.Pair{},
			EmbeddedExprs: $2,
			Src: yylex.(*Lexer).Source,
		}
	}
	| mapLBrace pairList comma kwargExpansionList RBRACE
	{
		$$ = &ast.MapLiteral{
			Token: $1.Literal,
			Pairs: $2,
			EmbeddedExprs: $4,
			Src: yylex.(*Lexer).Source,
		}
	}
	| mapLBrace pairList comma kwargExpansionList RET RBRACE
	{
		$$ = &ast.MapLiteral{
			Token: $1.Literal,
			Pairs: $2,
			EmbeddedExprs: $4,
			Src: yylex.(*Lexer).Source,
		}
	}
	| mapLBrace pairList comma kwargExpansionList comma RBRACE
	{
		$$ = &ast.MapLiteral{
			Token: $1.Literal,
			Pairs: $2,
			EmbeddedExprs: $4,
			Src: yylex.(*Lexer).Source,
		}
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
	| lBracket exprList RET RBRACKET
	{
		$$ = &ast.ArrLiteral{
			Token: $1.Literal,
			Elems: $2,
			Src: yylex.(*Lexer).Source,

		}
		yylex.(*Lexer).curRule = "arrLiteral -> lBracket exprList RBRACKET"
	}
	| lBracket exprList comma RBRACKET
	{
		$$ = &ast.ArrLiteral{
			Token: $1.Literal,
			Elems: $2,
			Src: yylex.(*Lexer).Source,

		}
		yylex.(*Lexer).curRule = "arrLiteral -> lBracket exprList RBRACKET"
	}

strLiteral
	: CHAR_STR
	{
		$$ = &ast.StrLiteral{
			Token: $1.Literal,
			Value: $1.Literal[1:],
			IsRaw: false,
			Src: yylex.(*Lexer).Source,
		}
	}

symLiteral
	: SYMBOL
	{
		$$ = &ast.SymLiteral{
			Token: $1.Literal,
			Value: $1.Literal[1:],
			Src: yylex.(*Lexer).Source,
		}
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
	| methodLBrace RBRACE
	{
		$$ = &ast.FuncLiteral{
			Token: $1.Literal,
			Args: ast.SelfIdentParamList(yylex.(*Lexer).Source).Args,
			Kwargs: map[*ast.Ident]ast.Expr{},
			Body: []ast.Stmt{},
			Src: yylex.(*Lexer).Source,
		}
		yylex.(*Lexer).curRule = "funcLiteral -> methodLBrace RBRACE"
	}
	| methodLBrace funcParams RBRACE
	{
		$$ = &ast.FuncLiteral{
			Token: $1.Literal,
			Args: $2.PrependSelf(yylex.(*Lexer).Source).Args,
			Kwargs: $2.Kwargs,
			Body: []ast.Stmt{},
			Src: yylex.(*Lexer).Source,
		}
		yylex.(*Lexer).curRule = "funcLiteral -> methodLBrace funcParams RBRACE"
	}
	| methodLBrace stmts RBRACE
	{
		$$ = &ast.FuncLiteral{
			Token: $1.Literal,
			Args: ast.SelfIdentParamList(yylex.(*Lexer).Source).Args,
			Kwargs: map[*ast.Ident]ast.Expr{},
			Body: $2,
			Src: yylex.(*Lexer).Source,
		}
		yylex.(*Lexer).curRule = "funcLiteral -> methodLBrace stmts RBRACE"
	}
	| methodLBrace funcParams stmts RBRACE
	{
		$$ = &ast.FuncLiteral{
			Token: $1.Literal,
			Args: $2.PrependSelf(yylex.(*Lexer).Source).Args,
			Kwargs: $2.Kwargs,
			Body: $3,
			Src: yylex.(*Lexer).Source,
		}
		yylex.(*Lexer).curRule = "funcLiteral -> methodLBrace funcParams stmts RBRACE"
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
	| OR
	{
		$$ = &ast.ParamList{
			Args: []*ast.Ident{},
			Kwargs: map[*ast.Ident]ast.Expr{},
		}
		yylex.(*Lexer).curRule = "funcParams -> OR"
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
	| callArgs funcLiteral
	{
		$$ = $1.AppendArg($2)
		yylex.(*Lexer).curRule = "callArgs -> callArgs funcLiteral"
	}
	| funcLiteral
	{
		$$ = ast.ExprToArgList($1)
		yylex.(*Lexer).curRule = "callArgs -> funcLiteral"
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
	| DOUBLE_SLASH
	{
		$$ = $1
		yylex.(*Lexer).curRule = "opMethod -> DOUBLE_SLASH"
	}
	| PERCENT
	{
		$$ = $1
		yylex.(*Lexer).curRule = "opMethod -> PERCENT"
	}
	| DOUBLE_STAR
	{
		$$ = $1
		yylex.(*Lexer).curRule = "opMethod -> DOUBLE_STAR"
	}
	| SPACESHIP
	{
		$$ = $1
		yylex.(*Lexer).curRule = "opMethod -> SPACESHIP"
	}
	| EQ
	{
		$$ = $1
		yylex.(*Lexer).curRule = "opMethod -> EQ"
	}
	| NEQ
	{
		$$ = $1
		yylex.(*Lexer).curRule = "opMethod -> NEQ"
	}
	| GE
	{
		$$ = $1
		yylex.(*Lexer).curRule = "opMethod -> GE"
	}
	| LE
	{
		$$ = $1
		yylex.(*Lexer).curRule = "opMethod -> LE"
	}
	| GT
	{
		$$ = $1
		yylex.(*Lexer).curRule = "opMethod -> GT"
	}
	| LT
	{
		$$ = $1
		yylex.(*Lexer).curRule = "opMethod -> LT"
	}
	| BIT_LSHIFT
	{
		$$ = $1
		yylex.(*Lexer).curRule = "opMethod -> BIT_LSHIFT"
	}
	| BIT_RSHIFT
	{
		$$ = $1
		yylex.(*Lexer).curRule = "opMethod -> BIT_RSHIFT"
	}
	| BIT_AND
	{
		$$ = $1
		yylex.(*Lexer).curRule = "opMethod -> BIT_AND"
	}
	| BIT_OR
	{
		$$ = $1
		yylex.(*Lexer).curRule = "opMethod -> BIT_OR"
	}
	| BIT_XOR
	{
		$$ = $1
		yylex.(*Lexer).curRule = "opMethod -> BIT_XOR"
	}
	| BIT_NOT
	{
		$$ = $1
		yylex.(*Lexer).curRule = "opMethod -> BIT_NOT"
	}
	| BANG
	{
		$$ = $1
		yylex.(*Lexer).curRule = "opMethod -> BANG"
	}
	| IADD
	{
		$$ = $1
		yylex.(*Lexer).curRule = "opMethod -> IADD"
	}
	| ISUB
	{
		$$ = $1
		yylex.(*Lexer).curRule = "opMethod -> ISUB"
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
	| argList comma kwargPair
	{
		$$ = $1.AppendKwarg($3.Key, $3.Val)
		yylex.(*Lexer).curRule = "argList -> argList comma pair"
	}
	| expr
	{
		$$ = ast.ExprToArgList($1)
		yylex.(*Lexer).curRule = "argList -> expr"
	}
	| kwargPair
	{
		$$ = ast.KwargPairToArgList($1)
		yylex.(*Lexer).curRule = "argList -> pair"
	}

paramList
	: paramList comma ident
	{
		$$ = $1.AppendArg($3)
		yylex.(*Lexer).curRule = "paramList -> paramList comma ident"
	}
	| paramList comma kwargPair
	{
		$$ = $1.AppendKwarg($3.Key, $3.Val)
		yylex.(*Lexer).curRule = "paramList -> paramList comma pair"
	}
	| ident
	{
		$$ = ast.IdentToParamList($1)
		yylex.(*Lexer).curRule = "paramList -> ident"
	}
	| kwargPair
	{
		$$ = ast.KwargPairToParamList($1)
		yylex.(*Lexer).curRule = "paramList -> pair"
	}

pairList
	: pairList comma pair
	{
		$$ = append($1, $3)
	}
	| pair
	{
		$$ = []*ast.Pair{$1}
	}

kwargExpansionList
	: kwargExpansionList comma DOUBLE_STAR expr
	{
		$$ = append($1, $4)
	}
	| DOUBLE_STAR expr
	{
		$$ = []ast.Expr{$2}
	}

kwargPair
	: ident COLON expr
	{
		$$ = &ast.KwargPair{Key: $1, Val: $3}
		yylex.(*Lexer).curRule = "kwargPair -> ident COLON expr"
	}

pair
	: expr COLON expr
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

mapLBrace
	: MAP_LBRACE
	{
		$$ = $1
		yylex.(*Lexer).curRule = "mapLBrace -> MAP_LBRACE"
	}
	| MAP_LBRACE RET
	{
		$$ = $1
		yylex.(*Lexer).curRule = "mapLBrace -> MAP_LBRACE RET"
	}

methodLBrace
	: METHOD_LBRACE
	{
		$$ = $1
		yylex.(*Lexer).curRule = "mapLBrace -> METHOD_LBRACE"
	}
	| METHOD_LBRACE RET
	{
		$$ = $1
		yylex.(*Lexer).curRule = "mapLBrace -> METHOD_LBRACE RET"
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

func tokenTypes() []simplexer.TokenType{
	// NOTE: make operators map to generate symbol regex easily
	// (for operator symbol such as '+ or '<=>)
	methodOps := map[string]string{
		"spaceship": `<=>`,
		"eq": `==`,
		"neq": `!=`,
		"ge": `>=`,
		"le": `<=`,
		"gt": `>`,
		"lt": `<`,
		"bitLShift": `<<`,
		"bitRShift": `>>`,
		"bitAnd": `/&`,
		"bitOr": `/\|`,
		"bitXor": `/\^`,
		"bitNot": `/~`,
		"bang": `!`,
		"plus": `\+`,
		"minus": `\-`,
		"star": `\*`,
		"doubleStar": `\*\*`,
		"slash": `/`,
		// NOTE: Do not use backquote! (otherwise commented out)
		"doubleSlash": "//",
		"percent": `%`,
		"iAdd": `\-%`,
		"iSub": `\+%`,
	}

	ident := `[a-zA-Z][a-zA-Z0-9_]*[!?]?`

	methodOpTokens := []string{}
	for _, op := range methodOps {
		methodOpTokens = append(methodOpTokens, op)
	}

	// sort by each token length (the longer, the earlier)
	sort.Slice(methodOpTokens, func(i, j int) bool {
		return len(methodOpTokens[i]) > len(methodOpTokens[j])
	})
	// NOTE: order is important!
	// in regex group with pipes, first match is selected (not the longest one!)
	// Therefore, a long token is never matched if there is a substring token
	// for that reason, sort methodOpTokens by token length
	// (e.g. : `'>>` should be tokenized [`'>>`], not [`'>`, `>`])

	// ident or private_ident or methodOps
	symbolable := fmt.Sprintf(`(%s|_+(%s)?|(%s))`,
		ident, ident, strings.Join(methodOpTokens, "|"))

	t := simplexer.NewRegexpTokenType

	// NOTE: order is important (the longer, the earlier)!
	// otherwise longer token is divided to shorter tokens unexpectedly
	// (e.g. : `>>` should be recognized one token (not `>` `>`))
	return []simplexer.TokenType{
		t(INT, `[0-9]+(\.[0-9]+)?`),
		t(CHAR_STR, `\?(\\[snt\\]|[^\r\n\\])`),
		// NOTE: comment(, which starts with "#") is included in RET
		// `#[^\n\r]*` is neseccery to lex final line comment (i.e. `#`)
		t(RET, `((#[^\n\r]*)?(\r|\n|\r\n)+|#[^\n\r]*)`),
		t(SYMBOL, "'"+symbolable),
		t(SPACESHIP, methodOps["spaceship"]),
		t(DOUBLE_STAR, methodOps["doubleStar"]),
		t(DOUBLE_SLASH, methodOps["doubleSlash"]),
		t(BIT_LSHIFT, methodOps["bitLShift"]),
		t(BIT_RSHIFT, methodOps["bitRShift"]),
		t(EQ, methodOps["eq"]),
		t(NEQ, methodOps["neq"]),
		t(GE, methodOps["ge"]),
		t(LE, methodOps["le"]),
		t(AND, `&&`),
		t(OR, `\|\|`),
		t(BIT_AND, methodOps["bitAnd"]),
		t(BIT_OR, methodOps["bitOr"]),
		t(BIT_XOR, methodOps["bitXor"]),
		t(BIT_NOT, methodOps["bitNot"]),
		t(IADD, methodOps["iAdd"]),
		t(ISUB, methodOps["iSub"]),
		t(MAP_LBRACE, `%\{`),
		t(METHOD_LBRACE, `m\{`),
		t(LPAREN, `\(`),
		t(RPAREN, `\)`),
		t(VERT, `\|`),
		t(LBRACE, `\{`),
		t(RBRACE, `\}`),
		t(LBRACKET, `\[`),
		t(RBRACKET, `\]`),
		t(COMMA, `,`),
		t(COLON, `:`),
		t(SEMICOLON, `;`),
		t(BANG, methodOps["bang"]),
		t(PLUS, methodOps["plus"]),
		t(MINUS, methodOps["minus"]),
		t(STAR, methodOps["star"]),
		t(SLASH, methodOps["slash"]),
		t(PERCENT, methodOps["percent"]),
		t(GT, methodOps["gt"]),
		t(LT, methodOps["lt"]),
		t(ADD_CHAIN, `[&~=]`),
		t(MAIN_CHAIN, `[\.@$]`),
		t(IDENT, ident),
		t(PRIVATE_IDENT, fmt.Sprintf(`_+(%s)?`, ident)),
	}
}

func NewLexer(reader io.Reader) *Lexer {
	l := simplexer.NewLexer(reader)
	l.TokenTypes = tokenTypes()
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
