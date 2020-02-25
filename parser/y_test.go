// Source code in this file is inherited and modified from
// "Writing an Interpreter in Go" https://interpreterbook.com/
// MIT License | Copyright (c) 2016-2017 Thorsten Ball
// see https://opensource.org/licenses/MIT

package parser

import (
	"../ast"
	"fmt"
	"strings"
	"testing"
)

// CAUTION: Capitalize test function names!

func TestInfixExpr(t *testing.T) {
	tests := []struct {
		input string
		left  int
		op    string
		right int
	}{
		{`5 + 2`, 5, "+", 2},
		{`5 - 2`, 5, "-", 2},
		{`5 * 2`, 5, "*", 2},
		{`5 / 2`, 5, "/", 2},
		{`5 ** 2`, 5, "**", 2},
		{`5 % 2`, 5, "%", 2},
		{`5 // 2`, 5, "//", 2},
		{`5 <=> 2`, 5, "<=>", 2},
		{`5 == 2`, 5, "==", 2},
		{`5 != 2`, 5, "!=", 2},
		{`5 <= 2`, 5, "<=", 2},
		{`5 >= 2`, 5, ">=", 2},
		{`5 > 2`, 5, ">", 2},
		{`5 < 2`, 5, "<", 2},
		{`5 << 2`, 5, "<<", 2},
		{`5 >> 2`, 5, ">>", 2},
		{`5 /& 2`, 5, "/&", 2},
		{`5 /| 2`, 5, "/|", 2},
		{`5 /^ 2`, 5, "/^", 2},
		{`5 && 2`, 5, "&&", 2},
		{`5 || 2`, 5, "||", 2},
		{`5+2`, 5, "+", 2}, // without space
	}

	for _, tt := range tests {
		program := testParse(t, tt.input)
		expr := testIfExprStmt(t, program)
		testInfixOperator(t, expr, tt.left, tt.op, tt.right)
	}
}

func TestInfixPrecedence(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{`5 + 2`, `(5 + 2)`},
		{`2 + 3 - 4`, `((2 + 3) - 4)`},
		{`2 - 3 + 4`, `((2 - 3) + 4)`},
		{`3 * 4 / 5`, `((3 * 4) / 5)`},
		{`3 / 4 * 5`, `((3 / 4) * 5)`},
		{`3 * 4 + 2`, `((3 * 4) + 2)`},
		{`3 + 4 * 2`, `(3 + (4 * 2))`},
		{`3+4*2`, `(3 + (4 * 2))`}, // without space
		{`3*4**2`, `(3 * (4 ** 2))`},
		{`3-4**2`, `(3 - (4 ** 2))`},
		{`3//4**2`, `(3 // (4 ** 2))`},
		{`3-4//2`, `(3 - (4 // 2))`},
		{`3%4**2`, `(3 % (4 ** 2))`},
		{`3-4%2`, `(3 - (4 % 2))`},
		{`3**4+2`, `((3 ** 4) + 2)`},
		{`3**4/2`, `((3 ** 4) / 2)`},
		{`3*2 == 2-2`, `((3 * 2) == (2 - 2))`},
		{`3*2 != 2-2`, `((3 * 2) != (2 - 2))`},
		{`3*2 > 2-2`, `((3 * 2) > (2 - 2))`},
		{`3*2 < 2-2`, `((3 * 2) < (2 - 2))`},
		{`3*2 <= 2-2`, `((3 * 2) <= (2 - 2))`},
		{`3*2 >= 2-2`, `((3 * 2) >= (2 - 2))`},
		{`3*2 <=> 2-2`, `((3 * 2) <=> (2 - 2))`},
		{`3+2 == 2*2`, `((3 + 2) == (2 * 2))`},
		{`3+2 == 2**2`, `((3 + 2) == (2 ** 2))`},
		{`3 == 4 && 5 != 6`, `((3 == 4) && (5 != 6))`},
		{`3 == 4 || 5 != 6`, `((3 == 4) || (5 != 6))`},
		{`3 && 4 + 2`, `(3 && (4 + 2))`},
		// NOTE: "&&" has higher precedence than "||" (same as other languages)
		// because "&&" is boolean mul, while "||" is boolean add
		{`3 && 4 || 2`, `((3 && 4) || 2)`},
		{`3 || 4 && 2`, `(3 || (4 && 2))`},
		{`3 << 1 == 6`, `((3 << 1) == 6)`},
		{`3 >> 1 == 1`, `((3 >> 1) == 1)`},
		{`3 >> 1 + 1`, `(3 >> (1 + 1))`},
		{`3 /& 1 + 1`, `(3 /& (1 + 1))`},
		{`3 /| 1 + 1`, `(3 /| (1 + 1))`},
		{`3 /^ 1 + 1`, `(3 /^ (1 + 1))`},
		{`3 /& 1 << 1`, `(3 /& (1 << 1))`},
		{`3 /| 1 << 1`, `(3 /| (1 << 1))`},
		{`3 /^ 1 << 1`, `(3 /^ (1 << 1))`},
		{`3 == 3 /& 1`, `(3 == (3 /& 1))`},
		{`3 == 3 /| 1`, `(3 == (3 /| 1))`},
		{`3 == 3 /^ 1`, `(3 == (3 /^ 1))`},
		// NOTE: "/&" has higher precedence than "/|" and "/^"
		// because "/&" is bitwise mul, while "/|", "/^" are bitwise add
		{`3 /& 2 /| 1`, `((3 /& 2) /| 1)`},
		{`3 /| 2 /& 1`, `(3 /| (2 /& 1))`},
		{`3 /& 2 /^ 1`, `((3 /& 2) /^ 1)`},
		{`3 /^ 2 /& 1`, `(3 /^ (2 /& 1))`},
		{`3 /| 2 /^ 1`, `((3 /| 2) /^ 1)`},
		{`3 /^ 2 /| 1`, `((3 /^ 2) /| 1)`},
	}

	for _, tt := range tests {
		program := testParse(t, tt.input)
		expr := testIfExprStmt(t, program)
		infixExpr, ok := expr.(*ast.InfixExpr)

		if !ok {
			t.Fatalf("expr is not ast.InfixExpr. got=%T", expr)
		}

		actual := infixExpr.String()
		if actual != tt.expected {
			t.Errorf("wrong precedence. expected=%s, got=%s",
				tt.expected, actual)
		}
	}
}

func TestPrefixExpression(t *testing.T) {
	tests := []struct {
		input string
		op    string
		right interface{}
	}{
		{`+5`, "+", 5},
		{`-10`, "-", 10},
		{`!1`, "!", 1},
		{`*1`, "*", 1}, // arr expansion
		// NOTE: bitwise not is "/~" (not "~") otherwise
		// conflict occurs in `~.a`
		// ("~" of ".a" or thoughtful scalar chain of prop "a"?)
		{`/~100`, "/~", 100},
	}

	for _, tt := range tests {
		program := testParse(t, tt.input)
		expr := testIfExprStmt(t, program)
		testPrefixOperator(t, expr, tt.op, tt.right)
	}
}

func TestPrefixPrecedence(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{`-3+1`, `((-3) + 1)`},
		{`!3-1`, `((!3) - 1)`},
		{`-3*1`, `((-3) * 1)`},
		{`-3**1`, `((-3) ** 1)`},
		{`--1`, `(-(-1))`},
		{`+-1`, `(+(-1))`},
		{`-1+-1`, `((-1) + (-1))`},
		{`-1---1`, `((-1) - (-(-1)))`},
	}

	for _, tt := range tests {
		program := testParse(t, tt.input)
		expr := testIfExprStmt(t, program)
		actual := expr.String()
		if actual != tt.expected {
			t.Errorf("wrong precedence. expected=%s, got=%s",
				tt.expected, actual)
		}
	}
}

func TestChainPrecedence(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{
			`-3.even?`,
			`(-3).even?()`,
		},
		{
			`2.p ** 3.q`,
			`(2.p() ** 3.q())`,
		},
		{
			`10.find {|a| a+b}.sum`,
			`10.find({|a| (a + b)}).sum()`,
		},
		{
			`10.find(1) {|a| a+b}.sum`,
			`10.find(1, {|a| (a + b)}).sum()`,
		},
		{
			`10.find(1, 2) {|a| a+b}.sum`,
			`10.find(1, 2, {|a| (a + b)}).sum()`,
		},
		{
			`10.find(1, 2) {|a| a+b} {|c| c+d}.sum`,
			`10.find(1, 2, {|a| (a + b)}, {|c| (c + d)}).sum()`,
		},
	}

	for _, tt := range tests {
		program := testParse(t, tt.input)
		expr := testIfExprStmt(t, program)
		actual := expr.String()
		if actual != tt.expected {
			t.Errorf("wrong precedence. expected=%s, got=%s",
				tt.expected, actual)
		}
	}
}

func TestIdentifier(t *testing.T) {
	tests := []string{
		"a",
		"Foo",
		"even?",
		"rand!",
		"_private",
		"__foo",
		"i0",
		"n0u1m2b3e4r",
		"snake_case",
		"snake_case_",
		"CamelCase",
		"pascalCase",
		"CONST",
		"v1234_56ab",
		"_",
		"___",
		//"_123", // not allowed
	}

	for _, tt := range tests {
		program := testParse(t, tt)
		expr := testIfExprStmt(t, program)
		testIdentifier(t, expr, tt)
	}
}

func TestArgIdentifier(t *testing.T) {
	tests := []string{
		`\1`,
		`\0`,
		`\2`,
		`\10`,
		`\999`,
		"\\",
		// `\` is syntax sugar of `\1`
		// NOTE: write `\` as escape otherwise editor syntax-highlighter breaks		{`\1`, `\1`},
	}

	for _, tt := range tests {
		program := testParse(t, tt)
		expr := testIfExprStmt(t, program)
		testArgIdent(t, expr, tt)
	}
}

func TestKwargIdentifier(t *testing.T) {
	tests := []string{
		`\_`,
		`\foo`,
		`\_foo`,
		`\foo_`,
		`\i0`,
		`\__`,
		`\even?`,
		`\fugafuga`,
	}

	for _, tt := range tests {
		program := testParse(t, tt)
		expr := testIfExprStmt(t, program)
		testKwargIdent(t, expr, tt)
	}
}

func TestSymLiteral(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{`'a`, "a"},
		{`'FooBar`, "FooBar"},
		{`'_hidden`, "_hidden"},
		{`'__foo`, "__foo"},
		{`'_f_o_o`, "_f_o_o"},
		{`'i0`, "i0"},
		{`'even?`, "even?"},
		{`'_`, "_"},
		{`'__`, "__"},
		{`'+`, "+"},
		{`'-`, "-"},
		{`'*`, "*"},
		{`'/`, "/"},
		{`'//`, "//"},
		{`'%`, "%"},
		{`'**`, "**"},
		{`'<=>`, "<=>"},
		{`'==`, "=="},
		{`'!=`, "!="},
		{`'<=`, "<="},
		{`'>=`, ">="},
		{`'<`, "<"},
		{`'>`, ">"},
		{`'+%`, "+%"},
		{`'-%`, "-%"},
		{`'!`, "!"},
		{`'/~`, "/~"},
		{`'/&`, "/&"},
		{`'/|`, "/|"},
		{`'/^`, "/^"},
		{`'<<`, "<<"},
		{`'>>`, ">>"},
	}

	for _, tt := range tests {
		program := testParse(t, tt.input)
		expr := testIfExprStmt(t, program)
		sym, ok := expr.(*ast.SymLiteral)
		if !ok {
			t.Fatalf("expr is not *ast.SymLiteral.got=%T", expr)
		}

		testSymbol(t, sym, tt.expected)
	}
}

func TestCharStrLiteral(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{`?a`, "a"},
		{`?b`, "b"},
		{`?I`, "I"},
		{`?1`, "1"},
		{`?_`, "_"},
		// espaced char is not evaluated in parser
		{`?\n`, `\n`},
		{`?\t`, `\t`},
		{`?\s`, `\s`},
		// `?\` is not allowed (for compatibility with other escaped form)
		{`?\\`, `\\`},
		{`? `, " "},
		{`?!`, "!"},
		{`??`, "?"},
		{`?"`, `"`},
		{`?'`, "'"},
		{`?.`, "."},
		{`?@`, "@"},
		{`?$`, "$"},
		{`?#`, "#"},
	}

	for _, tt := range tests {
		program := testParse(t, tt.input)
		expr := testIfExprStmt(t, program)
		str, ok := expr.(*ast.StrLiteral)
		if !ok {
			t.Fatalf("expr is not *ast.StrLiteral.got=%T", expr)
		}

		// IsRaw should be false (to evaluate escapes)
		testStr(t, str, tt.expected, false)
	}
}

