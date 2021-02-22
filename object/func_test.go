package object

import (
	"fmt"
	"testing"

	"github.com/Syuparn/pangaea/ast"
)

func TestFuncKind(t *testing.T) {
	obj := NewPanFunc(newMockFuncWrapper(), nil)
	if obj.Type() != FuncType {
		t.Fatalf("wrong type: expected=%s, got=%s", FuncType, obj.Type())
	}
}

func TestFuncInspect(t *testing.T) {
	tests := []struct {
		obj      *PanFunc
		expected string
	}{
		// AstFuncWrapper delegates to FuncComponent.String(), which works same as below
		{NewPanFunc(newMockFuncWrapperWithBody("|a| a + 1"), nil), "{|a| a + 1}"},
		{NewPanIter(newMockFuncWrapperWithBody("|a| a + 1"), nil), "<{|a| a + 1}>"},
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
		f            *PanFunc
		expected     PanObject
		expectedName string
	}{
		{
			NewPanFunc(newMockFuncWrapperWithBody("|foo| foo"), NewEnv()),
			BuiltInFuncObj,
			"BuiltInFuncObj",
		},
		{
			NewPanIter(newMockFuncWrapperWithBody("|foo| foo"), NewEnv()),
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
	var _ PanObject = NewPanFunc(newMockFuncWrapper(), nil)
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

func newMockFuncWrapperWithBody(b string) *mockFuncWrapper {
	return &mockFuncWrapper{body: b}
}

type mockFuncWrapper struct {
	body string
}

func (m *mockFuncWrapper) String() string {
	return m.body
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
