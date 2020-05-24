%{

package parser

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"math"
	"sort"
	"strconv"
	"strings"

	"github.com/macrat/simplexer"

	"../ast"
)

%}

%union{
    token *simplexer.Token
	recvAndChain *ast.RecvAndChain
    chain *ast.Chain
	ident *ast.Ident
	expr  ast.Expr
	funcComponent ast.FuncComponent
	funcComponentList []*ast.FuncComponent
	argList *ast.ArgList
	exprList []ast.Expr
	pair *ast.Pair
	pairList []*ast.Pair
	kwargPair *ast.KwargPair
	formerStrPiece *ast.FormerStrPiece
	stmt  ast.Stmt
	stmts []ast.Stmt
	program *ast.Program
}

%type<program> program
%type<stmts> stmts
%type<stmt> stmt exprStmt jumpStmt jumpIfStmt
%type<expr> expr infixExpr prefixExpr assignExpr callExpr embeddedStr indexExpr ifExpr
%type<expr> literal unitExpr
%type<expr> intLiteral floatLiteral
%type<expr> funcLiteral iterLiteral matchLiteral diamondLiteral
%type<expr> arrLiteral objLiteral mapLiteral
%type<expr> strLiteral symLiteral
%type<expr> rangeLiteral bareRange
%type<expr> listElem
%type<funcComponent> funcComponent formalFuncComponent
%type<funcComponentList> funcComponentList
%type<kwargPair> kwargPair
%type<pair> pair
%type<pairList> pairList
%type<argList> argList callArgs funcParams
%type<exprList> exprList kwargExpansionList
%type<ident> ident
%type<chain> chain
%type<recvAndChain> recvAndChain
%type<formerStrPiece> formerStrPiece
%type<token> opMethod breakLine
%type<token> comma
%type<token> lBrace lParen lBracket mapLBrace methodMapLBrace methodLBrace lIter methodLIter

%token<token> INT FLOAT HEX_INT BIN_INT OCT_INT EXP_FLOAT EXP_INT
%token<token> SYMBOL CHAR_STR BACKQUOTE_STR DOUBLEQUOTE_STR
%token<token> HEAD_STR_PIECE MID_STR_PIECE TAIL_STR_PIECE
%token<token> DOUBLE_STAR PLUS MINUS STAR SLASH BANG DOUBLE_SLASH PERCENT
%token<token> SPACESHIP EQ NEQ LT LE GT GE
%token<token> BIT_LSHIFT BIT_RSHIFT BIT_AND BIT_OR BIT_XOR BIT_NOT
%token<token> AND OR IADD ISUB
%token<token> ADD_CHAIN MAIN_CHAIN MULTILINE_ADD_CHAIN MULTILINE_MAIN_CHAIN
%token<token> IDENT PRIVATE_IDENT ARG_IDENT KWARG_IDENT
%token<token> LPAREN RPAREN COMMA COLON LBRACE RBRACE VERT LBRACKET RBRACKET CARET
%token<token> MAP_LBRACE METHOD_MAP_LBRACE METHOD_LBRACE LITER RITER METHOD_LITER DIAMOND
%token<token> RET SEMICOLON
%token<token> ASSIGN COMPOUND_ASSIGN RIGHT_ASSIGN
%token<token> IF ELSE
%token<token> RETURN RAISE YIELD

%left IF
%left ELSE
%left JUMP
%left JUMPIF
%left RIGHT_ASSIGN
%right ASSIGN COMPOUND_ASSIGN
%left OR
%left AND
%left SPACESHIP EQ NEQ LT LE GT GE
%left BIT_OR BIT_XOR
%left BIT_AND
%left BIT_LSHIFT BIT_RSHIFT
%left PLUS MINUS
%left STAR SLASH DOUBLE_SLASH PERCENT
%left DOUBLE_STAR
%left MULTILINE_ADD_CHAIN MULTILINE_MAIN_CHAIN
%left ADD_CHAIN MAIN_CHAIN
%left UNARY_OP
%left GROUPING
%left INDEXING

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
	| jumpStmt
	{
		$$ = $1
		yylex.(*Lexer).curRule = "stmt -> jumpStmt"
	}
	| jumpIfStmt
	{
		$$ = $1
		yylex.(*Lexer).curRule = "stmt -> jumpIfStmt"
	}

