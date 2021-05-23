// Source code in this file is inherited and modified from
// "Writing an Interpreter in Go" https://interpreterbook.com/
// MIT License | Copyright (c) 2016-2017 Thorsten Ball
// see https://opensource.org/licenses/MIT

package evaluator

import (
	"bytes"
	"fmt"
	"os"
	"strings"
	"testing"

	"github.com/Syuparn/pangaea/ast"
	"github.com/Syuparn/pangaea/object"
	"github.com/Syuparn/pangaea/parser"
	"github.com/Syuparn/pangaea/props"
)

func TestMain(m *testing.M) {
	// setup for name resolution
	ctn := NewPropContainer()
	// inject stubs
	ctn["Str_eval"] = object.NewNotImplementedErr("not implemented in evaluator")
	ctn["Str_evalEnv"] = object.NewNotImplementedErr("not implemented in evaluator")
	injectBuiltInProps(ctn)
	ret := m.Run()
	os.Exit(ret)
}

func injectBuiltInProps(ctn map[string]object.PanObject) {
	injectProps(object.BuiltInAssertionErr, props.AssertionErrProps, ctn)
	injectProps(object.BuiltInArrObj, props.ArrProps, ctn)
	injectProps(object.BuiltInBaseObj, props.BaseObjProps, ctn)
	injectProps(object.BuiltInComparableObj, props.ComparableProps, ctn)
	injectProps(object.BuiltInDiamondObj, props.DiamondProps, ctn)
	injectProps(object.BuiltInEitherObj, props.EitherProps, ctn)
	injectProps(object.BuiltInEitherValObj, props.EitherValProps, ctn)
	injectProps(object.BuiltInEitherErrObj, props.EitherErrProps, ctn)
	injectProps(object.BuiltInErrObj, props.ErrProps, ctn)
	injectProps(object.BuiltInFloatObj, props.FloatProps, ctn)
	injectProps(object.BuiltInFuncObj, props.FuncProps, ctn)
	injectProps(object.BuiltInIntObj, props.IntProps, ctn)
	injectProps(object.BuiltInIterObj, props.IterProps, ctn)
	injectProps(object.BuiltInIterableObj, props.IterableProps, ctn)
	injectProps(object.BuiltInKernelObj, props.KernelProps, ctn)
	injectProps(object.BuiltInMatchObj, props.MatchProps, ctn)
	injectProps(object.BuiltInMapObj, props.MapProps, ctn)
	injectProps(object.BuiltInNameErr, props.NameErrProps, ctn)
	injectProps(object.BuiltInNilObj, props.NilProps, ctn)
	injectProps(object.BuiltInNoPropErr, props.NoPropErrProps, ctn)
	injectProps(object.BuiltInNotImplementedErr, props.NotImplementedErrProps, ctn)
	injectProps(object.BuiltInNumObj, props.NumProps, ctn)
	injectProps(object.BuiltInObjObj, props.ObjProps, ctn)
	injectProps(object.BuiltInRangeObj, props.RangeProps, ctn)
	injectProps(object.BuiltInStopIterErr, props.StopIterErrProps, ctn)
	injectProps(object.BuiltInStrObj, props.StrProps, ctn)
	injectProps(object.BuiltInSyntaxErr, props.SyntaxErrProps, ctn)
	injectProps(object.BuiltInTypeErr, props.TypeErrProps, ctn)
	injectProps(object.BuiltInValueErr, props.ValueErrProps, ctn)
	injectProps(object.BuiltInWrappableObj, props.WrappableProps, ctn)
	injectProps(object.BuiltInZeroDivisionErr, props.ZeroDivisionErrProps, ctn)
}

func injectProps(
	obj *object.PanObj,
	props func(map[string]object.PanObject) map[string]object.PanObject,
	propContainer map[string]object.PanObject,
) {
	pairs := map[object.SymHash]object.Pair{}
	for k, v := range props(propContainer) {
		pair := object.Pair{
			Key:   object.NewPanStr(k),
			Value: v,
		}
		pairs[object.GetSymHash(k)] = pair
	}

	obj.AddPairs(&pairs)
}

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
		expected := object.NewPanInt(tt.expected)
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

func TestEvalIntAt(t *testing.T) {
	tests := []struct {
		input    string
		expected object.PanObject
	}{
		// n[i] returns ith bit of n
		{`5[0]`, object.NewPanInt(1)},
		{`5[1]`, object.NewPanInt(0)},
		{`5[2]`, object.NewPanInt(1)},
		{`5[3]`, object.NewPanInt(0)},
		{
			`6[0:3]`,
			object.NewPanArr(
				object.NewPanInt(0),
				object.NewPanInt(1),
				object.NewPanInt(1),
			),
		},
		// use descendant of arr for arg
		{
			`5.at([0].bear)`,
			object.NewPanInt(1),
		},
		// use descendant of int for index
		{
			`5[false]`,
			object.NewPanInt(1),
		},
		// use descendant of range for index
		{
			`5[(0:2).bear]`,
			object.NewPanArr(
				object.NewPanInt(1),
				object.NewPanInt(0),
			),
		},
		// if \1 is not int, return nil
		{
			`Int['at]({}, [1])`,
			object.BuiltInNil,
		},
		// if \2 is not arr, return nil
		{
			`1.at({})`,
			object.BuiltInNil,
		},
		// if key is insufficient, return nil
		{
			`1[]`,
			object.BuiltInNil,
		},
		// if index is out of range, return nil
		{
			`1[100]`,
			object.BuiltInNil,
		},
		// minus index
		{
			`1[-1]`,
			object.NewPanInt(0),
		},
		{
			`1[-100]`,
			object.BuiltInNil,
		},
		// if non-int value is passed, call parent's at method
		{
			`1['at]`,
			(*object.BuiltInIntObj.Pairs)[object.GetSymHash("at")].Value,
		},
	}

	for _, tt := range tests {
		actual := testEval(t, tt.input)
		testValue(t, actual, tt.expected)
	}
}

func TestEvalIntPrimep(t *testing.T) {
	tests := []struct {
		input    string
		expected object.PanObject
	}{
		{
			`7.prime?`,
			object.BuiltInTrue,
		},
		{
			`12.prime?`,
			object.BuiltInFalse,
		},
		{
			`-5.prime?`,
			object.BuiltInFalse,
		},
		// use descendant of int for recv
		{
			`true.prime?`,
			object.BuiltInFalse,
		},
		// if no args are passed, raise an error
		{
			`Int['prime?]()`,
			object.NewTypeErr("Int#prime? requires at least 1 arg"),
		},
		// if \1 is not int, raise an error
		{
			`Int['prime?]("a")`,
			object.NewTypeErr("\\1 must be int"),
		},
	}

	for _, tt := range tests {
		actual := testEval(t, tt.input)
		testValue(t, actual, tt.expected)
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
		expected := object.NewPanFloat(tt.expected)
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
		expected := object.NewPanStr(tt.expected)
		testPanStr(t, actual, expected)
	}
}

func TestEvalEmbeddedStr(t *testing.T) {
	tests := []struct {
		input    string
		expected object.PanObject
	}{
		{
			`"Hello, #{ "world" }!"`,
			object.NewPanStr("Hello, world!"),
		},
		// .S is called internal
		{
			`"1 + 1 = #{1 + 1}"`,
			object.NewPanStr("1 + 1 = 2"),
		},
		{
			`"arr: #{ [1, 2, "3"] }"`,
			object.NewPanStr(`arr: [1, 2, "3"]`),
		},
		// NOTE: currently obj and func literal cannot be embedded...
		{
			`obj := {a: 1}; "three: #{ obj.a + 2 if true }"`,
			object.NewPanStr(`three: 3`),
		},
		// multiple embedding
		{
			`"a#{'bbc[1:]}d#{'eiffel[0:3:2]}g"`,
			object.NewPanStr(`abcdefg`),
		},
		// embedding succeeds even if embedded.S is descendant of str
		{
			`{S: "elem".bear()}.{|e| "embedded: #{e}"}`,
			object.NewPanStr(`embedded: elem`),
		},
	}

	for _, tt := range tests {
		actual := testEval(t, tt.input)
		testValue(t, actual, tt.expected)
	}
}

func TestEvalStrAt(t *testing.T) {
	tests := []struct {
		input    string
		expected object.PanObject
	}{
		// s[i] returns ith rune of s
		{`'ab[0]`, object.NewPanStr("a")},
		{`'ab[1]`, object.NewPanStr("b")},
		{`'ab[2]`, object.BuiltInNil},
		{`'abcde[0:3]`, object.NewPanStr("abc")},
		// non-ascii string
		{`"日本語"[1]`, object.NewPanStr("本")},
		// use descendant of arr for arg
		{
			`'abc.at([0].bear)`,
			object.NewPanStr("a"),
		},
		// use descendant of int for index
		{
			`'abc[false]`,
			object.NewPanStr("a"),
		},
		// use descendant of range for index
		{
			`'abc[(0:2).bear]`,
			object.NewPanStr("ab"),
		},
		// if \1 is not str, return nil
		{
			`Str['at]({}, [1])`,
			object.BuiltInNil,
		},
		// if \2 is not arr, return nil
		{
			`"".at({})`,
			object.BuiltInNil,
		},
		// if key is insufficient, return nil
		{
			`""[]`,
			object.BuiltInNil,
		},
		// if index is out of range, return nil
		{
			`""[100]`,
			object.BuiltInNil,
		},
		// minus index
		{
			`"abc"[-1]`,
			object.NewPanStr("c"),
		},
		{
			`"abc"[-4]`,
			object.BuiltInNil,
		},
		// if non-int value is passed, call parent's at method
		{
			`""['at]`,
			(*object.BuiltInStrObj.Pairs)[object.GetSymHash("at")].Value,
		},
	}

	for _, tt := range tests {
		actual := testEval(t, tt.input)
		testValue(t, actual, tt.expected)
	}
}

func TestEvalStrLen(t *testing.T) {
	tests := []struct {
		input    string
		expected object.PanObject
	}{
		{
			`"".len`,
			object.NewPanInt(0),
		},
		{
			`"abc".len`,
			object.NewPanInt(3),
		},
		// length is counted by runes
		{
			`"ABCあいうえお".len`,
			object.NewPanInt(8),
		},
		// use descendant of str for recv
		{
			`"a".bear.len`,
			object.NewPanInt(1),
		},
		// if no args are passed, raise an error
		{
			`Str['len]()`,
			object.NewTypeErr("Str#len requires at least 1 arg"),
		},
		// if \1 is not str, raise an error
		{
			`Str['len](1)`,
			object.NewTypeErr("\\1 must be str"),
		},
	}

	for _, tt := range tests {
		actual := testEval(t, tt.input)
		testValue(t, actual, tt.expected)
	}
}

func TestEvalStrLc(t *testing.T) {
	tests := []struct {
		input    string
		expected object.PanObject
	}{
		{
			`"AbC".lc`,
			object.NewPanStr("abc"),
		},
		// non-alphabetical characters are ignored
		{
			`"AbC日本語".lc`,
			object.NewPanStr("abc日本語"),
		},
		// use descendant of str for recv
		{
			`"A".bear.lc`,
			object.NewPanStr("a"),
		},
		// if no args are passed, raise an error
		{
			`Str['lc]()`,
			object.NewTypeErr("Str#lc requires at least 1 arg"),
		},
		// if \1 is not str, raise an error
		{
			`Str['lc](1)`,
			object.NewTypeErr("\\1 must be str"),
		},
	}

	for _, tt := range tests {
		actual := testEval(t, tt.input)
		testValue(t, actual, tt.expected)
	}
}

func TestEvalStrMatch(t *testing.T) {
	tests := []struct {
		input    string
		expected object.PanObject
	}{
		{
			"`hello`.match(`l+`)",
			object.NewPanArr(
				object.NewPanStr("ll"),
			),
		},
		{
			"`--aabcbcbcxyyyzz--`.match(`a+((bc)*)x([yz]+)`)",
			object.NewPanArr(
				object.NewPanStr("aabcbcbcxyyyzz"),
				object.NewPanStr("bcbcbc"),
				object.NewPanStr("bc"),
				object.NewPanStr("yyyzz"),
			),
		},
		// no match
		{
			"`a`.match(`bc`)",
			object.NewPanArr(),
		},
		// use descendant of str for recv
		{
			"`abbbc`.bear.match(`b+`)",
			object.NewPanArr(
				object.NewPanStr("bbb"),
			),
		},
		// if no args are passed, raise an error
		{
			`Str.match`,
			object.NewTypeErr("Str#match requires at least 2 args"),
		},
		// if \1 is not str, raise an error
		{
			`Str['match](1, "a")`,
			object.NewTypeErr("\\1 must be str"),
		},
		// if \2 is not str, raise an error
		{
			`"a".match(1)`,
			object.NewTypeErr("\\2 must be str"),
		},
		// if \2 is illegal for regex syntax, raise an error
		{
			`"a".match("[")`,
			object.NewValueErr(`"[" is invalid regex pattern`),
		},
	}

	for _, tt := range tests {
		actual := testEval(t, tt.input)
		testValue(t, actual, tt.expected)
	}
}

func TestEvalStrSub(t *testing.T) {
	tests := []struct {
		input    string
		expected object.PanObject
	}{
		{
			`"rice".sub("r", "m")`,
			object.NewPanStr("mice"),
		},
		{
			"`secret: 1234`.sub(`\\d`, `x`)",
			object.NewPanStr("secret: xxxx"),
		},
		{
			`"abbc".sub("(b+)", "[$1]")`,
			object.NewPanStr("a[bb]c"),
		},
		// use descendant of str for recv
		{
			`"hoge".bear.sub("o", "i")`,
			object.NewPanStr("hige"),
		},
		// convert to lowercase
		{
			"`PANGAEA`.sub(`(A|E)`, `\\L$1\\E`)",
			object.NewPanStr("PaNGaea"),
		},
		// convert to uppercase
		{
			"`pangaea`.sub(`(a|e)`, `\\U$1\\E`)",
			object.NewPanStr("pAngAEA"),
		},
		// use both of \U and \L
		{
			"`SIng`.sub(`(Ing)`, `\\U$1\\E\\L$1\\E`)",
			object.NewPanStr("SINGing"),
		},
		// if no args are passed, raise an error
		{
			`Str.sub("a")`,
			object.NewTypeErr("Str#sub requires at least 3 args"),
		},
		// if \1 is not str, raise an error
		{
			`Str['sub](1, "a", "b")`,
			object.NewTypeErr("\\1 must be str"),
		},
		// if \2 is not str, raise an error
		{
			`"a".sub(2, "b")`,
			object.NewTypeErr("\\2 must be str"),
		},
		// if \3 is not str, raise an error
		{
			`"a".sub("b", 3)`,
			object.NewTypeErr("\\3 must be str"),
		},
		// if \2 is illegal for regex syntax, raise an error
		{
			`"a".sub("[", "b")`,
			object.NewValueErr(`"[" is invalid regex pattern`),
		},
		// if \3 is illegal for regex syntax, raise an error
		{
			`"a".sub("a", "$999999999999")`,
			object.NewValueErr(`"$999999999999" is invalid regex pattern`),
		},
	}

	for _, tt := range tests {
		actual := testEval(t, tt.input)
		testValue(t, actual, tt.expected)
	}
}

func TestEvalStrSymp(t *testing.T) {
	tests := []struct {
		input    string
		expected object.PanObject
	}{
		// symbol
		{
			`'mySymbol.sym?`,
			object.BuiltInTrue,
		},
		// able to make from symbol
		{
			`"var".sym?`,
			object.BuiltInTrue,
		},
		{
			`"_private".sym?`,
			object.BuiltInTrue,
		},
		{
			`?+.sym?`,
			object.BuiltInTrue,
		},
		// cannot be used as symbol
		{
			`"hello world".sym?`,
			object.BuiltInFalse,
		},
		{
			`"にほんご".sym?`,
			object.BuiltInFalse,
		},
		{
			`Str['sym?]()`,
			object.NewTypeErr("Str#sym? requires at least 1 arg"),
		},
		// if \1 is not str, raise an error
		{
			`Str['sym?](1)`,
			object.NewTypeErr("1 cannot be treated as str"),
		},
	}

	for _, tt := range tests {
		actual := testEval(t, tt.input)
		testValue(t, actual, tt.expected)
	}
}

func TestEvalStrUc(t *testing.T) {
	tests := []struct {
		input    string
		expected object.PanObject
	}{
		{
			`"AbC".uc`,
			object.NewPanStr("ABC"),
		},
		// non-alphabetical characters are ignored
		{
			`"AbC日本語".uc`,
			object.NewPanStr("ABC日本語"),
		},
		// use descendant of str for recv
		{
			`"a".bear.uc`,
			object.NewPanStr("A"),
		},
		// if no args are passed, raise an error
		{
			`Str['uc]()`,
			object.NewTypeErr("Str#uc requires at least 1 arg"),
		},
		// if \1 is not str, raise an error
		{
			`Str['uc](1)`,
			object.NewTypeErr("\\1 must be str"),
		},
	}

	for _, tt := range tests {
		actual := testEval(t, tt.input)
		testValue(t, actual, tt.expected)
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
			object.NewPanRange(
				object.NewPanInt(3),
				object.NewPanStr("s"),
				object.BuiltInFalse,
			),
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
			return object.NewPanStr(o)
		case int:
			return object.NewPanInt(int64(o))
		default:
			return object.BuiltInNil
		}
	}
	return object.NewPanRange(obj(start), obj(stop), obj(step))
}

func TestEvalRangeStart(t *testing.T) {
	tests := []struct {
		input    string
		expected object.PanObject
	}{
		{
			`(1:2:3).start`,
			object.NewPanInt(1),
		},
		{
			`(::'step).start`,
			object.BuiltInNil,
		},
		{
			`Range['start]()`,
			object.NewTypeErr("Range#start requires at least 1 arg"),
		},
		{
			`Range['start](1)`,
			object.NewTypeErr("1 cannot be treated as range"),
		},
	}

	for _, tt := range tests {
		actual := testEval(t, tt.input)
		testValue(t, actual, tt.expected)
	}
}

func TestEvalRangeStop(t *testing.T) {
	tests := []struct {
		input    string
		expected object.PanObject
	}{
		{
			`(1:2:3).stop`,
			object.NewPanInt(2),
		},
		{
			`('start::'step).stop`,
			object.BuiltInNil,
		},
		{
			`Range['stop]()`,
			object.NewTypeErr("Range#stop requires at least 1 arg"),
		},
		{
			`Range['stop](1)`,
			object.NewTypeErr("1 cannot be treated as range"),
		},
	}

	for _, tt := range tests {
		actual := testEval(t, tt.input)
		testValue(t, actual, tt.expected)
	}
}

func TestEvalRangeStep(t *testing.T) {
	tests := []struct {
		input    string
		expected object.PanObject
	}{
		{
			`(1:2:3).step`,
			object.NewPanInt(3),
		},
		{
			`('start:'stop).step`,
			object.BuiltInNil,
		},
		{
			`Range['step]()`,
			object.NewTypeErr("Range#step requires at least 1 arg"),
		},
		{
			`Range['step](1)`,
			object.NewTypeErr("1 cannot be treated as range"),
		},
	}

	for _, tt := range tests {
		actual := testEval(t, tt.input)
		testValue(t, actual, tt.expected)
	}
}

func TestEvalArrLiteral(t *testing.T) {
	tests := []struct {
		input    string
		expected *object.PanArr
	}{
		{
			`[]`,
			object.NewPanArr(),
		},
		{
			`[1]`,
			object.NewPanArr(
				object.NewPanInt(1),
			),
		},
		{
			`[2, 3]`,
			object.NewPanArr(
				object.NewPanInt(2),
				object.NewPanInt(3),
			),
		},
		// arr can contain different type elements
		{
			`["a", 4]`,
			object.NewPanArr(
				object.NewPanStr("a"),
				object.NewPanInt(4),
			),
		},
		// nested
		{
			`[[10]]`,
			object.NewPanArr(
				object.NewPanArr(
					object.NewPanInt(10),
				),
			),
		},
		// embedded
		{
			`[*[1]]`,
			object.NewPanArr(
				object.NewPanInt(1),
			),
		},
		{
			`[*[1], 2]`,
			object.NewPanArr(
				object.NewPanInt(1),
				object.NewPanInt(2),
			),
		},
		{
			`[1, *[2, 3], 4]`,
			object.NewPanArr(
				object.NewPanInt(1),
				object.NewPanInt(2),
				object.NewPanInt(3),
				object.NewPanInt(4),
			),
		},
		// embed descendant of arr (arr proto is embedded)
		{
			`[1, *([2].bear)]`,
			object.NewPanArr(
				object.NewPanInt(1),
				object.NewPanInt(2),
			),
		},
	}

	for _, tt := range tests {
		actual := testEval(t, tt.input)
		testPanArr(t, actual, tt.expected)
	}
}

func TestEvalArrAt(t *testing.T) {
	tests := []struct {
		input    string
		expected object.PanObject
	}{
		{
			`[5, 10][0]`,
			object.NewPanInt(5),
		},
		// if key is insufficient, return nil
		{
			`[5, 10][]`,
			object.BuiltInNil,
		},
		// if index is out of range, return nil
		{
			`[1, 2][100]`,
			object.BuiltInNil,
		},
		// if non-int value is passed, call parent's at method
		{
			`[1]['at]`,
			(*object.BuiltInArrObj.Pairs)[object.GetSymHash("at")].Value,
		},
		// minus index
		{
			`[1, 2, 3][-1]`,
			object.NewPanInt(3),
		},
		{
			`[1, 2, 3][-3]`,
			object.NewPanInt(1),
		},
		{
			`[1, 2, 3][-4]`,
			object.BuiltInNil,
		},
		// use descendant of arr for recv
		{
			`[1].bear[0]`,
			object.NewPanInt(1),
		},
		// use descendant of arr for arg
		{
			`[1].at([0].bear)`,
			object.NewPanInt(1),
		},
		// use descendant of int for index
		{
			`[1][false]`,
			object.NewPanInt(1),
		},
		// if \1 is not arr, return nil
		{
			`Arr['at]({}, [1])`,
			object.BuiltInNil,
		},
		// if \2 is not arr, return nil
		{
			`[1, 2, 3].at({})`,
			object.BuiltInNil,
		},
	}

	for _, tt := range tests {
		actual := testEval(t, tt.input)
		testValue(t, actual, tt.expected)
	}
}

func TestEvalArrAtWithRange(t *testing.T) {
	tests := []struct {
		input    string
		expected object.PanObject
	}{
		{
			`[0, 1, 2][0:1]`,
			object.NewPanArr(
				object.NewPanInt(0),
			),
		},
		{
			`[0, 1, 2][0:]`,
			object.NewPanArr(
				object.NewPanInt(0),
				object.NewPanInt(1),
				object.NewPanInt(2),
			),
		},
		{
			`[0, 1, 2][:2]`,
			object.NewPanArr(
				object.NewPanInt(0),
				object.NewPanInt(1),
			),
		},
		{
			`[0, 1, 2][::-1]`,
			object.NewPanArr(
				object.NewPanInt(2),
				object.NewPanInt(1),
				object.NewPanInt(0),
			),
		},
		{
			`[0, 1, 2][1::-1]`,
			object.NewPanArr(
				object.NewPanInt(1),
				object.NewPanInt(0),
			),
		},
		{
			`[0, 1, 2, 3, 4][:2:-1]`,
			object.NewPanArr(
				object.NewPanInt(4),
				object.NewPanInt(3),
			),
		},
		{
			`[0, 1, 2, 3, 4, 5][1::2]`,
			object.NewPanArr(
				object.NewPanInt(1),
				object.NewPanInt(3),
				object.NewPanInt(5),
			),
		},
		{
			`[0, 1, 2, 3, 4][:3:2]`,
			object.NewPanArr(
				object.NewPanInt(0),
				object.NewPanInt(2),
			),
		},
		{
			`[0, 1, 2, 3, 4, 5][1:5:2]`,
			object.NewPanArr(
				object.NewPanInt(1),
				object.NewPanInt(3),
			),
		},
		// use descendant of range for index
		{
			`[1, 2, 3][(0:2).bear]`,
			object.NewPanArr(
				object.NewPanInt(1),
				object.NewPanInt(2),
			),
		},
		// range is fixed even if start or step is out of range
		{
			`[0, 1, 2][-3:-2]`,
			object.NewPanArr(
				object.NewPanInt(0),
			),
		},
		{
			`[0, 1, 2][-100:2]`,
			object.NewPanArr(
				object.NewPanInt(0),
				object.NewPanInt(1),
			),
		},
		{
			`[0, 1, 2][-2:10000]`,
			object.NewPanArr(
				object.NewPanInt(1),
				object.NewPanInt(2),
			),
		},
	}

	for _, tt := range tests {
		actual := testEval(t, tt.input)
		testValue(t, actual, tt.expected)
	}
}

func TestEvalArrHasp(t *testing.T) {
	tests := []struct {
		input    string
		expected object.PanObject
	}{
		{
			`[1, 2].has?(2)`,
			object.BuiltInTrue,
		},
		{
			`[1].has?(2)`,
			object.BuiltInFalse,
		},
		// non-scalar can be used
		{
			`[[1, 2, 3], [4, 5]].has?([1, 2, 3])`,
			object.BuiltInTrue,
		},
		// if no args are passed, raise an error
		{
			`[].has?`,
			object.NewTypeErr("Arr#has? requires at least 2 args"),
		},
		// if \1 is not Arr, raise an error
		{
			`Arr['has?](1, 1)`,
			object.NewTypeErr("\\1 must be arr"),
		},
	}

	for _, tt := range tests {
		actual := testEval(t, tt.input)
		testValue(t, actual, tt.expected)
	}
}

