package object

import (
	"testing"
)

func TestMatchType(t *testing.T) {
	obj := PanMatch{}
	if obj.Type() != MATCH_TYPE {
		t.Fatalf("wrong type: expected=%s, got=%s", MATCH_TYPE, obj.Type())
	}
}

type MockMatchWrapper struct {
	str string
}

func (m *MockMatchWrapper) String() string {
	return m.str
}

func TestMatchInspect(t *testing.T) {
	tests := []struct {
		obj      PanMatch
		expected string
	}{
		// AstFuncWrapper delegates to FuncComponent.String(), which works same as below
		{PanMatch{&MockMatchWrapper{"%{|1| 2 |a| a * 2}"}}, "%{|1| 2 |a| a * 2}"},
	}

	for _, tt := range tests {
		if tt.obj.Inspect() != tt.expected {
			t.Errorf("wrong output: expected=%s, got=%s",
				tt.expected, tt.obj.Inspect())
		}
	}
}

func TestMatchProto(t *testing.T) {
	m := PanMatch{&MockMatchWrapper{"%{}"}}
	if m.Proto() != BuiltInMatchObj {
		t.Fatalf("Proto is not BuiltInMatchObj. got=%T (%+v)",
			m.Proto(), m.Proto())
	}
}

// checked by compiler (this function works nothing)
func testMatchIsPanObject() {
	var _ PanObject = &PanMatch{&MockMatchWrapper{"%{|1| 2 |foo| foo}"}}
}
