package object

import (
	"testing"
)

func TestFuncType(t *testing.T) {
	obj := PanFunc{}
	if obj.Type() != FUNC_TYPE {
		t.Fatalf("wrong type: expected=%s, got=%s", FUNC_TYPE, obj.Type())
	}
}

type MockFuncWrapper struct {
	str string
}

func (m *MockFuncWrapper) String() string {
	return m.str
}

func TestFuncInspect(t *testing.T) {
	tests := []struct {
		obj      PanFunc
		expected string
	}{
		// AstFuncWrapper delegates to FuncComponent.String(), which works same as below
		{PanFunc{&MockFuncWrapper{"|a| a + 1"}, FUNC_FUNC, nil}, "{|a| a + 1}"},
		{PanFunc{&MockFuncWrapper{"|a| a + 1"}, ITER_FUNC, nil}, "<{|a| a + 1}>"},
	}

	for _, tt := range tests {
		if tt.obj.Inspect() != tt.expected {
			t.Errorf("wrong output: expected=%s, got=%s",
				tt.expected, tt.obj.Inspect())
		}
	}
}

func TestFuncProto(t *testing.T) {
	f := PanFunc{&MockFuncWrapper{"|foo| foo"}, FUNC_FUNC, nil}
	if f.Proto() != BuiltInFuncObj {
		t.Fatalf("Proto is not BuiltInFuncObj. got=%T (%+v)",
			f.Proto(), f.Proto())
	}
}

// checked by compiler (this function works nothing)
func testFuncIsPanObject() {
	var _ PanObject = &PanFunc{&MockFuncWrapper{"|foo| foo"}, FUNC_FUNC, nil}
}
