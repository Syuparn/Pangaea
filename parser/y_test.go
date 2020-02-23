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
			t.Fatalf("expr is not *ast.SymLiteral.got=%T", sym)
		}

		testSymbol(t, sym, tt.expected)
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
			t.Fatalf("expr is not *ast.ObjLiteral.got=%T", obj)
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
			t.Fatalf("expr is not *ast.ObjLiteral.got=%T", obj)
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
			t.Fatalf("expr is not *ast.ObjLiteral.got=%T", obj)
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
	}

	for _, tt := range tests {
		program := testParse(t, tt.input)
		expr := testIfExprStmt(t, program)

		f, ok := expr.(*ast.FuncLiteral)
		if !ok {
			t.Fatalf("f is not *ast.FuncLiteral. got=%T", expr)
		}

		if len(f.Args) != len(tt.args) {
			t.Fatalf("wrong arity of args, expected=%d, got=%d",
				len(tt.args), len(f.Args))
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
		t.Fatalf(err.Error())
		t.FailNow()
	}

	if ast == nil {
		t.Fatalf("ast not generated.")
		t.FailNow()
	}

	return ast
}