func TestBackQuoteStrLiteral(t *testing.T) {
	tests := []string{
		// NOTE: each input is wrapped by ``
		// (because `` cannot be written in ``)
		``,
		`foo`,
		`12345`,
		`\n`,
		`Hello, world!`,
		`#comment?`,
		`break
		
		line
		s`,
		`"a"`,
		`#{1 + 1}`,
		`#{}`,
		`.hoge`,
		`1 + 1`,
		`_`,
	}

	for _, tt := range tests {
		// NOTE: each input is wrapped by ``
		// (because `` cannot be written in ``)
		input := "`" + tt + "`"
		program := testParse(t, input)
		expr := testIfExprStmt(t, program)
		str, ok := expr.(*ast.StrLiteral)
		if !ok {
			t.Fatalf("expr is not *ast.StrLiteral.got=%T", expr)
		}

		// IsRaw should be true
		testStr(t, str, tt, true)
	}
}

func TestDoubleQuoteStrLiteral(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{`""`, ""},
		{`"foo"`, "foo"},
		{`"12345"`, "12345"},
		{`"Hello, world!"`, "Hello, world!"},
		{`"comment?"`, "comment?"},
		// NOTE: escape is not evaluated in parser
		{`"\n"`, `\n`},
		{`"` + "``" + `"`, "``"},
		{`".hoge"`, ".hoge"},
		{`"1 + 1"`, "1 + 1"},
		{`"_"`, "_"},
	}

	for _, tt := range tests {
		program := testParse(t, tt.input)
		expr := testIfExprStmt(t, program)
		str, ok := expr.(*ast.StrLiteral)
		if !ok {
			t.Fatalf("expr is not *ast.StrLiteral.got=%T", expr)
		}

		// IsRaw should be false
		testStr(t, str, tt.expected, false)
	}
}

func TestEmbeddedStr(t *testing.T) {
	input := `"abc#{1}def#{1+1}ghi#{foo.bar}jkl"`

	program := testParse(t, input)
	expr := testIfExprStmt(t, program)
	embeddedStr, ok := expr.(*ast.EmbeddedStr)
	if !ok {
		t.Fatalf("expr is not *ast.EmbeddedStr. got=%T", expr)
	}

	if embeddedStr.Latter != "jkl" {
		t.Errorf("str is not `jkl`. got=`%s`", embeddedStr.Latter)
	}

	former1 := embeddedStr.Former

	if former1 == nil {
		t.Fatalf("former1 must not be nil.")
	}

	ce1, ok := former1.Expr.(*ast.PropCallExpr)

	if !ok {
		t.Fatalf("former1 is not *ast.CallExpr. got=%T", former1)
	}

	testIdentifier(t, ce1.Receiver, "foo")

	if ce1.ChainToken() != "." {
		t.Errorf("chain is not `.`. got=%s", ce1.ChainToken())
	}

	testIdentifier(t, ce1.Prop, "bar")

	if former1.Str != "ghi" {
		t.Errorf("str is not `ghi`. got=`%s`", former1.Str)
	}

	former2 := former1.Former

	if former2 == nil {
		t.Fatalf("former2 must not be nil.")
	}

	infix2, ok := former2.Expr.(*ast.InfixExpr)

	if !ok {
		t.Fatalf("former2 is not *ast.InfixExpr. got=%T", former2)
	}

	testInfixOperator(t, infix2, 1, "+", 1)

	if former2.Str != "def" {
		t.Errorf("str is not `def`. got=`%s`", former2.Str)
	}

	former3 := former2.Former

	if former3 == nil {
		t.Fatalf("former3 must not be nil.")
	}

	testLiteralExpr(t, former3.Expr, 1)

	if former3.Str != "abc" {
		t.Errorf("str is not `abc`. got=`%s`", former3.Str)
	}

	former4 := former3.Former

	if former4 != nil {
		t.Fatalf("former4 must be nil")
	}

	expectedStr := `"abc#{ 1 }def#{ (1 + 1) }ghi#{ foo.bar() }jkl"`

	if embeddedStr.String() != expectedStr {
		t.Errorf("wrong str output. expected=`\n%s\n`. got=`\n%s\n`",
			expectedStr, embeddedStr.String())
	}
}

func TestObjLiteral(t *testing.T) {
	tests := []struct {
		input    string
		keys     []string
		vals     []interface{}
		embedded []string
	}{
		{
			`{}`,
			[]string{},
			[]interface{}{},
			[]string{},
		},
		{
			`{'a: 2}`,
			[]string{"a"},
			[]interface{}{2},
			[]string{},
		},
		{
			`{**a}`,
			[]string{},
			[]interface{}{},
			[]string{"a"},
		},
		{
			`{'a: 2, 'b: 3, 'c: 4}`,
			[]string{"a", "b", "c"},
			[]interface{}{2, 3, 4},
			[]string{},
		},
		{
			`{'b: 3, 'd: 5, **a, **c}`,
			[]string{"b", "d"},
			[]interface{}{3, 5},
			[]string{"a", "c"},
		},
	}

	for _, tt := range tests {
		program := testParse(t, tt.input)
		expr := testIfExprStmt(t, program)
		obj, ok := expr.(*ast.ObjLiteral)
		if !ok {
			t.Fatalf("expr is not *ast.ObjLiteral.got=%T", expr)
		}

		if len(tt.keys) != len(obj.Pairs) {
			t.Fatalf("wrong number of elements. expected=%d, got=%d.",
				len(tt.keys), len(obj.Pairs))
		}

		for i, pair := range obj.Pairs {
			key, ok := pair.Key.(*ast.SymLiteral)
			if !ok {
				t.Errorf("obj.Pairs[%d].Key is not *ast.SymLiteral. got=%T",
					i, pair.Key)
			}

			testSymbol(t, key, tt.keys[i])

			val, ok := pair.Val.(ast.Expr)
			if !ok {
				t.Errorf("obj.Pairs[%d].Val is not ast.Expr. got=%T",
					i, pair.Val)
			}

			testLiteralExpr(t, val, tt.vals[i])
		}

		if len(tt.embedded) != len(obj.EmbeddedExprs) {
			t.Fatalf("wrong number of embedded. expected=%d, got=%d.",
				len(tt.embedded), len(obj.EmbeddedExprs))
		}

		for i, expr := range obj.EmbeddedExprs {
			testIdentifier(t, expr, tt.embedded[i])
		}
	}
}

func TestObjString(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{
			`{}`,
			`{}`,
		},
		{
			`{'a:1,'b:2}`,
			`{'a: 1, 'b: 2}`,
		},
		{
			`{'a:1,b:2}`,
			`{'a: 1, b: 2}`,
		},
		{
			`{'a: 1, 'b: 2,}`,
			`{'a: 1, 'b: 2}`,
		},
		{
			`{
				'a: 1,
				'b: 2,
			}`,
			`{'a: 1, 'b: 2}`,
		},
		{
			`{'a:1,**foo,**bar}`,
			`{'a: 1, **foo, **bar}`,
		},
		{
			`{'a: 1, **foo,}`,
			`{'a: 1, **foo}`,
		},
		{
			`{
				'a: 1,
				**foo,
				**bar,
			}`,
			`{'a: 1, **foo, **bar}`,
		},
	}

	for _, tt := range tests {
		program := testParse(t, tt.input)
		expr := testIfExprStmt(t, program)
		obj, ok := expr.(*ast.ObjLiteral)
		if !ok {
			t.Fatalf("expr is not *ast.ObjLiteral.got=%T", expr)
		}

		if obj.String() != tt.expected {
			t.Errorf("wrong string. expected=`\n%s\n`, got=`\n%s\n`",
				tt.expected, obj.String())
		}
	}
}

func TestObjBreaklines(t *testing.T) {
	tests := []struct {
		input    string
		keys     []string
		vals     []interface{}
		embedded []string
	}{
		{
			`{'a: 1,
			'b: 2}`,
			[]string{"a", "b"},
			[]interface{}{1, 2},
			[]string{},
		},
		{
			`{
			'a: 1, 'b: 2}`,
			[]string{"a", "b"},
			[]interface{}{1, 2},
			[]string{},
		},
		{
			`{
				'a: 1, 'b: 2
			}`,
			[]string{"a", "b"},
			[]interface{}{1, 2},
			[]string{},
		},
		{
			`{
				'a: 1,
				'b: 2
			}`,
			[]string{"a", "b"},
			[]interface{}{1, 2},
			[]string{},
		},
		{
			`{
				'a: 1,
				'b: 2,
			}`,
			[]string{"a", "b"},
			[]interface{}{1, 2},
			[]string{},
		},
		{
			`{'a: 1,
			'b: 2,}`,
			[]string{"a", "b"},
			[]interface{}{1, 2},
			[]string{},
		},
		{
			`{
				'a: 1,'b: 2,}`,
			[]string{"a", "b"},
			[]interface{}{1, 2},
			[]string{},
		},
		{
			`{'a: 1,'b: 2,}`,
			[]string{"a", "b"},
			[]interface{}{1, 2},
			[]string{},
		},
		{
			`{**a,
			**b}`,
			[]string{},
			[]interface{}{},
			[]string{"a", "b"},
		},
		{
			`{
				**a, **b}`,
			[]string{},
			[]interface{}{},
			[]string{"a", "b"},
		},
		{
			`{
				**a,**b
			}`,
			[]string{},
			[]interface{}{},
			[]string{"a", "b"},
		},
		{
			`{
				**a,
				**b
			}`,
			[]string{},
			[]interface{}{},
			[]string{"a", "b"},
		},
		{
			`{
				**a,
				**b,
			}`,
			[]string{},
			[]interface{}{},
			[]string{"a", "b"},
		},
		{
			`{**a,
			**b,
			}`,
			[]string{},
			[]interface{}{},
			[]string{"a", "b"},
		},
		{
			`{
			**a,**b,}`,
			[]string{},
			[]interface{}{},
			[]string{"a", "b"},
		},
		{
			`{**a,**b,}`,
			[]string{},
			[]interface{}{},
			[]string{"a", "b"},
		},
		{
			`{'foo: 1,
			**a, **b}`,
			[]string{"foo"},
			[]interface{}{1},
			[]string{"a", "b"},
		},
		{
			`{'foo: 1,
			**a, **b
			}`,
			[]string{"foo"},
			[]interface{}{1},
			[]string{"a", "b"},
		},
		{
			`{
			**a,
			}`,
			[]string{},
			[]interface{}{},
			[]string{"a"},
		},
		{
			`{
			'a: 1,
			}`,
			[]string{"a"},
			[]interface{}{1},
			[]string{},
		},
	}

	for _, tt := range tests {
		program := testParse(t, tt.input)
		expr := testIfExprStmt(t, program)
		obj, ok := expr.(*ast.ObjLiteral)
		if !ok {
			t.Fatalf("expr is not *ast.ObjLiteral.got=%T", expr)
		}

		if len(tt.keys) != len(obj.Pairs) {
			t.Fatalf("wrong number of elements. expected=%d, got=%d.",
				len(tt.keys), len(obj.Pairs))
		}

		for i, pair := range obj.Pairs {
			key, ok := pair.Key.(*ast.SymLiteral)
			if !ok {
				t.Errorf("obj.Pairs[%d].Key is not *ast.SymLiteral. got=%T",
					i, pair.Key)
			}

			testSymbol(t, key, tt.keys[i])

			val, ok := pair.Val.(ast.Expr)
			if !ok {
				t.Errorf("obj.Pairs[%d].Val is not ast.Expr. got=%T",
					i, pair.Val)
			}

			testLiteralExpr(t, val, tt.vals[i])
		}

		if len(tt.embedded) != len(obj.EmbeddedExprs) {
			t.Fatalf("wrong number of embedded. expected=%d, got=%d.",
				len(tt.embedded), len(obj.EmbeddedExprs))
		}

		for i, expr := range obj.EmbeddedExprs {
			testIdentifier(t, expr, tt.embedded[i])
		}
	}
}

