package object

import (
	"fmt"
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

func TestNewPanFunc(t *testing.T) {
	tests := []struct {
		f   FuncWrapper
		env *Env
	}{
		{newMockFuncWrapper(), NewEnv()},
	}

	for _, tt := range tests {
		actual := NewPanFunc(tt.f, tt.env)
		// NOTE: functions are not comparable
		if fmt.Sprintf("%v", actual.FuncWrapper) != fmt.Sprintf("%v", tt.f) {
			t.Errorf("wrong value. expected=%#v, got=%#v",
				tt.f, actual.FuncWrapper)
		}

		if actual.FuncKind != FuncFunc {
			t.Errorf("kind must be FuncFunc. got=%#v", actual.FuncKind)
		}

		if actual.Env != tt.env {
			t.Errorf("wrong env. expected=%#v, got=%#v",
				tt.env, actual.Env)
		}
	}
}

func TestNewPanIter(t *testing.T) {
	tests := []struct {
		f   FuncWrapper
		env *Env
	}{
		{newMockFuncWrapper(), NewEnv()},
	}

	for _, tt := range tests {
		actual := NewPanIter(tt.f, tt.env)
		// NOTE: functions are not comparable
		if fmt.Sprintf("%v", actual.FuncWrapper) != fmt.Sprintf("%v", tt.f) {
			t.Errorf("wrong value. expected=%#v, got=%#v",
				tt.f, actual.FuncWrapper)
		}

		if actual.FuncKind != IterFunc {
			t.Errorf("kind must be IterFunc. got=%#v", actual.FuncKind)
		}

		if actual.Env != tt.env {
			t.Errorf("wrong env. expected=%#v, got=%#v",
				tt.env, actual.Env)
		}
	}
}

func newMockFuncWrapper() *mockFuncWrapper {
	return &mockFuncWrapper{}
}

type mockFuncWrapper struct{}

func (m *mockFuncWrapper) String() string {
	return "mock"
}

func (m *mockFuncWrapper) Args() *PanArr {
	return NewPanArr()
}

func (m *mockFuncWrapper) Kwargs() *PanObj {
	return EmptyPanObjPtr()
}

func (m *mockFuncWrapper) Body() *[]ast.Stmt {
	return nil
}
