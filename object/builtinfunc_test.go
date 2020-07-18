package object

import (
	"testing"
)

func TestBuiltInType(t *testing.T) {
	f := func(e *Env, Kwargs *PanObj, args ...PanObject) PanObject { return args[0] }
	obj := PanBuiltIn{f}
	if obj.Type() != BUILTIN_TYPE {
		t.Fatalf("wrong type: expected=%s, got=%s", BUILTIN_TYPE, obj.Type())
	}
}

func TestBuiltInInspect(t *testing.T) {
	f := func(e *Env, Kwargs *PanObj, args ...PanObject) PanObject { return args[0] }
	obj := PanBuiltIn{f}
	expected := `{|| [builtin]}`
	if obj.Inspect() != expected {
		t.Errorf("wrong output. expected=%s, got=%s",
			expected, obj.Inspect())
	}
}

func TestBuiltInProto(t *testing.T) {
	f := func(e *Env, Kwargs *PanObj, args ...PanObject) PanObject { return args[0] }
	obj := PanBuiltIn{f}
	if obj.Proto() != BuiltInFuncObj {
		t.Fatalf("Proto is not BuiltInFuncObj. got=%T (%+v)",
			obj.Proto(), obj.Proto())
	}
}

// checked by compiler (this function works nothing)
func testBuiltInIsPanObject() {
	f := func(e *Env, Kwargs *PanObj, args ...PanObject) PanObject { return args[0] }
	var _ PanObject = &PanBuiltIn{f}
}