func TestMapLiteral(t *testing.T) {
	tests := []struct {
		input    string
		keys     []interface{}
		vals     []interface{}
		embedded []string
	}{
		{
			`%{}`,
			[]interface{}{},
			[]interface{}{},
			[]string{},
		},
		{
			`%{'a: 2}`,
			[]interface{}{"a"},
			[]interface{}{2},
			[]string{},
		},
		{
			`%{**a}`,
			[]interface{}{},
			[]interface{}{},
			[]string{"a"},
		},
		{
			`%{'a: 2, 'b: 3, 'c: 4}`,
			[]interface{}{"a", "b", "c"},
			[]interface{}{2, 3, 4},
			[]string{},
		},
		{
			`%{'b: 3, 'd: 5, **a, **c}`,
			[]interface{}{"b", "d"},
			[]interface{}{3, 5},
			[]string{"a", "c"},
		},
		{
			`%{1: 2}`,
			[]interface{}{1},
			[]interface{}{2},
			[]string{},
		},
		{
			`%{1: 2, 3: 4}`,
			[]interface{}{1, 3},
			[]interface{}{2, 4},
			[]string{},
		},
		{
			`%{1: 2, 'a: 3}`,
			[]interface{}{1, "a"},
			[]interface{}{2, 3},
			[]string{},
		},
		{
			`%{1: 2, 'a: 3, **foo}`,
			[]interface{}{1, "a"},
			[]interface{}{2, 3},
			[]string{"foo"},
		},
	}

	for _, tt := range tests {
		program := testParse(t, tt.input)
		expr := testIfExprStmt(t, program)
		obj, ok := expr.(*ast.MapLiteral)
		if !ok {
			t.Fatalf("expr is not *ast.MapLiteral.got=%T", expr)
		}

		if len(tt.keys) != len(obj.Pairs) {
			t.Fatalf("wrong number of elements. expected=%d, got=%d.",
				len(tt.keys), len(obj.Pairs))
		}

		for i, pair := range obj.Pairs {
			switch expectedKey := tt.keys[i].(type) {
			case string:
				key, ok := pair.Key.(*ast.SymLiteral)
				if !ok {
					t.Errorf("obj.Pairs[%d].Key is not *ast.SymLiteral. got=%T",
						i, pair.Key)
				}

				testSymbol(t, key, expectedKey)
			default:
				testLiteralExpr(t, pair.Key, expectedKey)
			}

			val, ok := pair.Val.(ast.Expr)
			if !ok {
				t.Errorf("obj.Pairs[%d].Val is not ast.Expr. got=%T",
					i, pair.Val)
			}

			testLiteralExpr(t, val, tt.vals[i])
		}

		if len(tt.embedded) != len(obj.EmbeddedExprs) {
			t.Fatalf("wrong number of embedded. expected=%d, got=%d.",
				len(tt.embedded), len(obj.EmbeddedExprs))
		}

		for i, expr := range obj.EmbeddedExprs {
			testIdentifier(t, expr, tt.embedded[i])
		}
	}
}

func TestMapString(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{
			`%{}`,
			`%{}`,
		},
		{
			`%{'a:1,'b:2}`,
			`%{'a: 1, 'b: 2}`,
		},
		{
			`%{'a:1,b:2}`,
			`%{'a: 1, b: 2}`,
		},
		{
			`%{'a: 1, 'b: 2,}`,
			`%{'a: 1, 'b: 2}`,
		},
		{
			`%{
				'a: 1,
				'b: 2,
			}`,
			`%{'a: 1, 'b: 2}`,
		},
		{
			`%{'a:1,**foo,**bar}`,
			`%{'a: 1, **foo, **bar}`,
		},
		{
			`%{'a: 1, **foo,}`,
			`%{'a: 1, **foo}`,
		},
		{
			`%{
				'a: 1,
				**foo,
				**bar,
			}`,
			`%{'a: 1, **foo, **bar}`,
		},
		{
			`%{1:2,3:4}`,
			`%{1: 2, 3: 4}`,
		},
	}

	for _, tt := range tests {
		program := testParse(t, tt.input)
		expr := testIfExprStmt(t, program)
		obj, ok := expr.(*ast.MapLiteral)
		if !ok {
			t.Fatalf("expr is not *ast.MapLiteral.got=%T", expr)
		}

		if obj.String() != tt.expected {
			t.Errorf("wrong string. expected=`\n%s\n`, got=`\n%s\n`",
				tt.expected, obj.String())
		}
	}
}

func TestMapBreaklines(t *testing.T) {
	tests := []struct {
		input    string
		keys     []string
		vals     []interface{}
		embedded []string
	}{
		{
			`%{'a: 1,
			'b: 2}`,
			[]string{"a", "b"},
			[]interface{}{1, 2},
			[]string{},
		},
		{
			`%{
			'a: 1, 'b: 2}`,
			[]string{"a", "b"},
			[]interface{}{1, 2},
			[]string{},
		},
		{
			`%{
				'a: 1, 'b: 2
			}`,
			[]string{"a", "b"},
			[]interface{}{1, 2},
			[]string{},
		},
		{
			`%{
				'a: 1,
				'b: 2
			}`,
			[]string{"a", "b"},
			[]interface{}{1, 2},
			[]string{},
		},
		{
			`%{
				'a: 1,
				'b: 2,
			}`,
			[]string{"a", "b"},
			[]interface{}{1, 2},
			[]string{},
		},
		{
			`%{'a: 1,
			'b: 2,}`,
			[]string{"a", "b"},
			[]interface{}{1, 2},
			[]string{},
		},
		{
			`%{
				'a: 1,'b: 2,}`,
			[]string{"a", "b"},
			[]interface{}{1, 2},
			[]string{},
		},
		{
			`%{'a: 1,'b: 2,}`,
			[]string{"a", "b"},
			[]interface{}{1, 2},
			[]string{},
		},
		{
			`%{**a,
			**b}`,
			[]string{},
			[]interface{}{},
			[]string{"a", "b"},
		},
		{
			`%{
				**a, **b}`,
			[]string{},
			[]interface{}{},
			[]string{"a", "b"},
		},
		{
			`%{
				**a,**b
			}`,
			[]string{},
			[]interface{}{},
			[]string{"a", "b"},
		},
		{
			`%{
				**a,
				**b
			}`,
			[]string{},
			[]interface{}{},
			[]string{"a", "b"},
		},
		{
			`%{
				**a,
				**b,
			}`,
			[]string{},
			[]interface{}{},
			[]string{"a", "b"},
		},
		{
			`%{**a,
			**b,
			}`,
			[]string{},
			[]interface{}{},
			[]string{"a", "b"},
		},
		{
			`%{
			**a,**b,}`,
			[]string{},
			[]interface{}{},
			[]string{"a", "b"},
		},
		{
			`%{**a,**b,}`,
			[]string{},
			[]interface{}{},
			[]string{"a", "b"},
		},
		{
			`%{'foo: 1,
			**a, **b}`,
			[]string{"foo"},
			[]interface{}{1},
			[]string{"a", "b"},
		},
		{
			`%{'foo: 1,
			**a, **b
			}`,
			[]string{"foo"},
			[]interface{}{1},
			[]string{"a", "b"},
		},
		{
			`%{
			**a,
			}`,
			[]string{},
			[]interface{}{},
			[]string{"a"},
		},
		{
			`%{
			'a: 1,
			}`,
			[]string{"a"},
			[]interface{}{1},
			[]string{},
		},
	}

	for _, tt := range tests {
		program := testParse(t, tt.input)
		expr := testIfExprStmt(t, program)
		obj, ok := expr.(*ast.MapLiteral)
		if !ok {
			t.Fatalf("expr is not *ast.MapLiteral.got=%T", expr)
		}

		if len(tt.keys) != len(obj.Pairs) {
			t.Fatalf("wrong number of elements. expected=%d, got=%d.",
				len(tt.keys), len(obj.Pairs))
		}

		for i, pair := range obj.Pairs {
			key, ok := pair.Key.(*ast.SymLiteral)
			if !ok {
				t.Errorf("obj.Pairs[%d].Key is not *ast.SymLiteral. got=%T",
					i, pair.Key)
			}

			testSymbol(t, key, tt.keys[i])

			val, ok := pair.Val.(ast.Expr)
			if !ok {
				t.Errorf("obj.Pairs[%d].Val is not ast.Expr. got=%T",
					i, pair.Val)
			}

			testLiteralExpr(t, val, tt.vals[i])
		}

		if len(tt.embedded) != len(obj.EmbeddedExprs) {
			t.Fatalf("wrong number of embedded. expected=%d, got=%d.",
				len(tt.embedded), len(obj.EmbeddedExprs))
		}

		for i, expr := range obj.EmbeddedExprs {
			testIdentifier(t, expr, tt.embedded[i])
		}
	}
}

func TestArrLiteral(t *testing.T) {
	tests := []struct {
		input string
		vals  []interface{}
	}{
		{`[]`, []interface{}{}},
		{`[1]`, []interface{}{1}},
		{`[1, 2, 3]`, []interface{}{1, 2, 3}},
	}

	for _, tt := range tests {
		program := testParse(t, tt.input)
		expr := testIfExprStmt(t, program)

		a, ok := expr.(*ast.ArrLiteral)
		if !ok {
			t.Fatalf("f is not *ast.ArrLiteral. got=%T", expr)
		}

		if len(a.Elems) != len(tt.vals) {
			t.Fatalf("number of elements is not %d. got=%d",
				len(tt.vals), len(a.Elems))
		}

		for i, elem := range a.Elems {
			testLiteralExpr(t, elem, tt.vals[i])
		}
	}
}

func TestArrString(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{
			`[1,2,3]`,
			`[1, 2, 3]`,
		},
		{
			`[]`,
			`[]`,
		},
		{
			`[1,]`,
			`[1]`,
		},
		{
			`[1,2,]`,
			`[1, 2]`,
		},
		{
			`[
				1,
				2,
			]`,
			`[1, 2]`,
		},
		{
			`[1, *a,]`,
			`[1, (*a)]`,
		},
	}

	for _, tt := range tests {
		program := testParse(t, tt.input)
		expr := testIfExprStmt(t, program)

		a, ok := expr.(*ast.ArrLiteral)
		if !ok {
			t.Fatalf("f is not *ast.ArrLiteral. got=%T", expr)
		}

		if a.String() != tt.expected {
			t.Errorf("wrong string. expected=`\n%s\n`, got=`\n%s\n`",
				tt.expected, a.String())
		}
	}
}

func TestArrBreakLines(t *testing.T) {
	tests := []struct {
		input string
		vals  []interface{}
	}{
		{
			`[1,]`,
			[]interface{}{1},
		},
		{
			`[
			1,
			]`,
			[]interface{}{1},
		},
		{
			`[1,
			2]`,
			[]interface{}{1, 2},
		},
		{
			`[1,
			2
			]`,
			[]interface{}{1, 2},
		},
		{
			`[
			 1, 2]`,
			[]interface{}{1, 2},
		},
		{
			`[
				1, 2
			]`,
			[]interface{}{1, 2},
		},
		{
			`[
				1,
				2
			]`,
			[]interface{}{1, 2},
		},
		{
			`[
				1,
				2,
			]`,
			[]interface{}{1, 2},
		},
		{
			`[1,
			2,
			]`,
			[]interface{}{1, 2},
		},
		{
			`[
			1,2,]`,
			[]interface{}{1, 2},
		},
		{
			`[1,2,]`,
			[]interface{}{1, 2},
		},
	}

	for _, tt := range tests {
		program := testParse(t, tt.input)
		expr := testIfExprStmt(t, program)

		a, ok := expr.(*ast.ArrLiteral)
		if !ok {
			t.Fatalf("f is not *ast.ArrLiteral. got=%T", expr)
		}

		if len(a.Elems) != len(tt.vals) {
			t.Fatalf("number of elements is not %d. got=%d",
				len(tt.vals), len(a.Elems))
		}

		for i, elem := range a.Elems {
			testLiteralExpr(t, elem, tt.vals[i])
		}
	}
}

