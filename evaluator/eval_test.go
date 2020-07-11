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

func TestEvalArrLiteral(t *testing.T) {
	tests := []struct {
		input    string
		expected *object.PanArr
	}{
		{
			`[]`,
			&object.PanArr{Elems: []object.PanObject{}},
		},
		{
			`[1]`,
			&object.PanArr{Elems: []object.PanObject{
				&object.PanInt{Value: 1},
			}},
		},
		{
			`[2, 3]`,
			&object.PanArr{Elems: []object.PanObject{
				&object.PanInt{Value: 2},
				&object.PanInt{Value: 3},
			}},
		},
		// arr can contain different type elements
		{
			`["a", 4]`,
			&object.PanArr{Elems: []object.PanObject{
				&object.PanStr{Value: "a"},
				&object.PanInt{Value: 4},
			}},
		},
		// nested
		{
			`[[10]]`,
			&object.PanArr{Elems: []object.PanObject{
				&object.PanArr{Elems: []object.PanObject{
					&object.PanInt{Value: 10},
				}},
			}},
		},
		// embedded
		{
			`[*[1]]`,
			&object.PanArr{Elems: []object.PanObject{
				&object.PanInt{Value: 1},
			}},
		},
		{
			`[*[1], 2]`,
			&object.PanArr{Elems: []object.PanObject{
				&object.PanInt{Value: 1},
				&object.PanInt{Value: 2},
			}},
		},
		{
			`[1, *[2, 3], 4]`,
			&object.PanArr{Elems: []object.PanObject{
				&object.PanInt{Value: 1},
				&object.PanInt{Value: 2},
				&object.PanInt{Value: 3},
				&object.PanInt{Value: 4},
			}},
		},
	}

	for _, tt := range tests {
		actual := testEval(t, tt.input)
		testPanArr(t, actual, tt.expected)
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
				object.Pair{
					Key:   &object.PanStr{Value: "a"},
					Value: &object.PanInt{Value: 1},
				},
			}),
		},
		{
			`{a: 1, b: 2}`,
			toPanObj([]object.Pair{
				object.Pair{
					Key:   &object.PanStr{Value: "a"},
					Value: &object.PanInt{Value: 1},
				},
				object.Pair{
					Key:   &object.PanStr{Value: "b"},
					Value: &object.PanInt{Value: 2},
				},
			}),
		},
		// dangling keys are ignored (in this example, `a: 3` is ignored)
		{
			`{a: 1, b: 2, a: 3}`,
			toPanObj([]object.Pair{
				object.Pair{
					Key:   &object.PanStr{Value: "a"},
					Value: &object.PanInt{Value: 1},
				},
				object.Pair{
					Key:   &object.PanStr{Value: "b"},
					Value: &object.PanInt{Value: 2},
				},
			}),
		},
		// Obj can contain multiple types
		{
			`{a: 1, b: "B"}`,
			toPanObj([]object.Pair{
				object.Pair{
					Key:   &object.PanStr{Value: "a"},
					Value: &object.PanInt{Value: 1},
				},
				object.Pair{
					Key:   &object.PanStr{Value: "b"},
					Value: &object.PanStr{Value: "B"},
				},
			}),
		},
		// nested
		{
			`{a: {b: 10}}`,
			toPanObj([]object.Pair{
				object.Pair{
					Key: &object.PanStr{Value: "a"},
					Value: toPanObj([]object.Pair{
						object.Pair{
							Key:   &object.PanStr{Value: "b"},
							Value: &object.PanInt{Value: 10},
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
				object.Pair{
					Key:   &object.PanStr{Value: "a"},
					Value: &object.PanInt{Value: 1},
				},
				object.Pair{
					Key:   &object.PanStr{Value: "b"},
					Value: &object.PanInt{Value: 2},
				},
			}),
		},
		{
			`{**{a: 1}, **{b: 2}}`,
			toPanObj([]object.Pair{
				object.Pair{
					Key:   &object.PanStr{Value: "a"},
					Value: &object.PanInt{Value: 1},
				},
				object.Pair{
					Key:   &object.PanStr{Value: "b"},
					Value: &object.PanInt{Value: 2},
				},
			}),
		},
	}

	for _, tt := range tests {
		actual := testEval(t, tt.input)
		testPanObj(t, actual, tt.expected)
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
					object.Pair{
						Key:   &object.PanStr{Value: "a"},
						Value: &object.PanInt{Value: 1},
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
					object.Pair{
						Key:   &object.PanStr{Value: "a"},
						Value: &object.PanInt{Value: 1},
					},
					object.Pair{
						Key:   &object.PanInt{Value: 2},
						Value: &object.PanInt{Value: 3},
					},
					object.Pair{
						Key:   object.BuiltInTrue,
						Value: &object.PanInt{Value: 5},
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
					object.Pair{
						Key:   &object.PanStr{Value: "a"},
						Value: &object.PanInt{Value: 1},
					},
					object.Pair{
						Key:   &object.PanStr{Value: "b"},
						Value: &object.PanInt{Value: 2},
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
					object.Pair{
						Key:   &object.PanStr{Value: "a"},
						Value: &object.PanInt{Value: 1},
					},
					object.Pair{
						Key:   &object.PanStr{Value: "b"},
						Value: &object.PanStr{Value: "B"},
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
					object.Pair{
						Key: &object.PanStr{Value: "a"},
						Value: toPanMap(
							[]object.Pair{
								object.Pair{
									Key:   &object.PanStr{Value: "b"},
									Value: &object.PanInt{Value: 10},
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
					object.Pair{
						Key:   &object.PanStr{Value: "a"},
						Value: &object.PanInt{Value: 1},
					},
					object.Pair{
						Key:   &object.PanStr{Value: "b"},
						Value: &object.PanInt{Value: 2},
					},
				},
				[]object.Pair{},
			),
		},
		{
			`%{**%{'a: 1}, **%{'b: 2}}`,
			toPanMap(
				[]object.Pair{
					object.Pair{
						Key:   &object.PanStr{Value: "a"},
						Value: &object.PanInt{Value: 1},
					},
					object.Pair{
						Key:   &object.PanStr{Value: "b"},
						Value: &object.PanInt{Value: 2},
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
					object.Pair{
						Key:   &object.PanStr{Value: "a"},
						Value: &object.PanInt{Value: 1},
					},
					object.Pair{
						Key:   &object.PanStr{Value: "b"},
						Value: &object.PanInt{Value: 2},
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
					object.Pair{
						Key:   &object.PanStr{Value: "a"},
						Value: &object.PanInt{Value: 1},
					},
				},
				[]object.Pair{
					object.Pair{
						Key: &object.PanArr{Elems: []object.PanObject{
							&object.PanInt{Value: 1},
							&object.PanInt{Value: 2},
						}},
						Value: &object.PanInt{Value: 3},
					},
					object.Pair{
						Key: toPanMap(
							[]object.Pair{
								object.Pair{
									Key:   &object.PanInt{Value: 4},
									Value: &object.PanInt{Value: 5},
								},
							},
							[]object.Pair{},
						),
						Value: &object.PanInt{Value: 6},
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
	pairMap := map[object.HashKey]object.Pair{}

	for _, pair := range pairs {
		panScalar, _ := pair.Key.(object.PanScalar)
		hash := panScalar.Hash()
		pairMap[hash] = pair
	}

	return &object.PanMap{
		Pairs:            &pairMap,
		NonHashablePairs: &nonHashablePairs,
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
					object.Pair{
						Key:   &object.PanStr{Value: "c"},
						Value: &object.PanInt{Value: 10},
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
					object.Pair{
						Key:   &object.PanStr{Value: "c"},
						Value: &object.PanInt{Value: 10},
					},
					object.Pair{
						Key:   &object.PanStr{Value: "d"},
						Value: &object.PanStr{Value: "e"},
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
		argArr = append(argArr, &object.PanStr{Value: arg})
	}

	funcWrapper := &FuncWrapperImpl{
		codeStr: str,
		args:    &object.PanArr{Elems: argArr},
		kwargs:  toPanObj(kwargs),
		// empty stmt (body is not tested)
		body: &[]ast.Stmt{},
	}

	return &object.PanFunc{
		FuncWrapper: funcWrapper,
		FuncType:    object.FUNC_FUNC,
		Env:         object.NewEnclosedEnv(env),
	}
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

func testPanArr(t *testing.T, actual object.PanObject, expected *object.PanArr) {
	if actual == nil {
		t.Fatalf("actual must not be nil. expected=%v(%T)", expected, expected)
	}

	if actual.Type() != object.ARR_TYPE {
		t.Fatalf("Type must be ARR_TYPE. got=%s", actual.Type())
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

	if actual.Type() != object.OBJ_TYPE {
		t.Fatalf("Type must be OBJ_TYPE. got=%s", actual.Type())
		return
	}

	obj, ok := actual.(*object.PanObj)
	if !ok {
		t.Fatalf("actual must be *object.PanObj. got=%T (%v)", actual, actual)
		return
	}

	if len(*obj.Pairs) != len(*expected.Pairs) {
		t.Fatalf("length must be %d (%v). got=%d (%v)",
			len(*expected.Pairs), *expected.Pairs, len(*obj.Pairs), *obj.Pairs)
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

	if actual.Type() != object.MAP_TYPE {
		t.Fatalf("Type must be MAP_TYPE. got=%s", actual.Type())
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
}

func testPanFunc(t *testing.T, actual object.PanObject, expected *object.PanFunc) {
	if actual == nil {
		t.Fatalf("actual must not be nil. expected=%v(%T)", expected, expected)
	}

	if actual.Type() != object.FUNC_TYPE {
		t.Fatalf("Type must be FUNC_TYPE. got=%s", actual.Type())
		return
	}

	obj, ok := actual.(*object.PanFunc)
	if !ok {
		t.Fatalf("actual must be *object.PanFunc. got=%T (%v)", actual, actual)
		return
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

func testEnv(t *testing.T, actual object.Env, expected object.Env) {
	if actual.Outer() != expected.Outer() {
		t.Fatalf("Outer is wrong. expected=%v(%p), got=%v(%p)",
			expected.Outer(), expected.Outer(),
			actual.Outer(), actual.Outer())
	}

	// compare vars in env
	testValue(t, actual.Items(), expected.Items())
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
	default:
		t.Fatalf("type of expected %T cannot be handled by testValue()", expected)
	}
}

func testEval(t *testing.T, input string) object.PanObject {
	return testEvalInEnv(t, input, object.NewEnv())
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
