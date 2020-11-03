package di

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/Syuparn/pangaea/ast"
	"github.com/Syuparn/pangaea/evaluator"
	"github.com/Syuparn/pangaea/object"
	"github.com/Syuparn/pangaea/parser"
)

func TestMain(m *testing.M) {
	// HACK: adjust relative file path to read native codes
	apath, _ := filepath.Abs("../")
	os.Chdir(apath)
	// setup for name resolution
	InjectBuiltInProps(object.NewEnvWithConsts())
	ret := m.Run()
	os.Exit(ret)
}

func testEval(t *testing.T, input string) object.PanObject {
	return testEvalInEnv(t, input, object.NewEnvWithConsts())
}

func testEvalInEnv(t *testing.T, input string, env *object.Env) object.PanObject {
	node := testParse(t, input)
	panObject := evaluator.Eval(node, env)
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
		t.Fatalf("length must be %d. got=%d", len(expected.Elems), len(obj.Elems))
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
	case *object.PanFunc:
		testPanFunc(t, actual, expected)
	case *object.PanBuiltIn:
		testPanBulitIn(t, actual, expected)
	default:
		t.Fatalf("type of expected %T cannot be handled by testValue()", expected)
	}
}