func TestNestedArrLiteral(t *testing.T) {
	input := `[1, *foo, [2, 3]]`

	program := testParse(t, input)
	expr := testIfExprStmt(t, program)

	a, ok := expr.(*ast.ArrLiteral)
	if !ok {
		t.Fatalf("a is not *ast.ArrLiteral. got=%T", expr)
	}

	if len(a.Elems) != 3 {
		t.Fatalf("number of elements is not 6. got=%d", len(a.Elems))
	}

	testLiteralExpr(t, a.Elems[0], 1)

	elem1, ok := a.Elems[1].(*ast.PrefixExpr)

	if !ok {
		t.Fatalf("a.Elems[1] is not *ast.PrefixExpr. got=%T", a.Elems[1])
	}

	if elem1.Operator != "*" {
		t.Errorf("a.Elems[1] does not have `*`. got=%s", elem1.Operator)
	}

	testIdentifier(t, elem1.Right, "foo")

	elem2, ok := a.Elems[2].(*ast.ArrLiteral)
	if !ok {
		t.Fatalf("a.Elems[2] is not *ast.ArrLiteral. got=%T", a.Elems[2])
	}

	testLiteralExpr(t, elem2.Elems[0], 2)
	testLiteralExpr(t, elem2.Elems[1], 3)

}

func TestCallArgBreakLines(t *testing.T) {
	tests := []struct {
		input  string
		args   []int
		kwargs map[string]int
	}{
		{
			`a.b(1, b: 2, 3)`,
			[]int{1, 3},
			map[string]int{"b": 2},
		},
		{
			`a.b(
			  1, b: 2, 3)`,
			[]int{1, 3},
			map[string]int{"b": 2},
		},
		{
			`a.b(1,
				b: 2, 3)`,
			[]int{1, 3},
			map[string]int{"b": 2},
		},
		{
			`a.b(1, b: 2,
				3)`,
			[]int{1, 3},
			map[string]int{"b": 2},
		},
		{
			`a.b(
				1,
				b: 2,
				3
			)`,
			[]int{1, 3},
			map[string]int{"b": 2},
		},
	}

	for _, tt := range tests {
		program := testParse(t, tt.input)
		expr := testIfExprStmt(t, program)

		f, ok := expr.(*ast.PropCallExpr)
		if !ok {
			t.Fatalf("f is not *ast.PropCallExpr. got=%T", expr)
		}

		if len(f.Args) != len(tt.args) {
			t.Fatalf("arity of args is not %d. got=%d",
				len(tt.args), len(f.Args))
		}

		if len(f.Kwargs) != len(tt.kwargs) {
			t.Fatalf("arity of kwargs is not %d. got=%d",
				len(tt.kwargs), len(f.Kwargs))
		}

		for i, expArg := range tt.args {
			testLiteralExpr(t, f.Args[i], expArg)
		}

		for ident, val := range f.Kwargs {
			name := ident.Token
			exp, ok := tt.kwargs[name]
			if ok {
				testLiteralExpr(t, val, exp)
			} else {
				t.Errorf("unexpected kwarg %s found.", name)
			}
		}
	}
}

func TestFuncArgBreakLines(t *testing.T) {
	tests := []struct {
		input  string
		args   []string
		kwargs map[string]int64
	}{
		{
			`{|a, b: 2, c|}`,
			[]string{"a", "c"},
			map[string]int64{"b": 2},
		},
		{
			`{|a,
			   b: 2, c|}`,
			[]string{"a", "c"},
			map[string]int64{"b": 2},
		},
		{
			`{
				|a, b: 2,
				 c|}`,
			[]string{"a", "c"},
			map[string]int64{"b": 2},
		},
	}

	for _, tt := range tests {
		program := testParse(t, tt.input)
		expr := testIfExprStmt(t, program)

		f, ok := expr.(*ast.FuncLiteral)
		if !ok {
			t.Fatalf("f is not *ast.FuncLiteral. got=%T", expr)
		}

		if len(f.Args) != len(tt.args) {
			t.Fatalf("arity of args is not %d. got=%d",
				len(tt.args), len(f.Args))
		}

		if len(f.Kwargs) != len(tt.kwargs) {
			t.Fatalf("arity of kwargs is not %d. got=%d",
				len(tt.kwargs), len(f.Kwargs))
		}

		for i, expArg := range tt.args {
			testIdentifier(t, f.Args[i], expArg)
		}

		for ident, val := range f.Kwargs {
			name := ident.Token
			exp, ok := tt.kwargs[name]
			if ok {
				testLiteralExpr(t, val, exp)
			} else {
				t.Errorf("unexpected kwarg %s found.", name)
			}
		}

	}
}

func TestFuncBodyBreakLines(t *testing.T) {
	tests := []struct {
		input  string
		bodies []int64
	}{
		{
			`{|a| 1;2}`,
			[]int64{1, 2},
		},
		{
			`{|a| 1
			2}`,
			[]int64{1, 2},
		},
		{
			`{|a|
			   1
			   2
			 }`,
			[]int64{1, 2},
		},
		{
			`{1;2}`,
			[]int64{1, 2},
		},
		{
			`{1
			2}`,
			[]int64{1, 2},
		},
		{
			`{
			   1
			   2
			 }`,
			[]int64{1, 2},
		},
	}

	for _, tt := range tests {
		program := testParse(t, tt.input)
		expr := testIfExprStmt(t, program)

		f, ok := expr.(*ast.FuncLiteral)
		if !ok {
			t.Fatalf("f is not *ast.FuncLiteral. got=%T", expr)
		}

		if len(f.Body) != len(tt.bodies) {
			t.Fatalf("body length is not 1. got=%d", len(f.Body))
		}

		for i, b := range f.Body {
			es, ok := b.(*ast.ExprStmt)
			if !ok {
				t.Fatalf("f.Body[%d] is not *ast.ExprStmt. got=%T",
					i, f.Body[0])
			}
			lit, ok := es.Expr.(*ast.IntLiteral)
			if !ok {
				t.Fatalf("f.Body[%d] does not have *ast.IntLiteral. got=%T",
					i, es.Expr)
			}

			testIntLiteral(t, lit, tt.bodies[i])
		}
	}
}

func TestPropCall(t *testing.T) {
	tests := []struct {
		input        string
		receiver     interface{}
		chainContext string
		chainArg     interface{}
		propName     string
	}{
		{`5.times`, 5, ".", nil, "times"},
		{`10@puts`, 10, "@", nil, "puts"},
		{`5@(10)puts`, 5, "@", 10, "puts"},
		{`10$add`, 10, "$", nil, "add"},
		{`5$(0)add`, 5, "$", 0, "add"},
		{`10$(0)+`, 10, "$", 0, "+"},
		{`5.foo`, 5, ".", nil, "foo"},
		{`5@foo`, 5, "@", nil, "foo"},
		{`5$foo`, 5, "$", nil, "foo"},
		{`5&.foo`, 5, "&.", nil, "foo"},
		{`5~.foo`, 5, "~.", nil, "foo"},
		{`5=.foo`, 5, "=.", nil, "foo"},
		{`5&@foo`, 5, "&@", nil, "foo"},
		{`5~@foo`, 5, "~@", nil, "foo"},
		{`5=@foo`, 5, "=@", nil, "foo"},
		{`5&$foo`, 5, "&$", nil, "foo"},
		{`5~$foo`, 5, "~$", nil, "foo"},
		{`5=$foo`, 5, "=$", nil, "foo"},
	}

	for _, tt := range tests {
		program := testParse(t, tt.input)
		expr := testIfExprStmt(t, program)

		callExpr, ok := expr.(*ast.PropCallExpr)
		if !ok {
			t.Fatalf("expr is not *ast.PropCallExpr. got=%T", expr)
		}

		testLiteralExpr(t, callExpr.Receiver, tt.receiver)
		testChainContext(t, callExpr, tt.chainContext, tt.chainArg)
		testIdentifier(t, callExpr.Prop, tt.propName)

	}
}

func TestRecursiveChain(t *testing.T) {
	input := `foo.bar()@(10)hoge$piyo$(1)fuga().puts`
	expectedVals := []struct {
		prop         string
		chainContext string
		chainArg     interface{}
	}{
		{"bar", ".", nil},
		{"hoge", "@", 10},
		{"piyo", "$", nil},
		{"fuga", "$", 1},
		{"puts", ".", nil},
	}

	program := testParse(t, input)
	expr := testIfExprStmt(t, program)

	callExpr, ok := expr.(*ast.PropCallExpr)
	if !ok {
		t.Fatalf("expr is not *ast.PropCallExpr. got=%T", expr)
	}

	pc := callExpr

	for i := len(expectedVals) - 1; i > 0; i-- {
		exp := expectedVals[i]
		testChainContext(t, pc, exp.chainContext, exp.chainArg)
		testIdentifier(t, pc.Prop, exp.prop)
		recv, ok := pc.Receiver.(*ast.PropCallExpr)
		if ok {
			pc = recv
		} else {
			t.Errorf("receiver is not *ast.PropCallExpr. got=%T(i=%d)",
				pc.Receiver, i)
		}
	}
}

func TestOpMethods(t *testing.T) {
	tests := []struct {
		input string
		op    string
	}{
		{`10.+(1)`, `+`},
		{`10.-(1)`, `-`},
		{`10.*(1)`, `*`},
		{`10./(1)`, `/`},
		{`10.//(1)`, `//`},
		{`10.%(1)`, `%`},
		{`10.**(1)`, `**`},
		{`10.<=>(1)`, `<=>`},
		{`10.==(1)`, `==`},
		{`10.!=(1)`, `!=`},
		{`10.<=(1)`, `<=`},
		{`10.>=(1)`, `>=`},
		{`10.<(1)`, `<`},
		{`10.>(1)`, `>`},
		{`10.>>(1)`, `>>`},
		{`10.<<(1)`, `<<`},
		{`10./&(1)`, `/&`},
		{`10./|(1)`, `/|`},
		{`10./^(1)`, `/^`},
		{`10./~(1)`, `/~`},
		{`10.!(1)`, `!`},
		{`10.+%(1)`, `+%`},
		{`10.-%(1)`, `-%`},
	}

	for _, tt := range tests {
		program := testParse(t, tt.input)
		expr := testIfExprStmt(t, program)

		callExpr, ok := expr.(*ast.PropCallExpr)
		if !ok {
			t.Fatalf("expr is not *ast.PropCallExpr. got=%T", expr)
		}

		testLiteralExpr(t, callExpr.Receiver, 10)
		testChainContext(t, callExpr, ".", nil)
		testIdentifier(t, callExpr.Prop, tt.op)

		if len(callExpr.Args) != 1 {
			t.Fatalf("arity must be 1. got=%d", len(callExpr.Args))
		}

		testLiteralExpr(t, callExpr.Args[0], 1)
	}
}

func TestArgOrders(t *testing.T) {
	tests := []struct {
		input   string
		args    []int
		kwargs  map[string]int
		printed string
	}{
		{
			`5.a(1)`,
			[]int{1},
			map[string]int{},
			`5.a(1)`,
		},
		{
			`5.a(1, 2)`,
			[]int{1, 2},
			map[string]int{},
			`5.a(1, 2)`,
		},
		{
			`5.a(1, foo:2)`, // without space
			[]int{1},
			map[string]int{"foo": 2},
			`5.a(1, foo: 2)`,
		},
		{
			`5.a(foo: 3, 1, 2)`,
			[]int{1, 2},
			map[string]int{"foo": 3},
			`5.a(1, 2, foo: 3)`,
		},
		{
			`5.a(1, foo: 3, 2)`,
			[]int{1, 2},
			map[string]int{"foo": 3},
			`5.a(1, 2, foo: 3)`,
		},
		{
			`5.a(1, 2, foo: 3)`,
			[]int{1, 2},
			map[string]int{"foo": 3},
			`5.a(1, 2, foo: 3)`,
		},
		{
			`5.a(1, i: 2, j: 3)`,
			[]int{1},
			map[string]int{"i": 2, "j": 3},
			`5.a(1, i: 2, j: 3)`,
		},
		{
			`5.a(1, j: 3, i: 2)`,
			[]int{1},
			map[string]int{"i": 2, "j": 3},
			`5.a(1, i: 2, j: 3)`,
		},
	}

	for _, tt := range tests {
		program := testParse(t, tt.input)
		expr := testIfExprStmt(t, program)

		callExpr, ok := expr.(*ast.PropCallExpr)
		if !ok {
			t.Fatalf("expr is not *ast.PropCallExpr. got=%T", expr)
		}

		if len(callExpr.Args) != len(tt.args) {
			t.Fatalf("wrong arity of args, expected=%d, got=%d",
				len(tt.args), len(callExpr.Args))
		}

		if len(callExpr.Kwargs) != len(tt.kwargs) {
			t.Fatalf("wrong arity of kwargs, expected=%d, got=%d",
				len(tt.kwargs), len(callExpr.Kwargs))
		}

		if callExpr.String() != tt.printed {
			t.Errorf("wrong output.expected=\n%s,\ngot=\n%s",
				tt.printed, callExpr.String())
		}

		for i, expArg := range tt.args {
			testLiteralExpr(t, callExpr.Args[i], expArg)
		}

		for ident, val := range callExpr.Kwargs {
			name := ident.Token
			exp, ok := tt.kwargs[name]
			if ok {
				testLiteralExpr(t, val, exp)
			} else {
				t.Errorf("unexpected kwarg %s found.", name)
			}
		}
	}
}

