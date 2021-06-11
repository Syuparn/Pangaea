package object

import (
	"fmt"
	"testing"
)

func TestBuiltInIterType(t *testing.T) {
	f := func(e *Env, Kwargs *PanObj, args ...PanObject) PanObject { return args[0] }
	obj := PanBuiltInIter{Fn: f, Env: NewEnv()}
	if obj.Type() != BuiltInIterType {
		t.Fatalf("wrong type: expected=%s, got=%s", BuiltInIterType, obj.Type())
	}
}

func TestBuiltInIterInspect(t *testing.T) {
	f := func(e *Env, Kwargs *PanObj, args ...PanObject) PanObject { return args[0] }
	obj := PanBuiltInIter{Fn: f, Env: NewEnv()}
	expected := `<{|| [builtin]}>`
	if obj.Inspect() != expected {
		t.Errorf("wrong output. expected=%s, got=%s",
			expected, obj.Inspect())
	}
}

func TestBuiltInIterRepr(t *testing.T) {
	f := func(e *Env, Kwargs *PanObj, args ...PanObject) PanObject { return args[0] }
	obj := PanBuiltInIter{Fn: f, Env: NewEnv()}
	expected := `<{|| [builtin]}>`
	if obj.Repr() != expected {
		t.Errorf("wrong output. expected=%s, got=%s",
			expected, obj.Inspect())
	}
}

func TestBuiltInIterProto(t *testing.T) {
	f := func(e *Env, Kwargs *PanObj, args ...PanObject) PanObject { return args[0] }
	obj := PanBuiltInIter{Fn: f, Env: NewEnv()}
	if obj.Proto() != BuiltInIterObj {
		t.Fatalf("Proto is not BuiltInIterObj. got=%T (%+v)",
			obj.Proto(), obj.Proto())
	}
}

// checked by compiler (this function works nothing)
func testBuiltInIterIsPanObject() {
	f := func(e *Env, Kwargs *PanObj, args ...PanObject) PanObject { return args[0] }
	var _ PanObject = NewPanBuiltInIter(f, NewEnv())
}

func TestNewPanBuiltInIter(t *testing.T) {
	tests := []struct {
		f   BuiltInFunc
		env *Env
	}{
		{
			func(*Env, *PanObj, ...PanObject) PanObject { return nil },
			NewEnv(),
		},
	}

	for _, tt := range tests {
		actual := NewPanBuiltInIter(tt.f, tt.env)

		// NOTE: functions are not comparable
		if fmt.Sprintf("%v", actual.Fn) != fmt.Sprintf("%v", tt.f) {
			t.Errorf("wrong func. expected=%#v, got=%#v",
				tt.f, actual.Fn)
		}

		if actual.Env != tt.env {
			t.Errorf("wrong env. expected=%#v, got=%#v",
				tt.env, actual.Env)
		}
	}
}
