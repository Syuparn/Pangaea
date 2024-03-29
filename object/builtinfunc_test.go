package object

import (
	"fmt"
	"testing"
)

func TestBuiltInType(t *testing.T) {
	f := func(e *Env, Kwargs *PanObj, args ...PanObject) PanObject { return args[0] }
	obj := NewPanBuiltInFunc(f)
	if obj.Type() != BuiltInType {
		t.Fatalf("wrong type: expected=%s, got=%s", BuiltInType, obj.Type())
	}
}

func TestBuiltInInspect(t *testing.T) {
	f := func(e *Env, Kwargs *PanObj, args ...PanObject) PanObject { return args[0] }
	obj := NewPanBuiltInFunc(f)
	expected := `{|| [builtin]}`
	if obj.Inspect() != expected {
		t.Errorf("wrong output. expected=%s, got=%s",
			expected, obj.Inspect())
	}
}

func TestBuiltInRepr(t *testing.T) {
	f := func(e *Env, Kwargs *PanObj, args ...PanObject) PanObject { return args[0] }
	obj := NewPanBuiltInFunc(f)
	expected := `{|| [builtin]}`
	if obj.Repr() != expected {
		t.Errorf("wrong output. expected=%s, got=%s",
			expected, obj.Inspect())
	}
}

func TestBuiltInProto(t *testing.T) {
	f := func(e *Env, Kwargs *PanObj, args ...PanObject) PanObject { return args[0] }
	obj := NewPanBuiltInFunc(f)
	if obj.Proto() != BuiltInFuncObj {
		t.Fatalf("Proto is not BuiltInFuncObj. got=%T (%+v)",
			obj.Proto(), obj.Proto())
	}
}

func TestBuiltInZero(t *testing.T) {
	f := func(e *Env, Kwargs *PanObj, args ...PanObject) PanObject { return args[0] }

	tests := []struct {
		name string
		obj  *PanBuiltIn
	}{
		{"f", NewPanBuiltInFunc(f)},
	}

	for _, tt := range tests {
		tt := tt // pin

		t.Run(tt.name, func(t *testing.T) {
			actual := tt.obj.Zero()

			if actual != tt.obj {
				t.Errorf("zero must be itself (%#v). got=%s (%#v)",
					tt.obj, actual.Repr(), actual)
			}
		})
	}
}

// checked by compiler (this function works nothing)
func testBuiltInIsPanObject() {
	f := func(e *Env, Kwargs *PanObj, args ...PanObject) PanObject { return args[0] }
	var _ PanObject = NewPanBuiltInFunc(f)
}

func TestNewPanBuiltInFunc(t *testing.T) {
	tests := []struct {
		f BuiltInFunc
	}{
		{func(*Env, *PanObj, ...PanObject) PanObject { return nil }},
	}

	for _, tt := range tests {
		actual := NewPanBuiltInFunc(tt.f)
		// NOTE: functions are not comparable
		if fmt.Sprintf("%v", actual.Fn) != fmt.Sprintf("%v", tt.f) {
			t.Errorf("wrong value. expected=%#v, got=%#v",
				tt.f, actual.Fn)
		}
	}
}