func TestCallWithArgs(t *testing.T) {
	tests := []struct {
		input    string
		receiver interface{}
		propName string
		args     []interface{}
	}{
		{`5.hi`, 5, "hi", []interface{}{}},
		{`5@hi`, 5, "hi", []interface{}{}},
		{`5$hi`, 5, "hi", []interface{}{}},
		{`5.hi(6)`, 5, "hi", []interface{}{6}},
		{`5@hi(6)`, 5, "hi", []interface{}{6}},
		{`5$hi(6)`, 5, "hi", []interface{}{6}},
		{`5&.hi(6)`, 5, "hi", []interface{}{6}},
		{`5&@hi(6)`, 5, "hi", []interface{}{6}},
		{`5&$hi(6)`, 5, "hi", []interface{}{6}},
		{`5~.hi(6)`, 5, "hi", []interface{}{6}},
		{`5~@hi(6)`, 5, "hi", []interface{}{6}},
		{`5~$hi(6)`, 5, "hi", []interface{}{6}},
		{`5=.hi(6)`, 5, "hi", []interface{}{6}},
		{`5=@hi(6)`, 5, "hi", []interface{}{6}},
		{`5=$hi(6)`, 5, "hi", []interface{}{6}},
	}

	// TODO: inplement test
	for _, tt := range tests {
		program := testParse(t, tt.input)
		expr := testIfExprStmt(t, program)

		callExpr, ok := expr.(*ast.PropCallExpr)
		if !ok {
			t.Fatalf("expr is not *ast.PropCallExpr. got=%T", expr)
		}

		if len(callExpr.Args) != len(tt.args) {
			t.Fatalf("wrong arity of args, expected=%d, got=%d",
				len(tt.args), len(callExpr.Args))
		}

		for i, expArg := range tt.args {
			testLiteralExpr(t, callExpr.Args[i], expArg)
		}
	}
}

func TestAnonPropCall(t *testing.T) {
	tests := []struct {
		input        string
		chainContext string
		chainArg     interface{}
		propName     string
	}{
		{`.times`, ".", nil, "times"},
		{`@puts`, "@", nil, "puts"},
		{`@(10)puts`, "@", 10, "puts"},
		{`$add`, "$", nil, "add"},
		{`$(0)add`, "$", 0, "add"},
		{`$(0)+`, "$", 0, "+"},
		{`.foo`, ".", nil, "foo"},
		{`@foo`, "@", nil, "foo"},
		{`$foo`, "$", nil, "foo"},
		{`&.foo`, "&.", nil, "foo"},
		{`~.foo`, "~.", nil, "foo"},
		{`=.foo`, "=.", nil, "foo"},
		{`&@foo`, "&@", nil, "foo"},
		{`~@foo`, "~@", nil, "foo"},
		{`=@foo`, "=@", nil, "foo"},
		{`&$foo`, "&$", nil, "foo"},
		{`~$foo`, "~$", nil, "foo"},
		{`=$foo`, "=$", nil, "foo"},
	}

	for _, tt := range tests {
		program := testParse(t, tt.input)
		expr := testIfExprStmt(t, program)

		callExpr, ok := expr.(*ast.PropCallExpr)
		if !ok {
			t.Fatalf("expr is not *ast.PropCallExpr. got=%T", expr)
		}

		testNil(t, callExpr.Receiver)
		testChainContext(t, callExpr, tt.chainContext, tt.chainArg)
		testIdentifier(t, callExpr.Prop, tt.propName)
	}
}

func TestAnonPropCallWithArgs(t *testing.T) {
	tests := []struct {
		input    string
		propName string
		args     []interface{}
	}{
		{`.hi`, "hi", []interface{}{}},
		{`@hi`, "hi", []interface{}{}},
		{`$hi`, "hi", []interface{}{}},
		{`.hi(6)`, "hi", []interface{}{6}},
		{`@hi(6)`, "hi", []interface{}{6}},
		{`$hi(6)`, "hi", []interface{}{6}},
		{`&.hi(6)`, "hi", []interface{}{6}},
		{`&@hi(6)`, "hi", []interface{}{6}},
		{`&$hi(6)`, "hi", []interface{}{6}},
		{`~.hi(6)`, "hi", []interface{}{6}},
		{`~@hi(6)`, "hi", []interface{}{6}},
		{`~$hi(6)`, "hi", []interface{}{6}},
		{`=.hi(6)`, "hi", []interface{}{6}},
		{`=@hi(6)`, "hi", []interface{}{6}},
		{`=$hi(6)`, "hi", []interface{}{6}},
	}

	for _, tt := range tests {
		program := testParse(t, tt.input)
		expr := testIfExprStmt(t, program)

		callExpr, ok := expr.(*ast.PropCallExpr)
		if !ok {
			t.Fatalf("expr is not *ast.PropCallExpr. got=%T", expr)
		}

		if len(callExpr.Args) != len(tt.args) {
			t.Fatalf("wrong arity of args, expected=%d, got=%d",
				len(tt.args), len(callExpr.Args))
		}

		testNil(t, callExpr.Receiver)

		for i, expArg := range tt.args {
			testLiteralExpr(t, callExpr.Args[i], expArg)
		}
	}
}

func TestAnonOpMethods(t *testing.T) {
	tests := []struct {
		input string
		op    string
	}{
		{`.+(1)`, `+`},
		{`.-(1)`, `-`},
		{`.*(1)`, `*`},
		{`./(1)`, `/`},
		{`.//(1)`, `//`},
		{`.%(1)`, `%`},
		{`.**(1)`, `**`},
		{`.<=>(1)`, `<=>`},
		{`.==(1)`, `==`},
		{`.!=(1)`, `!=`},
		{`.<=(1)`, `<=`},
		{`.>=(1)`, `>=`},
		{`.<(1)`, `<`},
		{`.>(1)`, `>`},
		{`.>>(1)`, `>>`},
		{`.<<(1)`, `<<`},
		{`./&(1)`, `/&`},
		{`./|(1)`, `/|`},
		{`./^(1)`, `/^`},
		{`./~(1)`, `/~`},
		{`.!(1)`, `!`},
		{`.+%(1)`, `+%`},
		{`.-%(1)`, `-%`},
	}

	for _, tt := range tests {
		program := testParse(t, tt.input)
		expr := testIfExprStmt(t, program)

		callExpr, ok := expr.(*ast.PropCallExpr)
		if !ok {
			t.Fatalf("expr is not *ast.PropCallExpr. got=%T", expr)
		}

		testNil(t, callExpr.Receiver)
		testChainContext(t, callExpr, ".", nil)
		testIdentifier(t, callExpr.Prop, tt.op)

		if len(callExpr.Args) != 1 {
			t.Fatalf("arity must be 1. got=%d", len(callExpr.Args))
		}

		testLiteralExpr(t, callExpr.Args[0], 1)
	}
}

func TestLiteralCall(t *testing.T) {
	tests := []struct {
		input        string
		receiver     interface{}
		chainContext string
		chainArg     interface{}
	}{
		{`5.{|a| 1}`, 5, ".", nil},
		{`10@{|a| 1}`, 10, "@", nil},
		{`5@(10){|a| 1}`, 5, "@", 10},
		{`10${|a| 1}`, 10, "$", nil},
		{`5$(0){|a| 1}`, 5, "$", 0},
		{`5.{|a| 1}`, 5, ".", nil},
		{`5@{|a| 1}`, 5, "@", nil},
		{`5${|a| 1}`, 5, "$", nil},
		{`5&.{|a| 1}`, 5, "&.", nil},
		{`5~.{|a| 1}`, 5, "~.", nil},
		{`5=.{|a| 1}`, 5, "=.", nil},
		{`5&@{|a| 1}`, 5, "&@", nil},
		{`5~@{|a| 1}`, 5, "~@", nil},
		{`5=@{|a| 1}`, 5, "=@", nil},
		{`5&${|a| 1}`, 5, "&$", nil},
		{`5~${|a| 1}`, 5, "~$", nil},
		{`5=${|a| 1}`, 5, "=$", nil},
	}

	for _, tt := range tests {
		program := testParse(t, tt.input)
		expr := testIfExprStmt(t, program)

		callExpr, ok := expr.(*ast.LiteralCallExpr)
		if !ok {
			t.Fatalf("expr is not *ast.LiteralCallExpr. got=%T", expr)
		}

		testLiteralExpr(t, callExpr.Receiver, tt.receiver)
		testChainContext(t, callExpr, tt.chainContext, tt.chainArg)
	}
}

func TestLiteralCallFunc(t *testing.T) {
	input := `a.{|foo, hoge: 3, bar| 1+2; 3}`

	program := testParse(t, input)
	expr := testIfExprStmt(t, program)

	callExpr, ok := expr.(*ast.LiteralCallExpr)
	if !ok {
		t.Fatalf("expr is not *ast.LiteralCallExpr. got=%T", expr)
	}

	f := callExpr.Func

	if len(f.Args) != 2 {
		t.Fatalf("wrong arity of args. expected=2, got=%d", len(f.Args))
	}

	testIdentifier(t, f.Args[0], "foo")
	testIdentifier(t, f.Args[1], "bar")

	if len(f.Kwargs) != 1 {
		t.Fatalf("wrong arity of kwargs. expected=1, got=%d", len(f.Kwargs))
	}

	for k, v := range f.Kwargs {
		testIdentifier(t, k, "hoge")
		testLiteralExpr(t, v, 3)
	}

	if len(f.Body) != 2 {
		t.Fatalf("wrong num of stmts. expected=2, got=%d", len(f.Body))
	}

	es0, ok := f.Body[0].(*ast.ExprStmt)
	if !ok {
		t.Fatalf("Body[0] is not ast*ExprStmt. got=%T", f.Body[0])
	}

	testInfixOperator(t, es0.Expr, 1, "+", 2)

	es1, ok := f.Body[1].(*ast.ExprStmt)
	if !ok {
		t.Fatalf("Body[1] is not ast*ExprStmt. got=%T", f.Body[1])
	}

	testLiteralExpr(t, es1.Expr, 3)
}

