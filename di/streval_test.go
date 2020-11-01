package di

import (
	"testing"

	"github.com/Syuparn/pangaea/object"
)

func TestEvalStrEval(t *testing.T) {
	tests := []struct {
		input    string
		expected object.PanObject
	}{
		// literal
		{
			`"1".eval`,
			object.NewPanInt(1),
		},
		// infix
		{
			`"'he + 'llo".eval`,
			object.NewPanStr("hello"),
		},
		// call
		{
			`"[1, 2, 3]@{|a| a * 2}".eval`,
			&object.PanArr{Elems: []object.PanObject{
				object.NewPanInt(2),
				object.NewPanInt(4),
				object.NewPanInt(6),
			}},
		},
		// assign
		{
			`"a := 3; a * 2".eval`,
			object.NewPanInt(6),
		},
		// multiple lines
		// NOTE: concat each line because `` cannot be written in ``
		{
			"`{|a|\n b := a * 2\n b * 2} => f\n f(2)`.eval",
			object.NewPanInt(8),
		},
	}

	for _, tt := range tests {
		actual := testEval(t, tt.input)
		testValue(t, actual, tt.expected)
	}
}

func TestEvalStrEvalError(t *testing.T) {
	tests := []struct {
		input    string
		expected *object.PanErr
	}{
		{
			`"+".eval`,
			object.NewSyntaxErr("failed to parse"),
		},
		{
			`"a".eval`,
			object.NewNameErr("name `a` is not defined"),
		},
	}

	for _, tt := range tests {
		actual := testEval(t, tt.input)
		testValue(t, actual, tt.expected)
	}
}
