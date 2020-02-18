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
			t.Fatalf(fmt.Sprintf("%+v", err))
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