func TestLiteralCallFuncArgs(t *testing.T) {
	tests := []struct {
		input   string
		args    []interface{}
		kwargs  map[string]interface{}
		printed string
	}{
		{
			`a.{|b| c}()`,
			[]interface{}{},
			map[string]interface{}{},
			`a.{|b| c}()`,
		},
		{
			`a.{|b| c}`,
			[]interface{}{},
			map[string]interface{}{},
			`a.{|b| c}()`,
		},
		{
			`a.m{c}`,
			[]interface{}{},
			map[string]interface{}{},
			`a.{|self| c}()`,
		},
		{
			`a.m{c}(1)`,
			[]interface{}{1},
			map[string]interface{}{},
			`a.{|self| c}(1)`,
		},
		{
			`a.{|b| c}(1)`,
			[]interface{}{1},
			map[string]interface{}{},
			`a.{|b| c}(1)`,
		},
		{
			`a.{|b| c}(1, foo)`,
			[]interface{}{1, "foo"},
			map[string]interface{}{},
			`a.{|b| c}(1, foo)`,
		},
		{
			`a.{|b| c}(foo: 2)`,
			[]interface{}{},
			map[string]interface{}{"foo": 2},
			`a.{|b| c}(foo: 2)`,
		},
		{
			`a.{|b| c}(1, foo, bar: 2)`,
			[]interface{}{1, "foo"},
			map[string]interface{}{"bar": 2},
			`a.{|b| c}(1, foo, bar: 2)`,
		},
		{
			`a.{|b| c}(1, bar: 2, foo)`,
			[]interface{}{1, "foo"},
			map[string]interface{}{"bar": 2},
			`a.{|b| c}(1, foo, bar: 2)`,
		},
		{
			`a.{|b| c}(1, f: 3, e, d: 2)`,
			[]interface{}{1, "e"},
			map[string]interface{}{"d": 2, "f": 3},
			`a.{|b| c}(1, e, d: 2, f: 3)`,
		},
		{
			`.{|b| c}()`,
			[]interface{}{},
			map[string]interface{}{},
			`.{|b| c}()`,
		},
		{
			`.{|b| c}`,
			[]interface{}{},
			map[string]interface{}{},
			`.{|b| c}()`,
		},
		{
			`.m{c}`,
			[]interface{}{},
			map[string]interface{}{},
			`.{|self| c}()`,
		},
		{
			`.m{c}(1)`,
			[]interface{}{1},
			map[string]interface{}{},
			`.{|self| c}(1)`,
		},
	}

	for _, tt := range tests {
		errPrefix := fmt.Sprintf("err in ```\n%s\n```\n", tt.input)

		program := testParse(t, tt.input)
		expr := testIfExprStmt(t, program)

		callExpr, ok := expr.(*ast.LiteralCallExpr)

		if !ok {
			t.Fatalf("%sexpr is not *ast.LiteralCallExpr. got=%T",
				errPrefix, expr)
		}
		if len(callExpr.Args) != len(tt.args) {
			t.Fatalf("%swrong arity of args, expected=%d, got=%d",
				errPrefix, len(tt.args), len(callExpr.Args))
		}

		if len(callExpr.Kwargs) != len(tt.kwargs) {
			t.Fatalf("%swrong arity of kwargs, expected=%d, got=%d",
				errPrefix, len(tt.kwargs), len(callExpr.Kwargs))
		}

		if callExpr.String() != tt.printed {
			t.Errorf("wrong output.expected=\n%s,\ngot=\n%s",
				tt.printed, callExpr.String())
		}

		for i, expArg := range tt.args {
			switch a := expArg.(type) {
			case string:
				testIdentifier(t, callExpr.Args[i], a)
			default:
				testLiteralExpr(t, callExpr.Args[i], a)
			}
		}

		for ident, val := range callExpr.Kwargs {
			name := ident.Token
			exp, ok := tt.kwargs[name]
			if ok {
				testLiteralExpr(t, val, exp)
			} else {
				t.Errorf("unexpected kwarg %s found.", name)
			}
		}
	}
}

func TestAnonLiteralCall(t *testing.T) {
	tests := []struct {
		input        string
		chainContext string
		chainArg     interface{}
	}{
		{`.{|a| 1}`, ".", nil},
		{`@{|a| 1}`, "@", nil},
		{`@(10){|a| 1}`, "@", 10},
		{`${|a| 1}`, "$", nil},
		{`$(0){|a| 1}`, "$", 0},
		{`.{|a| 1}`, ".", nil},
		{`@{|a| 1}`, "@", nil},
		{`${|a| 1}`, "$", nil},
		{`&.{|a| 1}`, "&.", nil},
		{`~.{|a| 1}`, "~.", nil},
		{`=.{|a| 1}`, "=.", nil},
		{`&@{|a| 1}`, "&@", nil},
		{`~@{|a| 1}`, "~@", nil},
		{`=@{|a| 1}`, "=@", nil},
		{`&${|a| 1}`, "&$", nil},
		{`~${|a| 1}`, "~$", nil},
		{`=${|a| 1}`, "=$", nil},
	}

	for _, tt := range tests {
		program := testParse(t, tt.input)
		expr := testIfExprStmt(t, program)

		callExpr, ok := expr.(*ast.LiteralCallExpr)
		if !ok {
			t.Fatalf("expr is not *ast.LiteralCallExpr. got=%T", expr)
		}

		testNil(t, callExpr.Receiver)
		testChainContext(t, callExpr, tt.chainContext, tt.chainArg)
	}
}

func TestVarCall(t *testing.T) {
	tests := []struct {
		input        string
		receiver     interface{}
		chainContext string
		chainArg     interface{}
		varName      string
	}{
		{`5.^times`, 5, ".", nil, "times"},
		{`10@^puts`, 10, "@", nil, "puts"},
		{`5@(10)^puts`, 5, "@", 10, "puts"},
		{`10$^add`, 10, "$", nil, "add"},
		{`5$(0)^add`, 5, "$", 0, "add"},
		{`5.^foo`, 5, ".", nil, "foo"},
		{`5@^foo`, 5, "@", nil, "foo"},
		{`5$^foo`, 5, "$", nil, "foo"},
		{`5&.^foo`, 5, "&.", nil, "foo"},
		{`5~.^foo`, 5, "~.", nil, "foo"},
		{`5=.^foo`, 5, "=.", nil, "foo"},
		{`5&@^foo`, 5, "&@", nil, "foo"},
		{`5~@^foo`, 5, "~@", nil, "foo"},
		{`5=@^foo`, 5, "=@", nil, "foo"},
		{`5&$^foo`, 5, "&$", nil, "foo"},
		{`5~$^foo`, 5, "~$", nil, "foo"},
		{`5=$^foo`, 5, "=$", nil, "foo"},
	}

	for _, tt := range tests {
		program := testParse(t, tt.input)
		expr := testIfExprStmt(t, program)

		callExpr, ok := expr.(*ast.VarCallExpr)
		if !ok {
			t.Fatalf("expr is not *ast.VarCallExpr. got=%T", expr)
		}

		testLiteralExpr(t, callExpr.Receiver, tt.receiver)
		testChainContext(t, callExpr, tt.chainContext, tt.chainArg)
		testIdentifier(t, callExpr.Var, tt.varName)
	}
}

func TestVarCallFuncArgs(t *testing.T) {
	tests := []struct {
		input   string
		args    []interface{}
		kwargs  map[string]interface{}
		printed string
	}{
		{
			`a.^foo()`,
			[]interface{}{},
			map[string]interface{}{},
			`a.^foo()`,
		},
		{
			`a.^foo`,
			[]interface{}{},
			map[string]interface{}{},
			`a.^foo()`,
		},
		{
			`a.^foo(1)`,
			[]interface{}{1},
			map[string]interface{}{},
			`a.^foo(1)`,
		},
		{
			`a.^foo(1, foo)`,
			[]interface{}{1, "foo"},
			map[string]interface{}{},
			`a.^foo(1, foo)`,
		},
		{
			`a.^foo(foo: 2)`,
			[]interface{}{},
			map[string]interface{}{"foo": 2},
			`a.^foo(foo: 2)`,
		},
		{
			`a.^foo(1, foo, bar: 2)`,
			[]interface{}{1, "foo"},
			map[string]interface{}{"bar": 2},
			`a.^foo(1, foo, bar: 2)`,
		},
		{
			`a.^foo(1, bar: 2, foo)`,
			[]interface{}{1, "foo"},
			map[string]interface{}{"bar": 2},
			`a.^foo(1, foo, bar: 2)`,
		},
		{
			`a.^foo(1, f: 3, e, d: 2)`,
			[]interface{}{1, "e"},
			map[string]interface{}{"d": 2, "f": 3},
			`a.^foo(1, e, d: 2, f: 3)`,
		},
		{
			`.^foo()`,
			[]interface{}{},
			map[string]interface{}{},
			`.^foo()`,
		},
		{
			`.^foo`,
			[]interface{}{},
			map[string]interface{}{},
			`.^foo()`,
		},
		{
			`.^foo(1)`,
			[]interface{}{1},
			map[string]interface{}{},
			`.^foo(1)`,
		},
	}

	for _, tt := range tests {
		errPrefix := fmt.Sprintf("err in ```\n%s\n```\n", tt.input)

		program := testParse(t, tt.input)
		expr := testIfExprStmt(t, program)

		callExpr, ok := expr.(*ast.VarCallExpr)

		if !ok {
			t.Fatalf("%sexpr is not *ast.VarCallExpr. got=%T",
				errPrefix, expr)
		}
		if len(callExpr.Args) != len(tt.args) {
			t.Fatalf("%swrong arity of args, expected=%d, got=%d",
				errPrefix, len(tt.args), len(callExpr.Args))
		}

		if len(callExpr.Kwargs) != len(tt.kwargs) {
			t.Fatalf("%swrong arity of kwargs, expected=%d, got=%d",
				errPrefix, len(tt.kwargs), len(callExpr.Kwargs))
		}

		if callExpr.String() != tt.printed {
			t.Errorf("wrong output.expected=\n%s,\ngot=\n%s",
				tt.printed, callExpr.String())
		}

		for i, expArg := range tt.args {
			switch a := expArg.(type) {
			case string:
				testIdentifier(t, callExpr.Args[i], a)
			default:
				testLiteralExpr(t, callExpr.Args[i], a)
			}
		}

		for ident, val := range callExpr.Kwargs {
			name := ident.Token
			exp, ok := tt.kwargs[name]
			if ok {
				testLiteralExpr(t, val, exp)
			} else {
				t.Errorf("unexpected kwarg %s found.", name)
			}
		}
	}
}

func TestAnonVarCall(t *testing.T) {
	tests := []struct {
		input        string
		chainContext string
		chainArg     interface{}
		varName      string
	}{
		{`.^times`, ".", nil, "times"},
		{`@^puts`, "@", nil, "puts"},
		{`@(10)^puts`, "@", 10, "puts"},
		{`$^add`, "$", nil, "add"},
		{`$(0)^add`, "$", 0, "add"},
		{`.^foo`, ".", nil, "foo"},
		{`@^foo`, "@", nil, "foo"},
		{`$^foo`, "$", nil, "foo"},
		{`&.^foo`, "&.", nil, "foo"},
		{`~.^foo`, "~.", nil, "foo"},
		{`=.^foo`, "=.", nil, "foo"},
		{`&@^foo`, "&@", nil, "foo"},
		{`~@^foo`, "~@", nil, "foo"},
		{`=@^foo`, "=@", nil, "foo"},
		{`&$^foo`, "&$", nil, "foo"},
		{`~$^foo`, "~$", nil, "foo"},
		{`=$^foo`, "=$", nil, "foo"},
	}

	for _, tt := range tests {
		program := testParse(t, tt.input)
		expr := testIfExprStmt(t, program)

		callExpr, ok := expr.(*ast.VarCallExpr)
		if !ok {
			t.Fatalf("expr is not *ast.VarCallExpr. got=%T", expr)
		}

		testNil(t, callExpr.Receiver)
		testChainContext(t, callExpr, tt.chainContext, tt.chainArg)
		testIdentifier(t, callExpr.Var, tt.varName)
	}
}

func TestMultipleLineChain(t *testing.T) {
	input := `
	foo
	  |.bar(1)
	  |@{|i| i*2}
	  |$(3)^hoge
	  |~.+
	`

	program := testParse(t, input)
	expr := testIfExprStmt(t, program)

	expectedStr := `foo.bar(1)@{|i| (i * 2)}()$(3)^hoge()~.+()`
	if expr.String() != expectedStr {
		t.Errorf("wrong output. expected=`\n%s\n`, got=`\n%s\n`",
			expectedStr, expr.String())
	}

	callExpr1, ok := expr.(*ast.PropCallExpr)
	if !ok {
		t.Fatalf("expr is not *ast.PropCallExpr. got=%T", expr)
	}
	testIdentifier(t, callExpr1.Prop, "+")
	testChainContext(t, callExpr1, "~.", nil)

	expr2 := callExpr1.Receiver
	if expr2 == nil {
		t.Fatalf("expr2 must not be nil.")
	}

	callExpr2, ok := expr2.(*ast.VarCallExpr)
	if !ok {
		t.Fatalf("expr2 is not *ast.VarCallExpr. got=%T", expr2)
	}
	testIdentifier(t, callExpr2.Var, "hoge")
	testChainContext(t, callExpr2, "$", 3)

	expr3 := callExpr2.Receiver
	if expr3 == nil {
		t.Fatalf("expr3 must not be nil.")
	}

	callExpr3, ok := expr3.(*ast.LiteralCallExpr)
	if !ok {
		t.Fatalf("expr3 is not *ast.LiteralCallExpr. got=%T", expr3)
	}

	if len(callExpr3.Func.Args) != 1 {
		t.Fatalf("callExpr3.Func must have 1 arg. got=%d",
			len(callExpr3.Func.Args))
	}
	testIdentifier(t, callExpr3.Func.Args[0], "i")
	testChainContext(t, callExpr3, "@", nil)

	expr4 := callExpr3.Receiver
	if expr4 == nil {
		t.Fatalf("expr4 must not be nil.")
	}

	callExpr4, ok := expr4.(*ast.PropCallExpr)
	if !ok {
		t.Fatalf("expr4 is not *ast.PropCallExpr. got=%T", expr4)
	}
	testIdentifier(t, callExpr4.Prop, "bar")
	testChainContext(t, callExpr4, ".", nil)

	if len(callExpr4.Args) != 1 {
		t.Fatalf("callExpr4 must have 1 arg. got=%d",
			len(callExpr4.Args))
	}
	testLiteralExpr(t, callExpr4.Args[0], 1)

	expr5 := callExpr4.Receiver
	if expr5 == nil {
		t.Fatalf("expr5 must not be nil.")
	}
	testIdentifier(t, expr5, "foo")
}