func TestEvalArrJoin(t *testing.T) {
	tests := []struct {
		input    string
		expected object.PanObject
	}{
		{
			`["hello", "world"].join(" ")`,
			object.NewPanStr("hello world"),
		},
		// non-str elems are stringified by .S
		{
			`[1, [2, 3], true, "five"].join("&")`,
			object.NewPanStr("1&[2, 3]&true&five"),
		},
		// if arity is insufficient, raise an error
		{
			`[].join`,
			object.NewTypeErr("Arr#join requires at least 2 args"),
		},
		// if \1 is not arr, raise an error
		{
			`Arr['join](1, " ")`,
			object.NewTypeErr("\\1 must be arr"),
		},
		// if \2 is not str, raise an error
		{
			`['a, 'b].join(1)`,
			object.NewTypeErr("\\2 must be str"),
		},
		// if (elem).S is not str, raise an error
		{
			`['a, {S: 1}].join(" ")`,
			object.NewTypeErr(`{"S": 1}.S returned non-str value 1`),
		},
	}

	for _, tt := range tests {
		actual := testEval(t, tt.input)
		testValue(t, actual, tt.expected)
	}
}

func TestEvalArrLen(t *testing.T) {
	tests := []struct {
		input    string
		expected object.PanObject
	}{
		{
			`[].len`,
			object.NewPanInt(0),
		},
		{
			`[5, 10].len`,
			object.NewPanInt(2),
		},
		// use descendant of arr for recv
		{
			`[1].bear.len`,
			object.NewPanInt(1),
		},
		// if no args are passed, raise an error
		{
			`Arr['len]()`,
			object.NewTypeErr("Arr#len requires at least 1 arg"),
		},
		// if \1 is not arr, raise an error
		{
			`Arr['len](1)`,
			object.NewTypeErr("\\1 must be arr"),
		},
	}

	for _, tt := range tests {
		actual := testEval(t, tt.input)
		testValue(t, actual, tt.expected)
	}
}

func TestEvalObjLiteral(t *testing.T) {
	tests := []struct {
		input    string
		expected *object.PanObj
	}{
		{
			`{}`,
			toPanObj([]object.Pair{}),
		},
		{
			`{a: 1}`,
			toPanObj([]object.Pair{
				{
					Key:   object.NewPanStr("a"),
					Value: object.NewPanInt(1),
				},
			}),
		},
		{
			`{a: 1, b: 2}`,
			toPanObj([]object.Pair{
				{
					Key:   object.NewPanStr("a"),
					Value: object.NewPanInt(1),
				},
				{
					Key:   object.NewPanStr("b"),
					Value: object.NewPanInt(2),
				},
			}),
		},
		// dangling keys are ignored (in this example, `a: 3` is ignored)
		{
			`{a: 1, b: 2, a: 3}`,
			toPanObj([]object.Pair{
				{
					Key:   object.NewPanStr("a"),
					Value: object.NewPanInt(1),
				},
				{
					Key:   object.NewPanStr("b"),
					Value: object.NewPanInt(2),
				},
			}),
		},
		// Obj can contain multiple types
		{
			`{a: 1, b: "B"}`,
			toPanObj([]object.Pair{
				{
					Key:   object.NewPanStr("a"),
					Value: object.NewPanInt(1),
				},
				{
					Key:   object.NewPanStr("b"),
					Value: object.NewPanStr("B"),
				},
			}),
		},
		// nested
		{
			`{a: {b: 10}}`,
			toPanObj([]object.Pair{
				{
					Key: object.NewPanStr("a"),
					Value: toPanObj([]object.Pair{
						{
							Key:   object.NewPanStr("b"),
							Value: object.NewPanInt(10),
						},
					}),
				},
			}),
		},
		// embedded
		// NOTE: embedded elems must be after normal pairs (`{**a, b: 2}` is syntax error)
		{
			`{a: 1, **{b: 2}}`,
			toPanObj([]object.Pair{
				{
					Key:   object.NewPanStr("a"),
					Value: object.NewPanInt(1),
				},
				{
					Key:   object.NewPanStr("b"),
					Value: object.NewPanInt(2),
				},
			}),
		},
		{
			`{**{a: 1}, **{b: 2}}`,
			toPanObj([]object.Pair{
				{
					Key:   object.NewPanStr("a"),
					Value: object.NewPanInt(1),
				},
				{
					Key:   object.NewPanStr("b"),
					Value: object.NewPanInt(2),
				},
			}),
		},
		// enable to use descendant of str for key (str proto is set)
		{
			`{"hoge".bear: 1}`,
			toPanObj([]object.Pair{
				{
					Key:   object.NewPanStr("hoge"),
					Value: object.NewPanInt(1),
				},
			}),
		},
	}

	for _, tt := range tests {
		actual := testEval(t, tt.input)
		testPanObj(t, actual, tt.expected)
	}
}

func TestEvalPinnedObjKey(t *testing.T) {
	tests := []struct {
		input    string
		expected *object.PanObj
	}{
		{
			`a := "A"; {^a: 1}`,
			toPanObj([]object.Pair{
				{
					Key:   object.NewPanStr("A"),
					Value: object.NewPanInt(1),
				},
			}),
		},
		// enable to use descendant of str for key (str proto is set)
		{
			`b := "B".bear; {^b: 1}`,
			toPanObj([]object.Pair{
				{
					Key:   object.NewPanStr("B"),
					Value: object.NewPanInt(1),
				},
			}),
		},
	}

	for _, tt := range tests {
		actual := testEval(t, tt.input)
		testPanObj(t, actual, tt.expected)
	}
}

func TestEvalBuiltInCallProp(t *testing.T) {
	tests := []struct {
		input    string
		expected object.PanObject
	}{
		// non-callable
		{
			`{}.callProp({a: 1}, 'a)`,
			object.NewPanInt(1),
		},
		// callable
		{
			`{}.callProp({a: {|| 2}}, 'a)`,
			object.NewPanInt(2),
		},
		// NOTE: first arg (`self`) is reciever itself! (`{a: m{|x| x}}`)
		{
			`{}.callProp({a: m{|x| x}}, 'a, 3)`,
			object.NewPanInt(3),
		},
		{
			`{}.callProp({a: m{|x, y, z: 1| [x, y, z]}}, 'a, 4, 5, z: 6)`,
			object.NewPanArr(
				object.NewPanInt(4),
				object.NewPanInt(5),
				object.NewPanInt(6),
			),
		},
	}

	for _, tt := range tests {
		actual := testEval(t, tt.input)
		testValue(t, actual, tt.expected)
	}
}

func TestEvalObjLiteralInvalidKey(t *testing.T) {
	tests := []struct {
		input    string
		expected *object.PanErr
	}{
		{
			`{1: 2}`,
			object.NewTypeErr("cannot use `1` as Obj key."),
		},
		{
			`{a: 1, 2: 3}`,
			object.NewTypeErr("cannot use `2` as Obj key."),
		},
	}

	for _, tt := range tests {
		actual := testEval(t, tt.input)
		testPanErr(t, actual, tt.expected)
	}
}

func toPanObj(pairs []object.Pair) *object.PanObj {
	pairMap := map[object.SymHash]object.Pair{}

	for _, pair := range pairs {
		panStr, _ := pair.Key.(*object.PanStr)
		symHash := object.GetSymHash(panStr.Value)
		pairMap[symHash] = pair
	}

	obj := object.PanObjInstance(&pairMap)
	return &obj
}

func TestEvalMapLiteral(t *testing.T) {
	tests := []struct {
		input    string
		expected *object.PanMap
	}{
		{
			`%{}`,
			toPanMap([]object.Pair{}, []object.Pair{}),
		},
		{
			`%{'a: 1}`,
			toPanMap(
				[]object.Pair{
					{
						Key:   object.NewPanStr("a"),
						Value: object.NewPanInt(1),
					},
				},
				[]object.Pair{},
			),
		},
		// map can contain keys other than str
		{
			`%{'a: 1, 2: 3, true: 5}`,
			toPanMap(
				[]object.Pair{
					{
						Key:   object.NewPanStr("a"),
						Value: object.NewPanInt(1),
					},
					{
						Key:   object.NewPanInt(2),
						Value: object.NewPanInt(3),
					},
					{
						Key:   object.BuiltInTrue,
						Value: object.NewPanInt(5),
					},
				},
				[]object.Pair{},
			),
		},
		// dangling keys are ignored (in this example, `a: 3` is ignored)
		{
			`%{'a: 1, 'b: 2, 'a: 3}`,
			toPanMap(
				[]object.Pair{
					{
						Key:   object.NewPanStr("a"),
						Value: object.NewPanInt(1),
					},
					{
						Key:   object.NewPanStr("b"),
						Value: object.NewPanInt(2),
					},
				},
				[]object.Pair{},
			),
		},
		// Obj can contain multiple types
		{
			`%{"a": 1, "b": "B"}`,
			toPanMap(
				[]object.Pair{
					{
						Key:   object.NewPanStr("a"),
						Value: object.NewPanInt(1),
					},
					{
						Key:   object.NewPanStr("b"),
						Value: object.NewPanStr("B"),
					},
				},
				[]object.Pair{},
			),
		},
		// nested
		{
			`%{'a: %{'b: 10}}`,
			toPanMap(
				[]object.Pair{
					{
						Key: object.NewPanStr("a"),
						Value: toPanMap(
							[]object.Pair{
								{
									Key:   object.NewPanStr("b"),
									Value: object.NewPanInt(10),
								},
							},
							[]object.Pair{},
						),
					},
				},
				[]object.Pair{},
			),
		},
		// embedded
		// NOTE: embedded elems must be after normal pairs
		// (`%{**a, b: 2}` is syntax error)
		{
			`%{'a: 1, **%{'b: 2}}`,
			toPanMap(
				[]object.Pair{
					{
						Key:   object.NewPanStr("a"),
						Value: object.NewPanInt(1),
					},
					{
						Key:   object.NewPanStr("b"),
						Value: object.NewPanInt(2),
					},
				},
				[]object.Pair{},
			),
		},
		// embedded non-scalar
		{
			`%{'a: 1, **%{[]: 2}}`,
			toPanMap(
				[]object.Pair{
					{
						Key:   object.NewPanStr("a"),
						Value: object.NewPanInt(1),
					},
				},
				[]object.Pair{
					{
						Key:   object.NewPanArr(),
						Value: object.NewPanInt(2),
					},
				},
			),
		},
		{
			`%{**%{'a: 1}, **%{'b: 2}}`,
			toPanMap(
				[]object.Pair{
					{
						Key:   object.NewPanStr("a"),
						Value: object.NewPanInt(1),
					},
					{
						Key:   object.NewPanStr("b"),
						Value: object.NewPanInt(2),
					},
				},
				[]object.Pair{},
			),
		},
		// obj can be embedded to map (but map cannot be embedded to obj)
		{
			`%{**%{'a: 1}, **{b: 2}}`,
			toPanMap(
				[]object.Pair{
					{
						Key:   object.NewPanStr("a"),
						Value: object.NewPanInt(1),
					},
					{
						Key:   object.NewPanStr("b"),
						Value: object.NewPanInt(2),
					},
				},
				[]object.Pair{},
			),
		},
		// map can contain non-scalar keys
		{
			`%{'a: 1, [1, 2]: 3, %{4: 5}: 6}`,
			toPanMap(
				[]object.Pair{
					{
						Key:   object.NewPanStr("a"),
						Value: object.NewPanInt(1),
					},
				},
				[]object.Pair{
					{
						Key: object.NewPanArr(
							object.NewPanInt(1),
							object.NewPanInt(2),
						),
						Value: object.NewPanInt(3),
					},
					{
						Key: toPanMap(
							[]object.Pair{
								{
									Key:   object.NewPanInt(4),
									Value: object.NewPanInt(5),
								},
							},
							[]object.Pair{},
						),
						Value: object.NewPanInt(6),
					},
				},
			),
		},
		// non-hashable dangling keys are also ignored ([1,2]: 4 is ignored)
		{
			`%{[1, 2]: 3, [1, 2]: 4}`,
			toPanMap(
				[]object.Pair{},
				[]object.Pair{
					{
						Key: object.NewPanArr(
							object.NewPanInt(1),
							object.NewPanInt(2),
						),
						Value: object.NewPanInt(3),
					},
				},
			),
		},
	}

	for _, tt := range tests {
		actual := testEval(t, tt.input)
		testPanMap(t, actual, tt.expected)
	}
}

func toPanMap(pairs []object.Pair, nonHashablePairs []object.Pair) *object.PanMap {
	return object.NewPanMap(append(pairs, nonHashablePairs...)...)
}

func TestEvalPinnedMapKey(t *testing.T) {
	tests := []struct {
		input    string
		expected *object.PanMap
	}{
		{
			`a := true; %{^a: 1}`,
			toPanMap(
				[]object.Pair{
					{
						Key:   object.BuiltInTrue,
						Value: object.NewPanInt(1),
					},
				},
				[]object.Pair{},
			),
		},
	}

	for _, tt := range tests {
		actual := testEval(t, tt.input)
		testValue(t, actual, tt.expected)
	}
}

func TestEvalUnpackingErr(t *testing.T) {
	tests := []struct {
		input    string
		expected *object.PanErr
	}{
		{
			`*[1]`,
			object.NewSyntaxErr("cannot use `*` unpacking outside of Arr."),
		},
		{
			`[*1]`,
			object.NewTypeErr("cannot use `*` unpacking for `1`"),
		},
		{
			`{a: 1, **[2]}`,
			object.NewTypeErr("cannot use `**` unpacking for `[2]`"),
		},
		{
			`%{'a: 1, **[3]}`,
			object.NewTypeErr("cannot use `**` unpacking for `[3]`"),
		},
	}

	for _, tt := range tests {
		actual := testEval(t, tt.input)
		testPanErr(t, actual, tt.expected)
	}
}

func TestEvalFuncLiteral(t *testing.T) {
	outerEnv := object.NewEnv()

	tests := []struct {
		input    string
		expected *object.PanFunc
	}{
		{
			`{||}`,
			toPanFunc(
				[]string{},
				[]object.Pair{},
				`|| `,
				outerEnv,
			),
		},
		{
			`{|a|}`,
			toPanFunc(
				[]string{"a"},
				[]object.Pair{},
				`|a| `,
				outerEnv,
			),
		},
		{
			`{|a, b|}`,
			toPanFunc(
				[]string{"a", "b"},
				[]object.Pair{},
				`|a, b| `,
				outerEnv,
			),
		},
		{
			`{|a, b, c:10|}`,
			toPanFunc(
				[]string{"a", "b"},
				[]object.Pair{
					{
						Key:   object.NewPanStr("c"),
						Value: object.NewPanInt(10),
					},
				},
				`|a, b, c: 10| `,
				outerEnv,
			),
		},
		{
			`{|a, b, c: 10, d:'e|}`,
			toPanFunc(
				[]string{"a", "b"},
				[]object.Pair{
					{
						Key:   object.NewPanStr("c"),
						Value: object.NewPanInt(10),
					},
					{
						Key:   object.NewPanStr("d"),
						Value: object.NewPanStr("e"),
					},
				},
				`|a, b, c: 10, d: 'e| `,
				outerEnv,
			),
		},
	}

	for _, tt := range tests {
		actual := testEvalInEnv(t, tt.input, outerEnv)
		testPanFunc(t, actual, tt.expected)
	}
}

func toPanFunc(
	args []string,
	kwargs []object.Pair,
	str string,
	env *object.Env,
) *object.PanFunc {
	argArr := []object.PanObject{}
	for _, arg := range args {
		argArr = append(argArr, object.NewPanStr(arg))
	}

	funcWrapper := &FuncWrapperImpl{
		codeStr: str,
		args:    object.NewPanArr(argArr...),
		kwargs:  toPanObj(kwargs),
		// empty stmt (body is not tested)
		body: &[]ast.Stmt{},
	}

	return object.NewPanFunc(
		funcWrapper,
		object.NewEnclosedEnv(env),
	)
}

func TestEvalFuncCall(t *testing.T) {
	tests := []struct {
		input    string
		expected object.PanObject
	}{
		{
			`{10}()`,
			object.NewPanInt(10),
		},
		{
			`{|x| x}(5)`,
			object.NewPanInt(5),
		},
		{
			`{|x, y| [x, y]}("x", "y")`,
			object.NewPanArr(
				object.NewPanStr("x"),
				object.NewPanStr("y"),
			),
		},
		{
			`{|foo, bar: "bar", baz| [foo, bar, baz]}("FOO", "BAZ")`,
			object.NewPanArr(
				object.NewPanStr("FOO"),
				object.NewPanStr("bar"),
				object.NewPanStr("BAZ"),
			),
		},
	}

	for _, tt := range tests {
		actual := testEval(t, tt.input)
		testValue(t, actual, tt.expected)
	}
}

func TestEvalEmptyFuncCall(t *testing.T) {
	tests := []struct {
		input    string
		expected object.PanObject
	}{
		{
			`{||}()`,
			object.BuiltInNil,
		},
	}

	for _, tt := range tests {
		actual := testEval(t, tt.input)
		testValue(t, actual, tt.expected)
	}
}

func TestEvalMultiLineFuncCall(t *testing.T) {
	tests := []struct {
		input    string
		expected object.PanObject
	}{
		{
			`{|a, b| a; b}(1, 2)`,
			object.NewPanInt(2),
		},
	}

	for _, tt := range tests {
		actual := testEval(t, tt.input)
		testValue(t, actual, tt.expected)
	}
}

func TestEvalBuiltInFuncCall(t *testing.T) {
	tests := []struct {
		input    string
		expected object.PanObject
	}{
		{
			`f := {||}['at]; f({a: 10}, ['a])`,
			object.NewPanInt(10),
		},
	}

	for _, tt := range tests {
		actual := testEval(t, tt.input)
		testValue(t, actual, tt.expected)
	}
}

func TestEvalBuiltInFuncInvalidArgErr(t *testing.T) {
	tests := []struct {
		input    string
		expected *object.PanErr
	}{
		{
			`{call: {|x| x}['call]}.call(1)`,
			object.NewTypeErr("{\"call\": {|| [builtin]}} is not callable."),
		},
		{
			`f := {||}['at]; f()`,
			object.NewTypeErr("Obj#at requires at least 2 args"),
		},
		{
			`f := {||}['at]; f(1)`,
			object.NewTypeErr("Obj#at requires at least 2 args"),
		},
	}

	for _, tt := range tests {
		actual := testEval(t, tt.input)
		testPanErr(t, actual, tt.expected)
	}
}

func TestEvalIterLiteral(t *testing.T) {
	outerEnv := object.NewEnv()

	tests := []struct {
		input    string
		expected *object.PanFunc
	}{
		{
			`<{||}>`,
			toPanIter(
				[]string{},
				[]object.Pair{},
				`|| `,
				object.NewEnclosedEnv(outerEnv),
			),
		},
		{
			`<{|a|}>`,
			toPanIter(
				[]string{"a"},
				[]object.Pair{},
				`|a| `,
				object.NewEnclosedEnv(outerEnv),
			),
		},
		{
			`<{|a, b|}>`,
			toPanIter(
				[]string{"a", "b"},
				[]object.Pair{},
				`|a, b| `,
				object.NewEnclosedEnv(outerEnv),
			),
		},
		{
			`<{|a, b, c:10|}>`,
			toPanIter(
				[]string{"a", "b"},
				[]object.Pair{
					{
						Key:   object.NewPanStr("c"),
						Value: object.NewPanInt(10),
					},
				},
				`|a, b, c: 10| `,
				object.NewEnclosedEnv(outerEnv),
			),
		},
		{
			`<{|a, b, c: 10, d:'e|}>`,
			toPanIter(
				[]string{"a", "b"},
				[]object.Pair{
					{
						Key:   object.NewPanStr("c"),
						Value: object.NewPanInt(10),
					},
					{
						Key:   object.NewPanStr("d"),
						Value: object.NewPanStr("e"),
					},
				},
				`|a, b, c: 10, d: 'e| `,
				object.NewEnclosedEnv(outerEnv),
			),
		},
	}

	for _, tt := range tests {
		actual := testEvalInEnv(t, tt.input, outerEnv)
		testPanFunc(t, actual, tt.expected)
	}
}

func toPanIter(
	args []string,
	kwargs []object.Pair,
	str string,
	env *object.Env,
) *object.PanFunc {
	argArr := []object.PanObject{}
	for _, arg := range args {
		argArr = append(argArr, object.NewPanStr(arg))
	}

	funcWrapper := &FuncWrapperImpl{
		codeStr: str,
		args:    object.NewPanArr(argArr...),
		kwargs:  toPanObj(kwargs),
		// empty stmt (body is not tested)
		body: &[]ast.Stmt{},
	}

	return object.NewPanIter(funcWrapper, env)
}

func TestEvalIterNew(t *testing.T) {
	outerEnv := object.NewEnv()

	tests := []struct {
		input    string
		expected object.PanObject
	}{
		// return new Iter with args set
		{
			`<{|a| yield a}>.new(10)`,
			toPanIter(
				[]string{"a"},
				[]object.Pair{},
				`|a| yield a`,
				toEnv(
					[]object.Pair{
						{
							Key:   object.NewPanStr("a"),
							Value: object.NewPanInt(10),
						},
					},
					[]object.Pair{},
					outerEnv,
				),
			),
		},
		{
			`<{|a, b, c: 'c| yield a}>.new('A, 'B, c: 'C)`,
			toPanIter(
				[]string{"a", "b"},
				[]object.Pair{
					{
						Key:   object.NewPanStr("c"),
						Value: object.NewPanStr("c"),
					},
				},
				`|a, b, c: 'c| yield a`,
				toEnv(
					[]object.Pair{
						{
							Key:   object.NewPanStr("a"),
							Value: object.NewPanStr("A"),
						},
						{
							Key:   object.NewPanStr("b"),
							Value: object.NewPanStr("B"),
						},
					},
					[]object.Pair{
						{
							Key:   object.NewPanStr("c"),
							Value: object.NewPanStr("C"),
						},
					},
					outerEnv,
				),
			),
		},
	}

	for _, tt := range tests {
		if outerEnv.Items().Inspect() != "{}" {
			t.Errorf("outerEnv must be empty. got=%s", outerEnv.Items().Inspect())
		}

		actual := testEvalInEnv(t, tt.input, outerEnv)
		testValue(t, actual, tt.expected)
	}
}

func toEnv(
	argPairs []object.Pair,
	kwargPairs []object.Pair,
	outer *object.Env,
) *object.Env {
	env := object.NewEnclosedEnv(outer)
	for i, pair := range argPairs {
		sym := object.GetSymHash(pair.Key.(*object.PanStr).Value)
		env.Set(sym, pair.Value)
		// \n
		argSym := object.GetSymHash(fmt.Sprintf(`\%d`, i+1))
		env.Set(argSym, pair.Value)
	}
	// \0
	argValues := []object.PanObject{}
	for _, pair := range argPairs {
		argValues = append(argValues, pair.Value)
	}
	env.Set(object.GetSymHash(`\0`), object.NewPanArr(argValues...))
	// \
	if len(argValues) > 0 {
		env.Set(object.GetSymHash("\\"), argValues[0])
	}

	for _, pair := range kwargPairs {
		sym := object.GetSymHash(pair.Key.(*object.PanStr).Value)
		env.Set(sym, pair.Value)
		// \hoge
		kwargSym := object.GetSymHash("\\" + pair.Key.(*object.PanStr).Value)
		env.Set(kwargSym, pair.Value)
	}
	// \_
	env.Set(object.GetSymHash("\\_"), toPanObj(kwargPairs))

	return env
}

func TestEvalDiamondLiteral(t *testing.T) {
	tests := []struct {
		input    string
		expected object.PanObject
	}{
		{
			`<>`,
			object.BuiltInDiamondObj,
		},
	}

	for _, tt := range tests {
		actual := testEval(t, tt.input)
		testValue(t, actual, tt.expected)
	}
}

func TestEvalDiamondProps(t *testing.T) {
	tests := []struct {
		stdin    string
		input    string
		expected object.PanObject
	}{
		// NOTE: <> can use all str props
		{
			`hoge`,
			`<>.uc`,
			object.NewPanStr("HOGE"),
		},
		{
			"HOGE\nFUGA",
			`<>.lc
			 <>.lc`,
			object.NewPanStr("fuga"),
		},
		// after reached to EOF, "" is read
		{
			"hoge\nfuga",
			`<>.uc
			 <>.uc
			 <>.uc`,
			object.NewPanStr(""),
		},
		// Diamond#S returns read line str
		{
			"hoge\nfuga",
			`<>.S+<>.S`,
			object.NewPanStr("hogefuga"),
		},
	}
	for _, tt := range tests {
		reader := strings.NewReader(tt.stdin)
		// setup IO
		env := object.NewEnvWithConsts()
		env.InjectIO(reader, os.Stdout)

		actual := testEvalInEnv(t, tt.input, env)
		testValue(t, actual, tt.expected)
	}
}

func TestEvalDiamondPropsErr(t *testing.T) {
	tests := []struct {
		input    string
		expected object.PanObject
	}{
		{
			`<>._missing`,
			object.NewTypeErr("_missing requires at least 2 args"),
		},
		{
			`<>['S]()`,
			object.NewTypeErr("Diamond#S requires at least 1 arg"),
		},
		{
			`<>._missing(1)`,
			object.NewTypeErr("\\2 must be prop name str"),
		},
		{
			`IO := 1; <>.uc`,
			object.NewTypeErr("name `IO` must be IO obj"),
		},
	}
	for _, tt := range tests {
		// setup IO
		env := object.NewEnvWithConsts()
		env.InjectIO(os.Stdin, os.Stdout)

		actual := testEvalInEnv(t, tt.input, env)
		testValue(t, actual, tt.expected)
	}
}

func TestEvalDiamondPropsErrIONotFound(t *testing.T) {
	tests := []struct {
		input    string
		expected object.PanObject
	}{
		{
			`<>.lc`,
			object.NewNameErr("`IO` is not found in env"),
		},
		{
			`<>.S`,
			object.NewNameErr("`IO` is not found in env"),
		},
		{
			`<>@{\}`,
			object.NewNameErr("`IO` is not found in env"),
		},
	}
	for _, tt := range tests {
		// without IO injection
		env := object.NewEnvWithConsts()
		actual := testEvalInEnv(t, tt.input, env)
		testValue(t, actual, tt.expected)
	}
}

