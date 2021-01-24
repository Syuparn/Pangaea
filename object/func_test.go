package object

import (
	"testing"

	"github.com/Syuparn/pangaea/ast"
)

func TestFuncKind(t *testing.T) {
	obj := PanFunc{}
	if obj.Type() != FuncType {
		t.Fatalf("wrong type: expected=%s, got=%s", FuncType, obj.Type())
	}
}

type MockFuncWrapper struct {
	str string
}

func (m *MockFuncWrapper) String() string {
	return m.str
}

func (m *MockFuncWrapper) Args() *PanArr {
	// return empty arr
	return &PanArr{}
}

func (m *MockFuncWrapper) Kwargs() *PanObj {
	// return empty obj
	obj, _ := PanObjInstancePtr(&map[SymHash]Pair{}).(*PanObj)
	return obj
}

func (m *MockFuncWrapper) Body() *[]ast.Stmt {
	// return empty stmt
	return &[]ast.Stmt{}
}

func TestFuncInspect(t *testing.T) {
	tests := []struct {
		obj      PanFunc
		expected string
	}{
		// AstFuncWrapper delegates to FuncComponent.String(), which works same as below
		{PanFunc{&MockFuncWrapper{"|a| a + 1"}, FuncFunc, nil}, "{|a| a + 1}"},
		{PanFunc{&MockFuncWrapper{"|a| a + 1"}, IterFunc, nil}, "<{|a| a + 1}>"},
	}

	for _, tt := range tests {
		if tt.obj.Inspect() != tt.expected {
			t.Errorf("wrong output: expected=%s, got=%s",
				tt.expected, tt.obj.Inspect())
		}
	}
}

func TestFuncProto(t *testing.T) {
	tests := []struct {
		f            PanFunc
		expected     PanObject
		expectedName string
	}{
		{
			PanFunc{&MockFuncWrapper{"|foo| foo"}, FuncFunc, nil},
			BuiltInFuncObj,
			"BuiltInFuncObj",
		},
		{
			PanFunc{&MockFuncWrapper{"|foo| foo"}, IterFunc, nil},
			BuiltInIterObj,
			"BuiltInIterObj",
		},
	}

	for _, tt := range tests {
		actual := tt.f.Proto()

		if actual != tt.expected {
			t.Fatalf("Proto is not %s. got=%T (%+v)",
				tt.expectedName, actual, actual)
		}
	}
}

// checked by compiler (this function works nothing)
func testFuncIsPanObject() {
	var _ PanObject = &PanFunc{&MockFuncWrapper{"|foo| foo"}, FuncFunc, nil}
}
