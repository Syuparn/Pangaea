package object

import (
	"os"
	"strings"
	"testing"
)

func TestIOType(t *testing.T) {
	obj := NewPanIO(os.Stdin, os.Stdout)
	if obj.Type() != IOType {
		t.Fatalf("wrong type: expected=%s, got=%s", IOType, obj.Type())
	}
}

func TestIOInspect(t *testing.T) {
	tests := []struct {
		obj      *PanIO
		expected string
	}{
		// keys are sorted so that Inspect() always returns same output
		{
			NewPanIO(os.Stdin, os.Stdout),
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

func TestIORepr(t *testing.T) {
	tests := []struct {
		obj      *PanIO
		expected string
	}{
		// keys are sorted so that Inspect() always returns same output
		{
			NewPanIO(os.Stdin, os.Stdout),
			`IO`,
		},
	}

	for _, tt := range tests {
		if tt.obj.Repr() != tt.expected {
			t.Errorf("wrong output: expected=%s, got=%s",
				tt.expected, tt.obj.Inspect())
		}
	}
}

func TestIOProto(t *testing.T) {
	o := NewPanIO(os.Stdin, os.Stdout)
	if o.Proto() != BuiltInIOObj {
		t.Fatalf("Proto is not BuiltInIOObj. got=%T (%+v)",
			o.Proto(), o.Proto())
	}
}

func TestNewLine(t *testing.T) {
	tests := []struct {
		in       string
		expected []string
	}{
		{
			"hoge",
			[]string{"hoge"},
		},
		// ignore last blank line
		{
			"hoge\n",
			[]string{"hoge"},
		},
		{
			"hoge\nfuga\n",
			[]string{"hoge", "fuga"},
		},
		{
			"hoge\nfuga\npiyo",
			[]string{"hoge", "fuga", "piyo"},
		},
		// read non-last blank line
		{
			"hoge\n\npiyo",
			[]string{"hoge", "", "piyo"},
		},
	}

	for _, tt := range tests {
		stdin := strings.NewReader(tt.in)
		o := NewPanIO(stdin, os.Stdout)

		for i, expected := range tt.expected {
			actual, ok := o.ReadLine()
			if !ok {
				t.Fatalf("line %d must not be empty", i)
			}

			if actual.Value != expected {
				t.Errorf("line must be %s. got=%s", expected, actual.Value)
			}
		}

		got, ok := o.ReadLine()
		if ok {
			t.Fatalf("line must be empty. got=%s", got.Value)
		}
	}
}

// checked by compiler (this function works nothing)
func testIOIsPanObject() {
	var _ PanObject = NewPanIO(os.Stdin, os.Stdout)
}