exprStmt
	: expr
	{
		$$ = &ast.ExprStmt{
			Token: "(exprStmt)",
			Expr: $1,
			Src: yylex.(*Lexer).Source,
		}
		yylex.(*Lexer).curRule = "exprStmt -> expr"
	}

jumpIfStmt
	: jumpStmt IF expr %prec JUMPIF
	{
		$$ = &ast.JumpIfStmt{
			JumpStmt: $1.(*ast.JumpStmt),
			Cond: $3,
			Src: yylex.(*Lexer).Source,
		}
	}

jumpStmt
	: RETURN expr %prec JUMP
	{
		$$ = &ast.JumpStmt{
			Token: $1.Literal,
			Val: $2,
			JumpType: ast.ReturnJump,
			Src: yylex.(*Lexer).Source,
		}
	}
	| RAISE expr %prec JUMP
	{
		$$ = &ast.JumpStmt{
			Token: $1.Literal,
			Val: $2,
			JumpType: ast.RaiseJump,
			Src: yylex.(*Lexer).Source,
		}
	}
	| YIELD expr %prec JUMP
	{
		$$ = &ast.JumpStmt{
			Token: $1.Literal,
			Val: $2,
			JumpType: ast.YieldJump,
			Src: yylex.(*Lexer).Source,
		}
	}


expr
	: unitExpr
	{
		$$ = $1
		yylex.(*Lexer).curRule = "expr -> unitExpr"
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
	| assignExpr
	{
		$$ = $1
		yylex.(*Lexer).curRule = "expr -> assignExpr"
	}
	| ifExpr
	{
		$$ = $1
		yylex.(*Lexer).curRule = "expr -> ifExpr"
	}

unitExpr
	: literal
	{
		$$ = $1
	}
	| embeddedStr
	{
		$$ = $1
	}
	| callExpr
	{
		$$ = $1
	}
	| indexExpr
	{
		$$ = $1
	}
	| LPAREN expr RPAREN %prec GROUPING
	{
		$$ = $2
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
			Src: yylex.(*Lexer).Source,
			IsPrivate: false,
			IdentAttr: ast.NormalIdent,
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
			IdentAttr: ast.NormalIdent,
		}
		yylex.(*Lexer).curRule = "ident -> PRIVATE_IDENT"
	}
	| ARG_IDENT
	{
		$$ = &ast.Ident{
			Token: $1.Literal,
			Value: $1.Literal,
			Src: yylex.(*Lexer).Source,
			IsPrivate: true,
			IdentAttr: ast.ArgIdent,
		}
		yylex.(*Lexer).curRule = "ident -> ARG_IDENT"
	}
	| KWARG_IDENT
	{
		$$ = &ast.Ident{
			Token: $1.Literal,
			Value: $1.Literal,
			Src: yylex.(*Lexer).Source,
			IsPrivate: true,
			IdentAttr: ast.KwargIdent,
		}
		yylex.(*Lexer).curRule = "ident -> KWARG_IDENT"
	}

literal
	: intLiteral
	{
		$$ = $1
	}
	| floatLiteral
	{
		$$ = $1
	}
	| funcLiteral
	{
		$$ = $1
	}
	| matchLiteral
	{
		$$ = $1
	}
	| iterLiteral
	{
		$$ = $1
	}
	| diamondLiteral
	{
		$$ = $1
	}
	| arrLiteral
	{
		$$ = $1
	}
	| objLiteral
	{
		$$ = $1
	}
	| mapLiteral
	{
		$$ = $1
	}
	| strLiteral
	{
		$$ = $1
	}
	| symLiteral
	{
		$$ = $1
	}
	| rangeLiteral
	{
		$$ = $1
	}

