package evaluator

import (
	"../object"
	"testing"
)

func TestFuncCall(t *testing.T) {
	tests := []struct {
		funcSrc  string
		kwargs   *object.PanObj
		args     []object.PanObject
		expected object.PanObject
	}{
		{
			`{5}`,
			toPanObj([]object.Pair{}),
			[]object.PanObject{},
			object.NewPanInt(5),
		},
		{
			`{|x| x}`,
			toPanObj([]object.Pair{}),
			[]object.PanObject{
				object.NewPanInt(10),
			},
			object.NewPanInt(10),
		},
		{
			`{|x, y| [x, y]}`,
			toPanObj([]object.Pair{}),
			[]object.PanObject{
				object.NewPanStr("x"),
				object.NewPanStr("y"),
			},
			&object.PanArr{Elems: []object.PanObject{
				object.NewPanStr("x"),
				object.NewPanStr("y"),
			}},
		},
	}

	for _, tt := range tests {
		evaluated := testEval(t, tt.funcSrc)
		f, ok := evaluated.(*object.PanFunc)
		if !ok {
			t.Fatalf("evaluated value is not PanFunc. got=%s", evaluated.Type())
		}

		// prepend f
		args := append([]object.PanObject{f}, tt.args...)

		ret := evalFuncCall(object.NewEnv(), tt.kwargs, args...)
		testValue(t, ret, tt.expected)
	}
}