func TestIndexExprRecv(t *testing.T) {
	// `expr[arg]` is syntax sugar of `expr.at([arg])`

	parseIndexRecv := func(input string) (ast.Expr, bool) {
		program := testParse(t, input)
		expr := testIfExprStmt(t, program)
		idxExpr, ok := expr.(*ast.PropCallExpr)
		if !ok {
			t.Fatalf("expr is not *ast.PropCallExpr. got=%T", expr)
			return nil, false
		}
		if !testIdentifier(t, idxExpr.Prop, "at") {
			return nil, false
		}
		if !testChainContext(t, idxExpr, ".", nil) {
			return nil, false
		}
		if len(idxExpr.Args) != 1 {
			t.Fatalf("arity must be 1. got=%d", len(idxExpr.Args))
			return nil, false
		}
		if _, ok := idxExpr.Args[0].(*ast.ArrLiteral); !ok {
			t.Fatalf("1st arg must be *ast.ArrLiteral. got=%T",
				idxExpr.Args[0])
			return nil, false
		}

		return idxExpr.Receiver, true
	}

	input0 := `a[0]`
	recv0, ok := parseIndexRecv(input0)
	if !ok {
		t.Fatalf("recv0 test failed")
	}
	testIdentifier(t, recv0, "a")

	input1 := `10[0]`
	recv1, ok := parseIndexRecv(input1)
	if !ok {
		t.Fatalf("recv1 test failed")
	}
	testLiteralExpr(t, recv1, 10)

	input2 := `"a"[0]`
	recv2, ok := parseIndexRecv(input2)
	if !ok {
		t.Fatalf("recv2 test failed")
	}
	testStr(t, recv2, "a", false)

	input3 := `[0,1,2][0]`
	recv3, ok := parseIndexRecv(input3)
	if !ok {
		t.Fatalf("recv3 test failed")
	}
	arr, ok := recv3.(*ast.ArrLiteral)
	if !ok {
		t.Fatalf("recv3 is not *ast.ArrLiteral. got=%T", recv3)
	}
	if len(arr.Elems) != 3 {
		t.Fatalf("length of recv3 must be 3. got=%d", len(arr.Elems))
	}
	for i, e := range arr.Elems {
		testLiteralExpr(t, e, i)
	}

	input4 := `{'a: 1}[0]`
	recv4, ok := parseIndexRecv(input4)
	if !ok {
		t.Fatalf("recv4 test failed")
	}
	obj, ok := recv4.(*ast.ObjLiteral)
	if !ok {
		t.Fatalf("recv4 is not *ast.ObjLiteral. got=%T", recv4)
	}
	if len(obj.Pairs) != 1 {
		t.Fatalf("length of recv4 must be 1. got=%d", len(obj.Pairs))
	}
	p := obj.Pairs[0]
	testSymbol(t, p.Key, "a")
	testLiteralExpr(t, p.Val, 1)

	input5 := `a.foo[0]`
	recv5, ok := parseIndexRecv(input5)
	if !ok {
		t.Fatalf("recv5 test failed")
	}
	ce, ok := recv5.(*ast.PropCallExpr)
	if !ok {
		t.Fatalf("recv5 is not *ast.PropCallExpr. got=%T", recv5)
	}
	testIdentifier(t, ce.Receiver, "a")
	testIdentifier(t, ce.Prop, "foo")
	if len(ce.Args) != 0 {
		t.Errorf("arity of ce must be 0. got=%d", len(ce.Args))
	}
}

func TestIndexExprArg(t *testing.T) {
	tests := []struct {
		input     string
		elemTypes []string
		elems     []interface{}
	}{
		{`a[1]`, []string{"Int"}, []interface{}{1}},
		{`a[2,3]`, []string{"Int", "Int"}, []interface{}{2, 3}},
		{`a[4,5,"a"]`, []string{"Int", "Int", "Str"}, []interface{}{4, 5, "a"}},
		{`a[[6,7]]`, []string{"Arr"}, []interface{}{[]int{6, 7}}},
		{`a[[8,9], 10]`, []string{"Arr", "Int"}, []interface{}{[]int{8, 9}, 10}},
		{`a[{'a: 11}]`, []string{"Obj"}, []interface{}{[]interface{}{"a", 11}}},
		{`a[b]`, []string{"Ident"}, []interface{}{"b"}},
	}

	for _, tt := range tests {
		program := testParse(t, tt.input)
		expr := testIfExprStmt(t, program)
		idxExpr, ok := expr.(*ast.PropCallExpr)
		if !ok {
			t.Fatalf("expr is not *ast.PropCallExpr. got=%T", expr)
		}

		testIdentifier(t, idxExpr.Prop, "at")
		testChainContext(t, idxExpr, ".", nil)

		if len(idxExpr.Args) != 1 {
			t.Fatalf("arity must be 1. got=%d", len(idxExpr.Args))
		}

		lst, ok := idxExpr.Args[0].(*ast.ArrLiteral)
		if !ok {
			t.Fatalf("1st arg must be *ast.ArrLiteral. got=%T",
				idxExpr.Args[0])
		}

		if len(lst.Elems) != len(tt.elems) {
			t.Fatalf("wrong length of elems. expected=%d. got=%d",
				len(tt.elems), len(lst.Elems))
		}

		for i, e := range lst.Elems {
			switch tt.elemTypes[i] {
			case "Int":
				testLiteralExpr(t, e, tt.elems[i])
			case "Str":
				str := tt.elems[i].(string)
				testStr(t, e, str, false)
			case "Ident":
				id := tt.elems[i].(string)
				testIdentifier(t, e, id)
			case "Arr":
				arr, ok := e.(*ast.ArrLiteral)
				if !ok {
					t.Fatalf("Elems[%d] in `%s` is not Arr. got=%T",
						i, tt.input, e)
				}
				ttArr := tt.elems[i].([]int)
				if len(arr.Elems) != len(ttArr) {
					t.Fatalf("wrong arr length. expected=%d, got=%d",
						len(ttArr), len(arr.Elems))
				}
				for j, arrElem := range arr.Elems {
					testLiteralExpr(t, arrElem, ttArr[j])
				}
			case "Obj":
				obj, ok := e.(*ast.ObjLiteral)
				if !ok {
					t.Fatalf("Elems[%d] in `%s` is not Obj. got=%T",
						i, tt.input, e)
				}
				if len(obj.Pairs) != 1 {
					t.Fatalf("length of obj.Pairs must be 1. got=%d(%s)",
						len(obj.Pairs), obj.String())
				}
				pair := obj.Pairs[0]
				testSymbol(t, pair.Key, "a")
				testLiteralExpr(t, pair.Val, 11)
			}
		}
	}
}

func TestIntLiteralExpr(t *testing.T) {
	tests := []struct {
		input    string
		expected int64
	}{
		// NOTE: minus is recognized as prefix
		{`5`, 5},
		{`100`, 100},
	}

	for _, tt := range tests {
		program := testParse(t, tt.input)
		expr := testIfExprStmt(t, program)

		if !testIntLiteral(t, expr, tt.expected) {
			return
		}
	}
}

func TestFuncLiteralArgs(t *testing.T) {
	tests := []struct {
		input   string
		args    []string
		kwargs  map[string]interface{}
		printed string
	}{
		// NOTE: `{}` is recognized as Obj
		{
			`{||}`,
			[]string{},
			map[string]interface{}{},
			`{|| }`,
		},
		{
			`{|| a}`,
			[]string{},
			map[string]interface{}{},
			`{|| a}`,
		},
		{
			`{|a| 1}`,
			[]string{"a"},
			map[string]interface{}{},
			`{|a| 1}`,
		},
		{
			`{1}`,
			[]string{},
			map[string]interface{}{},
			`{|| 1}`,
		},
		{
			`{|a, foo| 1}`,
			[]string{"a", "foo"},
			map[string]interface{}{},
			`{|a, foo| 1}`,
		},
		{
			`{|val: 1| val}`,
			[]string{},
			map[string]interface{}{"val": 1},
			`{|val: 1| val}`,
		},
		{
			`{|a, val: 1| val}`,
			[]string{"a"},
			map[string]interface{}{"val": 1},
			`{|a, val: 1| val}`,
		},
		{
			`{|val: 1, a| val}`,
			[]string{"a"},
			map[string]interface{}{"val": 1},
			`{|a, val: 1| val}`,
		},
		{
			`{|a, b, c: 1, d: 2| val}`,
			[]string{"a", "b"},
			map[string]interface{}{"c": 1, "d": 2},
			`{|a, b, c: 1, d: 2| val}`,
		},
		{
			`{|a, b, c: 1, d: 2| val}`,
			[]string{"a", "b"},
			map[string]interface{}{"c": 1, "d": 2},
			`{|a, b, c: 1, d: 2| val}`,
		},
		{
			`{|a, b, d: 2, c: 1| val}`,
			[]string{"a", "b"},
			map[string]interface{}{"c": 1, "d": 2},
			`{|a, b, c: 1, d: 2| val}`,
		},
		{
			`{|a, c: 1, b, d: 2| val}`,
			[]string{"a", "b"},
			map[string]interface{}{"c": 1, "d": 2},
			`{|a, b, c: 1, d: 2| val}`,
		},
		{
			`{|d: 2, c: 1, a, b| val}`,
			[]string{"a", "b"},
			map[string]interface{}{"c": 1, "d": 2},
			`{|a, b, c: 1, d: 2| val}`,
		},
		{
			`{|d: 2, a, c: 1, b| val}`,
			[]string{"a", "b"},
			map[string]interface{}{"c": 1, "d": 2},
			`{|a, b, c: 1, d: 2| val}`,
		},
		{
			`m{}`,
			[]string{"self"},
			map[string]interface{}{},
			`{|self| }`,
		},
		{
			`m{val}`,
			[]string{"self"},
			map[string]interface{}{},
			`{|self| val}`,
		},
		{
			`m{||}`,
			[]string{"self"},
			map[string]interface{}{},
			`{|self| }`,
		},
		{
			`m{|a|}`,
			[]string{"self", "a"},
			map[string]interface{}{},
			`{|self, a| }`,
		},
		{
			`m{|opt: 1|}`,
			[]string{"self"},
			map[string]interface{}{"opt": 1},
			`{|self, opt: 1| }`,
		},
		{
			`m{|opt: 1| val}`,
			[]string{"self"},
			map[string]interface{}{"opt": 1},
			`{|self, opt: 1| val}`,
		},
		{
			`m{|opt: 1, a|}`,
			[]string{"self", "a"},
			map[string]interface{}{"opt": 1},
			`{|self, a, opt: 1| }`,
		},
		{
			`m{|b: 1, a, c: 2|}`,
			[]string{"self", "a"},
			map[string]interface{}{"b": 1, "c": 2},
			`{|self, a, b: 1, c: 2| }`,
		},
	}

	for _, tt := range tests {
		program := testParse(t, tt.input)
		expr := testIfExprStmt(t, program)

		f, ok := expr.(*ast.FuncLiteral)
		if !ok {
			t.Fatalf("f is not *ast.FuncLiteral. got=%T", expr)
		}

		if len(f.Args) != len(tt.args) {
			t.Fatalf("wrong arity of args, expected=%d, got=%d in\n%s",
				len(tt.args), len(f.Args), tt.input)
		}

		if len(f.Kwargs) != len(tt.kwargs) {
			t.Fatalf("wrong arity of kwargs, expected=%d, got=%d",
				len(tt.kwargs), len(f.Kwargs))
		}

		for i, expArg := range tt.args {
			testIdentifier(t, f.Args[i], expArg)
		}

		for ident, val := range f.Kwargs {
			name := ident.Token
			exp, ok := tt.kwargs[name]
			if ok {
				testLiteralExpr(t, val, exp)
			} else {
				t.Errorf("unexpected kwarg %s found.", name)
			}
		}

		if f.String() != tt.printed {
			t.Errorf("wrong output.expected=\n%s,\ngot=\n%s",
				tt.printed, f.String())
		}
	}
}

