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

func TestEvalZeroAndOneToBuiltIn(t *testing.T) {
	tests := []struct {
		input    string
		expected *object.PanInt
	}{
		{"1", object.BuiltInOneInt},
		{"0", object.BuiltInZeroInt},
	}

	for _, tt := range tests {
		actual := testEval(t, tt.input)
		// check if cached objects (BuiltIn**Int) are used
		if actual != tt.expected {
			t.Errorf("wrong output. expected=%v(%p), got=%v(%p)",
				tt.expected, tt.expected, actual, actual)
		}
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

func TestEvalStrLiteral(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{`"a"`, "a"},
		{"`a`", "a"},
		{`"Hello, world!"`, "Hello, world!"},
		// symbol is also evaluated to PanStr
		{`'sym`, "sym"},
		{`'_hidden`, "_hidden"},
		{`'even?`, "even?"},
		{`'rand!`, "rand!"},
		{`'+`, "+"},
	}

	for _, tt := range tests {
		actual := testEval(t, tt.input)
		expected := &object.PanStr{Value: tt.expected}
		testPanStr(t, actual, expected)
	}
}

func TestEvalBoolLiteral(t *testing.T) {
	tests := []struct {
		input    string
		expected object.PanObject
	}{
		{`true`, object.BuiltInTrue},
		{`false`, object.BuiltInFalse},
	}

	for _, tt := range tests {
		actual := testEval(t, tt.input)
		// check if cached objects (BuiltIn**) are used
		if actual != tt.expected {
			t.Errorf("wrong output. expected=%v(%p), got=%v(%p)",
				tt.expected, tt.expected, actual, actual)
		}
		testPanBool(t, actual, tt.expected.(*object.PanBool))
	}
}

func TestEvalNilLiteral(t *testing.T) {
	actual := testEval(t, "nil")
	if actual != object.BuiltInNil {
		t.Errorf("wrong output. expected=object.BuiltInNil, got=%v", actual)
	}
}

func TestEvalRangeLiteral(t *testing.T) {
	tests := []struct {
		input    string
		expected *object.PanRange
	}{
		{
			`(1:2:3)`,
			toPanRange(1, 2, 3),
		},
		{
			`('a:'z:1)`,
			toPanRange("a", "z", 1),
		},
		// NOTE: parser error occurs in `(::)`
		{
			`(::'step)`,
			toPanRange(nil, nil, "step"),
		},
		{
			`(:'stop)`,
			toPanRange(nil, "stop", nil),
		},
		{
			`(:'stop:'step)`,
			toPanRange(nil, "stop", "step"),
		},
		{
			`('start:)`,
			toPanRange("start", nil, nil),
		},
		{
			`('start::'step)`,
			toPanRange("start", nil, "step"),
		},
		{
			`('start:'stop)`,
			toPanRange("start", "stop", nil),
		},
		{
			`('start:'stop:'step)`,
			toPanRange("start", "stop", "step"),
		},
		// multiple types
		{
			`(3:"s":false)`,
			&object.PanRange{
				Start: &object.PanInt{Value: 3},
				Stop:  &object.PanStr{Value: "s"},
				Step:  object.BuiltInFalse,
			},
		},
	}

	for _, tt := range tests {
		actual := testEval(t, tt.input)
		testPanRange(t, actual, tt.expected)
	}
}

func toPanRange(start, stop, step interface{}) *object.PanRange {
	obj := func(o interface{}) object.PanObject {
		switch o := o.(type) {
		case string:
			return &object.PanStr{Value: o}
		case int:
			return &object.PanInt{Value: int64(o)}
		default:
			return object.BuiltInNil
		}
	}
	return &object.PanRange{Start: obj(start), Stop: obj(stop), Step: obj(step)}
}

func testPanInt(t *testing.T, actual object.PanObject, expected *object.PanInt) {
	if actual == nil {
		t.Fatalf("actual must not be nil. expected=%v(%T)", expected, expected)
	}

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
	if actual == nil {
		t.Fatalf("actual must not be nil. expected=%v(%T)", expected, expected)
	}

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

func testPanStr(t *testing.T, actual object.PanObject, expected *object.PanStr) {
	if actual == nil {
		t.Fatalf("actual must not be nil. expected=%v(%T)", expected, expected)
	}

	if actual.Type() != object.STR_TYPE {
		t.Fatalf("Type must be STR_TYPE. got=%s", actual.Type())
		return
	}

	strObj, ok := actual.(*object.PanStr)
	if !ok {
		t.Fatalf("actual must be *object.PanStr. got=%T (%v)", actual, actual)
		return
	}

	if strObj.Value != expected.Value {
		t.Errorf("wrong value. expected=%s, got=%s", expected.Value, strObj.Value)
	}
}

func testPanBool(t *testing.T, actual object.PanObject, expected *object.PanBool) {
	if actual == nil {
		t.Fatalf("actual must not be nil. expected=%v(%T)", expected, expected)
	}

	if actual.Type() != object.BOOL_TYPE {
		t.Fatalf("Type must be BOOL_TYPE. got=%s", actual.Type())
		return
	}

	boolObj, ok := actual.(*object.PanBool)
	if !ok {
		t.Fatalf("actual must be *object.PanBool. got=%T (%v)", actual, actual)
		return
	}

	if boolObj.Value != expected.Value {
		t.Errorf("wrong value. expected=%t, got=%t", expected.Value, boolObj.Value)
	}
}

func testPanNil(t *testing.T, actual object.PanObject, expected *object.PanNil) {
	if actual == nil {
		t.Fatalf("actual must not be nil. expected=%v(%T)", expected, expected)
	}

	if actual.Type() != object.NIL_TYPE {
		t.Fatalf("Type must be NIL_TYPE. got=%s", actual.Type())
		return
	}

	_, ok := actual.(*object.PanNil)
	if !ok {
		t.Fatalf("actual must be *object.PanNil. got=%T (%v)", actual, actual)
		return
	}
}

func testPanRange(t *testing.T, actual object.PanObject, expected *object.PanRange) {
	if actual == nil {
		t.Fatalf("actual must not be nil. expected=%v(%T)", expected, expected)
	}

	if actual.Type() != object.RANGE_TYPE {
		t.Fatalf("Type must be RANGE_TYPE. got=%s", actual.Type())
		return
	}

	obj, ok := actual.(*object.PanRange)
	if !ok {
		t.Fatalf("actual must be *object.PanRange. got=%T (%v)", actual, actual)
		return
	}

	testValue(t, obj.Start, expected.Start)
	testValue(t, obj.Stop, expected.Stop)
	testValue(t, obj.Step, expected.Step)
}

func testValue(t *testing.T, actual object.PanObject, expected object.PanObject) {
	// switch to test_XX functions by expected type
	switch expected := expected.(type) {
	case *object.PanInt:
		testPanInt(t, actual, expected)
	case *object.PanFloat:
		testPanFloat(t, actual, expected)
	case *object.PanStr:
		testPanStr(t, actual, expected)
	case *object.PanBool:
		testPanBool(t, actual, expected)
	case *object.PanNil:
		testPanNil(t, actual, expected)
	case *object.PanRange:
		testPanRange(t, actual, expected)
	default:
		t.Fatalf("type of expected %T cannot be handled by testValue()", expected)
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