func TestEvalDiamondIter(t *testing.T) {
	tests := []struct {
		stdin    string
		input    string
		expected object.PanObject
	}{
		// <> reads all lines until EOF
		{
			"hoge\nfuga\npiyo",
			`<>@{\}`,
			object.NewPanArr(
				object.NewPanStr("hoge"),
				object.NewPanStr("fuga"),
				object.NewPanStr("piyo"),
			),
		},
		{
			"hoge\nfuga\npiyo\n",
			`<>@{\}`,
			object.NewPanArr(
				object.NewPanStr("hoge"),
				object.NewPanStr("fuga"),
				object.NewPanStr("piyo"),
			),
		},
	}
	for _, tt := range tests {
		reader := strings.NewReader(tt.stdin)
		// setup IO
		env := object.NewEnvWithConsts()
		env.InjectIO(reader, os.Stdout)

		actual := testEvalInEnv(t, tt.input, env)
		testValue(t, actual, tt.expected)
	}
}

func TestEvalYield(t *testing.T) {
	tests := []struct {
		input    string
		expected object.PanObject
	}{
		{
			`<{|a| yield a}>.new(1).next`,
			object.NewPanInt(1),
		},
		{
			`<{|a| yield a}>.new('a).next`,
			object.NewPanStr("a"),
		},
		// recur
		{
			`it := <{|a| yield a; recur(a+1)}>.new(1)
			 it.next
			 it.next
			 `,
			object.NewPanInt(2),
		},
		{
			`it := <{|a| yield a; recur(a+1)}>.new(1)
			 it.next
			 it.next
			 it.next
			`,
			object.NewPanInt(3),
		},
		{
			`it := <{|a, b| yield [a, b]; recur(a+1, b*2)}>.new(1, 2)
			 it.next
			 it.next
			 it.next
			`,
			object.NewPanArr(
				object.NewPanInt(3),
				object.NewPanInt(8),
			),
		},
	}

	for _, tt := range tests {
		actual := testEval(t, tt.input)
		testValue(t, actual, tt.expected)
	}
}

func TestEvalIterYieldIsIndependent(t *testing.T) {
	tests := []struct {
		input    string
		expected object.PanObject
	}{
		{
			`
			it := <{|a| yield a; recur(a+1)}>.new(1)
			it.next
			it.next
			it2 := it.new(1)
			it2.next`,
			object.NewPanInt(1),
		},
	}

	for _, tt := range tests {
		actual := testEval(t, tt.input)
		testValue(t, actual, tt.expected)
	}
}

func TestEvalYieldStopped(t *testing.T) {
	tests := []struct {
		input    string
		expected object.PanObject
	}{
		{
			`it := <{|a| yield a if a == 1; recur(a+1)}>.new(1)
			it.next`,
			object.NewPanInt(1),
		},
		{
			`it := <{|a| yield a if a == 1; recur(a+1)}>.new(1)
			it.next
			it.next`,
			object.NewStopIterErr("iter stopped"),
		},
	}

	for _, tt := range tests {
		actual := testEval(t, tt.input)
		testValue(t, actual, tt.expected)
	}
}

func TestEvalJumpStmt(t *testing.T) {
	tests := []struct {
		input    string
		expected object.PanObject
	}{
		{
			`{|i| yield i}(1)`,
			object.NewPanInt(1),
		},
		{
			`{|i| return i}(2)`,
			object.NewPanInt(2),
		},
		{
			`{|i| raise Err.new("new error")}(2)`,
			object.NewPanErr("new error"),
		},
	}
	for _, tt := range tests {
		actual := testEval(t, tt.input)
		testValue(t, actual, tt.expected)
	}
}

func TestEvalJumpIfStmt(t *testing.T) {
	tests := []struct {
		input    string
		expected object.PanObject
	}{
		{
			`{|i| yield i if true; i * 2}(1)`,
			object.NewPanInt(1),
		},
		// raise StopIterErr
		{
			`{|i| yield i if false; i * 2}(1)`,
			object.NewStopIterErr("iter stopped"),
		},
		{
			`{|i| return i if true; i * 2}(10)`,
			object.NewPanInt(10),
		},
		{
			`{|i| return i if false; i * 2}(10)`,
			object.NewPanInt(20),
		},
		{
			`{|i| raise Err.new("new error") if true; i}(2)`,
			object.NewPanErr("new error"),
		},
		{
			`{|i| raise Err.new("new error") if false; i}(2)`,
			object.NewPanInt(2),
		},
	}
	for _, tt := range tests {
		actual := testEval(t, tt.input)
		testValue(t, actual, tt.expected)
	}
}

func TestEvalJumpIfStmtWithNonBoolCond(t *testing.T) {
	tests := []struct {
		input    string
		expected object.PanObject
	}{
		// `(cond).B` is called to check cond is truthy/falsy
		{
			`{|i| yield i if 1; i * 2}(1)`,
			object.NewPanInt(1),
		},
		{
			`{|i| yield i if ""; i * 2}(1)`,
			object.NewStopIterErr("iter stopped"),
		},
		{
			`{|i| return i if [1,2,3]; i * 2}(10)`,
			object.NewPanInt(10),
		},
		{
			`{|i| return i if nil; i * 2}(10)`,
			object.NewPanInt(20),
		},
	}
	for _, tt := range tests {
		actual := testEval(t, tt.input)
		testValue(t, actual, tt.expected)
	}
}

func TestEvalStmtsAfterYieldAreEvaluated(t *testing.T) {
	tests := []struct {
		input          string
		expected       object.PanObject
		expectedStdOut string
	}{
		{
			`{|i| yield i; "evaluated!".p}(1)`,
			object.NewPanInt(1),
			"evaluated!\n",
		},
	}
	for _, tt := range tests {
		writer := &bytes.Buffer{}
		// setup IO
		env := object.NewEnvWithConsts()
		env.InjectIO(os.Stdin, writer)

		actual := testEvalInEnv(t, tt.input, env)
		testValue(t, actual, tt.expected)

		// check output
		output := writer.String()
		if output != tt.expectedStdOut {
			t.Errorf("wrong output. expected=`%s`, got=`%s`",
				tt.expectedStdOut, output)
		}
	}
}

func TestEvalJumpStmtJumpPrecedence(t *testing.T) {
	tests := []struct {
		input    string
		expected object.PanObject
	}{
		// jump precedence: last line < yield < return = raise
		{
			`{|i| yield i; i * 2}(1)`,
			object.NewPanInt(1),
		},
		{
			`{|i| yield i; return i * 2}(1)`,
			object.NewPanInt(2),
		},
		{
			`{|i| return i; yield i * 2}(3)`,
			object.NewPanInt(3),
		},
		{
			`{|i| return i; i * 2}(4)`,
			object.NewPanInt(4),
		},
		{
			`{|i| yield i; raise i * 2}(1)`,
			object.NewPanInt(2),
		},
		{
			`{|i| raise i; yield i * 2}(3)`,
			object.NewPanInt(3),
		},
		{
			`{|i| raise i; i * 2}(4)`,
			object.NewPanInt(4),
		},
	}

	for _, tt := range tests {
		actual := testEval(t, tt.input)
		testValue(t, actual, tt.expected)
	}
}

func TestEvalDefer(t *testing.T) {
	// defer is evaluated after jump (same as that in golang)
	tests := []struct {
		input          string
		expected       object.PanObject
		expectedStdOut string
	}{
		{
			`{|i| defer "one".p; i}(1)`,
			object.NewPanInt(1),
			"one\n",
		},
		{
			`{|i| defer "one".p; return i}(1)`,
			object.NewPanInt(1),
			"one\n",
		},
		{
			`{|i| defer "one".p; yield i}(1)`,
			object.NewPanInt(1),
			"one\n",
		},
		{
			`{|i| defer "one".p; defer "two".p; i}(1)`,
			object.NewPanInt(1),
			"one\ntwo\n",
		},
		// NOTE: only defers above last evaluated line are valid
		{
			`{|i| defer "one".p; return i; defer "two".p}(1)`,
			object.NewPanInt(1),
			"one\n",
		},
		{
			`{|i| defer "one".p; yield i; defer "two".p}(1)`,
			object.NewPanInt(1),
			"one\ntwo\n",
		},
		{
			`{|i| defer "one".p if true; i}(1)`,
			object.NewPanInt(1),
			"one\n",
		},
		{
			`{|i| defer "one".p if false; i}(1)`,
			object.NewPanInt(1),
			"",
		},
	}

	for _, tt := range tests {
		writer := &bytes.Buffer{}
		// setup IO
		env := object.NewEnvWithConsts()
		env.InjectIO(os.Stdin, writer)

		actual := testEvalInEnv(t, tt.input, env)
		testValue(t, actual, tt.expected)

		// check output
		output := writer.String()
		if output != tt.expectedStdOut {
			t.Errorf("wrong output. expected=`%s`, got=`%s`",
				tt.expectedStdOut, output)
		}
	}
}

func TestEvalRaise(t *testing.T) {
	// raise works like return but unwraps errWrappers
	tests := []struct {
		input    string
		expected object.PanObject
	}{
		{
			`{|i| raise 10; i}(5)`,
			object.NewPanInt(10),
		},
		{
			`{|i| raise Err.new("new error"); i}(5)`,
			object.NewPanErr("new error"),
		},
		// err wrapped by Either is unwrapped by raise
		{
			`5.try.fmap {Err.new("new error")}.A.{|val, err| raise err; val}`,
			object.NewPanErr("new error"),
		},
	}

	for _, tt := range tests {
		actual := testEval(t, tt.input)
		testValue(t, actual, tt.expected)
	}
}

func TestEvalRaiseIf(t *testing.T) {
	// raise works like return but unwraps errWrappers
	tests := []struct {
		input    string
		expected object.PanObject
	}{
		{
			`{|i| raise 10 if true; i}(5)`,
			object.NewPanInt(10),
		},
		{
			`{|i| raise Err.new("new error") if true; i}(5)`,
			object.NewPanErr("new error"),
		},
		// err wrapped by Either is unwrapped by raise
		{
			`5.try.fmap {Err.new("new error")}.A.{|val, err| raise err if true; val}`,
			object.NewPanErr("new error"),
		},
	}

	for _, tt := range tests {
		actual := testEval(t, tt.input)
		testValue(t, actual, tt.expected)
	}
}

func TestEvalArrIter(t *testing.T) {
	tests := []struct {
		input    string
		expected object.PanObject
	}{
		{
			`it := ['a, 'b]._iter
			 it.next`,
			object.NewPanStr("a"),
		},
		{
			`it := ['a, 'b]._iter
			 it.next
			 it.next`,
			object.NewPanStr("b"),
		},
		{
			`it := ['a, 'b]._iter
			it.next
			it.next
			it.next`,
			object.NewStopIterErr("iter stopped"),
		},
	}

	for _, tt := range tests {
		actual := testEval(t, tt.input)
		testValue(t, actual, tt.expected)
	}
}

func TestEvalIntIter(t *testing.T) {
	tests := []struct {
		input    string
		expected object.PanObject
	}{
		{
			`it := 2._iter
			 it.next`,
			object.NewPanInt(1),
		},
		{
			`it := 2._iter
			 it.next
			 it.next`,
			object.NewPanInt(2),
		},
		{
			`it := 2._iter
			it.next
			it.next
			it.next`,
			object.NewStopIterErr("iter stopped"),
		},
	}

	for _, tt := range tests {
		actual := testEval(t, tt.input)
		testValue(t, actual, tt.expected)
	}
}

func TestEvalStrIter(t *testing.T) {
	tests := []struct {
		input    string
		expected object.PanObject
	}{
		{
			`it := "Hi"._iter
			 it.next`,
			object.NewPanStr("H"),
		},
		{
			`it := "Hi"._iter
			 it.next
			 it.next`,
			object.NewPanStr("i"),
		},
		{
			`it := "Hi"._iter
			it.next
			it.next
			it.next`,
			object.NewStopIterErr("iter stopped"),
		},
		// each iteration yields rune (not byte!)
		{
			`it := "日本語"._iter
			 it.next`,
			object.NewPanStr("日"),
		},
		{
			`it := "日本語"._iter
			 it.next
			 it.next`,
			object.NewPanStr("本"),
		},
		{
			`it := "日本語"._iter
			 it.next
			 it.next
			 it.next
			 it.next`,
			object.NewStopIterErr("iter stopped"),
		},
	}

	for _, tt := range tests {
		actual := testEval(t, tt.input)
		testValue(t, actual, tt.expected)
	}
}

func TestEvalObjIter(t *testing.T) {
	tests := []struct {
		input    string
		expected object.PanObject
	}{
		{
			`it := {a: "A", b: "B"}._iter
			 it.next`,
			object.NewPanArr(
				object.NewPanStr("a"),
				object.NewPanStr("A"),
			),
		},
		{
			`it := {a: "A", b: "B"}._iter
			 it.next
			 it.next`,
			object.NewPanArr(
				object.NewPanStr("b"),
				object.NewPanStr("B"),
			),
		},
		// sorted in alphabetical order of keys
		{
			`it := {b: "B", a: "A"}._iter
			 it.next`,
			object.NewPanArr(
				object.NewPanStr("a"),
				object.NewPanStr("A"),
			),
		},
		{
			`it := {b: "B", a: "A"}._iter
			 it.next
			 it.next`,
			object.NewPanArr(
				object.NewPanStr("b"),
				object.NewPanStr("B"),
			),
		},
		{
			`it := {a: "A", b: "B"}._iter
			it.next
			it.next
			it.next`,
			object.NewStopIterErr("iter stopped"),
		},
		// private keys are ignored
		{
			`it := {
				_secret: "invisible",
				'+: "PLUS",
				"include spaces": "ignored",
			 }._iter
			 it.next`,
			object.NewStopIterErr("iter stopped"),
		},
	}

	for _, tt := range tests {
		actual := testEval(t, tt.input)
		testValue(t, actual, tt.expected)
	}
}

func TestEvalMapIter(t *testing.T) {
	tests := []struct {
		input    string
		expected object.PanObject
	}{
		{
			`it := %{true: 1, false: 2}._iter
			 it.next`,
			object.NewPanArr(
				object.BuiltInTrue,
				object.NewPanInt(1),
			),
		},
		{
			`it := %{true: 1, false: 2}._iter
			 it.next
			 it.next`,
			object.NewPanArr(
				object.BuiltInFalse,
				object.NewPanInt(2),
			),
		},
		{
			`it := %{true: 1}._iter
			 it.next
			 it.next`,
			object.NewStopIterErr("iter stopped"),
		},
		// non-scalar key
		{
			`it := %{[0]: 1}._iter
			 it.next`,
			object.NewPanArr(
				object.NewPanArr(
					object.NewPanInt(0),
				),
				object.NewPanInt(1),
			),
		},
	}

	for _, tt := range tests {
		actual := testEval(t, tt.input)
		testValue(t, actual, tt.expected)
	}
}

func TestEvalRangeIter(t *testing.T) {
	tests := []struct {
		input    string
		expected object.PanObject
	}{
		{
			`it := (0:2)._iter
			 it.next`,
			object.NewPanInt(0),
		},
		{
			`it := (0:2)._iter
			 it.next
			 it.next`,
			object.NewPanInt(1),
		},
		{
			`it := (0:2)._iter
			 it.next
			 it.next`,
			object.NewPanInt(1),
		},
		{
			`it := (0:2)._iter
			 it.next
			 it.next
			 it.next`,
			object.NewStopIterErr("iter stopped"),
		},
		{
			`it := (10:7:-2)._iter
			 it.next`,
			object.NewPanInt(10),
		},
		{
			`it := (10:7:-2)._iter
			 it.next
			 it.next`,
			object.NewPanInt(8),
		},
		{
			`it := (10:7:-2)._iter
			 it.next
			 it.next
			 it.next`,
			object.NewStopIterErr("iter stopped"),
		},
		// str
		{
			`it := ('a:'c)._iter
			 it.next`,
			object.NewPanStr("a"),
		},
		{
			`it := ('a:'c)._iter
			 it.next
			 it.next`,
			object.NewPanStr("b"),
		},
		{
			`it := ('a:'c)._iter
			 it.next
			 it.next
			 it.next`,
			object.NewStopIterErr("iter stopped"),
		},
	}

	for _, tt := range tests {
		actual := testEval(t, tt.input)
		testValue(t, actual, tt.expected)
	}
}

func TestEvalIterIter(t *testing.T) {
	tests := []struct {
		input    string
		expected object.PanObject
	}{
		{
			`it := <{|i| yield i; recur(i+1)}>.new(1)._iter
			 it.next`,
			object.NewPanInt(1),
		},
		{
			`it := <{|i| yield i; recur(i+1)}>.new(1)._iter
			 it.next
			 it.next`,
			object.NewPanInt(2),
		},
		{
			`it := <{|i| yield i if i < 2; recur(i+1)}>.new(1)._iter
			 it.next
			 it.next`,
			object.NewStopIterErr("iter stopped"),
		},
		// new iter generated by ._iter is independent of original one
		{
			`orig := <{|i| yield i; recur(i+1)}>.new(1)
			 it := orig._iter
			 orig.next
			 orig.next
			 it.next`,
			object.NewPanInt(1),
		},
		// errors
		{
			`it := <{}>['_iter]()`,
			object.NewTypeErr("Iter#_iter requires at least 1 arg"),
		},
		{
			`it := <{}>['_iter](1)`,
			object.NewTypeErr("1 cannot be treated as iter"),
		},
		{
			`it := <{}>['_iter]({|i| i})`,
			object.NewTypeErr("{|i| i} cannot be treated as iter"),
		},
	}

	for _, tt := range tests {
		actual := testEval(t, tt.input)
		testValue(t, actual, tt.expected)
	}
}

func TestEvalListChainPropCall(t *testing.T) {
	tests := []struct {
		input    string
		expected object.PanObject
	}{
		{
			`[1, 2]@+(1)`,
			object.NewPanArr(
				object.NewPanInt(2),
				object.NewPanInt(3),
			),
		},
		{
			`3@S`,
			object.NewPanArr(
				object.NewPanStr("1"),
				object.NewPanStr("2"),
				object.NewPanStr("3"),
			),
		},
		{
			`"日本語"@S`,
			object.NewPanArr(
				object.NewPanStr("日"),
				object.NewPanStr("本"),
				object.NewPanStr("語"),
			),
		},
		// NOTE: must new() to init params
		{
			`<{|i| yield i if i != 3; recur(i+1)}>.new(0)@S`,
			object.NewPanArr(
				object.NewPanStr("0"),
				object.NewPanStr("1"),
				object.NewPanStr("2"),
			),
		},
		{
			`(1:6:2)@S`,
			object.NewPanArr(
				object.NewPanStr("1"),
				object.NewPanStr("3"),
				object.NewPanStr("5"),
			),
		},
		{
			`{a: 1, b: 2}@at([0])`,
			object.NewPanArr(
				object.NewPanStr("a"),
				object.NewPanStr("b"),
			),
		},
		{
			`%{"a": 1, "b": 2, "c": 3}@at([0])`,
			object.NewPanArr(
				object.NewPanStr("a"),
				object.NewPanStr("b"),
				object.NewPanStr("c"),
			),
		},
		// if any element raises error, propagate it
		{
			`[1, 2, "str", 4]@+(1)`,
			object.NewTypeErr("1 cannot be treated as str"),
		},
	}

	for _, tt := range tests {
		actual := testEval(t, tt.input)
		testValue(t, actual, tt.expected)
	}
}

func TestEvalListChainPropCallIgnoresNil(t *testing.T) {
	tests := []struct {
		input    string
		expected object.PanObject
	}{
		{
			`[{name: 'Taro}, {name: nil}, {name: 'Jiro}]@name`,
			object.NewPanArr(
				object.NewPanStr("Taro"),
				object.NewPanStr("Jiro"),
			),
		},
	}

	for _, tt := range tests {
		actual := testEval(t, tt.input)
		testValue(t, actual, tt.expected)
	}
}

func TestEvalListChainLiteralCall(t *testing.T) {
	tests := []struct {
		input    string
		expected object.PanObject
	}{
		{
			`[1, 2]@{|i| i + 1}`,
			object.NewPanArr(
				object.NewPanInt(2),
				object.NewPanInt(3),
			),
		},
		{
			`3@{|i| i.S}`,
			object.NewPanArr(
				object.NewPanStr("1"),
				object.NewPanStr("2"),
				object.NewPanStr("3"),
			),
		},
		{
			`"日本語"@{|c| c + "!"}`,
			object.NewPanArr(
				object.NewPanStr("日!"),
				object.NewPanStr("本!"),
				object.NewPanStr("語!"),
			),
		},
		{
			`{a: 1, b: 2}@{|k, v| "#{k}: #{v}"}`,
			object.NewPanArr(
				object.NewPanStr("a: 1"),
				object.NewPanStr("b: 2"),
			),
		},
		{
			`%{true: "yes", false: "no"}@{|k, v| [k, v]}`,
			object.NewPanArr(
				object.NewPanArr(
					object.BuiltInTrue,
					object.NewPanStr("yes"),
				),
				object.NewPanArr(
					object.BuiltInFalse,
					object.NewPanStr("no"),
				),
			),
		},
		{
			`(1:6:2)@{|n| n+1}`,
			object.NewPanArr(
				object.NewPanInt(2),
				object.NewPanInt(4),
				object.NewPanInt(6),
			),
		},
		// if any element raises error, propagate it
		{
			`[2, 0, 3]@{|n| 6/n}`,
			object.NewZeroDivisionErr("cannot be divided by 0"),
		},
		// if chainarg raises error, propagate it
		{
			`[1]@(_){|n| n + 1}`,
			object.BuiltInNotImplemented,
		},
	}

	for _, tt := range tests {
		actual := testEval(t, tt.input)
		testValue(t, actual, tt.expected)
	}
}

func TestEvalListChainLiteralCallIgnoresNil(t *testing.T) {
	tests := []struct {
		input    string
		expected object.PanObject
	}{
		{
			`[{name: 'Taro}, {name: nil}, {name: 'Jiro}]@{|o| o.name}`,
			object.NewPanArr(
				object.NewPanStr("Taro"),
				object.NewPanStr("Jiro"),
			),
		},
	}

	for _, tt := range tests {
		actual := testEval(t, tt.input)
		testValue(t, actual, tt.expected)
	}
}

func TestEvalListChainVarCall(t *testing.T) {
	tests := []struct {
		input    string
		expected object.PanObject
	}{
		{
			`f := {|i| i + 1}; [1, 2]@^f`,
			object.NewPanArr(
				object.NewPanInt(2),
				object.NewPanInt(3),
			),
		},
		{
			`f := {|i| i.S}; 3@^f`,
			object.NewPanArr(
				object.NewPanStr("1"),
				object.NewPanStr("2"),
				object.NewPanStr("3"),
			),
		},
		{
			`f := {|c| c + "!"}; "日本語"@^f`,
			object.NewPanArr(
				object.NewPanStr("日!"),
				object.NewPanStr("本!"),
				object.NewPanStr("語!"),
			),
		},
		{
			`f := {|k, v| "#{k}: #{v}"}; {a: 1, b: 2}@^f`,
			object.NewPanArr(
				object.NewPanStr("a: 1"),
				object.NewPanStr("b: 2"),
			),
		},
		{
			`f := {|k, v| [k, v]}; %{true: "yes", false: "no"}@^f`,
			object.NewPanArr(
				object.NewPanArr(
					object.BuiltInTrue,
					object.NewPanStr("yes"),
				),
				object.NewPanArr(
					object.BuiltInFalse,
					object.NewPanStr("no"),
				),
			),
		},
		{
			`f := {|n| n+1}; (1:6:2)@^f`,
			object.NewPanArr(
				object.NewPanInt(2),
				object.NewPanInt(4),
				object.NewPanInt(6),
			),
		},
		// if any element raises error, propagate it
		{
			`f := {|n| 6/n}; [2, 0, 3]@^f`,
			object.NewZeroDivisionErr("cannot be divided by 0"),
		},
	}

	for _, tt := range tests {
		actual := testEval(t, tt.input)
		testValue(t, actual, tt.expected)
	}
}

func TestEvalListChainVarCallIgnoresNil(t *testing.T) {
	tests := []struct {
		input    string
		expected object.PanObject
	}{
		{
			`f := {|o| o.name}; [{name: 'Taro}, {name: nil}, {name: 'Jiro}]@^f`,
			object.NewPanArr(
				object.NewPanStr("Taro"),
				object.NewPanStr("Jiro"),
			),
		},
	}

	for _, tt := range tests {
		actual := testEval(t, tt.input)
		testValue(t, actual, tt.expected)
	}
}

func TestEvalReduceChainPropCall(t *testing.T) {
	tests := []struct {
		input    string
		expected object.PanObject
	}{
		{
			`[1, 2, 3]$(4)+`,
			object.NewPanInt(10),
		},
		// if chainarg is not set, acc is `nil`
		{
			`[nil]$==`,
			object.BuiltInTrue,
		},
	}

	for _, tt := range tests {
		actual := testEval(t, tt.input)
		testValue(t, actual, tt.expected)
	}
}

func TestEvalReduceChainPropCallRecvIsAcc(t *testing.T) {
	// reciever of prop is resolved to acc
	tests := []struct {
		input    string
		expected object.PanObject
	}{
		// "a"['+] is called
		{
			`[1]$("a")+`,
			object.NewTypeErr("1 cannot be treated as str"),
		},
	}

	for _, tt := range tests {
		actual := testEval(t, tt.input)
		testValue(t, actual, tt.expected)
	}
}

func TestEvalReduceChainLiteralCall(t *testing.T) {
	tests := []struct {
		input    string
		expected object.PanObject
	}{
		{
			`[1, 2, 3]$(0){|i, j| i + j}`,
			object.NewPanInt(6),
		},
		{
			`[1, 2, 3]$(4){|i, j| i + j}`,
			object.NewPanInt(10),
		},
		// if chainarg is not set, acc is `nil`
		{
			`[nil]${|i, j| i == j}`,
			object.BuiltInTrue,
		},
	}

	for _, tt := range tests {
		actual := testEval(t, tt.input)
		testValue(t, actual, tt.expected)
	}
}

func TestEvalReduceChainVarCall(t *testing.T) {
	tests := []struct {
		input    string
		expected object.PanObject
	}{
		{
			`f := {|i, j| i + j}; [1, 2, 3]$(0)^f`,
			object.NewPanInt(6),
		},
		{
			`f := {|i, j| i + j}; [1, 2, 3]$(4)^f`,
			object.NewPanInt(10),
		},
		// if chainarg is not set, acc is `nil`
		{
			`f := {|i, j| i == j}; [nil]$^f`,
			object.BuiltInTrue,
		},
	}

	for _, tt := range tests {
		actual := testEval(t, tt.input)
		testValue(t, actual, tt.expected)
	}
}

