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
				outerEnv,
			),
		},
		{
			`<{|a|}>`,
			toPanIter(
				[]string{"a"},
				[]object.Pair{},
				`|a| `,
				outerEnv,
			),
		},
		{
			`<{|a, b|}>`,
			toPanIter(
				[]string{"a", "b"},
				[]object.Pair{},
				`|a, b| `,
				outerEnv,
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
				outerEnv,
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
				outerEnv,
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
		Env:         object.NewEnclosedEnv(env),
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
			(*object.BuiltInObjObj.Pairs)[object.GetSymHash("at")].Value,
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
		t.Fatalf("Outer is wrong. expected=%v(%p), got=%v(%p)",
			expected.Outer(), expected.Outer(),
			actual.Outer(), actual.Outer())
	}

	// compare vars in env
	testValue(t, actual.Items(), expected.Items())
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