intLiteral
	: INT
	{
		// remove separator "_"s
		intStr := strings.Replace($1.Literal, "_", "", -1)
		n, _ := strconv.ParseInt(intStr, 10, 64)
		$$ = &ast.IntLiteral{
			Token: $1.Literal,
			Value: n,
			Src: yylex.(*Lexer).Source,
		}
	}
	| HEX_INT
	{
		// remove separator "_"s
		lit := strings.Replace($1.Literal, "_", "", -1)
		// remove prefix "0x"
		intStr := lit[2:]
		n, _ := strconv.ParseInt(intStr, 16, 64)
		$$ = &ast.IntLiteral{
			Token: $1.Literal,
			Value: n,
			Src: yylex.(*Lexer).Source,
		}
	}
	| OCT_INT
	{
		// remove separator "_"s
		lit := strings.Replace($1.Literal, "_", "", -1)
		// remove prefix "0o"
		intStr := lit[2:]
		n, _ := strconv.ParseInt(intStr, 8, 64)
		$$ = &ast.IntLiteral{
			Token: $1.Literal,
			Value: n,
			Src: yylex.(*Lexer).Source,
		}
	}
	| BIN_INT
	{
		// remove separator "_"s
		lit := strings.Replace($1.Literal, "_", "", -1)
		// remove prefix "0b"
		intStr := lit[2:]
		n, _ := strconv.ParseInt(intStr, 2, 64)
		$$ = &ast.IntLiteral{
			Token: $1.Literal,
			Value: n,
			Src: yylex.(*Lexer).Source,
		}
	}
	| EXP_INT
	{
		// remove separator "_"s
		lit := strings.Replace($1.Literal, "_", "", -1)
		// NOTE: ToLower is nesessary (to split by both e and E)
		toks := strings.Split(strings.ToLower(lit), "e")
		// NOTE: cast float to deal with minus exp (i.e. `100e-2 == 1`)
		val, _ := strconv.ParseFloat(toks[0], 64)
		// NOTE: cannot use ParseInt (math.Pow requires float)
		exp, _ := strconv.ParseFloat(toks[1], 64)
		$$ = &ast.IntLiteral{
			Token: $1.Literal,
			Value: int64(val * math.Pow(10, exp)),
			Src: yylex.(*Lexer).Source,
		}
	}

floatLiteral
	: FLOAT
	{
		// remove separator "_"s
		floatStr := strings.Replace($1.Literal, "_", "", -1)
		n, _ := strconv.ParseFloat(floatStr, 64)
		$$ = &ast.FloatLiteral{
			Token: $1.Literal,
			Value: n,
			Src: yylex.(*Lexer).Source,
		}
	}
	| EXP_FLOAT
	{
		// remove separator "_"s
		lit := strings.Replace($1.Literal, "_", "", -1)
		// NOTE: ToLower is nesessary (to split by both e and E)
		toks := strings.Split(strings.ToLower(lit), "e")
		val, _ := strconv.ParseFloat(toks[0], 64)
		exp, _ := strconv.ParseFloat(toks[1], 64)
		$$ = &ast.FloatLiteral{
			Token: $1.Literal,
			Value: float64(val * math.Pow(10, exp)),
			Src: yylex.(*Lexer).Source,
		}
	} 

