package object

import (
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
	var _ PanObject = &PanBuiltInIter{Fn: f, Env: NewEnv()}
}