func TestEvalIfExpr(t *testing.T) {
	tests := []struct {
		input    string
		expected object.PanObject
	}{
		{
			`1 if true`,
			object.NewPanInt(1),
		},
		{
			`1 if false`,
			object.BuiltInNil,
		},
		{
			`10 if true else 5`,
			object.NewPanInt(10),
		},
		{
			`10 if false else 5`,
			object.NewPanInt(5),
		},
		// cond other than bool
		{
			`'t if 100 else 'f`,
			object.NewPanStr("t"),
		},
		{
			`'t if [] else 'f`,
			object.NewPanStr("f"),
		},
	}

	for _, tt := range tests {
		actual := testEval(t, tt.input)
		testValue(t, actual, tt.expected)
	}
}

func TestEvalInfixOr(t *testing.T) {
	tests := []struct {
		input    string
		expected object.PanObject
	}{
		{
			`true || false`,
			object.BuiltInTrue,
		},
		{
			`false || true`,
			object.BuiltInTrue,
		},
		{
			`false || false`,
			object.BuiltInFalse,
		},
		{
			`true || true`,
			object.BuiltInTrue,
		},
		{
			`'a || 10`,
			object.NewPanStr("a"),
		},
		{
			`false || 10`,
			object.NewPanInt(10),
		},
		// evaluated to right if left is falsy
		{
			`[] || 1`,
			object.NewPanInt(1),
		},
		{
			`"" || "b"`,
			object.NewPanStr("b"),
		},
		// shortcut
		{
			`5 || _`,
			object.NewPanInt(5),
		},
		{
			`_ || 5`,
			object.BuiltInNotImplemented,
		},
	}

	for _, tt := range tests {
		actual := testEval(t, tt.input)
		testValue(t, actual, tt.expected)
	}
}

func TestEvalInfixAnd(t *testing.T) {
	tests := []struct {
		input    string
		expected object.PanObject
	}{
		{
			`true && false`,
			object.BuiltInFalse,
		},
		{
			`false && true`,
			object.BuiltInFalse,
		},
		{
			`false && false`,
			object.BuiltInFalse,
		},
		{
			`true && true`,
			object.BuiltInTrue,
		},
		{
			`false && 10`,
			object.BuiltInFalse,
		},
		// evaluated to right if left is truthy
		{
			`['a] && 1`,
			object.NewPanInt(1),
		},
		{
			`"abc" && "b"`,
			object.NewPanStr("b"),
		},
		// shortcut
		{
			`"" && _`,
			object.NewPanStr(""),
		},
		{
			`_ && ""`,
			object.BuiltInNotImplemented,
		},
	}

	for _, tt := range tests {
		actual := testEval(t, tt.input)
		testValue(t, actual, tt.expected)
	}
}

func TestEvalAssign(t *testing.T) {
	tests := []struct {
		input       string
		expected    object.PanObject
		expectedEnv *object.PanObj
	}{
		{
			`a := 'A`,
			object.NewPanStr("A"),
			toPanObj([]object.Pair{
				{
					Key:   object.NewPanStr("a"),
					Value: object.NewPanStr("A"),
				},
			}),
		},
		{
			`'A => a`,
			object.NewPanStr("A"),
			toPanObj([]object.Pair{
				{
					Key:   object.NewPanStr("a"),
					Value: object.NewPanStr("A"),
				},
			}),
		},
		{
			`a := 'A; a`,
			object.NewPanStr("A"),
			toPanObj([]object.Pair{
				{
					Key:   object.NewPanStr("a"),
					Value: object.NewPanStr("A"),
				},
			}),
		},
		{
			`a := 5; b := 10; [a, b]`,
			object.NewPanArr(
				object.NewPanInt(5),
				object.NewPanInt(10),
			),
			toPanObj([]object.Pair{
				{
					Key:   object.NewPanStr("a"),
					Value: object.NewPanInt(5),
				},
				{
					Key:   object.NewPanStr("b"),
					Value: object.NewPanInt(10),
				},
			}),
		},
		{
			`a := b := 2; [a, b]`,
			object.NewPanArr(
				object.NewPanInt(2),
				object.NewPanInt(2),
			),
			toPanObj([]object.Pair{
				{
					Key:   object.NewPanStr("a"),
					Value: object.NewPanInt(2),
				},
				{
					Key:   object.NewPanStr("b"),
					Value: object.NewPanInt(2),
				},
			}),
		},
		{
			`"hi" => a; a`,
			object.NewPanStr("hi"),
			toPanObj([]object.Pair{
				{
					Key:   object.NewPanStr("a"),
					Value: object.NewPanStr("hi"),
				},
			}),
		},
		{
			`3 => c => d; [c, d]`,
			object.NewPanArr(
				object.NewPanInt(3),
				object.NewPanInt(3),
			),
			toPanObj([]object.Pair{
				{
					Key:   object.NewPanStr("c"),
					Value: object.NewPanInt(3),
				},
				{
					Key:   object.NewPanStr("d"),
					Value: object.NewPanInt(3),
				},
			}),
		},
		// right value is not assigned if error raises
		{
			`a := 0 / 0`,
			object.NewZeroDivisionErr("cannot be divided by 0"),
			toPanObj([]object.Pair{}),
		},
		{
			`0 / 0 => b`,
			object.NewZeroDivisionErr("cannot be divided by 0"),
			toPanObj([]object.Pair{}),
		},
	}

	for _, tt := range tests {
		// NOTE: use empty env without consts to make tt simple
		env := object.NewEnv()
		actual := testEvalInEnv(t, tt.input, env)
		testValue(t, actual, tt.expected)
		testValue(t, env.Items(), tt.expectedEnv)
	}
}

func TestEvalAssignShadowingConsts(t *testing.T) {
	tests := []struct {
		input    string
		expected object.PanObject
	}{
		{
			`true := 100; true`,
			object.NewPanInt(100),
		},
	}

	for _, tt := range tests {
		actual := testEval(t, tt.input)
		testValue(t, actual, tt.expected)
	}
}

func TestEvalArgVars(t *testing.T) {
	tests := []struct {
		input    string
		expected object.PanObject
	}{
		{
			`{\1}("a", "b")`,
			object.NewPanStr("a"),
		},
		{
			`{\2}("a", "b")`,
			object.NewPanStr("b"),
		},
		// NameErr if number exceeds arity
		{
			`{\3}("a", "b")`,
			object.NewNameErr("name `\\3` is not defined"),
		},
		// `\` is syntax sugar of `\1`
		{
			`{\}("one", "two")`,
			object.NewPanStr("one"),
		},
		// `\0` is arr of all args
		{
			`{\0}("one", "two")`,
			object.NewPanArr(
				object.NewPanStr("one"),
				object.NewPanStr("two"),
			),
		},
		// NOTE: \0 does not contain kwargs
		{
			`{\0}("arg", kwarg: "kwarg")`,
			object.NewPanArr(
				object.NewPanStr("arg"),
			),
		},
	}

	for _, tt := range tests {
		actual := testEval(t, tt.input)
		testValue(t, actual, tt.expected)
	}
}

func TestEvalKwargVars(t *testing.T) {
	tests := []struct {
		input    string
		expected object.PanObject
	}{
		{
			`{\a}(a: "A", b: "B")`,
			object.NewPanStr("A"),
		},
		{
			`{\b}(a: "A", b: "B")`,
			object.NewPanStr("B"),
		},
		// NameErr if kwarg is not found
		{
			`{\c}(a: "A", b: "B")`,
			object.NewNameErr("name `\\c` is not defined"),
		},
		// `\_` is obj of all kwargs
		{
			`{\_}(one: 1, two: 2)`,
			toPanObj([]object.Pair{
				{
					Key:   object.NewPanStr("one"),
					Value: object.NewPanInt(1),
				},
				{
					Key:   object.NewPanStr("two"),
					Value: object.NewPanInt(2),
				},
			}),
		},
		// NOTE: \0 does not contain kwargs
		{
			`{\0}("arg", kwarg: "kwarg")`,
			object.NewPanArr(
				object.NewPanStr("arg"),
			),
		},
	}

	for _, tt := range tests {
		actual := testEval(t, tt.input)
		testValue(t, actual, tt.expected)
	}
}

func TestEvalProto(t *testing.T) {
	tests := []struct {
		input    string
		expected object.PanObject
	}{
		{
			`Obj.proto`,
			object.BuiltInBaseObj,
		},
		{
			`[].proto`,
			object.BuiltInArrObj,
		},
		{
			`1.0.proto`,
			object.BuiltInFloatObj,
		},
		{
			`{||}.proto`,
			object.BuiltInFuncObj,
		},
		{
			`<{||}>.proto`,
			object.BuiltInIterObj,
		},
		{
			`1.proto`,
			object.BuiltInIntObj,
		},
		{
			`%{}.proto`,
			object.BuiltInMapObj,
		},
		{
			`nil.proto`,
			object.BuiltInNilObj,
		},
		{
			`{}.proto`,
			object.BuiltInObjObj,
		},
		{
			`'a.proto`,
			object.BuiltInStrObj,
		},
		{
			`true.proto`,
			object.BuiltInOneInt,
		},
		{
			`false.proto`,
			object.BuiltInZeroInt,
		},
	}

	for _, tt := range tests {
		actual := testEval(t, tt.input)
		testValue(t, actual, tt.expected)
	}
}

func TestEvalBearProto(t *testing.T) {
	tests := []struct {
		input string
	}{
		{`obj := {a: 1};  obj.bear({x: 2}).proto == obj`},
		{`arr := [1, 2];  arr.bear({x: 2}).proto == arr`},
		{`num := 2;       num.bear({x: 2}).proto == num`},
		{`map := %{1: 2}; map.bear({x: 2}).proto == map`},
	}
	for _, tt := range tests {
		actual := testEval(t, tt.input)
		testValue(t, actual, object.BuiltInTrue)
	}
}

func TestEvalBearErr(t *testing.T) {
	tests := []struct {
		input    string
		expected object.PanObject
	}{
		{
			`{}.bear(1)`,
			object.NewTypeErr("BaseObj#bear requires obj literal src"),
		},
		{
			`f := {}['bear]; f()`,
			object.NewTypeErr("BaseObj#bear requires at least 1 arg"),
		},
	}
	for _, tt := range tests {
		actual := testEval(t, tt.input)
		testValue(t, actual, tt.expected)
	}
}

func TestEvalBearContents(t *testing.T) {
	tests := []struct {
		input    string
		expected object.PanObject
	}{
		{
			`Obj.bear({a: 1})`,
			toPanObj([]object.Pair{
				{
					Key:   object.NewPanStr("a"),
					Value: object.NewPanInt(1),
				},
			}),
		},
		// if no arg, use empty obj
		{
			`Obj.bear`,
			toPanObj([]object.Pair{}),
		},
	}
	for _, tt := range tests {
		actual := testEval(t, tt.input)
		testValue(t, actual, tt.expected)
	}
}

func TestEvalObjKeys(t *testing.T) {
	tests := []struct {
		input    string
		expected object.PanObject
	}{
		{
			`{}.keys`,
			object.NewPanArr(),
		},
		{
			`{a: 1, b: 2}.keys`,
			object.NewPanArr(
				object.NewPanStr("a"),
				object.NewPanStr("b"),
			),
		},
		// private keys are ignored
		{
			`{a: 1, _b: 2, '+: 3, "with space": 4}.keys`,
			object.NewPanArr(
				object.NewPanStr("a"),
			),
		},
		// with private keys (private keys follow public keys)
		{
			`{a: 1, _b: 2, '+: 3, "with space": 4}.keys(private?: true)`,
			object.NewPanArr(
				object.NewPanStr("a"),
				object.NewPanStr("+"),
				object.NewPanStr("_b"),
				object.NewPanStr("with space"),
			),
		},
		// child of obj
		{
			`{a: 1}.bear({b: 2}).keys`,
			object.NewPanArr(
				object.NewPanStr("b"),
			),
		},
	}

	for _, tt := range tests {
		actual := testEval(t, tt.input)
		testValue(t, actual, tt.expected)
	}
}

func TestEvalObjValues(t *testing.T) {
	tests := []struct {
		input    string
		expected object.PanObject
	}{
		{
			`{}.values`,
			object.NewPanArr(),
		},
		{
			`{a: "A", b: "B"}.values`,
			object.NewPanArr(
				object.NewPanStr("A"),
				object.NewPanStr("B"),
			),
		},
		// private keys are ignored
		{
			`{
			   a: "A",
			   _b: "_B",
			   '+: "PLUS",
			   "with space": "WITH SPACE",
			 }.values`,
			object.NewPanArr(
				object.NewPanStr("A"),
			),
		},
		// with private keys (private keys follow public keys)
		{
			`{
			   a: "A",
			   _b: "_B",
			   '+: "PLUS",
			   "with space": "WITH SPACE",
			 }.values(private?: true)`,
			object.NewPanArr(
				object.NewPanStr("A"),
				object.NewPanStr("PLUS"),
				object.NewPanStr("_B"),
				object.NewPanStr("WITH SPACE"),
			),
		},
		// child of obj
		{
			`{a: 1}.bear({b: 2}).values`,
			object.NewPanArr(
				object.NewPanInt(2),
			),
		},
	}

	for _, tt := range tests {
		actual := testEval(t, tt.input)
		testValue(t, actual, tt.expected)
	}
}

func TestEvalObjItems(t *testing.T) {
	tests := []struct {
		input    string
		expected object.PanObject
	}{
		{
			`{}.items`,
			object.NewPanArr(),
		},
		{
			`{a: "A", b: "B"}.items`,
			object.NewPanArr(
				object.NewPanArr(
					object.NewPanStr("a"),
					object.NewPanStr("A"),
				),
				object.NewPanArr(
					object.NewPanStr("b"),
					object.NewPanStr("B"),
				),
			),
		},
		// private keys are ignored
		{
			`{
			   a: "A",
			   _b: "_B",
			   '+: "PLUS",
			   "with space": "WITH SPACE",
			 }.items`,
			object.NewPanArr(
				object.NewPanArr(
					object.NewPanStr("a"),
					object.NewPanStr("A"),
				),
			),
		},
		// with private keys (private keys follow public keys)
		{
			`{
			   a: "A",
			   _b: "_B",
			   '+: "PLUS",
			   "with space": "WITH SPACE",
			 }.items(private?: true)`,
			object.NewPanArr(
				object.NewPanArr(
					object.NewPanStr("a"),
					object.NewPanStr("A"),
				),
				object.NewPanArr(
					object.NewPanStr("+"),
					object.NewPanStr("PLUS"),
				),
				object.NewPanArr(
					object.NewPanStr("_b"),
					object.NewPanStr("_B"),
				),
				object.NewPanArr(
					object.NewPanStr("with space"),
					object.NewPanStr("WITH SPACE"),
				),
			),
		},
		// child of obj
		{
			`{a: 1}.bear({b: 2}).items`,
			object.NewPanArr(
				object.NewPanArr(
					object.NewPanStr("b"),
					object.NewPanInt(2),
				),
			),
		},
	}

	for _, tt := range tests {
		actual := testEval(t, tt.input)
		testValue(t, actual, tt.expected)
	}
}

func TestEvalMapKeys(t *testing.T) {
	tests := []struct {
		input    string
		expected object.PanObject
	}{
		{
			`%{}.keys`,
			object.NewPanArr(),
		},
		{
			`%{true: 1}.keys`,
			object.NewPanArr(
				object.BuiltInTrue,
			),
		},
		{
			`%{[0]: "zero"}.keys`,
			object.NewPanArr(
				object.NewPanArr(
					object.NewPanInt(0),
				),
			),
		},
		{
			`%{true: 1, "false": 2}.keys`,
			object.NewPanArr(
				object.BuiltInTrue,
				object.NewPanStr("false"),
			),
		},
		{
			`%{[0]: "zero", [1]: "one"}.keys`,
			object.NewPanArr(
				object.NewPanArr(
					object.NewPanInt(0),
				),
				object.NewPanArr(
					object.NewPanInt(1),
				),
			),
		},
		// non-hashables follow hashables
		{
			`%{[0]: "zero", 1: "one", "2": "two", [3]: "three"}.keys`,
			object.NewPanArr(
				object.NewPanInt(1),
				object.NewPanStr("2"),
				object.NewPanArr(
					object.NewPanInt(0),
				),
				object.NewPanArr(
					object.NewPanInt(3),
				),
			),
		},
	}

	for _, tt := range tests {
		actual := testEval(t, tt.input)
		testValue(t, actual, tt.expected)
	}
}

func TestEvalMapValues(t *testing.T) {
	tests := []struct {
		input    string
		expected object.PanObject
	}{
		{
			`%{}.values`,
			object.NewPanArr(),
		},
		{
			`%{true: 1}.values`,
			object.NewPanArr(
				object.NewPanInt(1),
			),
		},
		{
			`%{true: 1, false: 0}.values`,
			object.NewPanArr(
				object.NewPanInt(1),
				object.NewPanInt(0),
			),
		},
		{
			`%{[0]: "zero"}.values`,
			object.NewPanArr(
				object.NewPanStr("zero"),
			),
		},
		{
			`%{[0]: "zero", [1]: "one"}.values`,
			object.NewPanArr(
				object.NewPanStr("zero"),
				object.NewPanStr("one"),
			),
		},
		// non-hashables follow hashables
		{
			`%{[0]: "zero", 1: "one", "2": "two", [3]: "three"}.values`,
			object.NewPanArr(
				object.NewPanStr("one"),
				object.NewPanStr("two"),
				object.NewPanStr("zero"),
				object.NewPanStr("three"),
			),
		},
	}

	for _, tt := range tests {
		actual := testEval(t, tt.input)
		testValue(t, actual, tt.expected)
	}
}

func TestEvalMapItems(t *testing.T) {
	tests := []struct {
		input    string
		expected object.PanObject
	}{
		{
			`%{}.items`,
			object.NewPanArr(),
		},
		{
			`%{true: 1}.items`,
			object.NewPanArr(
				object.NewPanArr(
					object.BuiltInTrue,
					object.NewPanInt(1),
				),
			),
		},
		{
			`%{true: 1, false: 0}.items`,
			object.NewPanArr(
				object.NewPanArr(
					object.BuiltInTrue,
					object.NewPanInt(1),
				),
				object.NewPanArr(
					object.BuiltInFalse,
					object.NewPanInt(0),
				),
			),
		},
		{
			`%{[0]: "zero"}.items`,
			object.NewPanArr(
				object.NewPanArr(
					object.NewPanArr(
						object.NewPanInt(0),
					),
					object.NewPanStr("zero"),
				),
			),
		},
		{
			`%{[0]: "zero", [1]: "one"}.items`,
			object.NewPanArr(
				object.NewPanArr(
					object.NewPanArr(
						object.NewPanInt(0),
					),
					object.NewPanStr("zero"),
				),
				object.NewPanArr(
					object.NewPanArr(
						object.NewPanInt(1),
					),
					object.NewPanStr("one"),
				),
			),
		},
		// non-hashables follow hashables
		{
			`%{[0]: "zero", 1: "one", "2": "two", [3]: "three"}.items`,
			object.NewPanArr(
				object.NewPanArr(
					object.NewPanInt(1),
					object.NewPanStr("one"),
				),
				object.NewPanArr(
					object.NewPanStr("2"),
					object.NewPanStr("two"),
				),
				object.NewPanArr(
					object.NewPanArr(
						object.NewPanInt(0),
					),
					object.NewPanStr("zero"),
				),
				object.NewPanArr(
					object.NewPanArr(
						object.NewPanInt(3),
					),
					object.NewPanStr("three"),
				),
			),
		},
	}

	for _, tt := range tests {
		actual := testEval(t, tt.input)
		testValue(t, actual, tt.expected)
	}
}

func TestEvalStringify(t *testing.T) {
	tests := []struct {
		input    string
		expected object.PanObject
	}{
		{
			`[1].S`,
			object.NewPanStr("[1]"),
		},
		// unnecessary zeros are omitted
		{
			`1.0.S`,
			object.NewPanStr("1.0"),
		},
		{
			`{|x| x}.S`,
			object.NewPanStr("{|x| x}"),
		},
		{
			`10.S`,
			object.NewPanStr("10"),
		},
		{
			`%{'a: 1}.S`,
			object.NewPanStr(`%{"a": 1}`),
		},
		{
			`nil.S`,
			object.NewPanStr(`nil`),
		},
		{
			`{a: 1}.S`,
			object.NewPanStr(`{"a": 1}`),
		},
		{
			`(1:2).S`,
			object.NewPanStr("(1:2:nil)"),
		},
		// str is not quoted
		{
			`'a.S`,
			object.NewPanStr("a"),
		},
		{
			`true.S`,
			object.NewPanStr("true"),
		},
		{
			`false.S`,
			object.NewPanStr("false"),
		},
	}

	for _, tt := range tests {
		actual := testEval(t, tt.input)
		testValue(t, actual, tt.expected)
	}
}

func TestEvalBoolify(t *testing.T) {
	tests := []struct {
		input    string
		expected object.PanObject
	}{
		// arr
		{
			`[1].B`,
			object.BuiltInTrue,
		},
		{
			`[].B`,
			object.BuiltInFalse,
		},
		// float
		{
			`1.0.B`,
			object.BuiltInTrue,
		},
		{
			`0.0.B`,
			object.BuiltInFalse,
		},
		// func (always true)
		{
			`{|x| x}.B`,
			object.BuiltInTrue,
		},
		{
			`{}['keys].B`,
			object.BuiltInTrue,
		},
		// iter (always true)
		{
			`<{|x| x}>.B`,
			object.BuiltInTrue,
		},
		// int
		{
			`10.B`,
			object.BuiltInTrue,
		},
		{
			`0.B`,
			object.BuiltInFalse,
		},
		{
			`true.B`,
			object.BuiltInTrue,
		},
		{
			`false.B`,
			object.BuiltInFalse,
		},
		// map
		{
			`%{'a: 1}.B`,
			object.BuiltInTrue,
		},
		{
			`%{[1]: 1}.B`,
			object.BuiltInTrue,
		},
		{
			`%{}.B`,
			object.BuiltInFalse,
		},
		// nil
		{
			`nil.B`,
			object.BuiltInFalse,
		},
		// obj
		{
			`{a: 1}.B`,
			object.BuiltInTrue,
		},
		{
			`{}.B`,
			object.BuiltInFalse,
		},
		// range (always true)
		{
			`(1:2).B`,
			object.BuiltInTrue,
		},
		// str
		{
			`'a.B`,
			object.BuiltInTrue,
		},
		{
			`"".B`,
			object.BuiltInFalse,
		},
	}

	for _, tt := range tests {
		actual := testEval(t, tt.input)
		testValue(t, actual, tt.expected)
	}
}

func TestEvalIntify(t *testing.T) {
	tests := []struct {
		input    string
		expected object.PanObject
	}{
		{
			`"1".I`,
			object.NewPanInt(1),
		},
		{
			`"-5".I`,
			object.NewPanInt(-5),
		},
	}

	for _, tt := range tests {
		actual := testEval(t, tt.input)
		testValue(t, actual, tt.expected)
	}
}

func TestEvalIntifyErr(t *testing.T) {
	tests := []struct {
		input    string
		expected object.PanObject
	}{
		{
			`Str['I]()`,
			object.NewTypeErr("Str#I requires at least 1 arg"),
		},
		{
			`"1"['I](1)`,
			object.NewTypeErr("1 cannot be treated as str"),
		},
		{
			`"-5a".I`,
			object.NewValueErr("\"-5a\" cannot be converted into int"),
		},
	}

	for _, tt := range tests {
		actual := testEval(t, tt.input)
		testValue(t, actual, tt.expected)
	}
}

func TestEvalFloatify(t *testing.T) {
	tests := []struct {
		input    string
		expected object.PanObject
	}{
		// int
		{
			`1.F`,
			object.NewPanFloat(1.0),
		},
		// float
		{
			`4.0.F`,
			object.NewPanFloat(4.0),
		},
		// str
		{
			`"1.0".F`,
			object.NewPanFloat(1.0),
		},
		{
			`"-5.25".F`,
			object.NewPanFloat(-5.25),
		},
	}

	for _, tt := range tests {
		actual := testEval(t, tt.input)
		testValue(t, actual, tt.expected)
	}
}

func TestEvalFloatifyErr(t *testing.T) {
	tests := []struct {
		input    string
		expected object.PanObject
	}{
		{
			`1['F]()`,
			object.NewTypeErr("Num#F requires at least 1 arg"),
		},
		{
			`1['F]("a")`,
			object.NewTypeErr("\"a\" cannot be treated as num"),
		},
		{
			`Str['F]()`,
			object.NewTypeErr("Str#F requires at least 1 arg"),
		},
		{
			`"1"['F](2)`,
			object.NewTypeErr("2 cannot be treated as str"),
		},
		{
			`"-5.a".F`,
			object.NewValueErr("\"-5.a\" cannot be converted into float"),
		},
	}

	for _, tt := range tests {
		actual := testEval(t, tt.input)
		testValue(t, actual, tt.expected)
	}
}

func TestEvalRepr(t *testing.T) {
	tests := []struct {
		input    string
		expected object.PanObject
	}{
		{
			`[1].repr`,
			object.NewPanStr("[1]"),
		},
		// precise value is shown
		{
			`1.0.repr`,
			object.NewPanStr("1.000000"),
		},
		{
			`{|x| x}.repr`,
			object.NewPanStr("{|x| x}"),
		},
		{
			`10.repr`,
			object.NewPanStr("10"),
		},
		{
			`%{'a: 1}.repr`,
			object.NewPanStr(`%{"a": 1}`),
		},
		{
			`nil.repr`,
			object.NewPanStr(`nil`),
		},
		{
			`{a: 1}.repr`,
			object.NewPanStr(`{"a": 1}`),
		},
		{
			`(1:2).repr`,
			object.NewPanStr("(1:2:nil)"),
		},
		// str is quoted
		{
			`'a.repr`,
			object.NewPanStr(`"a"`),
		},
		{
			`true.repr`,
			object.NewPanStr("true"),
		},
		{
			`false.repr`,
			object.NewPanStr("false"),
		},
		// return name instead if self has '_name and it is str
		{
			`{_name: "myObj"}.repr`,
			object.NewPanStr("myObj"),
		},
		{
			`{_name: 1}.repr`,
			object.NewPanStr(`{"_name": 1}`),
		},
		// proto's _name is not referred
		{
			`{_name: "myObj"}.bear.repr`,
			object.NewPanStr(`{}`),
		},
	}

	for _, tt := range tests {
		actual := testEval(t, tt.input)
		testValue(t, actual, tt.expected)
	}
}

