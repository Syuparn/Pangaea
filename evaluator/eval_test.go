// Source code in this file is inherited and modified from
// "Writing an Interpreter in Go" https://interpreterbook.com/
// MIT License | Copyright (c) 2016-2017 Thorsten Ball
// see https://opensource.org/licenses/MIT

package evaluator

import (
	"../ast"
	"../object"
	"../parser"
	"bytes"
	"fmt"
	"os"
	"strings"
	"testing"
)

func TestMain(m *testing.M) {
	// setup for name resolution
	InjectBuiltInProps()
	ret := m.Run()
	os.Exit(ret)
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

func TestEvalArrAt(t *testing.T) {
	tests := []struct {
		input    string
		expected object.PanObject
	}{
		{
			`[5, 10][0]`,
			&object.PanInt{Value: 5},
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
			&object.PanInt{Value: 3},
		},
		{
			`[1, 2, 3][-3]`,
			&object.PanInt{Value: 1},
		},
		{
			`[1, 2, 3][-4]`,
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
			&object.PanArr{Elems: []object.PanObject{
				&object.PanInt{Value: 0},
			}},
		},
		{
			`[0, 1, 2][0:]`,
			&object.PanArr{Elems: []object.PanObject{
				&object.PanInt{Value: 0},
				&object.PanInt{Value: 1},
				&object.PanInt{Value: 2},
			}},
		},
		{
			`[0, 1, 2][:2]`,
			&object.PanArr{Elems: []object.PanObject{
				&object.PanInt{Value: 0},
				&object.PanInt{Value: 1},
			}},
		},
		{
			`[0, 1, 2][::-1]`,
			&object.PanArr{Elems: []object.PanObject{
				&object.PanInt{Value: 2},
				&object.PanInt{Value: 1},
				&object.PanInt{Value: 0},
			}},
		},
		{
			`[0, 1, 2][1::-1]`,
			&object.PanArr{Elems: []object.PanObject{
				&object.PanInt{Value: 1},
				&object.PanInt{Value: 0},
			}},
		},
		{
			`[0, 1, 2, 3, 4][:2:-1]`,
			&object.PanArr{Elems: []object.PanObject{
				&object.PanInt{Value: 4},
				&object.PanInt{Value: 3},
			}},
		},
		{
			`[0, 1, 2, 3, 4, 5][1::2]`,
			&object.PanArr{Elems: []object.PanObject{
				&object.PanInt{Value: 1},
				&object.PanInt{Value: 3},
				&object.PanInt{Value: 5},
			}},
		},
		{
			`[0, 1, 2, 3, 4][:3:2]`,
			&object.PanArr{Elems: []object.PanObject{
				&object.PanInt{Value: 0},
				&object.PanInt{Value: 2},
			}},
		},
		{
			`[0, 1, 2, 3, 4, 5][1:5:2]`,
			&object.PanArr{Elems: []object.PanObject{
				&object.PanInt{Value: 1},
				&object.PanInt{Value: 3},
			}},
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

func TestEvalBuiltInCallProp(t *testing.T) {
	tests := []struct {
		input    string
		expected object.PanObject
	}{
		// non-callable
		{
			`{}.callProp({a: 1}, 'a)`,
			&object.PanInt{Value: 1},
		},
		// callable
		{
			`{}.callProp({a: {|| 2}}, 'a)`,
			&object.PanInt{Value: 2},
		},
		// NOTE: first arg (`self`) is reciever itself! (`{a: m{|x| x}}`)
		{
			`{}.callProp({a: m{|x| x}}, 'a, 3)`,
			&object.PanInt{Value: 3},
		},
		{
			`{}.callProp({a: m{|x, y, z: 1| [x, y, z]}}, 'a, 4, 5, z: 6)`,
			&object.PanArr{Elems: []object.PanObject{
				&object.PanInt{Value: 4},
				&object.PanInt{Value: 5},
				&object.PanInt{Value: 6},
			}},
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
		// non-hashable dangling keys are also ignored ([1,2]: 4 is ignored)
		{
			`%{[1, 2]: 3, [1, 2]: 4}`,
			toPanMap(
				[]object.Pair{},
				[]object.Pair{
					object.Pair{
						Key: &object.PanArr{Elems: []object.PanObject{
							&object.PanInt{Value: 1},
							&object.PanInt{Value: 2},
						}},
						Value: &object.PanInt{Value: 3},
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

func TestEvalFuncCall(t *testing.T) {
	tests := []struct {
		input    string
		expected object.PanObject
	}{
		{
			`{10}()`,
			&object.PanInt{Value: 10},
		},
		{
			`{|x| x}(5)`,
			&object.PanInt{Value: 5},
		},
		{
			`{|x, y| [x, y]}("x", "y")`,
			&object.PanArr{Elems: []object.PanObject{
				&object.PanStr{Value: "x"},
				&object.PanStr{Value: "y"},
			}},
		},
		{
			`{|foo, bar: "bar", baz| [foo, bar, baz]}("FOO", "BAZ")`,
			&object.PanArr{Elems: []object.PanObject{
				&object.PanStr{Value: "FOO"},
				&object.PanStr{Value: "bar"},
				&object.PanStr{Value: "BAZ"},
			}},
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
			&object.PanInt{Value: 2},
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
			&object.PanInt{Value: 10},
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
			object.NewTypeErr("`{\"call\": {|| [builtin]}}` is not callable."),
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
					object.Pair{
						Key:   &object.PanStr{Value: "c"},
						Value: &object.PanInt{Value: 10},
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
		FuncType:    object.ITER_FUNC,
		Env:         env,
	}
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
						object.Pair{
							Key:   &object.PanStr{Value: "a"},
							Value: &object.PanInt{Value: 10},
						},
					},
					outerEnv,
				),
			),
		},
		{
			`<{|a, b, c: 'c| yield a}>.new('A, 'B, c: 'C)`,
			toPanIter(
				[]string{"a", "b"},
				[]object.Pair{
					object.Pair{
						Key:   &object.PanStr{Value: "c"},
						Value: &object.PanStr{Value: "c"},
					},
				},
				`|a, b, c: 'c| yield a`,
				toEnv(
					[]object.Pair{
						object.Pair{
							Key:   &object.PanStr{Value: "a"},
							Value: &object.PanStr{Value: "A"},
						},
						object.Pair{
							Key:   &object.PanStr{Value: "b"},
							Value: &object.PanStr{Value: "B"},
						},
						object.Pair{
							Key:   &object.PanStr{Value: "c"},
							Value: &object.PanStr{Value: "C"},
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

func toEnv(pairs []object.Pair, outer *object.Env) *object.Env {
	env := object.NewEnclosedEnv(outer)
	for _, pair := range pairs {
		sym := object.GetSymHash(pair.Key.(*object.PanStr).Value)
		env.Set(sym, pair.Value)
	}
	return env
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
			&object.PanStr{Value: "a"},
		},
		// recur
		{
			`it := <{|a| yield a; recur(a+1)}>.new(1)
			 it.next
			 it.next
			 `,
			&object.PanInt{Value: 2},
		},
		{
			`it := <{|a| yield a; recur(a+1)}>.new(1)
			 it.next
			 it.next
			 it.next
			`,
			&object.PanInt{Value: 3},
		},
		{
			`it := <{|a, b| yield [a, b]; recur(a+1, b*2)}>.new(1, 2)
			 it.next
			 it.next
			 it.next
			`,
			&object.PanArr{Elems: []object.PanObject{
				&object.PanInt{Value: 3},
				&object.PanInt{Value: 8},
			}},
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
		// TODO: raise
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
		// TODO: raise
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
		// jump precedence: last line < yield < return
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

func TestEvalArrIter(t *testing.T) {
	tests := []struct {
		input    string
		expected object.PanObject
	}{
		{
			`it := ['a, 'b]._iter
			 it.next`,
			&object.PanStr{Value: "a"},
		},
		{
			`it := ['a, 'b]._iter
			 it.next
			 it.next`,
			&object.PanStr{Value: "b"},
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
			&object.PanStr{Value: "H"},
		},
		{
			`it := "Hi"._iter
			 it.next
			 it.next`,
			&object.PanStr{Value: "i"},
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
			&object.PanStr{Value: "日"},
		},
		{
			`it := "日本語"._iter
			 it.next
			 it.next`,
			&object.PanStr{Value: "本"},
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

func TestEvalListChainPropCall(t *testing.T) {
	tests := []struct {
		input    string
		expected object.PanObject
	}{
		{
			`[1, 2]@+(1)`,
			&object.PanArr{Elems: []object.PanObject{
				object.NewPanInt(2),
				object.NewPanInt(3),
			}},
		},
		{
			`3@S`,
			&object.PanArr{Elems: []object.PanObject{
				&object.PanStr{Value: "1"},
				&object.PanStr{Value: "2"},
				&object.PanStr{Value: "3"},
			}},
		},
		{
			`"日本語"@S`,
			&object.PanArr{Elems: []object.PanObject{
				&object.PanStr{Value: "日"},
				&object.PanStr{Value: "本"},
				&object.PanStr{Value: "語"},
			}},
		},
		// TODO: check obj/map/range
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
			&object.PanArr{Elems: []object.PanObject{
				&object.PanStr{Value: "Taro"},
				&object.PanStr{Value: "Jiro"},
			}},
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
			&object.PanArr{Elems: []object.PanObject{
				object.NewPanInt(2),
				object.NewPanInt(3),
			}},
		},
		{
			`3@{|i| i.S}`,
			&object.PanArr{Elems: []object.PanObject{
				&object.PanStr{Value: "1"},
				&object.PanStr{Value: "2"},
				&object.PanStr{Value: "3"},
			}},
		},
		{
			`"日本語"@{|c| c + "!"}`,
			&object.PanArr{Elems: []object.PanObject{
				&object.PanStr{Value: "日!"},
				&object.PanStr{Value: "本!"},
				&object.PanStr{Value: "語!"},
			}},
		},
		// TODO: check obj/map/range
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
			&object.PanArr{Elems: []object.PanObject{
				&object.PanStr{Value: "Taro"},
				&object.PanStr{Value: "Jiro"},
			}},
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
			`[1, 2, 3]$(0)+`,
			object.NewPanInt(6),
		},
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
			&object.PanStr{Value: "t"},
		},
		{
			`'t if [] else 'f`,
			&object.PanStr{Value: "f"},
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
			&object.PanStr{Value: "A"},
			toPanObj([]object.Pair{
				object.Pair{
					Key:   &object.PanStr{Value: "a"},
					Value: &object.PanStr{Value: "A"},
				},
			}),
		},
		{
			`'A => a`,
			&object.PanStr{Value: "A"},
			toPanObj([]object.Pair{
				object.Pair{
					Key:   &object.PanStr{Value: "a"},
					Value: &object.PanStr{Value: "A"},
				},
			}),
		},
		{
			`a := 'A; a`,
			&object.PanStr{Value: "A"},
			toPanObj([]object.Pair{
				object.Pair{
					Key:   &object.PanStr{Value: "a"},
					Value: &object.PanStr{Value: "A"},
				},
			}),
		},
		{
			`a := 5; b := 10; [a, b]`,
			&object.PanArr{Elems: []object.PanObject{
				&object.PanInt{Value: 5},
				&object.PanInt{Value: 10},
			}},
			toPanObj([]object.Pair{
				object.Pair{
					Key:   &object.PanStr{Value: "a"},
					Value: &object.PanInt{Value: 5},
				},
				object.Pair{
					Key:   &object.PanStr{Value: "b"},
					Value: &object.PanInt{Value: 10},
				},
			}),
		},
		{
			`a := b := 2; [a, b]`,
			&object.PanArr{Elems: []object.PanObject{
				&object.PanInt{Value: 2},
				&object.PanInt{Value: 2},
			}},
			toPanObj([]object.Pair{
				object.Pair{
					Key:   &object.PanStr{Value: "a"},
					Value: &object.PanInt{Value: 2},
				},
				object.Pair{
					Key:   &object.PanStr{Value: "b"},
					Value: &object.PanInt{Value: 2},
				},
			}),
		},
		{
			`"hi" => a; a`,
			&object.PanStr{Value: "hi"},
			toPanObj([]object.Pair{
				object.Pair{
					Key:   &object.PanStr{Value: "a"},
					Value: &object.PanStr{Value: "hi"},
				},
			}),
		},
		{
			`3 => c => d; [c, d]`,
			&object.PanArr{Elems: []object.PanObject{
				&object.PanInt{Value: 3},
				&object.PanInt{Value: 3},
			}},
			toPanObj([]object.Pair{
				object.Pair{
					Key:   &object.PanStr{Value: "c"},
					Value: &object.PanInt{Value: 3},
				},
				object.Pair{
					Key:   &object.PanStr{Value: "d"},
					Value: &object.PanInt{Value: 3},
				},
			}),
		},
	}

	for _, tt := range tests {
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
			&object.PanInt{Value: 100},
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
				object.Pair{
					Key:   &object.PanStr{Value: "a"},
					Value: &object.PanInt{Value: 1},
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

func TestEvalStringify(t *testing.T) {
	tests := []struct {
		input    string
		expected object.PanObject
	}{
		{
			`[1].S`,
			&object.PanStr{Value: "[1]"},
		},
		// unnecessary zeros are omitted
		{
			`1.0.S`,
			&object.PanStr{Value: "1.0"},
		},
		{
			`{|x| x}.S`,
			&object.PanStr{Value: "{|x| x}"},
		},
		{
			`10.S`,
			&object.PanStr{Value: "10"},
		},
		{
			`%{'a: 1}.S`,
			&object.PanStr{Value: `%{"a": 1}`},
		},
		{
			`nil.S`,
			&object.PanStr{Value: `nil`},
		},
		{
			`{a: 1}.S`,
			&object.PanStr{Value: `{"a": 1}`},
		},
		{
			`(1:2).S`,
			&object.PanStr{Value: "(1:2:nil)"},
		},
		// str is not quoted
		{
			`'a.S`,
			&object.PanStr{Value: "a"},
		},
		{
			`true.S`,
			&object.PanStr{Value: "true"},
		},
		{
			`false.S`,
			&object.PanStr{Value: "false"},
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

func TestEvalRepr(t *testing.T) {
	tests := []struct {
		input    string
		expected object.PanObject
	}{
		{
			`[1].repr`,
			&object.PanStr{Value: "[1]"},
		},
		// precise value is shown
		{
			`1.0.repr`,
			&object.PanStr{Value: "1.000000"},
		},
		{
			`{|x| x}.repr`,
			&object.PanStr{Value: "{|x| x}"},
		},
		{
			`10.repr`,
			&object.PanStr{Value: "10"},
		},
		{
			`%{'a: 1}.repr`,
			&object.PanStr{Value: `%{"a": 1}`},
		},
		{
			`nil.repr`,
			&object.PanStr{Value: `nil`},
		},
		{
			`{a: 1}.repr`,
			&object.PanStr{Value: `{"a": 1}`},
		},
		{
			`(1:2).repr`,
			&object.PanStr{Value: "(1:2:nil)"},
		},
		// str is quoted
		{
			`'a.repr`,
			&object.PanStr{Value: `"a"`},
		},
		{
			`true.S`,
			&object.PanStr{Value: "true"},
		},
		{
			`false.S`,
			&object.PanStr{Value: "false"},
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

func TestEvalPropChain(t *testing.T) {
	tests := []struct {
		input    string
		expected object.PanObject
	}{
		{
			`{a: 5, b: 10}.a`,
			&object.PanInt{Value: 5},
		},
		{
			`{a: 5, b: 10}.b`,
			&object.PanInt{Value: 10},
		},
		// call method
		{
			`{a: {|| 2}}.a`,
			&object.PanInt{Value: 2},
		},
		{
			`{a: m{|x| x}}.a(3)`,
			&object.PanInt{Value: 3},
		},
		{
			`{a: m{|x, y| [x, y]}}.a("one", "two")`,
			&object.PanArr{Elems: []object.PanObject{
				&object.PanStr{Value: "one"},
				&object.PanStr{Value: "two"},
			}},
		},
		{
			`{a: m{|x, y: "y"| [x, y]}}.a("x")`,
			&object.PanArr{Elems: []object.PanObject{
				&object.PanStr{Value: "x"},
				&object.PanStr{Value: "y"},
			}},
		},
		{
			`{a: m{|x, y: "y"| [x, y]}}.a("x", y: "Y")`,
			&object.PanArr{Elems: []object.PanObject{
				&object.PanStr{Value: "x"},
				&object.PanStr{Value: "Y"},
			}},
		},
		// if args are insufficient, they are padded by nil
		{
			`{a: m{|x, y| [x, y]}}.a("X")`,
			&object.PanArr{Elems: []object.PanObject{
				&object.PanStr{Value: "X"},
				object.BuiltInNil,
			}},
		},
		// if too many args are passed, they are just ignored
		{
			`{a: m{|x| x}}.a("arg", "needless", "extra", "args")`,
			&object.PanStr{Value: "arg"},
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
			&object.PanInt{Value: 5},
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

func TestEvalMapAt(t *testing.T) {
	tests := []struct {
		input    string
		expected object.PanObject
	}{
		{
			`%{'a: 5, 'b: 10}['a]`,
			&object.PanInt{Value: 5},
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
			&object.PanStr{Value: "one"},
		},
		{
			`%{nil: "nil"}[nil]`,
			&object.PanStr{Value: "nil"},
		},
		{
			`%{[10]: "tenArr"}[[10]]`,
			&object.PanStr{Value: "tenArr"},
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
			&object.PanArr{Elems: []object.PanObject{
				object.NewPanInt(1),
				object.NewPanInt(2),
				object.NewPanInt(1),
				object.NewPanInt(2),
			}},
		},
		// if arity is more than 1, arr is extracted to each param
		{
			`[1, 2].{|i, j| i + j}`,
			object.NewPanInt(3),
		},
		// if args are insufficient, they are padded by nil
		{
			`'X.{|i, j| [i, j]}`,
			&object.PanArr{Elems: []object.PanObject{
				&object.PanStr{Value: "X"},
				object.BuiltInNil,
			}},
		},
		// if too many args are passed, they are just ignored
		{
			`[2, 3, "extra"].{|x, y| x + y}`,
			object.NewPanInt(5),
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
		{
			`nil&.a`,
			object.BuiltInNil,
		},
		{
			`[1, nil, 3]&@{|x| x * 2}`,
			&object.PanArr{Elems: []object.PanObject{
				object.NewPanInt(2),
				object.NewPanInt(6),
			}},
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
			&object.PanArr{Elems: []object.PanObject{
				object.NewPanInt(1),
				&object.PanStr{Value: "b"},
			}},
		},
		{
			`[2, nil, 3, nil, 4]~$(1){|acc, i| acc * i}`,
			object.NewPanInt(24),
		},
		// propcall
		{
			`{a: nil}~.a`,
			toPanObj([]object.Pair{
				object.Pair{
					Key:   &object.PanStr{Value: "a"},
					Value: object.BuiltInNil,
				},
			}),
		},
		{
			`[{a: 1}, {a: nil}]~@a`,
			&object.PanArr{Elems: []object.PanObject{
				object.NewPanInt(1),
				toPanObj([]object.Pair{
					object.Pair{
						Key:   &object.PanStr{Value: "a"},
						Value: object.BuiltInNil,
					},
				}),
			}},
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

func TestNameErr(t *testing.T) {
	tests := []struct {
		input    string
		expected *object.PanErr
	}{
		{
			`a`,
			object.NewNameErr("name `a` is not defined."),
		},
		// multiple lines
		{
			`1; two; 3`,
			object.NewNameErr("name `two` is not defined."),
		},
		// in arr
		{
			`[A]`,
			object.NewNameErr("name `A` is not defined."),
		},
		// in arr expansion
		{
			`[*ae]`,
			object.NewNameErr("name `ae` is not defined."),
		},
		// in obj
		{
			`{key: o}`,
			object.NewNameErr("name `o` is not defined."),
		},
		// in obj expansion
		{
			`{**oe}`,
			object.NewNameErr("name `oe` is not defined."),
		},
		// in map key
		{
			`%{key: 1}`,
			object.NewNameErr("name `key` is not defined."),
		},
		// in map val
		{
			`%{1: val}`,
			object.NewNameErr("name `val` is not defined."),
		},
		// in map expansion
		{
			`%{**me}`,
			object.NewNameErr("name `me` is not defined."),
		},
		// in func call
		{
			`{|a| fc}(1)`,
			object.NewNameErr("name `fc` is not defined."),
		},
		// in arg of func call
		{
			`{|a| 10}(afc)`,
			object.NewNameErr("name `afc` is not defined."),
		},
		// in kwarg of func call
		{
			`{|a: 1| 10}(a: kwfc)`,
			object.NewNameErr("name `kwfc` is not defined."),
		},
		// TODO: err in iter call
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
			&object.PanInt{Value: 2},
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

func testPanInt(t *testing.T, actual object.PanObject, expected *object.PanInt) {
	if actual == nil {
		t.Fatalf("actual must not be nil. expected=%v(%T)", expected, expected)
	}

	if actual.Type() != object.INT_TYPE {
		t.Fatalf("Type must be INT_TYPE(%s). got=%s(%s)",
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

	if actual.Type() != object.FLOAT_TYPE {
		t.Fatalf("Type must be FLOAT_TYPE(%s). got=%s(%s)",
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

	if actual.Type() != object.STR_TYPE {
		t.Fatalf("Type must be STR_TYPE(%s). got=%s(%s)",
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

	if actual.Type() != object.BOOL_TYPE {
		t.Fatalf("Type must be BOOL_TYPE(%s). got=%s(%s)",
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

	if actual.Type() != object.NIL_TYPE {
		t.Fatalf("Type must be NIL_TYPE(%s). got=%s(%s)",
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

	if actual.Type() != object.RANGE_TYPE {
		t.Fatalf("Type must be RANGE_TYPE(%s). got=%s(%s)",
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

	if actual.Type() != object.ARR_TYPE {
		t.Fatalf("Type must be ARR_TYPE(%s). got=%s(%s)",
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

	if actual.Type() != object.OBJ_TYPE {
		t.Fatalf("Type must be OBJ_TYPE(%s). got=%s(%s)",
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

	if actual.Type() != object.MAP_TYPE {
		t.Fatalf("Type must be MAP_TYPE(%s). got=%s(%s)",
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

	if actual.Type() != object.FUNC_TYPE {
		t.Fatalf("Type must be FUNC_TYPE(%s). got=%s(%s)",
			expected.Inspect(), actual.Type(), actual.Inspect())
		return
	}

	obj, ok := actual.(*object.PanFunc)
	if !ok {
		t.Fatalf("actual must be *object.PanFunc. got=%T (%v)", actual, actual)
		return
	}

	if obj.FuncType != expected.FuncType {
		t.Errorf("FuncType must be %d. got=%d",
			expected.FuncType, obj.FuncType)
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

	if actual.Type() != object.BUILTIN_TYPE {
		t.Fatalf("Type must be BUILTIN_TYPE(`%s`). got=%s(`%s`)",
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

	if actual.Type() != object.ERR_TYPE {
		t.Fatalf("Type must be ERR_TYPE(`%s`). got=%s(`%s`)",
			expected.Inspect(), actual.Type(), actual.Inspect())
		return
	}

	e, ok := actual.(*object.PanErr)
	if !ok {
		t.Fatalf("actual must be *object.PanErr. got=%T (%v)", actual, actual)
		return
	}

	if e.ErrType != expected.ErrType {
		t.Errorf("ErrType must be %s. got=%s", expected.ErrType, e.ErrType)
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

func testEval(t *testing.T, input string) object.PanObject {
	return testEvalInEnv(t, input, object.NewEnvWithConsts())
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
