// Source code in this file is inherited and modified from
// "Writing an Interpreter in Go" https://interpreterbook.com/
// MIT License | Copyright (c) 2016-2017 Thorsten Ball
// see https://opensource.org/licenses/MIT

package evaluator

import (
	"../ast"
	"../object"
	"../parser"
	"fmt"
	"strings"
	"testing"
)

func TestEvalIntLiteral(t *testing.T) {
	tests := []struct {
		input    string
		expected int64
	}{
		{"5", 5},
		{"12345", 12345},
		{"1", 1},
		{"0", 0},
		// special forms
		{`0x10`, 16},
		{`0o10`, 8},
		{`0b10`, 2},
		{`1e3`, 1000},
		{`1_0`, 10},
		// minus values
		{"-5", -5},
		{"-0", 0},
		{`-0x10`, -16},
		{`-0o10`, -8},
		{`-0b10`, -2},
		{`-1e3`, -1000},
		{`-1_0`, -10},
	}

	for _, tt := range tests {
		actual := testEval(t, tt.input)
		expected := &object.PanInt{Value: tt.expected}
		testPanInt(t, actual, expected)
	}
}

func TestEvalFloatLiteral(t *testing.T) {
	tests := []struct {
		input    string
		expected float64
	}{
		{"5.0", 5.0},
		{"12.345", 12.345},
		{"-1.4", -1.4},
		// special forms
		{"1_0.2", 10.2},
		{"1.2e1", 12.0},
		{"1.2e-1", 0.12},
		{".25", 0.25},
		{"-1_0.2", -10.2},
		{"-1.2e1", -12.0},
		{"-1.2e-1", -0.12},
		{"-.25", -0.25},
	}

	for _, tt := range tests {
		actual := testEval(t, tt.input)
		expected := &object.PanFloat{Value: tt.expected}
		testPanFloat(t, actual, expected)
	}
}

func testPanInt(t *testing.T, actual object.PanObject, expected *object.PanInt) {
	if actual.Type() != object.INT_TYPE {
		t.Fatalf("Type must be INT_TYPE. got=%s", actual.Type())
		return
	}

	intObj, ok := actual.(*object.PanInt)
	if !ok {
		t.Fatalf("actual must be *object.PanInt. got=%T (%v)", actual, actual)
		return
	}

	if intObj.Value != expected.Value {
		t.Errorf("wrong value. expected=%d, got=%d", expected.Value, intObj.Value)
	}
}

func testPanFloat(t *testing.T, actual object.PanObject, expected *object.PanFloat) {
	if actual.Type() != object.FLOAT_TYPE {
		t.Fatalf("Type must be FLOAT_TYPE. got=%s", actual.Type())
		return
	}

	floatObj, ok := actual.(*object.PanFloat)
	if !ok {
		t.Fatalf("actual must be *object.PanFloat. got=%T (%v)", actual, actual)
		return
	}

	if floatObj.Value != expected.Value {
		t.Errorf("wrong value. expected=%f, got=%f", expected.Value, floatObj.Value)
	}
}

func testEval(t *testing.T, input string) object.PanObject {
	node := testParse(t, input)
	panObject := Eval(node, object.NewEnv())
	if panObject == nil {
		t.Fatalf("Eval() returned nothing (input=`%s`)", input)
	}
	return panObject
}

func testParse(t *testing.T, input string) *ast.Program {
	node, err := parser.Parse(strings.NewReader(input))
	if err != nil {
		msg := fmt.Sprintf("%v\nOccurred in input ```\n%s\n```",
			err.Error(), input)
		t.Fatalf(msg)
		t.FailNow()
	}

	if node == nil {
		t.Fatalf("ast not generated.")
		t.FailNow()
	}

	return node
}