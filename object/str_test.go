package object

import (
	"testing"
)

func TestStrType(t *testing.T) {
	strObj := PanStr{"hello"}
	if strObj.Type() != STR_TYPE {
		t.Fatalf("wrong type: expected=%s, got=%s", STR_TYPE, strObj.Type())
	}
}

func TestStrInspect(t *testing.T) {
	tests := []struct {
		obj      PanStr
		expected string
	}{
		{PanStr{"hello"}, "hello"},
		{PanStr{"_foo"}, "_foo"},
		{PanStr{"a i u e o"}, "a i u e o"},
		{PanStr{`\a`}, `\a`},
	}

	for _, tt := range tests {
		if tt.obj.Inspect() != tt.expected {
			t.Errorf("wrong output: expected=%s, got=%s",
				tt.expected, tt.obj.Inspect())
		}
	}
}

func TestStrProto(t *testing.T) {
	s := PanStr{"foo"}
	if s.Proto() != builtInStrObj {
		t.Fatalf("Proto is not BuiltinStrObj. got=%T (%+v)",
			s.Proto(), s.Proto())
	}
}

// checked by compiler (this function works nothing)
func testStrIsPanObject() {
	var _ PanObject = &PanStr{"FOO"}
}
