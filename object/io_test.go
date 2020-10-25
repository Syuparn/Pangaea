package object

import (
	"os"
	"testing"
)

func TestIOType(t *testing.T) {
	obj := PanIO{In: os.Stdin, Out: os.Stdout}
	if obj.Type() != IOType {
		t.Fatalf("wrong type: expected=%s, got=%s", IOType, obj.Type())
	}
}

func TestIOInspect(t *testing.T) {
	tests := []struct {
		obj      PanIO
		expected string
	}{
		// keys are sorted so that Inspect() always returns same output
		{
			PanIO{In: os.Stdin, Out: os.Stdout},
			`IO`,
		},
	}

	for _, tt := range tests {
		if tt.obj.Inspect() != tt.expected {
			t.Errorf("wrong output: expected=%s, got=%s",
				tt.expected, tt.obj.Inspect())
		}
	}
}

func TestIOProto(t *testing.T) {
	o := PanIO{In: os.Stdin, Out: os.Stdout}
	if o.Proto() != BuiltInIOObj {
		t.Fatalf("Proto is not BuiltInIOObj. got=%T (%+v)",
			o.Proto(), o.Proto())
	}
}

// checked by compiler (this function works nothing)
func testIOIsPanObject() {
	var _ PanObject = &PanIO{In: os.Stdin, Out: os.Stdout}
}
