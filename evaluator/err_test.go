package evaluator

import (
	"strings"
	"testing"

	"github.com/Syuparn/pangaea/object"
)

func TestStackTrace(t *testing.T) {
	tests := []struct {
		input    string
		expected []string
	}{
		{
			`*1`,
			[]string{
				`"<string>" line: 1, col: 2`,
				`*1`,
			},
		},
	}

	for _, tt := range tests {
		actual := testEval(t, tt.input)

		e, ok := actual.(*object.PanErr)

		if !ok {
			t.Fatalf("must be evaluated to Err. got=%T(%v)",
				actual, actual)
		}

		expected := strings.Join(tt.expected, "\n")

		if e.StackTrace != expected {
			t.Errorf("stacktrace must be ```\n%s\n```. got=```\n%s\n```",
				expected, e.StackTrace)
		}

	}
}