func TestEvalNamesOfConsts(t *testing.T) {
	tests := []struct {
		input    string
		expected object.PanObject
	}{
		{
			`Int._name`,
			object.NewPanStr("Int"),
		},
		{
			`Float._name`,
			object.NewPanStr("Float"),
		},
		{
			`Num._name`,
			object.NewPanStr("Num"),
		},
		{
			`Nil._name`,
			object.NewPanStr("Nil"),
		},
		{
			`Str._name`,
			object.NewPanStr("Str"),
		},
		{
			`Arr._name`,
			object.NewPanStr("Arr"),
		},
		{
			`Range._name`,
			object.NewPanStr("Range"),
		},
		{
			`Func._name`,
			object.NewPanStr("Func"),
		},
		{
			`Iter._name`,
			object.NewPanStr("Iter"),
		},
		{
			`Iterable._name`,
			object.NewPanStr("Iterable"),
		},
		{
			`Comparable._name`,
			object.NewPanStr("Comparable"),
		},
		{
			`Wrappable._name`,
			object.NewPanStr("Wrappable"),
		},
		{
			`Match._name`,
			object.NewPanStr("Match"),
		},
		{
			`Obj._name`,
			object.NewPanStr("Obj"),
		},
		{
			`BaseObj._name`,
			object.NewPanStr("BaseObj"),
		},
		{
			`Map._name`,
			object.NewPanStr("Map"),
		},
		{
			`Diamond._name`,
			object.NewPanStr("Diamond"),
		},
		{
			`Kernel._name`,
			object.NewPanStr("Kernel"),
		},
		{
			`Either._name`,
			object.NewPanStr("Either"),
		},
		{
			`EitherErr._name`,
			object.NewPanStr("EitherErr"),
		},
		{
			`EitherVal._name`,
			object.NewPanStr("EitherVal"),
		},
		{
			`Err._name`,
			object.NewPanStr("Err"),
		},
		{
			`AssertionErr._name`,
			object.NewPanStr("AssertionErr"),
		},
		{
			`NameErr._name`,
			object.NewPanStr("NameErr"),
		},
		{
			`NoPropErr._name`,
			object.NewPanStr("NoPropErr"),
		},
		{
			`NotImplementedErr._name`,
			object.NewPanStr("NotImplementedErr"),
		},
		{
			`StopIterErr._name`,
			object.NewPanStr("StopIterErr"),
		},
		{
			`SyntaxErr._name`,
			object.NewPanStr("SyntaxErr"),
		},
		{
			`TypeErr._name`,
			object.NewPanStr("TypeErr"),
		},
		{
			`ValueErr._name`,
			object.NewPanStr("ValueErr"),
		},
		{
			`ZeroDivisionErr._name`,
			object.NewPanStr("ZeroDivisionErr"),
		},
	}

	for _, tt := range tests {
		actual := testEval(t, tt.input)
		testValue(t, actual, tt.expected)
	}
}

func TestEvalPrint(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		// print obj.S results
		{
			`[1].p`,
			"[1]\n",
		},
		{
			`1.0.p`,
			"1.0\n",
		},
		{
			`{|x| x}.p`,
			"{|x| x}\n",
		},
		{
			`10.p`,
			"10\n",
		},
		{
			`%{'a: 1}.p`,
			"%{\"a\": 1}\n",
		},
		{
			`nil.p`,
			"nil\n",
		},
		{
			`{a: 1}.p`,
			"{\"a\": 1}\n",
		},
		{
			`(1:2).p`,
			"(1:2:nil)\n",
		},
		// str is not quoted
		{
			`'a.p`,
			"a\n",
		},
		{
			`true.p`,
			"true\n",
		},
		{
			`false.p`,
			"false\n",
		},
		// set-up end
		{
			`"abc".p(end: ",")`,
			"abc,",
		},
	}

	for _, tt := range tests {
		writer := &bytes.Buffer{}
		// setup IO
		env := object.NewEnvWithConsts()
		env.InjectIO(os.Stdin, writer)

		actual := testEvalInEnv(t, tt.input, env)
		// p returns nil
		testValue(t, actual, object.BuiltInNil)

		// check output
		output := writer.String()
		if output != tt.expected {
			t.Errorf("wrong output. expected=`%s`, got=`%s`",
				tt.expected, output)
		}
	}
}

func TestEvalPrintErr(t *testing.T) {
	tests := []struct {
		input    string
		expected object.PanObject
	}{
		{
			`f := {}['p]; f()`,
			object.NewTypeErr("Obj#p requires at least 1 arg"),
		},
		{
			`BaseObj.p`,
			object.NewNoPropErr("property `p` is not defined."),
		},
		{
			`IO := 1; 'a.p`,
			object.NewTypeErr("name `IO` is not io object"),
		},
		{
			`{S: 1}.p`,
			object.NewTypeErr(`\1.S must be str`),
		},
		{
			`"abc".p(end: 1)`,
			object.NewTypeErr(`end must be str`),
		},
	}

	for _, tt := range tests {
		writer := &bytes.Buffer{}
		// setup IO
		env := object.NewEnvWithConsts()
		env.InjectIO(os.Stdin, writer)

		actual := testEvalInEnv(t, tt.input, env)
		testValue(t, actual, tt.expected)

		// check output is empty
		output := writer.String()
		if output != "" {
			t.Errorf("output must be empty. got=`%s`", output)
		}
	}
}

func TestEvalPrintErrIfNoIO(t *testing.T) {
	tests := []struct {
		input    string
		expected object.PanObject
	}{
		{
			`1.p`,
			object.NewNameErr("name `IO` is not defined."),
		},
	}

	for _, tt := range tests {
		// const `IO` is not set up
		actual := testEval(t, tt.input)
		testValue(t, actual, tt.expected)
	}
}

func TestEvalScalarPropChain(t *testing.T) {
	tests := []struct {
		input    string
		expected object.PanObject
	}{
		{
			`{a: 5, b: 10}.a`,
			object.NewPanInt(5),
		},
		{
			`{a: 5, b: 10}.b`,
			object.NewPanInt(10),
		},
		// call method
		{
			`{a: {|| 2}}.a`,
			object.NewPanInt(2),
		},
		{
			`{a: m{|x| x}}.a(3)`,
			object.NewPanInt(3),
		},
		{
			`{a: m{|x, y| [x, y]}}.a("one", "two")`,
			object.NewPanArr(
				object.NewPanStr("one"),
				object.NewPanStr("two"),
			),
		},
		{
			`{a: m{|x, y: "y"| [x, y]}}.a("x")`,
			object.NewPanArr(
				object.NewPanStr("x"),
				object.NewPanStr("y"),
			),
		},
		{
			`{a: m{|x, y: "y"| [x, y]}}.a("x", y: "Y")`,
			object.NewPanArr(
				object.NewPanStr("x"),
				object.NewPanStr("Y"),
			),
		},
		// if args are insufficient, they are padded by nil
		{
			`{a: m{|x, y| [x, y]}}.a("X")`,
			object.NewPanArr(
				object.NewPanStr("X"),
				object.BuiltInNil,
			),
		},
		// if too many args are passed, they are just ignored
		{
			`{a: m{|x| x}}.a("arg", "needless", "extra", "args")`,
			object.NewPanStr("arg"),
		},
	}

	for _, tt := range tests {
		actual := testEval(t, tt.input)
		testValue(t, actual, tt.expected)
	}
}

func TestEvalScalarPropChainWithIter(t *testing.T) {
	env := object.NewEnv()

	tests := []struct {
		input    string
		expected object.PanObject
	}{
		// unlike func, iter is not called
		{
			`{a: <{1}>}.a`,
			toPanIter(
				[]string{},
				[]object.Pair{},
				"|| 1",
				object.NewEnclosedEnv(env),
			),
		},
	}

	for _, tt := range tests {
		actual := testEvalInEnv(t, tt.input, env)
		testValue(t, actual, tt.expected)
	}
}

func TestEvalArgUnpack(t *testing.T) {
	tests := []struct {
		input    string
		expected object.PanObject
	}{
		{
			`{|i, j, k| [i, j, k]}(*[1, 2, 3])`,
			object.NewPanArr(
				object.NewPanInt(1),
				object.NewPanInt(2),
				object.NewPanInt(3),
			),
		},
		// with other prefix
		{
			`{|i, j, k| [i, j, k]}(*[1, 2], !3)`,
			object.NewPanArr(
				object.NewPanInt(1),
				object.NewPanInt(2),
				object.BuiltInFalse,
			),
		},
	}

	for _, tt := range tests {
		actual := testEval(t, tt.input)
		testValue(t, actual, tt.expected)
	}
}

func TestEvalKwargUnpack(t *testing.T) {
	tests := []struct {
		input    string
		expected object.PanObject
	}{
		{
			`{|a: 1, b: 2| [a, b]}(**{a: 5, b: 10})`,
			object.NewPanArr(
				object.NewPanInt(5),
				object.NewPanInt(10),
			),
		},
		{
			`{|a: 1, b: 2| [a, b]}(a: 3, **{a: 6, b: 9})`,
			object.NewPanArr(
				object.NewPanInt(3),
				object.NewPanInt(9),
			),
		},
	}

	for _, tt := range tests {
		actual := testEval(t, tt.input)
		testValue(t, actual, tt.expected)
	}
}

func TestEvalUnpackError(t *testing.T) {
	tests := []struct {
		input    string
		expected object.PanObject
	}{
		{
			`{|a: 1, b: 2| [a, b]}(**c)`,
			object.NewNameErr("name `c` is not defined"),
		},
		{
			`{|a: 1, b: 2| [a, b]}(**[])`,
			object.NewTypeErr("cannot use `**` unpacking for `[]`"),
		},
		{
			`{|a, b| [a, b]}(*c)`,
			object.NewNameErr("name `c` is not defined"),
		},
		{
			`{|a, b| [a, b]}(*{})`,
			object.NewTypeErr("cannot use `*` unpacking for `{}`"),
		},
	}

	for _, tt := range tests {
		actual := testEval(t, tt.input)
		testValue(t, actual, tt.expected)
	}
}

func TestEvalAnonPropChain(t *testing.T) {
	tests := []struct {
		input    string
		expected object.PanObject
	}{
		{
			`{|x| .a}({a: 5})`,
			object.NewPanInt(5),
		},
		{
			`{|x| @a}([{a: 1}, {a: 2}, {a: 3}])`,
			object.NewPanArr(
				object.NewPanInt(1),
				object.NewPanInt(2),
				object.NewPanInt(3),
			),
		},
		{
			`{|x| $(0)+}([1, 2, 3])`,
			object.NewPanInt(6),
		},
		{
			`.a`,
			object.NewNameErr("name `\\1` is not defined"),
		},
	}

	for _, tt := range tests {
		actual := testEval(t, tt.input)
		testValue(t, actual, tt.expected)
	}
}

func TestEvalAnonLiteralChain(t *testing.T) {
	tests := []struct {
		input    string
		expected object.PanObject
	}{
		{
			`{|x| .{|y| y + 2}}(3)`,
			object.NewPanInt(5),
		},
		{
			`{|nums| @{|n| n.S}}([1, 2, 3])`,
			object.NewPanArr(
				object.NewPanStr("1"),
				object.NewPanStr("2"),
				object.NewPanStr("3"),
			),
		},
		{
			`{|nums| $(0){|acc, i| acc + i}}([1, 2, 3])`,
			object.NewPanInt(6),
		},
		{
			`.{|x| x}`,
			object.NewNameErr("name `\\1` is not defined"),
		},
	}

	for _, tt := range tests {
		actual := testEval(t, tt.input)
		testValue(t, actual, tt.expected)
	}
}

func TestEvalObjAt(t *testing.T) {
	tests := []struct {
		input    string
		expected object.PanObject
	}{
		{
			`{a: 5, b: 10}['a]`,
			object.NewPanInt(5),
		},
		// if key is insufficient, return nil
		{
			`{a: 5, b: 10}[]`,
			object.BuiltInNil,
		},
		// if key is not found, return nil
		{
			`{a: 5, b: 10}['c]`,
			object.BuiltInNil,
		},
		// trace prototype chain
		{
			`{a: 5, b: 10}['at]`,
			(*object.BuiltInBaseObj.Pairs)[object.GetSymHash("at")].Value,
		},
		// use descendant of arr for arg
		{
			`{a: 1}.at(['a].bear)`,
			object.NewPanInt(1),
		},
		// use descendant of str for index
		{
			`{a: 1}['a.bear]`,
			object.NewPanInt(1),
		},
		// if \1 is not obj, return nil
		{
			`Obj['at]([], [1])`,
			object.BuiltInNil,
		},
		// if \2 is not arr, return nil
		{
			`{}.at({})`,
			object.BuiltInNil,
		},
	}

	for _, tt := range tests {
		actual := testEval(t, tt.input)
		testValue(t, actual, tt.expected)
	}
}

func TestEvalIndexErr(t *testing.T) {
	tests := []struct {
		input    string
		expected object.PanObject
	}{
		{
			`{}.at`,
			object.NewTypeErr("Obj#at requires at least 2 args"),
		},
		{
			`[].at`,
			object.NewTypeErr("Arr#at requires at least 2 args"),
		},
		{
			`"".at`,
			object.NewTypeErr("Str#at requires at least 2 args"),
		},
		{
			`1.at`,
			object.NewTypeErr("Int#at requires at least 2 args"),
		},
		// invalid index range
		{
			`[1,2,3][1:5:0]`,
			object.NewValueErr("cannot use 0 for range step"),
		},
	}

	for _, tt := range tests {
		actual := testEval(t, tt.input)
		testValue(t, actual, tt.expected)
	}
}

func TestEvalNoPropErr(t *testing.T) {
	tests := []struct {
		input    string
		expected *object.PanErr
	}{
		{
			`{a: 1}.b`,
			object.NewNoPropErr("property `b` is not defined."),
		},
		// case sensitive
		{
			`{A: 1}.a`,
			object.NewNoPropErr("property `a` is not defined."),
		},
	}

	for _, tt := range tests {
		actual := testEval(t, tt.input)
		testPanErr(t, actual, tt.expected)
	}
}

func TestEvalMissing(t *testing.T) {
	tests := []struct {
		input    string
		expected object.PanObject
	}{
		// if no prop found is proto chain, _missing is called
		{
			`{_missing: 'default}.hoge`,
			object.NewPanStr("default"),
		},
		{
			`[{_missing: 'a}, {_missing: 'b}]@hoge`,
			object.NewPanArr(
				object.NewPanStr("a"),
				object.NewPanStr("b"),
			),
		},
		{
			`['A, 'B]$({_missing: m{|name, i| "#{name} not found"+i}})+`,
			object.NewPanStr("+ not foundAB"),
		},
		// prop name and args are passed
		{
			`{_missing: m{|name, arg| [name, arg]}}.hoge(1)`,
			object.NewPanArr(
				object.NewPanStr("hoge"),
				object.NewPanInt(1),
			),
		},
		// proto _missing works
		{
			`{_missing: 'proto}.bear.hoge`,
			object.NewPanStr("proto"),
		},
	}

	for _, tt := range tests {
		actual := testEval(t, tt.input)
		testValue(t, actual, tt.expected)
	}
}

func TestEvalMapAt(t *testing.T) {
	tests := []struct {
		input    string
		expected object.PanObject
	}{
		{
			`%{'a: 5, 'b: 10}['a]`,
			object.NewPanInt(5),
		},
		// if key is insufficient, return nil
		{
			`%{'a: 5, 'b: 10}[]`,
			object.BuiltInNil,
		},
		// if key is not found, return nil
		{
			`%{'a: 5, 'b: 10}['c]`,
			object.BuiltInNil,
		},
		// trace prototype chain
		{
			`%{'a: 5, 'b: 10}['at]`,
			(*object.BuiltInMapObj.Pairs)[object.GetSymHash("at")].Value,
		},
		// index other than str can also be used
		{
			`%{1: "one"}[1]`,
			object.NewPanStr("one"),
		},
		{
			`%{nil: "nil"}[nil]`,
			object.NewPanStr("nil"),
		},
		{
			`%{[10]: "tenArr"}[[10]]`,
			object.NewPanStr("tenArr"),
		},
		// use descendant of arr for arg
		{
			`%{0: 1}.at([0].bear)`,
			object.NewPanInt(1),
		},
		// if \1 is not map, return nil
		{
			`Map['at]({}, [1])`,
			object.BuiltInNil,
		},
		// if \2 is not arr, return nil
		{
			`%{}.at({})`,
			object.BuiltInNil,
		},
		// if key is insufficient, return nil
		{
			`%{}[]`,
			object.BuiltInNil,
		},
	}

	for _, tt := range tests {
		actual := testEval(t, tt.input)
		testValue(t, actual, tt.expected)
	}
}

func TestEvalMapLen(t *testing.T) {
	tests := []struct {
		input    string
		expected object.PanObject
	}{
		{
			`%{}.len`,
			object.NewPanInt(0),
		},
		{
			`%{"a": 1}.len`,
			object.NewPanInt(1),
		},
		// with non-scalar key
		{
			`%{"a": 1, [2]: 3}.len`,
			object.NewPanInt(2),
		},
		// use descendant of str for recv
		{
			`%{1: 2}.bear.len`,
			object.NewPanInt(1),
		},
		// if no args are passed, raise an error
		{
			`Map['len]()`,
			object.NewTypeErr("Map#len requires at least 1 arg"),
		},
		// if \1 is not map, raise an error
		{
			`Map['len](1)`,
			object.NewTypeErr("\\1 must be map"),
		},
	}

	for _, tt := range tests {
		actual := testEval(t, tt.input)
		testValue(t, actual, tt.expected)
	}
}

func TestEvalLiteralCall(t *testing.T) {
	tests := []struct {
		input    string
		expected object.PanObject
	}{
		{
			`3.{|i| i * 2}`,
			object.NewPanInt(6),
		},
		{
			`[1, 2].{|i| i * 2}`,
			object.NewPanArr(
				object.NewPanInt(1),
				object.NewPanInt(2),
				object.NewPanInt(1),
				object.NewPanInt(2),
			),
		},
		// if arity is more than 1, arr is extracted to each param
		{
			`[1, 2].{|i, j| i + j}`,
			object.NewPanInt(3),
		},
		// if arity is more than 1, descendant of arr is extracted to each param
		{
			`[1, 2].bear.{|i, j| i + j}`,
			object.NewPanInt(3),
		},
		// if args are insufficient, they are padded by nil
		{
			`'X.{|i, j| [i, j]}`,
			object.NewPanArr(
				object.NewPanStr("X"),
				object.BuiltInNil,
			),
		},
		// if too many args are passed, they are just ignored
		{
			`[2, 3, "extra"].{|x, y| x + y}`,
			object.NewPanInt(5),
		},
		// if error raises while kwarg evaluation, propagate it
		{
			`[1].{|x, y: 6 / 0| x + y}`,
			object.NewZeroDivisionErr("cannot be divided by 0"),
		},
	}

	for _, tt := range tests {
		actual := testEval(t, tt.input)
		testValue(t, actual, tt.expected)
	}
}

func TestEvalLiteralCallLiteralProxy(t *testing.T) {
	// _literalProxy prop works as proxy for literal chain
	tests := []struct {
		input    string
		expected object.PanObject
	}{
		// NOTE: do not use literal-chain(or var-chain) in _literalProxy! (infinite loop)
		{
			`3.bear({_literalProxy: m{|f| f(f(self))}}).{|i| i * 2}`,
			object.NewPanInt(12),
		},
		// _literalProxy in proto can also work
		{
			`'hi.bear({_literalProxy: m{|f| f(f(self))}}).bear.{|i| i * 2}`,
			object.NewPanStr("hihihihi"),
		},
	}

	for _, tt := range tests {
		actual := testEval(t, tt.input)
		testValue(t, actual, tt.expected)
	}
}

func TestEvalVarCall(t *testing.T) {
	tests := []struct {
		input    string
		expected object.PanObject
	}{
		{
			`f := {|i| i * 2}; 3.^f`,
			object.NewPanInt(6),
		},
		{
			`f := {|i| i * 2}; [1, 2].^f`,
			object.NewPanArr(
				object.NewPanInt(1),
				object.NewPanInt(2),
				object.NewPanInt(1),
				object.NewPanInt(2),
			),
		},
		// built-in func
		{
			`f := Int['-%]; 3.^f`,
			object.NewPanInt(-3),
		},
		// error
		{
			`f := 2; 3.^f`,
			object.NewTypeErr("literal call must be func"),
		},
		{
			`3.^f`,
			object.NewNameErr("name `f` is not defined"),
		},
		// if arity is more than 1, arr is extracted to each param
		{
			`f := {|i, j| i + j}; [1, 2].^f`,
			object.NewPanInt(3),
		},
		// if args are insufficient, they are padded by nil
		{
			`f := {|i, j| [i, j]}; "X".^f`,
			object.NewPanArr(
				object.NewPanStr("X"),
				object.BuiltInNil,
			),
		},
		// if too many args are passed, they are just ignored
		{
			`f := {|x, y| x + y}; [2, 3, "extra"].^f`,
			object.NewPanInt(5),
		},
		// if descendant of func is passed, its proto is called
		{
			`f := {|i| i * 2}.bear; 5.^f`,
			object.NewPanInt(10),
		},
		// if descendant of built-in func is passed, its proto is called
		{
			`f := Int['-%].bear; 3.^f`,
			object.NewPanInt(-3),
		},
		// if recv has call prop, use it
		// NOTE: arguments [f, *args] are passed to 'call
		// (self in 'call func is f itself)
		{
			`f := {call: m{|i| i + 4}}; 3.^f`,
			object.NewPanInt(7),
		},
		{
			`f := {a: "A", call: m{|s| .a + s}}; "B".^f`,
			object.NewPanStr("AB"),
		},
		// built-in cannot be used because signature is wrong (no addtional self)
		{
			`f := {call: Int['-%]}; 2.^f`,
			object.NewTypeErr("prop 'call in varcall must be func"),
		},
		// if proto has call prop, use it
		{
			`f := {call: m{|i| i + 6}}.bear; 3.^f`,
			object.NewPanInt(9),
		},
		// if call is not func, raise error
		{
			`f := {call: 1}; 3.^f`,
			object.NewTypeErr("prop 'call in varcall must be func"),
		},
		// if an error occurred while calling 'call, raise it
		{
			`f := {call: _}; 3.^f`,
			object.BuiltInNotImplemented,
		},
		// _missing is not substitute of call
		{
			`f := {_missing: m{|propName, x| x * 2}}; 3.^f`,
			object.NewTypeErr("literal call must be func"),
		},
	}

	for _, tt := range tests {
		actual := testEval(t, tt.input)
		testValue(t, actual, tt.expected)
	}
}

func TestEvalLonelyChain(t *testing.T) {
	tests := []struct {
		input    string
		expected object.PanObject
	}{
		// propcall
		{
			`nil&.a`,
			object.BuiltInNil,
		},
		{
			`[nil, {a: 1}]&@a`,
			object.NewPanArr(
				object.NewPanInt(1),
			),
		},
		// literalcall
		{
			`nil&.{|x| x.a}`,
			object.BuiltInNil,
		},
		{
			`[1, nil, 3]&@{|x| x * 2}`,
			object.NewPanArr(
				object.NewPanInt(2),
				object.NewPanInt(6),
			),
		},
	}

	for _, tt := range tests {
		actual := testEval(t, tt.input)
		testValue(t, actual, tt.expected)
	}
}

func TestEvalThoughtfulChain(t *testing.T) {
	tests := []struct {
		input    string
		expected object.PanObject
	}{
		// literalcall
		{
			`1~.{nil}`,
			object.NewPanInt(1),
		},
		{
			`['a, 'b]~@{|x| {a: 1}[x]}`,
			object.NewPanArr(
				object.NewPanInt(1),
				object.NewPanStr("b"),
			),
		},
		{
			`[2, 'a, 3, 'b, 4]~$(1){|acc, i| acc * i}`,
			object.NewPanInt(24),
		},
		// propcall
		{
			`{a: nil}~.a`,
			toPanObj([]object.Pair{
				{
					Key:   object.NewPanStr("a"),
					Value: object.BuiltInNil,
				},
			}),
		},
		{
			`[{a: 1}, {a: nil}]~@a`,
			object.NewPanArr(
				object.NewPanInt(1),
				toPanObj([]object.Pair{
					{
						Key:   object.NewPanStr("a"),
						Value: object.BuiltInNil,
					},
				}),
			),
		},
		{
			`[2, nil, 3, nil, 4]~$(1)*`,
			object.NewPanInt(24),
		},
		// avoid error
		{
			`{}~.a`,
			toPanObj([]object.Pair{}),
		},
	}

	for _, tt := range tests {
		actual := testEval(t, tt.input)
		testValue(t, actual, tt.expected)
	}
}

func TestEvalStrictChain(t *testing.T) {
	tests := []struct {
		input    string
		expected object.PanObject
	}{
		// literalcall
		{
			`[1, nil]=@{|i| i}`,
			object.NewPanArr(
				object.NewPanInt(1),
				object.BuiltInNil,
			),
		},
		// propcall
		{
			`[{a: 1}, {a: nil}]=@a`,
			object.NewPanArr(
				object.NewPanInt(1),
				object.BuiltInNil,
			),
		},
	}

	for _, tt := range tests {
		actual := testEval(t, tt.input)
		testValue(t, actual, tt.expected)
	}
}

