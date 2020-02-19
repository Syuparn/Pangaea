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
		{`3 * 4 / 5`, `((3 * 4) / 5)`},
		{`3 * 4 + 2`, `((3 * 4) + 2)`},
		{`3 + 4 * 2`, `(3 + (4 * 2))`},
		{`3+4*2`, `(3 + (4 * 2))`}, // without space
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
	}

	// TODO: inplement test
	_ = tests
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
			`{||}`,
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

		if f.String() != tt.printed {
			t.Errorf("wrong output.expected=\n%s,\ngot=\n%s",
				tt.printed, f.String())
		}
	}
}

// TODO: check exprs in body
// TODO: check multiple stmts func

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

func testInfixOperator(t *testing.T, expr ast.Expr,
	left int, op string, right int) {
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
		t.Errorf("il not *ast.IntegerLiteral. got=%T", il)
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
	// HACK: catch yacc error by recover
	// (because yacc cannot return error)
	defer func() {
		if err := recover(); err != nil {
			m := fmt.Sprintf("error occured while parsing\n%s", input)
			t.Fatalf(m + "\n" + fmt.Sprintf("%+v", err))
			t.FailNow()
		}
	}()

	ast, err := Parse(strings.NewReader(input))
	if err != nil {
		t.Fatalf(err.Error())
		t.FailNow()
	}
	return ast
}
