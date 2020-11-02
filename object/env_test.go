package object

import (
	"testing"
)

func TestNewEnv(t *testing.T) {
	env := NewEnv()
	items := env.Items()

	if items.Type() != ObjType {
		t.Fatalf("wrong type: expected=%s, got=%s",
			ObjType, items.Type())
	}

	if items.Inspect() != "{}" {
		t.Fatalf("items must be empty({}). got=`%s`",
			items.Inspect())
	}
}

func TestEnvGetAndSet(t *testing.T) {
	env := NewEnv()
	obj := &PanInt{100}
	env.Set(GetSymHash("myInt"), obj)

	got, ok := env.Get(GetSymHash("myInt"))
	if !ok {
		t.Fatalf("element myInt must be set.")
	}

	if got != obj {
		t.Errorf("wrong value. expected=%s, got=%s",
			obj.Inspect(), got.Inspect())
	}

	if env.Items().Inspect() != `{"myInt": 100}` {
		t.Errorf("Items() are wrong. expected=%s, got=%s",
			`{"myInt": 100}`, env.Items().Inspect())
	}
}

func TestEnclosedEnv(t *testing.T) {
	outer := NewEnv()
	inner := NewEnclosedEnv(outer)
	if inner.Outer() != outer {
		t.Fatalf("Outer() must be Env outer. expected=%v, got=%v",
			outer, inner.Outer())
	}
}

func TestEnvWithConsts(t *testing.T) {
	tests := []struct {
		constName string
		expected  PanObject
	}{
		{"Int", BuiltInIntObj},
		{"Float", BuiltInFloatObj},
		{"Num", BuiltInNumObj},
		{"Nil", BuiltInNilObj},
		{"Str", BuiltInStrObj},
		{"Arr", BuiltInArrObj},
		{"Range", BuiltInRangeObj},
		{"Func", BuiltInFuncObj},
		{"Iter", BuiltInIterObj},
		{"Match", BuiltInMatchObj},
		{"Obj", BuiltInObjObj},
		{"BaseObj", BuiltInBaseObj},
		{"Map", BuiltInMapObj},
		{"true", BuiltInTrue},
		{"false", BuiltInFalse},
		{"nil", BuiltInNil},
		{"Err", BuiltInErrObj},
		{"AssertionErr", BuiltInAssertionErr},
		{"NameErr", BuiltInNameErr},
		{"NoPropErr", BuiltInNoPropErr},
		{"NotImplementedErr", BuiltInNotImplementedErr},
		{"StopIterErr", BuiltInStopIterErr},
		{"SyntaxErr", BuiltInSyntaxErr},
		{"TypeErr", BuiltInTypeErr},
		{"ValueErr", BuiltInValueErr},
		{"ZeroDivisionErr", BuiltInZeroDivisionErr},
		{"_", BuiltInNotImplemented},
	}

	env := NewEnvWithConsts()
	for _, tt := range tests {
		actual, ok := env.Get(GetSymHash(tt.constName))
		if !ok {
			t.Fatalf("element %s must be in env", tt.constName)
		}

		if actual != tt.expected {
			t.Errorf("element %s must be %s. got=%s",
				tt.constName, tt.expected, actual)
		}
	}
}

func TestGetInOuter(t *testing.T) {
	outer := NewEnv()
	inner := NewEnclosedEnv(outer)

	obj := &PanInt{100}
	outer.Set(GetSymHash("myInt"), obj)

	found, ok := inner.Get(GetSymHash("myInt"))

	if !ok {
		t.Fatalf("element myInt must be found.")
	}

	if found != obj {
		t.Errorf("wrong value. expected=%s, got=%s",
			obj.Inspect(), found.Inspect())
	}
}