func TestNameErr(t *testing.T) {
	tests := []struct {
		input    string
		expected *object.PanErr
	}{
		{
			`a`,
			object.NewNameErr("name `a` is not defined"),
		},
		// multiple lines
		{
			`1; two; 3`,
			object.NewNameErr("name `two` is not defined"),
		},
		// in arr
		{
			`[A]`,
			object.NewNameErr("name `A` is not defined"),
		},
		// in arr expansion
		{
			`[*ae]`,
			object.NewNameErr("name `ae` is not defined"),
		},
		// in obj
		{
			`{key: o}`,
			object.NewNameErr("name `o` is not defined"),
		},
		// in obj expansion
		{
			`{**oe}`,
			object.NewNameErr("name `oe` is not defined"),
		},
		// in map key
		{
			`%{key: 1}`,
			object.NewNameErr("name `key` is not defined"),
		},
		// in map val
		{
			`%{1: val}`,
			object.NewNameErr("name `val` is not defined"),
		},
		// in map expansion
		{
			`%{**me}`,
			object.NewNameErr("name `me` is not defined"),
		},
		// in func call
		{
			`{|a| fc}(1)`,
			object.NewNameErr("name `fc` is not defined"),
		},
		// in arg of func call
		{
			`{|a| 10}(afc)`,
			object.NewNameErr("name `afc` is not defined"),
		},
		// in kwarg of func call
		{
			`{|a: 1| 10}(a: kwfc)`,
			object.NewNameErr("name `kwfc` is not defined"),
		},
		// in iter call
		{
			`<{|a| ic}>.new(1).next`,
			object.NewNameErr("name `ic` is not defined"),
		},
		// in arg of iter call
		{
			`{|a| 10}.new(aic)`,
			object.NewNameErr("name `aic` is not defined"),
		},
		// in kwarg of iter call
		{
			`{|a: 1| 10}(a: kwic)`,
			object.NewNameErr("name `kwic` is not defined"),
		},
	}

	for _, tt := range tests {
		actual := testEval(t, tt.input)
		testPanErr(t, actual, tt.expected)
	}
}

func TestEvalConsts(t *testing.T) {
	tests := []struct {
		input    string
		expected object.PanObject
	}{
		{
			`Int`,
			object.BuiltInIntObj,
		},
		{
			`Float`,
			object.BuiltInFloatObj,
		},
		{
			`Num`,
			object.BuiltInNumObj,
		},
		{
			`Nil`,
			object.BuiltInNilObj,
		},
		{
			`Str`,
			object.BuiltInStrObj,
		},
		{
			`Arr`,
			object.BuiltInArrObj,
		},
		{
			`Range`,
			object.BuiltInRangeObj,
		},
		{
			`Func`,
			object.BuiltInFuncObj,
		},
		{
			`Iter`,
			object.BuiltInIterObj,
		},
		{
			`Match`,
			object.BuiltInMatchObj,
		},
		{
			`Obj`,
			object.BuiltInObjObj,
		},
		{
			`BaseObj`,
			object.BuiltInBaseObj,
		},
		{
			`Map`,
			object.BuiltInMapObj,
		},
		{
			`true`,
			object.BuiltInTrue,
		},
		{
			`false`,
			object.BuiltInFalse,
		},
		{
			`nil`,
			object.BuiltInNil,
		},
		{
			`Err`,
			object.BuiltInErrObj,
		},
		{
			`AssertionErr`,
			object.BuiltInAssertionErr,
		},
		{
			`NameErr`,
			object.BuiltInNameErr,
		},
		{
			`NoPropErr`,
			object.BuiltInNoPropErr,
		},
		{
			`NotImplementedErr`,
			object.BuiltInNotImplementedErr,
		},
		{
			`StopIterErr`,
			object.BuiltInStopIterErr,
		},
		{
			`SyntaxErr`,
			object.BuiltInSyntaxErr,
		},
		{
			`TypeErr`,
			object.BuiltInTypeErr,
		},
		{
			`ValueErr`,
			object.BuiltInValueErr,
		},
		{
			`ZeroDivisionErr`,
			object.BuiltInZeroDivisionErr,
		},
		{
			`_`,
			object.NewNotImplementedErr("Not implemented"),
		},
	}

	for _, tt := range tests {
		actual := testEval(t, tt.input)
		testValue(t, actual, tt.expected)
	}
}

func TestEvalInfix(t *testing.T) {
	tests := []struct {
		input    string
		expected object.PanObject
	}{
		{
			`1 == 1`,
			object.BuiltInTrue,
		},
		{
			`1 + 1`,
			object.NewPanInt(2),
		},
	}

	for _, tt := range tests {
		actual := testEval(t, tt.input)
		testValue(t, actual, tt.expected)
	}
}

func TestEvalInfixIntEq(t *testing.T) {
	tests := []struct {
		input    string
		expected object.PanObject
	}{
		{
			`10 == 10`,
			object.BuiltInTrue,
		},
		{
			`10 == 5`,
			object.BuiltInFalse,
		},
		{
			`1 == 1.0`,
			object.BuiltInFalse,
		},
		{
			`1 == "1"`,
			object.BuiltInFalse,
		},
		// ancestor of int is also comparable
		{
			`1 == true`,
			object.BuiltInTrue,
		},
		{
			`0 == false`,
			object.BuiltInTrue,
		},
		{
			`2.bear == 2`,
			object.BuiltInTrue,
		},
	}

	for _, tt := range tests {
		actual := testEval(t, tt.input)
		testValue(t, actual, tt.expected)
	}
}

func TestEvalInfixIntNeq(t *testing.T) {
	tests := []struct {
		input    string
		expected object.PanObject
	}{
		{
			`10 != 10`,
			object.BuiltInFalse,
		},
		{
			`10 != 5`,
			object.BuiltInTrue,
		},
		{
			`1 != 1.0`,
			object.BuiltInTrue,
		},
		{
			`1 != "1"`,
			object.BuiltInTrue,
		},
		// ancestor of int is also comparable
		{
			`1 != true`,
			object.BuiltInFalse,
		},
		{
			`0 != true`,
			object.BuiltInTrue,
		},
		{
			`2.bear != 2`,
			object.BuiltInFalse,
		},
	}

	for _, tt := range tests {
		actual := testEval(t, tt.input)
		testValue(t, actual, tt.expected)
	}
}

func TestEvalInfixFloatEq(t *testing.T) {
	tests := []struct {
		input    string
		expected object.PanObject
	}{
		{
			`1.0 == 1.0`,
			object.BuiltInTrue,
		},
		{
			`1.0 == 1.1`,
			object.BuiltInFalse,
		},
		{
			`1.0 == 'a`,
			object.BuiltInFalse,
		},
		// ancestors are also comparable
		{
			`1.0.bear == 1.0`,
			object.BuiltInTrue,
		},
	}

	for _, tt := range tests {
		actual := testEval(t, tt.input)
		testValue(t, actual, tt.expected)
	}
}

func TestEvalInfixNilEq(t *testing.T) {
	tests := []struct {
		input    string
		expected object.PanObject
	}{
		{
			`nil == nil`,
			object.BuiltInTrue,
		},
		{
			`nil == 'nil`,
			object.BuiltInFalse,
		},
		// ancestors are also comparable
		{
			`nil.bear == nil`,
			object.BuiltInTrue,
		},
		// Nil is treated as nil
		{
			`nil == Nil`,
			object.BuiltInTrue,
		},
		{
			`Nil == Nil`,
			object.BuiltInTrue,
		},
	}

	for _, tt := range tests {
		actual := testEval(t, tt.input)
		testValue(t, actual, tt.expected)
	}
}

func TestEvalInfixStrEq(t *testing.T) {
	tests := []struct {
		input    string
		expected object.PanObject
	}{
		{
			`'a == 'a`,
			object.BuiltInTrue,
		},
		{
			`'a == "a"`,
			object.BuiltInTrue,
		},
		{
			`'a == 'b`,
			object.BuiltInFalse,
		},
		{
			`'a == 'A`,
			object.BuiltInFalse,
		},
		{
			`'a == 1`,
			object.BuiltInFalse,
		},
	}

	for _, tt := range tests {
		actual := testEval(t, tt.input)
		testValue(t, actual, tt.expected)
	}
}

func TestEvalInfixArrEq(t *testing.T) {
	tests := []struct {
		input    string
		expected object.PanObject
	}{
		// true if all elements are equivalent
		{
			`[] == []`,
			object.BuiltInTrue,
		},
		{
			`[1, "a"] == [1, "a"]`,
			object.BuiltInTrue,
		},
		{
			`[] == 'a`,
			object.BuiltInFalse,
		},
		{
			`[1, 2] == [1]`,
			object.BuiltInFalse,
		},
		{
			`[1, 2] == [1, 2, 3]`,
			object.BuiltInFalse,
		},
		{
			`[1, 2] == [1, "2"]`,
			object.BuiltInFalse,
		},
		{
			`[1, [2, 3]] == [1, [2, 3]]`,
			object.BuiltInTrue,
		},
		{
			`[1, [2, 3]] == [1, [2, 4]]`,
			object.BuiltInFalse,
		},
		// ancestors are also comparable
		{
			`[2].bear == [2]`,
			object.BuiltInTrue,
		},
		// Arr is treated as []
		{
			`[] == Arr`,
			object.BuiltInTrue,
		},
		{
			`Arr == Arr`,
			object.BuiltInTrue,
		},
	}

	for _, tt := range tests {
		actual := testEval(t, tt.input)
		testValue(t, actual, tt.expected)
	}
}

func TestEvalInfixRangeEq(t *testing.T) {
	tests := []struct {
		input    string
		expected object.PanObject
	}{
		// true if all of start, stop, step are same
		{
			`(1:2) == (1:2)`,
			object.BuiltInTrue,
		},
		{
			`(1:2) == (1:3)`,
			object.BuiltInFalse,
		},
		{
			`(1:2:-1) == (1:2:-1)`,
			object.BuiltInTrue,
		},
		{
			`(1:2:-1) == (0:2:-1)`,
			object.BuiltInFalse,
		},
		{
			`(1:2:-1) == (1:3:-1)`,
			object.BuiltInFalse,
		},
		{
			`(1:2:-1) == (1:2:-2)`,
			object.BuiltInFalse,
		},
		{
			`(1:2) == (1:2:nil)`,
			object.BuiltInTrue,
		},
		// ancestors are also comparable
		{
			`(1:2).bear == (1:2)`,
			object.BuiltInTrue,
		},
	}

	for _, tt := range tests {
		actual := testEval(t, tt.input)
		testValue(t, actual, tt.expected)
	}
}

func TestEvalInfixFuncEq(t *testing.T) {
	tests := []struct {
		input    string
		expected object.PanObject
	}{
		// true only if Inspect() is same
		{
			`f := {|x| x}; f == f`,
			object.BuiltInTrue,
		},
		{
			`{|x| x} == {|x| x}`,
			object.BuiltInTrue,
		},
		{
			`{|x| x} == 1`,
			object.BuiltInFalse,
		},
		{
			`{|x| x} == {|x| 1}`,
			object.BuiltInFalse,
		},
		{
			`
			{|x|
				x
			} == {|x| x}
			`,
			object.BuiltInTrue,
		},
		// even though funcs always return the same value, they seems to be different
		// if src is different
		{
			`
			{|x| 1; x} == {|x| x}
			`,
			object.BuiltInFalse,
		},
		// ancestors are also comparable
		{
			`{|x| x}.bear == {|x| x}`,
			object.BuiltInTrue,
		},
	}

	for _, tt := range tests {
		actual := testEval(t, tt.input)
		testValue(t, actual, tt.expected)
	}
}

func TestEvalInfixObjEq(t *testing.T) {
	tests := []struct {
		input    string
		expected object.PanObject
	}{
		// true if all pairs are same
		{
			`{} == {}`,
			object.BuiltInTrue,
		},
		{
			`{} == "a"`,
			object.BuiltInFalse,
		},
		{
			`{a: 1, b: 2} == {a: 1, b: 2}`,
			object.BuiltInTrue,
		},
		{
			`{a: 1, b: 2} == {b: 2, a: 1}`,
			object.BuiltInTrue,
		},
		{
			`{a: 1, b: 2} == {a: 1, b: 1}`,
			object.BuiltInFalse,
		},
		{
			`{a: 1, b: 2, c: 3} == {a: 1, b: 2}`,
			object.BuiltInFalse,
		},
		{
			`{a: 1, b: 2} == {a: 1}`,
			object.BuiltInFalse,
		},
		{
			`{a: 1, b: {c: 3}} == {a: 1, b: {c: 3}}`,
			object.BuiltInTrue,
		},
		{
			`{a: 1, b: {c: 3}} == {a: 1, b: {c: 4}}`,
			object.BuiltInFalse,
		},
		// `==` is defined in BaseObj
		{
			`BaseObj == BaseObj`,
			object.BuiltInTrue,
		},
		// check ansestor
		// NOTE: == comparison is nothing to do with proto hierarchy!
		// ancestors are also comparable
		{
			`{}.bear({a: 1}) == {a: 1}`,
			object.BuiltInTrue,
		},
		{
			`{a: 1}.bear({b: 2}) == {a: 1}`,
			object.BuiltInFalse,
		},
	}

	for _, tt := range tests {
		actual := testEval(t, tt.input)
		testValue(t, actual, tt.expected)
	}
}

func TestEvalInfixMapEq(t *testing.T) {
	tests := []struct {
		input    string
		expected object.PanObject
	}{
		// true if all pairs are same
		{
			`%{} == %{}`,
			object.BuiltInTrue,
		},
		{
			`%{} == "a"`,
			object.BuiltInFalse,
		},
		{
			`%{} == {}`,
			object.BuiltInFalse,
		},
		{
			`%{0: 1, nil: 2} == %{0: 1, nil: 2}`,
			object.BuiltInTrue,
		},
		{
			`%{0: 1, nil: 2} == %{nil: 2, 0: 1}`,
			object.BuiltInTrue,
		},
		{
			`%{0: 1, nil: 2} == %{0: 1, nil: 3}`,
			object.BuiltInFalse,
		},
		{
			`%{0: 1, nil: 2, 'extra: 3} == %{0: 1, nil: 2}`,
			object.BuiltInFalse,
		},
		{
			`%{0: 1, 2: %{3: 4}} == %{0: 1, 2: %{3: 4}}`,
			object.BuiltInTrue,
		},
		{
			`%{0: 1, 2: %{3: 4}} == %{0: 1, 2: %{3: 5}}`,
			object.BuiltInFalse,
		},
		// map with non-hashable keys
		{
			`%{[1]: 1, [2]: 2} == %{[1]: 1, [2]: 2}`,
			object.BuiltInTrue,
		},
		{
			`%{[1]: 1, [2]: 2} == %{[2]: 2, [1]: 1}`,
			object.BuiltInTrue,
		},
		{
			`%{[1]: 1, [2]: 2} == %{[1]: 1, [2]: 3}`,
			object.BuiltInFalse,
		},
		{
			`%{[1]: 1, [2]: 2} == %{[2]: 3, [1]: 1}`,
			object.BuiltInFalse,
		},
		{
			`%{[1]: 1, [2]: 2, [3]: 3} == %{[1]: 1, [2]: 2}`,
			object.BuiltInFalse,
		},
		{
			`%{[1]: 1, [2]: 2} == %{[1]: 1}`,
			object.BuiltInFalse,
		},
		{
			`%{'a: 0, [1]: 1} == %{'a: 0, [1]: 1}`,
			object.BuiltInTrue,
		},
		// ancestors are also comparable
		{
			`%{'a: 0}.bear == %{'a: 0}`,
			object.BuiltInTrue,
		},
		// Map is treated as %{}
		{
			`%{} == Map`,
			object.BuiltInTrue,
		},
		{
			`Map == Map`,
			object.BuiltInTrue,
		},
	}

	for _, tt := range tests {
		actual := testEval(t, tt.input)
		testValue(t, actual, tt.expected)
	}
}

func TestEvalInfixBuiltInFuncEq(t *testing.T) {
	tests := []struct {
		input    string
		expected object.PanObject
	}{
		// true if pointers are same
		{
			`BaseObj['==] == BaseObj['==]`,
			object.BuiltInTrue,
		},
		{
			`BaseObj['==] == BaseObj['at]`,
			object.BuiltInFalse,
		},
	}

	for _, tt := range tests {
		actual := testEval(t, tt.input)
		testValue(t, actual, tt.expected)
	}
}

func TestEvalInfixErrEq(t *testing.T) {
	tests := []struct {
		input    string
		expected object.PanObject
	}{
		{
			`Err == Err`,
			object.BuiltInTrue,
		},
		{
			`AssertionErr == AssertionErr`,
			object.BuiltInTrue,
		},
		{
			`TypeErr == ValueErr`,
			object.BuiltInFalse,
		},
		// this is called only by errWrappers
		{
			`1.try.fmap {\/0}.err == 1.try.fmap {\/0}.err`,
			object.BuiltInTrue,
		},
		{
			`1.try.fmap {\/0}.err == 1.try.fmap {a}.err`,
			object.BuiltInFalse,
		},
		// only message is different
		{
			`1.try.fmap {raise Err.new("a")}.err == 1.try.fmap {raise Err.new("b")}.err`,
			object.BuiltInFalse,
		},
		// only type is different
		{
			`1.try.fmap {raise TypeErr.new("a")}.err == 1.try.fmap {raise ValueErr.new("a")}.err`,
			object.BuiltInFalse,
		},
	}

	for _, tt := range tests {
		actual := testEval(t, tt.input)
		testValue(t, actual, tt.expected)
	}
}

func TestEvalInfixConstsEq(t *testing.T) {
	tests := []struct {
		input    string
		expected object.PanObject
	}{
		{
			`Arr == Arr`,
			object.BuiltInTrue,
		},
		{
			`BaseObj == BaseObj`,
			object.BuiltInTrue,
		},
		{
			`Float == Float`,
			object.BuiltInTrue,
		},
		{
			`Func == Func`,
			object.BuiltInTrue,
		},
		{
			`Int == Int`,
			object.BuiltInTrue,
		},
		{
			`Iter == Iter`,
			object.BuiltInTrue,
		},
		{
			`Map == Map`,
			object.BuiltInTrue,
		},
		{
			`Nil == Nil`,
			object.BuiltInTrue,
		},
		{
			`Obj == Obj`,
			object.BuiltInTrue,
		},
		{
			`Range == Range`,
			object.BuiltInTrue,
		},
		{
			`Str == Str`,
			object.BuiltInTrue,
		},
	}

	for _, tt := range tests {
		actual := testEval(t, tt.input)
		testValue(t, actual, tt.expected)
	}
}

func TestEvalInfixIntSpaceship(t *testing.T) {
	tests := []struct {
		input    string
		expected object.PanObject
	}{
		{
			`2 <=> 2`,
			object.NewPanInt(0),
		},
		{
			`2 <=> 3`,
			object.NewPanInt(-1),
		},
		{
			`3 <=> 2`,
			object.NewPanInt(1),
		},
		// ancestor of int is also comparable
		{
			`1 <=> true`,
			object.NewPanInt(0),
		},
	}

	for _, tt := range tests {
		actual := testEval(t, tt.input)
		testValue(t, actual, tt.expected)
	}
}

func TestEvalInfixStrSpaceship(t *testing.T) {
	tests := []struct {
		input    string
		expected object.PanObject
	}{
		{
			`'b <=> 'b`,
			object.NewPanInt(0),
		},
		{
			`'a <=> 'b`,
			object.NewPanInt(-1),
		},
		{
			`'c <=> 'b`,
			object.NewPanInt(1),
		},
		// ancestor of int is also comparable
		{
			`'b <=> 'b.bear`,
			object.NewPanInt(0),
		},
	}

	for _, tt := range tests {
		actual := testEval(t, tt.input)
		testValue(t, actual, tt.expected)
	}
}

func TestEvalInfixFloatSpaceship(t *testing.T) {
	tests := []struct {
		input    string
		expected object.PanObject
	}{
		{
			`1.0 <=> 1.0`,
			object.NewPanInt(0),
		},
		{
			`1.5 <=> 2.0`,
			object.NewPanInt(-1),
		},
		{
			`2.5 <=> 2.0`,
			object.NewPanInt(1),
		},
		// ancestor of float is also comparable
		{
			`1.0 <=> 1.0.bear`,
			object.NewPanInt(0),
		},
	}

	for _, tt := range tests {
		actual := testEval(t, tt.input)
		testValue(t, actual, tt.expected)
	}
}

func TestEvalInfixIntAdd(t *testing.T) {
	tests := []struct {
		input    string
		expected object.PanObject
	}{
		{
			`1 + 1`,
			object.NewPanInt(2),
		},
		{
			`-5 + 10`,
			object.NewPanInt(5),
		},
		// decendant of int can be added
		{
			`3 + true`,
			object.NewPanInt(4),
		},
		// nil is treated as 0
		{
			`3 + nil`,
			object.NewPanInt(3),
		},
		// Int is treated as 0
		{
			`3 + Int`,
			object.NewPanInt(3),
		},
		{
			`Int + 3`,
			object.NewPanInt(3),
		},
	}

	for _, tt := range tests {
		actual := testEval(t, tt.input)
		testValue(t, actual, tt.expected)
	}
}

func TestEvalInfixIntAddErr(t *testing.T) {
	tests := []struct {
		input    string
		expected object.PanObject
	}{
		{
			`1 + []`,
			object.NewTypeErr("[] cannot be treated as int"),
		},
		{
			`1['+]({}, 2)`,
			object.NewTypeErr("{} cannot be treated as int"),
		},
		{
			`1.+`,
			object.NewTypeErr("+ requires at least 2 args"),
		},
	}

	for _, tt := range tests {
		actual := testEval(t, tt.input)
		testValue(t, actual, tt.expected)
	}
}

func TestEvalInfixIntSub(t *testing.T) {
	tests := []struct {
		input    string
		expected object.PanObject
	}{
		{
			`3 - 1`,
			object.NewPanInt(2),
		},
		{
			`5 - -10`,
			object.NewPanInt(15),
		},
		// decendant of int can be subbed
		{
			`3 - true`,
			object.NewPanInt(2),
		},
		// nil is treated as 0
		{
			`3 - nil`,
			object.NewPanInt(3),
		},
		// Int is treated as 0
		{
			`3 - Int`,
			object.NewPanInt(3),
		},
		{
			`Int - 3`,
			object.NewPanInt(-3),
		},
	}

	for _, tt := range tests {
		actual := testEval(t, tt.input)
		testValue(t, actual, tt.expected)
	}
}

func TestEvalInfixIntSubErr(t *testing.T) {
	tests := []struct {
		input    string
		expected object.PanObject
	}{
		{
			`1 - []`,
			object.NewTypeErr("[] cannot be treated as int"),
		},
		{
			`1['-]({}, 2)`,
			object.NewTypeErr("{} cannot be treated as int"),
		},
		{
			`1.-`,
			object.NewTypeErr("- requires at least 2 args"),
		},
	}

	for _, tt := range tests {
		actual := testEval(t, tt.input)
		testValue(t, actual, tt.expected)
	}
}

func TestEvalInfixIntMul(t *testing.T) {
	tests := []struct {
		input    string
		expected object.PanObject
	}{
		{
			`2 * 3`,
			object.NewPanInt(6),
		},
		// decendant of int can be added
		{
			`3 * true`,
			object.NewPanInt(3),
		},
		// nil is treated as 1
		{
			`3 * nil`,
			object.NewPanInt(3),
		},
		// Int is treated as 0
		{
			`3 * Int`,
			object.NewPanInt(0),
		},
		{
			`Int * 3`,
			object.NewPanInt(0),
		},
	}

	for _, tt := range tests {
		actual := testEval(t, tt.input)
		testValue(t, actual, tt.expected)
	}
}

func TestEvalInfixIntMulErr(t *testing.T) {
	tests := []struct {
		input    string
		expected object.PanObject
	}{
		{
			`1 * []`,
			object.NewTypeErr("[] cannot be treated as int"),
		},
		{
			`1['*]({}, 2)`,
			object.NewTypeErr("{} cannot be treated as int"),
		},
		{
			`1.*`,
			object.NewTypeErr("* requires at least 2 args"),
		},
	}

	for _, tt := range tests {
		actual := testEval(t, tt.input)
		testValue(t, actual, tt.expected)
	}
}

func TestEvalInfixIntDiv(t *testing.T) {
	tests := []struct {
		input    string
		expected object.PanObject
	}{
		{
			`10 / 5`,
			object.NewPanFloat(2.0),
		},
		{
			`-5 / 4`,
			object.NewPanFloat(-1.25),
		},
		// decendant of int can be devided
		{
			`3 / true`,
			object.NewPanFloat(3.0),
		},
		// nil is treated as 1
		{
			`5 / nil`,
			object.NewPanFloat(5.0),
		},
		// Int is treated as 0
		{
			`Int / 3`,
			object.NewPanFloat(0.0),
		},
	}

	for _, tt := range tests {
		actual := testEval(t, tt.input)
		testValue(t, actual, tt.expected)
	}
}

func TestEvalInfixIntDivErr(t *testing.T) {
	tests := []struct {
		input    string
		expected object.PanObject
	}{
		{
			`1 / []`,
			object.NewTypeErr("[] cannot be treated as int"),
		},
		{
			`1['/]({}, 2)`,
			object.NewTypeErr("{} cannot be treated as int"),
		},
		{
			`1./`,
			object.NewTypeErr("/ requires at least 2 args"),
		},
		{
			`1 / 0`,
			object.NewZeroDivisionErr("cannot be divided by 0"),
		},
	}

	for _, tt := range tests {
		actual := testEval(t, tt.input)
		testValue(t, actual, tt.expected)
	}
}

func TestEvalInfixIntFloorDiv(t *testing.T) {
	tests := []struct {
		input    string
		expected object.PanObject
	}{
		{
			`10 // 5`,
			object.NewPanInt(2),
		},
		{
			`-5 // 2`,
			object.NewPanInt(-3),
		},
		// decendant of int can be devided
		{
			`3 // true`,
			object.NewPanInt(3),
		},
		// nil is treated as 1
		{
			`5 // nil`,
			object.NewPanInt(5),
		},
		// Int is treated as 0
		{
			`Int // 3`,
			object.NewPanInt(0),
		},
	}

	for _, tt := range tests {
		actual := testEval(t, tt.input)
		testValue(t, actual, tt.expected)
	}
}

