package di

import (
	"testing"

	"github.com/Syuparn/pangaea/object"
)

func TestEvalStrEvalEnv(t *testing.T) {
	tests := []struct {
		input    string
		expected object.PanObject
	}{
		// empty
		{
			`"1".evalEnv`,
			toPanObj([]object.Pair{}),
		},
		// single
		{
			`"a := 1".evalEnv`,
			toPanObj([]object.Pair{
				{
					Key:   object.NewPanStr("a"),
					Value: object.NewPanInt(1),
				},
			}),
		},
		// multiple
		{
			`"x := 3; y := x * 2".evalEnv`,
			toPanObj([]object.Pair{
				{
					Key:   object.NewPanStr("x"),
					Value: object.NewPanInt(3),
				},
				{
					Key:   object.NewPanStr("y"),
					Value: object.NewPanInt(6),
				},
			}),
		},
		{
			`"x := true".evalEnv`,
			toPanObj([]object.Pair{
				{
					Key:   object.NewPanStr("x"),
					Value: object.BuiltInTrue,
				},
			}),
		},
	}

	for _, tt := range tests {
		actual := testEval(t, tt.input)
		testValue(t, actual, tt.expected)
	}
}

func TestEvalStrEvalEnvError(t *testing.T) {
	tests := []struct {
		input    string
		expected *object.PanErr
	}{
		{
			`"+".evalEnv`,
			object.NewSyntaxErr("failed to parse"),
		},
		{
			`"a".evalEnv`,
			object.NewNameErr("name `a` is not defined"),
		},
		{
			`Str['evalEnv]()`,
			object.NewTypeErr("evalEnv requires at least 1 arg"),
		},
		{
			`Str['evalEnv](1)`,
			object.NewTypeErr("\\1 must be str"),
		},
	}

	for _, tt := range tests {
		actual := testEval(t, tt.input)
		testValue(t, actual, tt.expected)
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