func TestFuncLiteralBody(t *testing.T) {
	tests := []struct {
		input    string
		bodyType string
		body     interface{}
	}{
		{`{|| 2}`, "literal", 2},
		{`{|a| a}`, "ident", "a"},
		{`{|a: 1| 1+1}`, "infix", []interface{}{1, "+", 1}},
		{`m{2}`, "literal", 2},
		{`m{|| 2}`, "literal", 2},
		{`m{|a| a}`, "ident", "a"},
		{`m{|a: 1| 1+1}`, "infix", []interface{}{1, "+", 1}},
	}

	for _, tt := range tests {
		program := testParse(t, tt.input)
		expr := testIfExprStmt(t, program)

		f, ok := expr.(*ast.FuncLiteral)
		if !ok {
			t.Fatalf("f is not *ast.FuncLiteral. got=%T", expr)
		}

		if len(f.Body) != 1 {
			t.Fatalf("body length is not 1. got=%d", len(f.Body))
		}

		es, ok := f.Body[0].(*ast.ExprStmt)
		if !ok {
			t.Fatalf("f.Body[0] is not *ast.ExprStmt. got=%T", f.Body[0])
		}

		e := es.Expr
		switch tt.bodyType {
		case "literal":
			testLiteralExpr(t, e, tt.body)
		case "ident":
			s, _ := tt.body.(string)
			testIdentifier(t, e, s)
		case "infix":
			arr, _ := tt.body.([]interface{})
			left, _ := arr[0].(int)
			op, _ := arr[1].(string)
			right, _ := arr[2].(int)
			testInfixOperator(t, e, left, op, right)
		}
	}
}

func TestFuncLiteralBodies(t *testing.T) {
	input := `
	{|a, b|
	  2
	  a
	  1 + 1
	}
	`

	tests := []struct {
		bodyType string
		body     interface{}
	}{
		{"literal", 2},
		{"ident", "a"},
		{"infix", []interface{}{1, "+", 1}},
	}

	program := testParse(t, input)
	expr := testIfExprStmt(t, program)

	f, ok := expr.(*ast.FuncLiteral)
	if !ok {
		t.Fatalf("f is not *ast.FuncLiteral. got=%T", expr)
	}

	if len(f.Body) != len(tests) {
		t.Fatalf("body length is not %d. got=%d",
			len(tests), len(f.Body))
	}

	for i, tt := range tests {
		es, ok := f.Body[i].(*ast.ExprStmt)
		if !ok {
			t.Fatalf("f.Body[%d] is not *ast.ExprStmt. got=%T",
				i, f.Body[i])
		}

		e := es.Expr
		switch tt.bodyType {
		case "literal":
			testLiteralExpr(t, e, tt.body)
		case "ident":
			s, _ := tt.body.(string)
			testIdentifier(t, e, s)
		case "infix":
			arr, _ := tt.body.([]interface{})
			left, _ := arr[0].(int)
			op, _ := arr[1].(string)
			right, _ := arr[2].(int)
			testInfixOperator(t, e, left, op, right)
		}
	}
}

func TestComment(t *testing.T) {
	tests := []struct {
		input   string
		vals    []interface{}
		printed string
	}{
		{
			`#`,
			[]interface{}{},
			"",
		},
		{
			`#foo`,
			[]interface{}{},
			"",
		},
		{
			`1 #foo`,
			[]interface{}{1},
			"1",
		},
		{
			`1#foo`,
			[]interface{}{1},
			"1",
		},
		{
			`1 ##f#o#o`,
			[]interface{}{1},
			"1",
		},
		{
			`1 #+2`,
			[]interface{}{1},
			"1",
		},
		{
			`a#foo`,
			[]interface{}{"a"},
			"a",
		},
		{
			`1 #foo
			#2 * 3
			4`,
			[]interface{}{1, 4},
			"1\n4",
		},
		{
			`#foo
			4`,
			[]interface{}{4},
			"4",
		},
		{
			`4
			#foo`,
			[]interface{}{4},
			"4",
		},
	}

	for _, tt := range tests {
		program := testParse(t, tt.input)

		if len(program.Stmts) != len(tt.vals) {
			t.Fatalf("wrong length of stmts in`\n%s\n`. expected=%d, got=%d",
				tt.input, len(tt.vals), len(program.Stmts))
		}

		if program.String() != tt.printed {
			t.Errorf("wrong output of String. expected=`\n%s\n`. got=`\n%s\n`.",
				tt.printed, program.String())
		}

		for i, stmt := range program.Stmts {
			exprStmt, ok := stmt.(*ast.ExprStmt)
			if !ok {
				t.Fatalf("stmts[%d] is not *ast.ExprStmt. got=%T", i, stmt)
			}

			switch val := tt.vals[i].(type) {
			case string:
				testIdentifier(t, exprStmt.Expr, val)
			default:
				testLiteralExpr(t, exprStmt.Expr, val)
			}
		}
	}
}

func testChainContext(t *testing.T, ce ast.CallExpr, expContext string,
	expArg interface{}) bool {
	if ce.ChainToken() != expContext {
		t.Errorf("chain is not %s. got=%s", expContext, ce.ChainToken())
		return false
	}

	if expArg == nil {
		return testNil(t, ce.ChainArg())
	}
	return testLiteralExpr(t, ce.ChainArg(), expArg)
}

func testIdentifier(t *testing.T, expr ast.Expr, expected string) bool {
	ident, ok := expr.(*ast.Ident)
	if !ok {
		t.Errorf("exp not *ast.Ident. got=%T", expr)
		return false
	}

	if ident.IdentAttr != ast.NormalIdent {
		t.Errorf("ident.IdentAttr not ast.NormalIdent. got=%T",
			ident.IdentAttr)
		return false
	}

	if ident.Value != expected {
		t.Errorf("ident.Value not %s. got=%s", expected, ident.Value)
		return false
	}

	if ident.TokenLiteral() != expected {
		t.Errorf("ident.TokenLiteral() not %s. got=%s",
			expected, ident.TokenLiteral())
		return false
	}

	return true
}

func testArgIdent(t *testing.T, expr ast.Expr, expected string) bool {
	ident, ok := expr.(*ast.Ident)
	if !ok {
		t.Errorf("exp not *ast.Ident. got=%T", expr)
		return false
	}

	if ident.IdentAttr != ast.ArgIdent {
		t.Errorf("ident.IdentAttr not ast.ArgIdent. got=%T",
			ident.IdentAttr)
		return false
	}

	if ident.Value != expected {
		t.Errorf("ident.Value not %s. got=%s", expected, ident.Value)
		return false
	}

	if ident.TokenLiteral() != expected {
		t.Errorf("ident.TokenLiteral() not %s. got=%s",
			expected, ident.TokenLiteral())
		return false
	}

	return true
}

func testKwargIdent(t *testing.T, expr ast.Expr, expected string) bool {
	ident, ok := expr.(*ast.Ident)
	if !ok {
		t.Errorf("exp not *ast.Ident. got=%T", expr)
		return false
	}

	if ident.IdentAttr != ast.KwargIdent {
		t.Errorf("ident.IdentAttr not ast.KwargIdent. got=%T",
			ident.IdentAttr)
		return false
	}

	if ident.Value != expected {
		t.Errorf("ident.Value not %s. got=%s", expected, ident.Value)
		return false
	}

	if ident.TokenLiteral() != expected {
		t.Errorf("ident.TokenLiteral() not %s. got=%s",
			expected, ident.TokenLiteral())
		return false
	}

	return true
}

func testPrefixOperator(t *testing.T, expr ast.Expr,
	op string, right interface{}) {
	prefixExpr, ok := expr.(*ast.PrefixExpr)

	if !ok {
		t.Fatalf("expr is not ast.PrefixExpr. got=%T", expr)
	}

	testLiteralExpr(t, prefixExpr.Right, right)

	if prefixExpr.Operator != op {
		t.Errorf("operator is not '%s'. got=%s", op, prefixExpr.Operator)
	}
}

func testInfixOperator(t *testing.T, expr ast.Expr,
	left interface{}, op string, right interface{}) {
	infixExpr, ok := expr.(*ast.InfixExpr)

	if !ok {
		t.Fatalf("expr is not ast.InfixExpr. got=%T", expr)
	}

	testLiteralExpr(t, infixExpr.Left, left)
	testLiteralExpr(t, infixExpr.Right, right)

	if infixExpr.Operator != op {
		t.Errorf("operator is not '%s'. got=%s", op, infixExpr.Operator)
	}
}

func testIfExprStmt(t *testing.T, program *ast.Program) ast.Expr {
	if len(program.Stmts) != 1 {
		t.Fatalf("program does not contain %d statements. got=%d",
			1, len(program.Stmts))
	}

	stmt, ok := program.Stmts[0].(*ast.ExprStmt)
	if !ok {
		t.Fatalf("program.Stmts[0] is not ast.ExprStmt. got=%T",
			program.Stmts[0])
	}

	return stmt.Expr
}

func testLiteralExpr(t *testing.T, exp ast.Expr, expected interface{}) bool {
	switch v := expected.(type) {
	case int:
		return testIntLiteral(t, exp, int64(v))
	case int64:
		return testIntLiteral(t, exp, v)
		//case string:
		//	return testIdentifier(t, exp, v)
		//case bool:
		//	return testBoolean(t, exp, v)
	}
	t.Errorf("type of exp not expected. got=%T", exp)
	return false
}

func testStr(t *testing.T, expr ast.Expr,
	expected string, isRaw bool) bool {
	sl, ok := expr.(*ast.StrLiteral)
	if !ok {
		t.Errorf("sl not *ast.StrLiteral. got=%T", expr)
		return false
	}

	if sl.Value != expected {
		t.Errorf("sl.Value not %s. got=%v", expected, sl.Value)
		return false
	}

	if sl.IsRaw != isRaw {
		t.Errorf("sl.IsRaw should be %v. got=%v", isRaw, sl.IsRaw)
		return false
	}

	return true
}

func testSymbol(t *testing.T, expr ast.Expr, expected string) bool {
	sl, ok := expr.(*ast.SymLiteral)
	if !ok {
		t.Errorf("sl not *ast.SymLiteral. got=%T", expr)
		return false
	}

	if sl.Value != expected {
		t.Errorf("il.Value not %s. got=%v", expected, sl.Value)
		return false
	}

	return true
}

func testNil(t *testing.T, val interface{}) bool {
	if val != nil {
		t.Errorf("type of value is not nil. got=%T", val)
		return false
	}
	return true
}

func testIntLiteral(t *testing.T, ex ast.Expr, expected int64) bool {
	il, ok := ex.(*ast.IntLiteral)

	if !ok {
		t.Errorf("il not *ast.IntLiteral. got=%T", ex)
		return false
	}

	if il.Value != expected {
		t.Errorf("il.Value not %d. got=%d", expected, il.Value)
		return false
	}

	if il.TokenLiteral() != fmt.Sprintf("%d", expected) {
		t.Errorf("il.TokenLiteral() not %d. got=%s", expected,
			il.TokenLiteral())
		return false
	}

	return true
}

func testParse(t *testing.T, input string) *ast.Program {
	ast, err := Parse(strings.NewReader(input))
	if err != nil {
		msg := fmt.Sprintf("%v\nOccurred in input ```\n%s\n```",
			err.Error(), input)
		t.Fatalf(msg)
		t.FailNow()
	}

	if ast == nil {
		t.Fatalf("ast not generated.")
		t.FailNow()
	}

	return ast
}