func TestEvalInfixIntFloorDivErr(t *testing.T) {
	tests := []struct {
		input    string
		expected object.PanObject
	}{
		{
			`1 // []`,
			object.NewTypeErr("[] cannot be treated as int"),
		},
		{
			`1['//]({}, 2)`,
			object.NewTypeErr("{} cannot be treated as int"),
		},
		{
			`1.//`,
			object.NewTypeErr("// requires at least 2 args"),
		},
		{
			`1 // 0`,
			object.NewZeroDivisionErr("cannot be divided by 0"),
		},
	}

	for _, tt := range tests {
		actual := testEval(t, tt.input)
		testValue(t, actual, tt.expected)
	}
}

func TestEvalInfixStrAdd(t *testing.T) {
	tests := []struct {
		input    string
		expected object.PanObject
	}{
		{
			`'a + 'b`,
			object.NewPanStr("ab"),
		},
		{
			`"にほ" + "んご"`,
			object.NewPanStr("にほんご"),
		},
		// decendant of str can be added
		{
			`"c" + "d".bear`,
			object.NewPanStr("cd"),
		},
		// nil is treated as ""
		{
			`"abc" + nil`,
			object.NewPanStr("abc"),
		},
	}

	for _, tt := range tests {
		actual := testEval(t, tt.input)
		testValue(t, actual, tt.expected)
	}
}

func TestEvalInfixStrAddErr(t *testing.T) {
	tests := []struct {
		input    string
		expected object.PanObject
	}{
		{
			`"a" + []`,
			object.NewTypeErr("[] cannot be treated as str"),
		},
		{
			`""['+]({}, "b")`,
			object.NewTypeErr("{} cannot be treated as str"),
		},
		{
			`"a".+`,
			object.NewTypeErr("+ requires at least 2 args"),
		},
	}

	for _, tt := range tests {
		actual := testEval(t, tt.input)
		testValue(t, actual, tt.expected)
	}
}

func TestEvalInfixStrMul(t *testing.T) {
	tests := []struct {
		input    string
		expected object.PanObject
	}{
		{
			`"w" * 3`,
			object.NewPanStr("www"),
		},
	}

	for _, tt := range tests {
		actual := testEval(t, tt.input)
		testValue(t, actual, tt.expected)
	}
}

func TestEvalInfixStrMulErr(t *testing.T) {
	tests := []struct {
		input    string
		expected object.PanObject
	}{
		{
			`"a" * -1`,
			object.NewValueErr("-1 is not positive"),
		},
		{
			`"w" * []`,
			object.NewTypeErr("[] cannot be treated as int"),
		},
		{
			`'x['*]({}, 2)`,
			object.NewTypeErr("{} cannot be treated as str"),
		},
		{
			`?y.*`,
			object.NewTypeErr("* requires at least 2 args"),
		},
	}

	for _, tt := range tests {
		actual := testEval(t, tt.input)
		testValue(t, actual, tt.expected)
	}
}

func TestEvalInfixStrDiv(t *testing.T) {
	tests := []struct {
		input    string
		expected object.PanObject
	}{
		{
			`"hello world" / " "`,
			object.NewPanArr(
				object.NewPanStr("hello"),
				object.NewPanStr("world"),
			),
		},
		{
			`"a  b    c " / " "`,
			object.NewPanArr(
				object.NewPanStr("a"),
				object.NewPanStr("b"),
				object.NewPanStr("c"),
			),
		},
		// use regex
		{
			"`abc1de23f` / `\\d`",
			object.NewPanArr(
				object.NewPanStr("abc"),
				object.NewPanStr("de"),
				object.NewPanStr("f"),
			),
		},
		{
			"`123 \n45  \t6` / `\\s`",
			object.NewPanArr(
				object.NewPanStr("123"),
				object.NewPanStr("45"),
				object.NewPanStr("6"),
			),
		},
		// if no separators found, return whole string
		{
			"`a` / `x`",
			object.NewPanArr(
				object.NewPanStr("a"),
			),
		},
		// decendant of str can be added
		{
			`"a,b" / ",".bear`,
			object.NewPanArr(
				object.NewPanStr("a"),
				object.NewPanStr("b"),
			),
		},
		// if \2 is illegal for regex syntax, raise an error
		{
			"`a` / `[`",
			object.NewValueErr(`"[" is invalid regex pattern`),
		},
	}

	for _, tt := range tests {
		actual := testEval(t, tt.input)
		testValue(t, actual, tt.expected)
	}
}

func TestEvalInfixStrDivErr(t *testing.T) {
	tests := []struct {
		input    string
		expected object.PanObject
	}{
		{
			`"w" / []`,
			object.NewTypeErr("[] cannot be treated as str"),
		},
		{
			`'x['/]({}, "a")`,
			object.NewTypeErr("{} cannot be treated as str"),
		},
		{
			`?y./`,
			object.NewTypeErr("/ requires at least 2 args"),
		},
	}

	for _, tt := range tests {
		actual := testEval(t, tt.input)
		testValue(t, actual, tt.expected)
	}
}

func TestEvalInfixFloatAdd(t *testing.T) {
	tests := []struct {
		input    string
		expected object.PanObject
	}{
		{
			`1.0 + 1.0`,
			object.NewPanFloat(2.0),
		},
		{
			`-5.0 + 10.0`,
			object.NewPanFloat(5.0),
		},
		// decendant of int can be added
		{
			`3.0 + 1.0.bear`,
			object.NewPanFloat(4.0),
		},
		// nil is treated as 0.0
		{
			`3.0 + nil`,
			object.NewPanFloat(3.0),
		},
		// Float is treated as 0
		{
			`3.0 + Float`,
			object.NewPanFloat(3.0),
		},
		{
			`Float + 3.0`,
			object.NewPanFloat(3.0),
		},
	}

	for _, tt := range tests {
		actual := testEval(t, tt.input)
		testValue(t, actual, tt.expected)
	}
}

func TestEvalInfixFloatAddErr(t *testing.T) {
	tests := []struct {
		input    string
		expected object.PanObject
	}{
		{
			`1.0 + []`,
			object.NewTypeErr("[] cannot be treated as float"),
		},
		{
			`0.0['+]({}, 1.0)`,
			object.NewTypeErr("{} cannot be treated as float"),
		},
		{
			`1.0.+`,
			object.NewTypeErr("+ requires at least 2 args"),
		},
	}

	for _, tt := range tests {
		actual := testEval(t, tt.input)
		testValue(t, actual, tt.expected)
	}
}

func TestEvalInfixFloatSub(t *testing.T) {
	tests := []struct {
		input    string
		expected object.PanObject
	}{
		{
			`3.0 - 1.0`,
			object.NewPanFloat(2.0),
		},
		{
			`5.0 - -10.0`,
			object.NewPanFloat(15.0),
		},
		// decendant of int can be added
		{
			`4.0 - 1.0.bear`,
			object.NewPanFloat(3.0),
		},
		// nil is treated as 0.0
		{
			`3.0 - nil`,
			object.NewPanFloat(3.0),
		},
		// Float is treated as 0
		{
			`3.0 - Float`,
			object.NewPanFloat(3.0),
		},
		{
			`Float - 3.0`,
			object.NewPanFloat(-3.0),
		},
	}

	for _, tt := range tests {
		actual := testEval(t, tt.input)
		testValue(t, actual, tt.expected)
	}
}

func TestEvalInfixFloatSubErr(t *testing.T) {
	tests := []struct {
		input    string
		expected object.PanObject
	}{
		{
			`1.0 - []`,
			object.NewTypeErr("[] cannot be treated as float"),
		},
		{
			`0.0['-]({}, 1.0)`,
			object.NewTypeErr("{} cannot be treated as float"),
		},
		{
			`1.0.-`,
			object.NewTypeErr("- requires at least 2 args"),
		},
	}

	for _, tt := range tests {
		actual := testEval(t, tt.input)
		testValue(t, actual, tt.expected)
	}
}

func TestEvalInfixArrAdd(t *testing.T) {
	tests := []struct {
		input    string
		expected object.PanObject
	}{
		{
			`[1] + [2]`,
			object.NewPanArr(
				object.NewPanInt(1),
				object.NewPanInt(2),
			),
		},
		// decendant of arr can be added
		{
			`[1] + [3].bear`,
			object.NewPanArr(
				object.NewPanInt(1),
				object.NewPanInt(3),
			),
		},
		// nil is treated as []
		{
			`[5] + nil`,
			object.NewPanArr(
				object.NewPanInt(5),
			),
		},
		// Arr is treated as []
		{
			`[5] + Arr`,
			object.NewPanArr(
				object.NewPanInt(5),
			),
		},
	}

	for _, tt := range tests {
		actual := testEval(t, tt.input)
		testValue(t, actual, tt.expected)
	}
}

func TestEvalInfixArrAddErr(t *testing.T) {
	tests := []struct {
		input    string
		expected object.PanObject
	}{
		{
			`[] + 1`,
			object.NewTypeErr("1 cannot be treated as arr"),
		},
		{
			`[]['+]({}, [1])`,
			object.NewTypeErr("{} cannot be treated as arr"),
		},
		{
			`[].+`,
			object.NewTypeErr("+ requires at least 2 args"),
		},
	}

	for _, tt := range tests {
		actual := testEval(t, tt.input)
		testValue(t, actual, tt.expected)
	}
}

func TestEvalInfixNilAdd(t *testing.T) {
	tests := []struct {
		input    string
		expected object.PanObject
	}{
		// anything can be added to nil (and nil works as zero value).
		{
			`nil + 1`,
			object.NewPanInt(1),
		},
		{
			`nil + "a"`,
			object.NewPanStr("a"),
		},
		{
			`nil + 1.0`,
			object.NewPanFloat(1.0),
		},
		{
			`nil + [5]`,
			object.NewPanArr(
				object.NewPanInt(5),
			),
		},
	}

	for _, tt := range tests {
		actual := testEval(t, tt.input)
		testValue(t, actual, tt.expected)
	}
}

func TestEvalInfixNilAddErr(t *testing.T) {
	tests := []struct {
		input    string
		expected object.PanObject
	}{
		{
			`nil.+`,
			object.NewTypeErr("+ requires at least 2 args"),
		},
	}

	for _, tt := range tests {
		actual := testEval(t, tt.input)
		testValue(t, actual, tt.expected)
	}
}

func TestEvalInfixIntMod(t *testing.T) {
	tests := []struct {
		input    string
		expected object.PanObject
	}{
		{
			`10 % 2`,
			object.NewPanInt(0),
		},
		{
			`10 % 3`,
			object.NewPanInt(1),
		},
		// decendant of int can be added
		{
			`4 % true`,
			object.NewPanInt(0),
		},
	}

	for _, tt := range tests {
		actual := testEval(t, tt.input)
		testValue(t, actual, tt.expected)
	}
}

func TestEvalInfixIntModErr(t *testing.T) {
	tests := []struct {
		input    string
		expected object.PanObject
	}{
		{
			`10 % 0`,
			object.NewZeroDivisionErr("cannot be divided by 0"),
		},
		{
			`1 % []`,
			object.NewTypeErr("[] cannot be treated as int"),
		},
		{
			`1['%]({}, 2)`,
			object.NewTypeErr("{} cannot be treated as int"),
		},
		{
			`1.%`,
			object.NewTypeErr("% requires at least 2 args"),
		},
	}

	for _, tt := range tests {
		actual := testEval(t, tt.input)
		testValue(t, actual, tt.expected)
	}
}

func TestEvalPrefix(t *testing.T) {
	tests := []struct {
		input    string
		expected object.PanObject
	}{
		{
			`!false`,
			object.BuiltInTrue,
		},
		// NOTE: negative number (like `-1`) is treated as literal by parser
		{
			`a := 1; -a`,
			object.NewPanInt(-1),
		},
		{
			`/~1`,
			object.NewPanInt(-2),
		},
	}

	for _, tt := range tests {
		actual := testEval(t, tt.input)
		testValue(t, actual, tt.expected)
	}
}

func TestEvalPrefixNot(t *testing.T) {
	tests := []struct {
		input    string
		expected object.PanObject
	}{
		{
			`!false`,
			object.BuiltInTrue,
		},
		{
			`!true`,
			object.BuiltInFalse,
		},
		{
			`!0`,
			object.BuiltInTrue,
		},
		{
			`!100`,
			object.BuiltInFalse,
		},
		{
			`!""`,
			object.BuiltInTrue,
		},
		{
			`!'a`,
			object.BuiltInFalse,
		},
	}

	for _, tt := range tests {
		actual := testEval(t, tt.input)
		testValue(t, actual, tt.expected)
	}
}

func TestEvalErrType(t *testing.T) {
	tests := []struct {
		input    string
		expected object.PanObject
	}{
		{
			`{}.try.fmap {Err.new("err")}.err.type`,
			object.BuiltInErrObj,
		},
		{
			`{}.try.fmap {AssertionErr.new("err")}.err.type`,
			object.BuiltInAssertionErr,
		},
		{
			`{}.try.fmap {NameErr.new("err")}.err.type`,
			object.BuiltInNameErr,
		},
		{
			`{}.try.fmap {NoPropErr.new("err")}.err.type`,
			object.BuiltInNoPropErr,
		},
		{
			`{}.try.fmap {NotImplementedErr.new("err")}.err.type`,
			object.BuiltInNotImplementedErr,
		},
		{
			`{}.try.fmap {StopIterErr.new("err")}.err.type`,
			object.BuiltInStopIterErr,
		},
		{
			`{}.try.fmap {SyntaxErr.new("err")}.err.type`,
			object.BuiltInSyntaxErr,
		},
		{
			`{}.try.fmap {TypeErr.new("err")}.err.type`,
			object.BuiltInTypeErr,
		},
		{
			`{}.try.fmap {ValueErr.new("err")}.err.type`,
			object.BuiltInValueErr,
		},
		{
			`{}.try.fmap {ZeroDivisionErr.new("err")}.err.type`,
			object.BuiltInZeroDivisionErr,
		},
		// raise error
		{
			`{}.try.fmap {Err.new("err")}.err['type]()`,
			object.NewTypeErr("Err#type requires at least 1 arg"),
		},
	}

	for _, tt := range tests {
		actual := testEval(t, tt.input)
		testValue(t, actual, tt.expected)
	}
}

func TestEvalErrMsg(t *testing.T) {
	tests := []struct {
		input    string
		expected object.PanObject
	}{
		{
			`{}.try.fmap {Err.new("new error")}.err.msg`,
			object.NewPanStr("new error"),
		},
		// raise error
		{
			`{}.try.fmap {Err.new("new error")}.err['msg]()`,
			object.NewTypeErr("Err#msg requires at least 1 arg"),
		},
		{
			`{}.try.fmap {Err.new("new error")}.err['msg](1)`,
			object.NewTypeErr("1 cannot be treated as err"),
		},
	}

	for _, tt := range tests {
		actual := testEval(t, tt.input)
		testValue(t, actual, tt.expected)
	}
}

func TestEvalErrConstructor(t *testing.T) {
	tests := []struct {
		input    string
		expected object.PanObject
	}{
		{
			`Err.new("new error")`,
			object.NewPanErr("new error"),
		},
		// args are converted to str by .S
		{
			`Err.new(1)`,
			object.NewPanErr("1"),
		},
		{
			`Err.new({a:1})`,
			object.NewPanErr(`{"a": 1}`),
		},
		{
			`Err.new()`,
			object.NewPanErr("nil"),
		},
	}

	for _, tt := range tests {
		actual := testEval(t, tt.input)
		testValue(t, actual, tt.expected)
	}
}

func TestEvalAssertionErrConstructor(t *testing.T) {
	tests := []struct {
		input    string
		expected object.PanObject
	}{
		{
			`AssertionErr.new("new error")`,
			object.NewAssertionErr("new error"),
		},
		// args are converted to str by .S
		{
			`AssertionErr.new(1)`,
			object.NewAssertionErr("1"),
		},
		{
			`AssertionErr.new({a:1})`,
			object.NewAssertionErr(`{"a": 1}`),
		},
		{
			`AssertionErr.new()`,
			object.NewAssertionErr("nil"),
		},
	}

	for _, tt := range tests {
		actual := testEval(t, tt.input)
		testValue(t, actual, tt.expected)
	}
}

func TestEvalNameErrConstructor(t *testing.T) {
	tests := []struct {
		input    string
		expected object.PanObject
	}{
		{
			`NameErr.new("new error")`,
			object.NewNameErr("new error"),
		},
		// args are converted to str by .S
		{
			`NameErr.new(1)`,
			object.NewNameErr("1"),
		},
		{
			`NameErr.new({a:1})`,
			object.NewNameErr(`{"a": 1}`),
		},
		{
			`NameErr.new()`,
			object.NewNameErr("nil"),
		},
	}

	for _, tt := range tests {
		actual := testEval(t, tt.input)
		testValue(t, actual, tt.expected)
	}
}

func TestEvalNoPropErrConstructor(t *testing.T) {
	tests := []struct {
		input    string
		expected object.PanObject
	}{
		{
			`NoPropErr.new("new error")`,
			object.NewNoPropErr("new error"),
		},
		// args are converted to str by .S
		{
			`NoPropErr.new(1)`,
			object.NewNoPropErr("1"),
		},
		{
			`NoPropErr.new({a:1})`,
			object.NewNoPropErr(`{"a": 1}`),
		},
		{
			`NoPropErr.new()`,
			object.NewNoPropErr("nil"),
		},
	}

	for _, tt := range tests {
		actual := testEval(t, tt.input)
		testValue(t, actual, tt.expected)
	}
}

func TestEvalNotImplementedErrConstructor(t *testing.T) {
	tests := []struct {
		input    string
		expected object.PanObject
	}{
		{
			`NotImplementedErr.new("new error")`,
			object.NewNotImplementedErr("new error"),
		},
		// args are converted to str by .S
		{
			`NotImplementedErr.new(1)`,
			object.NewNotImplementedErr("1"),
		},
		{
			`NotImplementedErr.new({a:1})`,
			object.NewNotImplementedErr(`{"a": 1}`),
		},
		{
			`NotImplementedErr.new()`,
			object.NewNotImplementedErr("nil"),
		},
	}

	for _, tt := range tests {
		actual := testEval(t, tt.input)
		testValue(t, actual, tt.expected)
	}
}

func TestEvalStopIterErrConstructor(t *testing.T) {
	tests := []struct {
		input    string
		expected object.PanObject
	}{
		{
			`StopIterErr.new("new error")`,
			object.NewStopIterErr("new error"),
		},
		// args are converted to str by .S
		{
			`StopIterErr.new(1)`,
			object.NewStopIterErr("1"),
		},
		{
			`StopIterErr.new({a:1})`,
			object.NewStopIterErr(`{"a": 1}`),
		},
		{
			`StopIterErr.new()`,
			object.NewStopIterErr("nil"),
		},
	}

	for _, tt := range tests {
		actual := testEval(t, tt.input)
		testValue(t, actual, tt.expected)
	}
}

func TestEvalSyntaxErrConstructor(t *testing.T) {
	tests := []struct {
		input    string
		expected object.PanObject
	}{
		{
			`SyntaxErr.new("new error")`,
			object.NewSyntaxErr("new error"),
		},
		// args are converted to str by .S
		{
			`SyntaxErr.new(1)`,
			object.NewSyntaxErr("1"),
		},
		{
			`SyntaxErr.new({a:1})`,
			object.NewSyntaxErr(`{"a": 1}`),
		},
		{
			`SyntaxErr.new()`,
			object.NewSyntaxErr("nil"),
		},
	}

	for _, tt := range tests {
		actual := testEval(t, tt.input)
		testValue(t, actual, tt.expected)
	}
}

func TestEvalTypeErrConstructor(t *testing.T) {
	tests := []struct {
		input    string
		expected object.PanObject
	}{
		{
			`TypeErr.new("new error")`,
			object.NewTypeErr("new error"),
		},
		// args are converted to str by .S
		{
			`TypeErr.new(1)`,
			object.NewTypeErr("1"),
		},
		{
			`TypeErr.new({a:1})`,
			object.NewTypeErr(`{"a": 1}`),
		},
		{
			`TypeErr.new()`,
			object.NewTypeErr("nil"),
		},
	}

	for _, tt := range tests {
		actual := testEval(t, tt.input)
		testValue(t, actual, tt.expected)
	}
}

func TestEvalValueErrConstructor(t *testing.T) {
	tests := []struct {
		input    string
		expected object.PanObject
	}{
		{
			`ValueErr.new("new error")`,
			object.NewValueErr("new error"),
		},
		// args are converted to str by .S
		{
			`ValueErr.new(1)`,
			object.NewValueErr("1"),
		},
		{
			`ValueErr.new({a:1})`,
			object.NewValueErr(`{"a": 1}`),
		},
		{
			`ValueErr.new()`,
			object.NewValueErr("nil"),
		},
	}

	for _, tt := range tests {
		actual := testEval(t, tt.input)
		testValue(t, actual, tt.expected)
	}
}

func TestEvalZeroDivisionErrConstructor(t *testing.T) {
	tests := []struct {
		input    string
		expected object.PanObject
	}{
		{
			`ZeroDivisionErr.new("new error")`,
			object.NewZeroDivisionErr("new error"),
		},
		// args are converted to str by .S
		{
			`ZeroDivisionErr.new(1)`,
			object.NewZeroDivisionErr("1"),
		},
		{
			`ZeroDivisionErr.new({a:1})`,
			object.NewZeroDivisionErr(`{"a": 1}`),
		},
		{
			`ZeroDivisionErr.new()`,
			object.NewZeroDivisionErr("nil"),
		},
	}

	for _, tt := range tests {
		actual := testEval(t, tt.input)
		testValue(t, actual, tt.expected)
	}
}

func TestEvalArrConstructor(t *testing.T) {
	tests := []struct {
		input    string
		expected object.PanObject
	}{
		{
			`Arr.new([1])`,
			object.NewPanArr(
				object.NewPanInt(1),
			),
		},
		{
			`Arr.new([1, 2].bear({child?: true}))`,
			object.NewPanArr(
				object.NewPanInt(1),
				object.NewPanInt(2),
			),
		},
		// errors
		{
			`Arr.new`,
			object.NewTypeErr("Arr#new requires at least 2 args"),
		},
		{
			`Arr.new(2)`,
			object.NewTypeErr("2 cannot be treated as arr"),
		},
	}

	for _, tt := range tests {
		actual := testEval(t, tt.input)
		testValue(t, actual, tt.expected)
	}
}

func TestEvalArrCall(t *testing.T) {
	tests := []struct {
		input    string
		expected object.PanObject
	}{
		{
			`Arr()`,
			object.NewPanArr(),
		},
		{
			`Arr(1)`,
			object.NewPanArr(
				object.NewPanInt(1),
			),
		},
		{
			`Arr('a, true)`,
			object.NewPanArr(
				object.NewPanStr("a"),
				object.BuiltInTrue,
			),
		},
		// errors
		{
			`Arr['call]()`,
			object.NewTypeErr("Arr#call requires at least 1 arg"),
		},
	}

	for _, tt := range tests {
		actual := testEval(t, tt.input)
		testValue(t, actual, tt.expected)
	}
}

func TestEvalFloatConstructor(t *testing.T) {
	tests := []struct {
		input    string
		expected object.PanObject
	}{
		{
			`Float.new(1.5)`,
			object.NewPanFloat(1.5),
		},
		{
			`Float.new(-2.5.bear({child?: true}))`,
			object.NewPanFloat(-2.5),
		},
		{
			`Float.new(2)`,
			object.NewPanFloat(2.0),
		},
		// errors
		{
			`Float.new`,
			object.NewTypeErr("Float#new requires at least 2 args"),
		},
		{
			`Float.new("a")`,
			object.NewTypeErr(`"a" cannot be treated as float`),
		},
	}

	for _, tt := range tests {
		actual := testEval(t, tt.input)
		testValue(t, actual, tt.expected)
	}
}

func TestEvalFuncConstructor(t *testing.T) {
	outerEnv := object.NewEnvWithConsts()

	tests := []struct {
		input    string
		expected object.PanObject
	}{
		{
			`Func.new({||})`,
			toPanFunc(
				[]string{},
				[]object.Pair{},
				`|| `,
				outerEnv,
			),
		},
		{
			`Func.new({||}.bear({child?: true}))`,
			toPanFunc(
				[]string{},
				[]object.Pair{},
				`|| `,
				outerEnv,
			),
		},
		// errors
		{
			`Func.new`,
			object.NewTypeErr("Func#new requires at least 2 args"),
		},
		{
			`Func.new(2)`,
			object.NewTypeErr("2 cannot be treated as func"),
		},
	}

	for _, tt := range tests {
		actual := testEvalInEnv(t, tt.input, outerEnv)
		testValue(t, actual, tt.expected)
	}
}

func TestEvalIntConstructor(t *testing.T) {
	tests := []struct {
		input    string
		expected object.PanObject
	}{
		{
			`Int.new(1)`,
			object.NewPanInt(1),
		},
		{
			`Int.new(false)`,
			object.NewPanInt(0),
		},
		{
			`Int.new(2.9)`,
			object.NewPanInt(2),
		},
		{
			`Int.new(-2.9)`,
			object.NewPanInt(-2),
		},
		// errors
		{
			`Int.new`,
			object.NewTypeErr("Int#new requires at least 2 args"),
		},
		{
			`Int.new("a")`,
			object.NewTypeErr(`"a" cannot be treated as int`),
		},
	}

	for _, tt := range tests {
		actual := testEval(t, tt.input)
		testValue(t, actual, tt.expected)
	}
}

func TestEvalNilConstructor(t *testing.T) {
	tests := []struct {
		input    string
		expected object.PanObject
	}{
		{
			`Nil.new(nil)`,
			object.BuiltInNil,
		},
		// NOTE: not error! (because insufficient args are filled with nil)
		{
			`Nil.new`,
			object.BuiltInNil,
		},
		// errors
		{
			`Nil.new(2)`,
			object.NewTypeErr("2 cannot be treated as nil"),
		},
	}

	for _, tt := range tests {
		actual := testEval(t, tt.input)
		testValue(t, actual, tt.expected)
	}
}