ifExpr
	: expr IF expr
	{
		$$ = &ast.IfExpr{
			Token: $2.Literal,
			Cond: $3,
			Then: $1,
			Else: nil,
			Src: yylex.(*Lexer).Source,
		}
	}
	| expr IF expr ELSE expr %prec ELSE
	{
		// NOTE: to refrain shift/reduce conflict, else has higher prec than if
		// `a if b if c else d` means `((a if b) if c else d)`
		$$ = &ast.IfExpr{
			Token: $2.Literal,
			Cond: $3,
			Then: $1,
			Else: $5,
			Src: yylex.(*Lexer).Source,
		}
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
		// HACK: convert -(number) to literal
		// TOFIX: deal with this process in lexer
		switch r := $2.(type) {
			case *ast.IntLiteral:
			$$ = &ast.IntLiteral{
				Token: "-" + r.Token,
				Value: -r.Value,
				Src: yylex.(*Lexer).Source,
			}
			case *ast.FloatLiteral:
			$$ = &ast.FloatLiteral{
				Token: "-" + r.Token,
				Value: -r.Value,
				Src: yylex.(*Lexer).Source,
			}
			default:
			$$ = &ast.PrefixExpr{
				Token: $1.Literal,
				Operator: $1.Literal,
				Right: $2,
				Src: yylex.(*Lexer).Source,
			}
		}		
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

assignExpr
	: ident ASSIGN expr
	{
		$$ = &ast.AssignExpr{
			Token: $2.Literal,
			Left: $1,
			Right: $3,
			Src: yylex.(*Lexer).Source,
		}
	}
	| ident COMPOUND_ASSIGN expr
	{
		op := $2.Literal[:len($2.Literal)-1]
		ie := &ast.InfixExpr{
			Token: op,
			Left: $1,
			Operator: op,
			Right: $3,
			Src: yylex.(*Lexer).Source,
		}
		$$ = &ast.AssignExpr{
			Token: ":=",
			Left: $1,
			Right: ie,
			Src: yylex.(*Lexer).Source,
		}
	}
	| expr RIGHT_ASSIGN ident
	{
		// NOTE: "Left" and "Right" are reversed!
		$$ = &ast.AssignExpr{
			Token: ":=",
			Left: $3,
			Right: $1,
			Src: yylex.(*Lexer).Source,
		}
	}

indexExpr
	: unitExpr arrLiteral %prec INDEXING
	{
		atIdent := &ast.Ident{
			Token: "at",
			Value: "at",
			Src: yylex.(*Lexer).Source,
			IsPrivate: false,
			IdentAttr: ast.NormalIdent,
		}
		$$ = &ast.PropCallExpr{
			Token: $1.TokenLiteral(),
			Chain: ast.MakeChain("", ".", nil),
			Receiver: $1,
			Prop: atIdent,
			Args: []ast.Expr{$2},
			Kwargs: map[*ast.Ident]ast.Expr{},
			Src: yylex.(*Lexer).Source,
		}
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
	| BACKQUOTE_STR
	{
		$$ = &ast.StrLiteral{
			Token: $1.Literal,
			Value: $1.Literal[1:len($1.Literal)-1],
			IsRaw: true,
			Src: yylex.(*Lexer).Source,
		}
	}
	| DOUBLEQUOTE_STR
	{
		$$ = &ast.StrLiteral{
			Token: $1.Literal,
			Value: $1.Literal[1:len($1.Literal)-1],
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

rangeLiteral
	: LPAREN bareRange RPAREN
	{
		$$ = $2
	}

bareRange
	: expr COLON expr COLON expr
	{
		$$ = &ast.RangeLiteral{
			Token: $2.Literal,
			Start: $1,
			Stop: $3,
			Step: $5,
			Src: yylex.(*Lexer).Source,
		}
	}
	| expr COLON expr
	{
		$$ = &ast.RangeLiteral{
			Token: $2.Literal,
			Start: $1,
			Stop: $3,
			Step: nil,
			Src: yylex.(*Lexer).Source,
		}
	}
	| expr COLON COLON expr
	{
		$$ = &ast.RangeLiteral{
			Token: $2.Literal,
			Start: $1,
			Stop: nil,
			Step: $4,
			Src: yylex.(*Lexer).Source,
		}
	}
	| COLON expr COLON expr
	{
		$$ = &ast.RangeLiteral{
			Token: $1.Literal,
			Start: nil,
			Stop: $2,
			Step: $4,
			Src: yylex.(*Lexer).Source,
		}
	}
	| expr COLON
	{
		$$ = &ast.RangeLiteral{
			Token: $2.Literal,
			Start: $1,
			Stop: nil,
			Step: nil,
			Src: yylex.(*Lexer).Source,
		}
	}
	| COLON expr
	{
		$$ = &ast.RangeLiteral{
			Token: $1.Literal,
			Start: nil,
			Stop: $2,
			Step: nil,
			Src: yylex.(*Lexer).Source,
		}
	}
	| COLON COLON expr
	{
		$$ = &ast.RangeLiteral{
			Token: $1.Literal,
			Start: nil,
			Stop: nil,
			Step: $3,
			Src: yylex.(*Lexer).Source,
		}
	}

embeddedStr
	: formerStrPiece TAIL_STR_PIECE
	{
		$$ = &ast.EmbeddedStr{
			Token: $1.Token,
			Former: $1,
			Latter: $2.Literal[1:len($2.Literal)-1],
			Src: yylex.(*Lexer).Source,
		}
	}

formerStrPiece
	: formerStrPiece MID_STR_PIECE expr
	{
		$$ = &ast.FormerStrPiece{
			Token: $1.Token,
			Former: $1,
			Str: $2.Literal[1:len($2.Literal)-2],
			Expr: $3,
		}
	}
	| HEAD_STR_PIECE expr
	{
		$$ = &ast.FormerStrPiece{
			Token: $1.Literal,
			Former: nil,
			Str: $1.Literal[1:len($1.Literal)-2],
			Expr: $2,
		}
	}

funcLiteral
	: lBrace funcComponent RBRACE
	{
		$$ = &ast.FuncLiteral{
			Token: $1.Literal,
			FuncComponent: $2,
			Src: yylex.(*Lexer).Source,
		}
	}
	| methodLBrace RBRACE
	{
		$$ = &ast.FuncLiteral{
			Token: $1.Literal,
			Src: yylex.(*Lexer).Source,
			FuncComponent: ast.FuncComponent{
				Args: ast.SelfIdentArgList(yylex.(*Lexer).Source).Args,
				Kwargs: map[*ast.Ident]ast.Expr{},
				Body: []ast.Stmt{},
				Src: yylex.(*Lexer).Source,
			},
		}
	}
	| methodLBrace funcComponent RBRACE
	{
		$$ = &ast.FuncLiteral{
			Token: $1.Literal,
			FuncComponent: *$2.PrependSelf(yylex.(*Lexer).Source),
			Src: yylex.(*Lexer).Source,
		}
	}

iterLiteral
	: lIter RITER
	{
		$$ = &ast.IterLiteral{
			Token: $1.Literal,
			Src: yylex.(*Lexer).Source,
			FuncComponent: ast.FuncComponent{
				Args: []ast.Expr{},
				Kwargs: map[*ast.Ident]ast.Expr{},
				Body: []ast.Stmt{},
				Src: yylex.(*Lexer).Source,
			},
		}
	}
	| lIter funcComponent RITER
	{
		$$ = &ast.IterLiteral{
			Token: $1.Literal,
			Src: yylex.(*Lexer).Source,
			FuncComponent: $2,
		}
	}
	| methodLIter RITER
	{
		$$ = &ast.IterLiteral{
			Token: $1.Literal,
			Src: yylex.(*Lexer).Source,
			FuncComponent: ast.FuncComponent{
				Args: ast.SelfIdentArgList(yylex.(*Lexer).Source).Args,
				Kwargs: map[*ast.Ident]ast.Expr{},
				Body: []ast.Stmt{},
				Src: yylex.(*Lexer).Source,
			},
		}
	}
	| methodLIter funcComponent RITER
	{
		$$ = &ast.IterLiteral{
			Token: $1.Literal,
			Src: yylex.(*Lexer).Source,
			FuncComponent: *$2.PrependSelf(yylex.(*Lexer).Source),
		}
	}

matchLiteral
	: mapLBrace funcComponentList RBRACE
	{
		$$ = &ast.MatchLiteral{
			Token: $1.Literal,
			Patterns: $2,
			Src: yylex.(*Lexer).Source,
		}
	}
	| methodMapLBrace funcComponentList RBRACE
	{
		patterns := []*ast.FuncComponent{}
		for _, p := range $2 {
			patterns = append(patterns, p.PrependSelf(yylex.(*Lexer).Source))
		}

		$$ = &ast.MatchLiteral{
			Token: $1.Literal,
			Patterns: patterns,
			Src: yylex.(*Lexer).Source,
		}
	}

funcComponentList
	: funcComponentList comma formalFuncComponent
	{
		// NOTE: assigning is nesessary because $3 is passed by reference
		// which means address of $3 is the last match of funcComponentList
		// (same as for loop)
		comp := $3
		$$ = append($1, &comp)
	}
	| formalFuncComponent
	{
		comp := $1
		$$ = []*ast.FuncComponent{&comp}
	}

funcComponent
	: funcParams
	{
		$$ = ast.FuncComponent{
			Args: $1.Args,
			Kwargs: $1.Kwargs,
			Body: []ast.Stmt{},
			Src: yylex.(*Lexer).Source,
		}
	}
	| stmts
	{
		$$ = ast.FuncComponent{
			Args: []ast.Expr{},
			Kwargs: map[*ast.Ident]ast.Expr{},
			Body: $1,
			Src: yylex.(*Lexer).Source,
		}
	}
	| formalFuncComponent
	{
		$$ = $1
	}

formalFuncComponent
	: funcParams stmts
	{
		$$ = ast.FuncComponent{
			Args: $1.Args,
			Kwargs: $1.Kwargs,
			Body: $2,
			Src: yylex.(*Lexer).Source,
		}
	}

diamondLiteral
	: DIAMOND
	{
		$$ = &ast.DiamondLiteral{
			Token: $1.Literal,
			Src: yylex.(*Lexer).Source,
		}
	}

funcParams
	: VERT VERT
	{
		$$ = &ast.ArgList{
			Args: []ast.Expr{},
			Kwargs: map[*ast.Ident]ast.Expr{},
		}
	}
	| OR
	{
		$$ = &ast.ArgList{
			Args: []ast.Expr{},
			Kwargs: map[*ast.Ident]ast.Expr{},
		}
	}
	| VERT argList VERT
	{
		$$ = $2
	}
	| VERT argList RET VERT
	{
		$$ = $2
	}
	|  VERT VERT RET
	{
		$$ = &ast.ArgList{
			Args: []ast.Expr{},
			Kwargs: map[*ast.Ident]ast.Expr{},
		}
	}
	|  OR RET
	{
		$$ = &ast.ArgList{
			Args: []ast.Expr{},
			Kwargs: map[*ast.Ident]ast.Expr{},
		}
	}
	| VERT argList VERT RET
	{
		$$ = $2
	}
	| VERT argList RET VERT RET
	{
		$$ = $2
	}

callExpr
	: recvAndChain ident
	{
		$$ = &ast.PropCallExpr{
			Token: "(propCall)",
			Chain: $1.Chain,
			Receiver: $1.Recv,
			Prop: $2,
			Args: []ast.Expr{},
			Kwargs: map[*ast.Ident]ast.Expr{},
			Src: yylex.(*Lexer).Source,
		}
	}
	| recvAndChain ident callArgs
	{
		$$ = &ast.PropCallExpr{
			Token: "(propCall)",
			Chain: $1.Chain,
			Receiver: $1.Recv,
			Prop: $2,
			Args: $3.Args,
			Kwargs: $3.Kwargs,
			Src: yylex.(*Lexer).Source,
		}
	}
	| recvAndChain opMethod
	{
		opIdent := &ast.Ident{
			Token: $2.Literal,
			Value: $2.Literal,
			Src: yylex.(*Lexer).Source,
			IsPrivate: true,
		}
		$$ = &ast.PropCallExpr{
			Token: "(propCall)",
			Chain: $1.Chain,
			Receiver: $1.Recv,
			Prop: opIdent,
			Args: []ast.Expr{},
			Kwargs: map[*ast.Ident]ast.Expr{},
			Src: yylex.(*Lexer).Source,
		}
	}
	| recvAndChain opMethod callArgs
	{
		opIdent := &ast.Ident{
			Token: $2.Literal,
			Value: $2.Literal,
			Src: yylex.(*Lexer).Source,
			IsPrivate: true,
		}
		$$ = &ast.PropCallExpr{
			Token: "(propCall)",
			Chain: $1.Chain,
			Receiver: $1.Recv,
			Prop: opIdent,
			Args: $3.Args,
			Kwargs: $3.Kwargs,
			Src: yylex.(*Lexer).Source,
		}
	}
	| recvAndChain funcLiteral
	{
		$$ = &ast.LiteralCallExpr{
			Token: "(literalCall)",
			Chain: $1.Chain,
			Receiver: $1.Recv,
			Func: $2.(*ast.FuncLiteral),
			Args: []ast.Expr{},
			Kwargs: map[*ast.Ident]ast.Expr{},
			Src: yylex.(*Lexer).Source,
		}
	}
	| recvAndChain funcLiteral callArgs
	{
		$$ = &ast.LiteralCallExpr{
			Token: "(literalCall)",
			Chain: $1.Chain,
			Receiver: $1.Recv,
			Func: $2.(*ast.FuncLiteral),
			Args: $3.Args,
			Kwargs: $3.Kwargs,
			Src: yylex.(*Lexer).Source,
		}
	}
	| recvAndChain CARET ident
	{
		$$ = &ast.VarCallExpr{
			Token: "(varCall)",
			Chain: $1.Chain,
			Receiver: $1.Recv,
			Var: $3,
			Args: []ast.Expr{},
			Kwargs: map[*ast.Ident]ast.Expr{},
			Src: yylex.(*Lexer).Source,
		}
	}
	| recvAndChain CARET ident callArgs
	{
		$$ = &ast.VarCallExpr{
			Token: "(varCall)",
			Chain: $1.Chain,
			Receiver: $1.Recv,
			Var: $3,
			Args: $4.Args,
			Kwargs: $4.Kwargs,
			Src: yylex.(*Lexer).Source,
		}
	}
	| ident callArgs
	{
		callIdent := &ast.Ident{
			Token: "call",
			Value: "call",
			Src: yylex.(*Lexer).Source,
			IsPrivate: false,
			IdentAttr: ast.NormalIdent,
		}
		$$ = &ast.PropCallExpr{
			Token: "(propCall)",
			Chain: ast.MakeChain("", ".", nil),
			Receiver: $1,
			Prop: callIdent,
			Args: $2.Args,
			Kwargs: $2.Kwargs,
			Src: yylex.(*Lexer).Source,
		}
	}
	| funcLiteral callArgs
	{
		callIdent := &ast.Ident{
			Token: "call",
			Value: "call",
			Src: yylex.(*Lexer).Source,
			IsPrivate: false,
			IdentAttr: ast.NormalIdent,
		}
		$$ = &ast.PropCallExpr{
			Token: "(propCall)",
			Chain: ast.MakeChain("", ".", nil),
			Receiver: $1,
			Prop: callIdent,
			Args: $2.Args,
			Kwargs: $2.Kwargs,
			Src: yylex.(*Lexer).Source,
		}
	}

recvAndChain
	: expr chain
	{
		$$ = &ast.RecvAndChain{
			Recv: $1,
			Chain: $2,
		}
	}
	| chain
	{
		$$ = &ast.RecvAndChain{
			Recv: nil,
			Chain: $1,
		}
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
	| MULTILINE_ADD_CHAIN MAIN_CHAIN
	{
		ac := string($1.Literal[len($1.Literal)-1])
		$$ = ast.MakeChain(ac, $2.Literal, nil)
	}
	| MULTILINE_MAIN_CHAIN
	{
		mc := string($1.Literal[len($1.Literal)-1])
		$$ = ast.MakeChain("", mc, nil)
	}
	| MULTILINE_MAIN_CHAIN lParen expr RPAREN
	{
		mc := string($1.Literal[len($1.Literal)-1])
		$$ = ast.MakeChain("", mc, $3)
	}
	| MULTILINE_ADD_CHAIN MAIN_CHAIN lParen expr RPAREN
	{
		ac := string($1.Literal[len($1.Literal)-1])
		$$ = ast.MakeChain(ac, $2.Literal, $4)
	}

exprList
	: exprList comma listElem
	{
		$$ = append($1, $3)
	}
	| listElem
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

listElem
	: expr
	{
		$$ = $1
	}
	| bareRange
	{
		$$ = $1
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
	}
	| CARET ident COLON expr
	{
		pinned := &ast.PinnedIdent{*$2}
		$$ = &ast.Pair{Key: pinned, Val: $4}
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

methodMapLBrace
	: METHOD_MAP_LBRACE
	{
		$$ = $1
	}
	| METHOD_MAP_LBRACE RET
	{
		$$ = $1
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

lIter
	: LITER
	{
		$$ = $1
	}
	| LITER RET
	{
		$$ = $1
	}

methodLIter
	: METHOD_LITER
	{
		$$ = $1
	}
	| METHOD_LITER RET
	{
		$$ = $1
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

	// NOTE: unary and comparison ops cannot be used for compound assign
	// `&&` and `||` are not methodops but can be used for compound assign
	compoundAssign := `(<<|>>|/&|/\||/\^|\+|\-|\*|\*\*|/|//|%|&&|\|\|)=`

	ident := `[a-zA-Z][a-zA-Z0-9_]*[!?]?`
	// NOTE: comment(, which starts with "#") is included in RET
	// `#[^\n\r]*` is neseccery to lex final line comment (i.e. `#`)
	comment := `#[^\n\r]*`
	retChar := `(\r|\n|\r\n)`
	ret := fmt.Sprintf(`(([ \t]*(%s)?%s)+|%s)`, comment, retChar, comment)
	
	// NOTE: lexer deals with multiline chain
	// (if parser does, shift/reduce conflict occurs)
	keepChainRet := fmt.Sprintf(`([ \t]*(%s)?%s)+[ \t]*\|`, comment, retChar)

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
		t(KWARG_IDENT, fmt.Sprintf(`\\(%s|_+(%s)?)`, ident, ident)),
		t(ARG_IDENT, `\\(0|[1-9][0-9]*)?`),
		t(EXP_FLOAT, `([0-9][0-9_]*[0-9]|[0-9]*)\.([0-9][0-9_]*[0-9]|[0-9]+)[eE]-?[0-9]+`),
		t(FLOAT, `([0-9][0-9_]*[0-9]|[0-9]*)\.([0-9][0-9_]*[0-9]|[0-9]+)`),
		t(HEX_INT, `0[xX]([0-9a-fA-F][0-9a-fA-F_]*[0-9a-fA-F]|[0-9a-fA-F]+)`),
		t(OCT_INT, `0[oO]([0-7][0-7_]*[0-7]|[0-7]+)`),
		t(BIN_INT, `0[bB]([01][01_]*[01]|[01]+)`),
		t(EXP_INT, `([0-9][0-9_]*[0-9]|[0-9]+)[eE]-?[0-9]+`),
		t(INT, `([0-9][0-9_]*[0-9]|[0-9]+)`),
		t(CHAR_STR, `\?(\\[snt\\]|[^\r\n\\])`),
		t(BACKQUOTE_STR, "`[^`]*`"),
		t(HEAD_STR_PIECE, `"[^\"\n\r#]*#\{`),
		t(MID_STR_PIECE, `\}[^\"\n\r#]*#\{`),
		t(TAIL_STR_PIECE, `\}[^\"\n\r#]*"`),
		t(DOUBLEQUOTE_STR, `"[^\"\n\r]*"`),
		// NOTE: lexer deals with multiline chain
		// (if parser does, shift/reduce conflict occurs)
		t(MULTILINE_ADD_CHAIN, fmt.Sprintf(`%s[&~=]`, keepChainRet)),
		t(MULTILINE_MAIN_CHAIN, fmt.Sprintf(`%s[\.@$]`, keepChainRet)),
		// NOTE: comment(, which starts with "#") is included in RET
		// `#[^\n\r]*` is neseccery to lex final line comment (i.e. `#`)
		t(RET, ret),
		t(COMPOUND_ASSIGN, compoundAssign),
		t(SYMBOL, "'"+symbolable),
		t(SPACESHIP, methodOps["spaceship"]),
		t(ASSIGN, `:=`),
		t(RIGHT_ASSIGN, `=>`),
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
		t(DIAMOND, `<>`),
		t(METHOD_LITER, `m<\{`),
		t(LITER, `<\{`),
		t(RITER, `\}>`),
		t(METHOD_MAP_LBRACE, `m%\{`),
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
		t(CARET, `\^`),
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
		t(IF, `if`),
		t(ELSE, `else`),
		t(RETURN, `return`),
		t(YIELD, `yield`),
		t(RAISE, `raise`),
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