func TestEvalObjConstructor(t *testing.T) {
	tests := []struct {
		input    string
		expected object.PanObject
	}{
		{
			`Obj.new({a: 1})`,
			toPanObj([]object.Pair{
				{Key: object.NewPanStr("a"), Value: object.NewPanInt(1)},
			}),
		},
		// proto info is ignored
		{
			`Obj.new({a: 1}.bear({child?: true}))`,
			toPanObj([]object.Pair{
				{Key: object.NewPanStr("child?"), Value: object.BuiltInTrue},
			}),
		},
		// obj other than PanObj is converted to PanObj
		{
			`Obj.new(10)`,
			toPanObj([]object.Pair{
				{Key: object.NewPanStr("_value"), Value: object.NewPanInt(10)},
			}),
		},
		// if nil is passed, empty obj is returned
		{
			`Obj.new(nil)`,
			toPanObj([]object.Pair{}),
		},
		// not error! (insufficient args are filled with nil)
		{
			`Obj.new`,
			toPanObj([]object.Pair{}),
		},
	}

	for _, tt := range tests {
		actual := testEval(t, tt.input)
		testValue(t, actual, tt.expected)
	}
}

func TestEvalRangeConstructor(t *testing.T) {
	tests := []struct {
		input    string
		expected object.PanObject
	}{
		{
			`Range.new(1)`,
			toPanRange(1, nil, nil),
		},
		{
			`Range.new(2, 5)`,
			toPanRange(2, 5, nil),
		},
		{
			`Range.new('e, 'a, -1)`,
			toPanRange("e", "a", -1),
		},
		// errors
		{
			`Range.new`,
			object.NewTypeErr("Range#new requires at least 2 args"),
		},
	}

	for _, tt := range tests {
		actual := testEval(t, tt.input)
		testValue(t, actual, tt.expected)
	}
}

func TestEvalStrConstructor(t *testing.T) {
	tests := []struct {
		input    string
		expected object.PanObject
	}{
		{
			`Str.new('a)`,
			object.NewPanStr("a"),
		},
		{
			`Str.new('b.bear({child?: true}))`,
			object.NewPanStr("b"),
		},
		// non-string arg is converted to str by .S
		{
			`Str.new({a: 1})`,
			object.NewPanStr(`{"a": 1}`),
		},
		// errors
		{
			`Str.new`,
			object.NewTypeErr("Str#new requires at least 2 args"),
		},
	}

	for _, tt := range tests {
		actual := testEval(t, tt.input)
		testValue(t, actual, tt.expected)
	}
}

func TestEvalAssert(t *testing.T) {
	tests := []struct {
		input    string
		expected object.PanObject
	}{
		// NOTE: do not pass reciever to assert! (it is a function)
		{
			`Kernel['assert](true)`,
			object.BuiltInNil,
		},
		{
			`assert(true)`,
			object.BuiltInNil,
		},
		{
			`assert(false)`,
			object.NewAssertionErr("false is not truthy"),
		},
		{
			`assert("")`,
			object.NewAssertionErr(`"" is not truthy`),
		},
		{
			`assert(1 == 2)`,
			object.NewAssertionErr(`false is not truthy`),
		},
		{
			`assert()`,
			object.NewTypeErr("assert requires at least 1 arg"),
		},
	}

	for _, tt := range tests {
		actual := testEval(t, tt.input)
		testValue(t, actual, tt.expected)
	}
}

func TestEvalAssertEq(t *testing.T) {
	tests := []struct {
		input    string
		expected object.PanObject
	}{
		// NOTE: do not pass reciever to assertEq! (it is a function)
		{
			`Kernel['assertEq]('a, 'a)`,
			object.BuiltInNil,
		},
		{
			`Kernel['assertEq](1 + 1, 2)`,
			object.BuiltInNil,
		},
		{
			`assertEq(true, true)`,
			object.BuiltInNil,
		},
		{
			`assertEq(1 + 1, 3)`,
			object.NewAssertionErr("2 != 3"),
		},
		{
			`assertEq({a: 1, b: 3}, {b: 2, a: 1})`,
			object.NewAssertionErr(`{"a": 1, "b": 3} != {"a": 1, "b": 2}`),
		},
		{
			`assertEq(1)`,
			object.NewTypeErr("assertEq requires at least 2 args"),
		},
	}

	for _, tt := range tests {
		actual := testEval(t, tt.input)
		testValue(t, actual, tt.expected)
	}
}

func TestEvalAssertRaises(t *testing.T) {
	tests := []struct {
		input    string
		expected object.PanObject
	}{
		// NOTE: do not pass receiver to assertRaises! (it is a function)
		{
			`assertRaises(ValueErr, "error") {raise ValueErr.new("error")}`,
			object.BuiltInNil,
		},
		{
			`assertRaises(ValueErr, "error") {1}`,
			object.NewAssertionErr("error must be raised"),
		},
		{
			`assertRaises(ValueErr, "error") {raise TypeErr.new("error")}`,
			object.NewAssertionErr("wrong type: TypeErr != ValueErr"),
		},
		{
			`assertRaises(ValueErr, "error") {raise ValueErr.new("hoge")}`,
			object.NewAssertionErr(`wrong msg: "hoge" != "error"`),
		},
		{
			`assertRaises(1)`,
			object.NewTypeErr("assertRaises requires at least 3 args"),
		},
		{
			`assertRaises(Err, 1) {raise Err.new("hoge")}`,
			object.NewTypeErr("1 cannot be treated as str"),
		},
		// if args[2] is not callable
		{
			`assertRaises(Err, "msg", 1)`,
			object.NewAssertionErr("error must be raised"),
		},
	}

	for _, tt := range tests {
		actual := testEval(t, tt.input)
		testValue(t, actual, tt.expected)
	}
}

func TestEvalTry(t *testing.T) {
	tests := []struct {
		input    string
		expected object.PanObject
	}{
		{
			`1.try`,
			toEitherVal([]object.Pair{
				{Key: object.NewPanStr("_value"), Value: object.NewPanInt(1)},
			}),
		},
		{
			`"".try`,
			toEitherVal([]object.Pair{
				{Key: object.NewPanStr("_value"), Value: object.NewPanStr("")},
			}),
		},
	}

	for _, tt := range tests {
		actual := testEval(t, tt.input)
		testValue(t, actual, tt.expected)
	}
}

func TestEvalFmap(t *testing.T) {
	tests := []struct {
		input    string
		expected object.PanObject
	}{
		// success
		{
			`1.try.fmap {\ + 2}`,
			toEitherVal([]object.Pair{
				{Key: object.NewPanStr("_value"), Value: object.NewPanInt(3)},
			}),
		},
		{
			`1.try.fmap {\ + 2}.fmap {\ + 4}`,
			toEitherVal([]object.Pair{
				{Key: object.NewPanStr("_value"), Value: object.NewPanInt(7)},
			}),
		},
		// failure
		{
			`1.try.fmap {\ / 0}`,
			toEitherErr(object.NewZeroDivisionErr("cannot be divided by 0")),
		},
		{
			`1.try.fmap {\ / 0}.fmap {\ - 1}`,
			toEitherErr(object.NewZeroDivisionErr("cannot be divided by 0")),
		},
		// raises errors (highly unrecommended because it is not caught)
		{
			`1.try.fmap`,
			object.NewTypeErr("EitherVal#fmap requires at least 2 args"),
		},
		{
			`1.try.fmap {\ / 0}.fmap`,
			object.NewTypeErr("EitherErr#fmap requires at least 2 args"),
		},
		{
			`1.try['fmap]({}) {\ + 1}`,
			object.NewTypeErr("{} cannot be treated as EitherVal"),
		},
		// Either['fmap] is abstract method
		{
			`Either.fmap {.keys}`,
			object.BuiltInNotImplemented,
		},
		// return Val(nil) if f does not have prop `call`
		{
			`1.try.fmap(2)`,
			toEitherVal([]object.Pair{
				{Key: object.NewPanStr("_value"), Value: object.BuiltInNil},
			}),
		},
	}

	for _, tt := range tests {
		actual := testEval(t, tt.input)
		testValue(t, actual, tt.expected)
	}
}

func TestEvalEitherA(t *testing.T) {
	tests := []struct {
		input    string
		expected object.PanObject
	}{
		// success
		{
			`1.try.fmap {\ + 2}.A`,
			object.NewPanArr(
				object.NewPanInt(3),
				object.BuiltInNil,
			),
		},
		// failure
		{
			`1.try.fmap {\ / 0}.A`,
			object.NewPanArr(
				object.BuiltInNil,
				object.WrapErr(object.NewZeroDivisionErr("cannot be divided by 0")),
			),
		},
		// raises errors (highly unrecommended because it is not caught)
		{
			`1.try.fmap {\ + 2}['A]()`,
			object.NewTypeErr("EitherVal#A requires at least 1 arg"),
		},
		{
			`1.try.fmap {\ / 0}['A]()`,
			object.NewTypeErr("EitherErr#A requires at least 1 arg"),
		},
		{
			`1.try.fmap {\ + 2}['A]([])`,
			object.NewTypeErr("[] cannot be treated as EitherVal"),
		},
		{
			`1.try.fmap {\ + 2}['A]({})`,
			object.NewTypeErr("{} cannot be treated as EitherVal"),
		},
		{
			`1.try.fmap {\ / 0}['A]()`,
			object.NewTypeErr("EitherErr#A requires at least 1 arg"),
		},
		{
			`1.try.fmap {\ / 0}['A]([])`,
			object.NewTypeErr("[] cannot be treated as EitherErr"),
		},
		{
			`1.try.fmap {\ / 0}['A]({})`,
			object.NewTypeErr("{} cannot be treated as EitherErr"),
		},
		// Either['A] is abstract method
		{
			`Either.A`,
			object.BuiltInNotImplemented,
		},
	}

	for _, tt := range tests {
		actual := testEval(t, tt.input)
		testValue(t, actual, tt.expected)
	}
}

func TestEvalEitherVal(t *testing.T) {
	tests := []struct {
		input    string
		expected object.PanObject
	}{
		// success
		{
			`1.try.fmap {\ + 2}.val`,
			object.NewPanInt(3),
		},
		// failure
		{
			`1.try.fmap {\ / 0}.val`,
			object.BuiltInNil,
		},
		// raises errors (highly unrecommended because it is not caught)
		{
			`1.try.fmap {\ + 2}['val]()`,
			object.NewTypeErr("EitherVal#val requires at least 1 arg"),
		},
		{
			`1.try.fmap {\ / 0}['val]()`,
			object.NewTypeErr("EitherErr#val requires at least 1 arg"),
		},
		{
			`1.try.fmap {\ + 2}['val]([])`,
			object.NewTypeErr("[] cannot be treated as EitherVal"),
		},
		{
			`1.try.fmap {\ + 2}['val]({})`,
			object.NewTypeErr("{} cannot be treated as EitherVal"),
		},
		// Either['val] is abstract method
		{
			`Either.val`,
			object.BuiltInNotImplemented,
		},
	}

	for _, tt := range tests {
		actual := testEval(t, tt.input)
		testValue(t, actual, tt.expected)
	}
}

func TestEvalEitherErr(t *testing.T) {
	tests := []struct {
		input    string
		expected object.PanObject
	}{
		// success
		{
			`1.try.fmap {\ + 2}.err`,
			object.BuiltInNil,
		},
		// failure
		{
			`1.try.fmap {\ / 0}.err`,
			object.WrapErr(object.NewZeroDivisionErr("cannot be divided by 0")),
		},
		// raises errors (highly unrecommended because it is not caught)
		{
			`1.try.fmap {\ + 2}['err]()`,
			object.NewTypeErr("EitherVal#err requires at least 1 arg"),
		},
		{
			`1.try.fmap {\ / 0}['err]()`,
			object.NewTypeErr("EitherErr#err requires at least 1 arg"),
		},
		{
			`1.try.fmap {\ / 0}['err]([])`,
			object.NewTypeErr("[] cannot be treated as EitherErr"),
		},
		{
			`1.try.fmap {\ / 0}['err]({})`,
			object.NewTypeErr("{} cannot be treated as EitherErr"),
		},
		// Either['err] is abstract method
		{
			`Either.err`,
			object.BuiltInNotImplemented,
		},
	}

	for _, tt := range tests {
		actual := testEval(t, tt.input)
		testValue(t, actual, tt.expected)
	}
}

func TestEvalEitherOr(t *testing.T) {
	tests := []struct {
		input    string
		expected object.PanObject
	}{
		// success
		{
			`1.try.fmap {\ + 2}.or("failed")`,
			object.NewPanInt(3),
		},
		// failure
		{
			`1.try.fmap {\ / 0}.or("failed")`,
			object.NewPanStr("failed"),
		},
		// raises errors (highly unrecommended because it is not caught)
		{
			`1.try.fmap {\ + 2}['or]()`,
			object.NewTypeErr("EitherVal#or requires at least 2 args"),
		},
		{
			`1.try.fmap {\ / 0}['or]()`,
			object.NewTypeErr("EitherErr#or requires at least 2 args"),
		},
		{
			`1.try.fmap {\ + 1}['or]([], "")`,
			object.NewTypeErr("[] cannot be treated as EitherVal"),
		},
		{
			`1.try.fmap {\ + 1}['or]({}, "")`,
			object.NewTypeErr("{} cannot be treated as EitherVal"),
		},
		// Either['or] is abstract method
		{
			`Either.or`,
			object.BuiltInNotImplemented,
		},
	}

	for _, tt := range tests {
		actual := testEval(t, tt.input)
		testValue(t, actual, tt.expected)
	}
}

func TestEvalErrWrapperS(t *testing.T) {
	tests := []struct {
		input    string
		expected object.PanObject
	}{
		// failure
		{
			`1.try.fmap {a}.A[1].S`,
			object.NewPanStr("[NameErr: name `a` is not defined]"),
		},
	}

	for _, tt := range tests {
		actual := testEval(t, tt.input)
		testValue(t, actual, tt.expected)
	}
}

func toEitherVal(pairs []object.Pair) *object.PanObj {
	pairMap := map[object.SymHash]object.Pair{}

	for _, pair := range pairs {
		panStr, _ := pair.Key.(*object.PanStr)
		symHash := object.GetSymHash(panStr.Value)
		pairMap[symHash] = pair
	}

	obj := object.NewPanObj(&pairMap, object.BuiltInEitherValObj)
	return obj
}

func toEitherErr(e *object.PanErr) *object.PanObj {
	pairMap := map[object.SymHash]object.Pair{}

	const Key = "_error"
	symHash := object.GetSymHash(Key)
	pairMap[symHash] = object.Pair{
		Key:   object.NewPanStr(Key),
		Value: object.WrapErr(e),
	}

	obj := object.NewPanObj(&pairMap, object.BuiltInEitherErrObj)
	return obj
}

func testPanInt(t *testing.T, actual object.PanObject, expected *object.PanInt) {
	if actual == nil {
		t.Fatalf("actual must not be nil. expected=%v(%T)", expected, expected)
	}

	if actual.Type() != object.IntType {
		t.Fatalf("Type must be IntType(%s). got=%s(%s)",
			expected.Inspect(), actual.Type(), actual.Inspect())
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

	if actual.Type() != object.FloatType {
		t.Fatalf("Type must be FloatType(%s). got=%s(%s)",
			expected.Inspect(), actual.Type(), actual.Inspect())
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

	if actual.Type() != object.StrType {
		t.Fatalf("Type must be StrType(%s). got=%s(%s)",
			expected.Inspect(), actual.Type(), actual.Inspect())
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

	if actual.Type() != object.BoolType {
		t.Fatalf("Type must be BoolType(%s). got=%s(%s)",
			expected.Inspect(), actual.Type(), actual.Inspect())
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

	if actual.Type() != object.NilType {
		t.Fatalf("Type must be NilType(%s). got=%s(%s)",
			expected.Inspect(), actual.Type(), actual.Inspect())
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

	if actual.Type() != object.RangeType {
		t.Fatalf("Type must be RangeType(%s). got=%s(%s)",
			expected.Inspect(), actual.Type(), actual.Inspect())
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

func testPanArr(t *testing.T, actual object.PanObject, expected *object.PanArr) {
	if actual == nil {
		t.Fatalf("actual must not be nil. expected=%v(%T)", expected, expected)
	}

	if actual.Type() != object.ArrType {
		t.Fatalf("Type must be ArrType(%s). got=%s(%s)",
			expected.Inspect(), actual.Type(), actual.Inspect())
		return
	}

	obj, ok := actual.(*object.PanArr)
	if !ok {
		t.Fatalf("actual must be *object.PanArr. got=%T (%v)", actual, actual)
		return
	}

	if len(obj.Elems) != len(expected.Elems) {
		t.Fatalf("length must be %d(%s). got=%d(%s)",
			len(expected.Elems), expected.Inspect(),
			len(obj.Elems), obj.Inspect())
		return
	}

	for i, act := range obj.Elems {
		testValue(t, act, expected.Elems[i])
	}
}

func testPanObj(t *testing.T, actual object.PanObject, expected *object.PanObj) {
	if actual == nil {
		t.Fatalf("actual must not be nil. expected=%v(%T)", expected, expected)
	}

	if actual.Type() != object.ObjType {
		t.Fatalf("Type must be ObjType(%s). got=%s(%s)",
			expected.Inspect(), actual.Type(), actual.Inspect())
		return
	}

	obj, ok := actual.(*object.PanObj)
	if !ok {
		t.Fatalf("actual must be *object.PanObj. got=%T (%v)", actual, actual)
		return
	}

	if len(*obj.Pairs) != len(*expected.Pairs) {
		t.Fatalf("length must be %d (%s). got=%d (%s)",
			len(*expected.Pairs), expected.Inspect(), len(*obj.Pairs), obj.Inspect())
		return
	}

	if obj.Proto() != expected.Proto() {
		t.Errorf("Proto must be same. expected=%v(%T), got=%v(%T)",
			expected.Proto(), expected.Proto(), obj.Proto(), obj.Proto())
	}

	for key, pair := range *expected.Pairs {
		actPair, ok := (*obj.Pairs)[key]
		if !ok {
			t.Errorf("key %v(%T) not found", pair.Key, pair.Key)
			continue
		}

		testValue(t, actPair.Key, pair.Key)
		testValue(t, actPair.Value, pair.Value)
	}
}

func testPanMap(t *testing.T, actual object.PanObject, expected *object.PanMap) {
	if actual == nil {
		t.Fatalf("actual must not be nil. expected=%v(%T)", expected, expected)
	}

	if actual.Type() != object.MapType {
		t.Fatalf("Type must be MapType(%s). got=%s(%s)",
			expected.Inspect(), actual.Type(), actual.Inspect())
		return
	}

	obj, ok := actual.(*object.PanMap)
	if !ok {
		t.Fatalf("actual must be *object.PanMap. got=%T (%v)", actual, actual)
		return
	}

	if len(*obj.HashKeys) != len(*expected.HashKeys) {
		t.Fatalf("length must be %d (%s). got=%d (%s)",
			len(*expected.HashKeys), expected.Inspect(),
			len(*obj.HashKeys), obj.Inspect())
		return
	}

	for i, h := range *expected.HashKeys {
		actHash := (*obj.HashKeys)[i]

		if actHash != h {
			t.Errorf("wrong hashkeys[%d]. expected=%v, got=%v",
				i, h, actHash)
		}
	}

	if len(*obj.Pairs) != len(*expected.Pairs) {
		t.Fatalf("length must be %d (%s). got=%d (%s)",
			len(*expected.Pairs), expected.Inspect(),
			len(*obj.Pairs), obj.Inspect())
		return
	}

	for key, pair := range *expected.Pairs {
		actPair, ok := (*obj.Pairs)[key]
		if !ok {
			t.Errorf("key %v(%T) not found", pair.Key, pair.Key)
			continue
		}

		testValue(t, actPair.Key, pair.Key)
		testValue(t, actPair.Value, pair.Value)
	}

	if len(*obj.NonHashablePairs) != len(*expected.NonHashablePairs) {
		t.Fatalf("nonHashablePair length must be %d (%s). got=%d (%s)",
			len(*expected.Pairs), expected.Inspect(),
			len(*obj.Pairs), obj.Inspect())
		return
	}

	for i, pair := range *expected.NonHashablePairs {
		actPair := (*obj.NonHashablePairs)[i]
		testValue(t, actPair.Key, pair.Key)
		testValue(t, actPair.Value, pair.Value)
	}
}

func testPanFunc(t *testing.T, actual object.PanObject, expected *object.PanFunc) {
	if actual == nil {
		t.Fatalf("actual must not be nil. expected=%v(%T)", expected, expected)
	}

	if actual.Type() != object.FuncType {
		t.Fatalf("Type must be FuncType(%s). got=%s(%s)",
			expected.Inspect(), actual.Type(), actual.Inspect())
		return
	}

	obj, ok := actual.(*object.PanFunc)
	if !ok {
		t.Fatalf("actual must be *object.PanFunc. got=%T (%v)", actual, actual)
		return
	}

	if obj.FuncKind != expected.FuncKind {
		t.Errorf("FuncKind must be %d. got=%d",
			expected.FuncKind, obj.FuncKind)
	}

	testEnv(t, *obj.Env, *expected.Env)
	testFuncComponent(t, obj.FuncWrapper, expected.FuncWrapper)
}

func testFuncComponent(
	t *testing.T,
	actual object.FuncWrapper,
	expected object.FuncWrapper,
) {
	if actual.String() != expected.String() {
		t.Errorf("String() must be `%s`. got=`%s`",
			expected.String(), actual.String())
	}

	testValue(t, actual.Args(), expected.Args())
	testValue(t, actual.Kwargs(), expected.Kwargs())
}

func testPanBulitIn(t *testing.T, actual object.PanObject, expected *object.PanBuiltIn) {
	if actual == nil {
		t.Fatalf("actual must not be nil. expected=%v(%T)", expected, expected)
	}

	if actual.Type() != object.BuiltInType {
		t.Fatalf("Type must be BuiltInType(`%s`). got=%s(`%s`)",
			expected.Inspect(), actual.Type(), actual.Inspect())
		return
	}

	f := actual.(*object.PanBuiltIn)
	actualPtr := fmt.Sprintf("%p", f.Fn)
	expectedPtr := fmt.Sprintf("%p", expected.Fn)

	if actualPtr != expectedPtr {
		t.Errorf("Fn must be %s. got=%s", expectedPtr, actualPtr)
	}
}

func testEnv(t *testing.T, actual object.Env, expected object.Env) {
	if actual.Outer() != expected.Outer() {
		t.Fatalf("Outer is wrong. expected=%s(%p), got=%s(%p)",
			inspectEnv(expected.Outer()), expected.Outer(),
			inspectEnv(actual.Outer()), actual.Outer())
	}

	// compare vars in env
	testValue(t, actual.Items(), expected.Items())
}

func inspectEnv(e *object.Env) string {
	if e == nil {
		return "{nil}"
	}
	return e.Items().Inspect()
}

func testPanErr(t *testing.T, actual object.PanObject, expected *object.PanErr) {
	if actual == nil {
		t.Fatalf("actual must not be nil. expected=%v(%T)", expected, expected)
	}

	if actual.Type() != object.ErrType {
		t.Fatalf("Type must be ErrType(`%s`). got=%s(`%s`)",
			expected.Inspect(), actual.Type(), actual.Inspect())
		return
	}

	e, ok := actual.(*object.PanErr)
	if !ok {
		t.Fatalf("actual must be *object.PanErr. got=%T (%v)", actual, actual)
		return
	}

	if e.ErrKind != expected.ErrKind {
		t.Errorf("ErrKind must be %s. got=%s", expected.ErrKind, e.ErrKind)
	}

	if e.Inspect() != expected.Inspect() {
		t.Errorf("wrong msg. expected=`\n%s\n`. got=`\n%s\n`",
			expected.Inspect(), e.Inspect())
	}

	if e.Proto() != expected.Proto() {
		t.Errorf("proto must be %v(%s). got=%v(%s)",
			expected.Proto(), expected.Proto().Inspect(),
			e.Proto(), e.Proto().Inspect())
	}
}

func testPanErrWrapper(t *testing.T, actual object.PanObject, expected *object.PanErrWrapper) {
	if actual == nil {
		t.Fatalf("actual must not be nil. expected=%v(%T)", expected, expected)
	}

	if actual.Type() != object.ErrWrapperType {
		t.Fatalf("Type must be ErrWrapperType(`%s`). got=%s(`%s`)",
			expected.Inspect(), actual.Type(), actual.Inspect())
		return
	}

	e, ok := actual.(*object.PanErrWrapper)
	if !ok {
		t.Fatalf("actual must be *object.PanErrWrapper. got=%T (%v)", actual, actual)
		return
	}

	if e.ErrKind != expected.ErrKind {
		t.Errorf("ErrKind must be %s. got=%s", expected.ErrKind, e.ErrKind)
	}

	if e.Inspect() != expected.Inspect() {
		t.Errorf("wrong msg. expected=`\n%s\n`. got=`\n%s\n`",
			expected.Inspect(), e.Inspect())
	}

	if e.Proto() != expected.Proto() {
		t.Errorf("proto must be %v(%s). got=%v(%s)",
			expected.Proto(), expected.Proto().Inspect(),
			e.Proto(), e.Proto().Inspect())
	}
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
	case *object.PanArr:
		testPanArr(t, actual, expected)
	case *object.PanObj:
		testPanObj(t, actual, expected)
	case *object.PanMap:
		testPanMap(t, actual, expected)
	case *object.PanErr:
		testPanErr(t, actual, expected)
	case *object.PanErrWrapper:
		testPanErrWrapper(t, actual, expected)
	case *object.PanFunc:
		testPanFunc(t, actual, expected)
	case *object.PanBuiltIn:
		testPanBulitIn(t, actual, expected)
	default:
		t.Fatalf("type of expected %T cannot be handled by testValue()", expected)
	}
}

func testEval(t *testing.T, input string) object.PanObject {
	// NOTE: props in Kernel can be accessed directly in top-level
	env := object.NewEnvWithConsts()
	env.InjectFrom(object.BuiltInKernelObj)

	return testEvalInEnv(t, input, env)
}

func testEvalInEnv(t *testing.T, input string, env *object.Env) object.PanObject {
	node := testParse(t, input)
	panObject := Eval(node, env)
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
