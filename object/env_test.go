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
	obj := NewPanInt(100)
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

func TestInjectFrom(t *testing.T) {
	value1 := NewPanInt(1)
	value2 := NewPanInt(2)

	obj := PanObjInstancePtr(&map[SymHash]Pair{
		GetSymHash("a"): {Key: NewPanStr("a"), Value: value1},
		GetSymHash("b"): {Key: NewPanStr("b"), Value: value2},
	}).(*PanObj)

	env := NewEnv()
	env.InjectFrom(obj)

	if len(env.Store) != 2 {
		t.Fatalf("env must have 2 vars. got=%d", len(env.Store))
	}

	actual1, ok := env.Get(GetSymHash("a"))
	if !ok {
		t.Fatalf("`a` must be in env.")
	}
	if actual1 != value1 {
		t.Errorf("a must be %s. got=%s", value1.Inspect(), actual1.Inspect())
	}

	actual2, ok := env.Get(GetSymHash("b"))
	if !ok {
		t.Fatalf("`b` must be in env.")
	}
	if actual2 != value2 {
		t.Errorf("b must be %s. got=%s", value2.Inspect(), actual2.Inspect())
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
		{"Iterable", BuiltInIterableObj},
		{"Comparable", BuiltInComparableObj},
		{"Wrappable", BuiltInWrappableObj},
		{"Match", BuiltInMatchObj},
		{"Obj", BuiltInObjObj},
		{"BaseObj", BuiltInBaseObj},
		{"Map", BuiltInMapObj},
		{"Diamond", BuiltInDiamondObj},
		{"Kernel", BuiltInKernelObj},
		{"Either", BuiltInEitherObj},
		{"EitherVal", BuiltInEitherValObj},
		{"EitherErr", BuiltInEitherErrObj},
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

func TestCopiedEnv(t *testing.T) {
	env := NewEnv()
	obj := NewPanInt(100)
	env.Set(GetSymHash("myInt"), obj)

	copiedEnv := NewCopiedEnv(env)

	got, ok := copiedEnv.Get(GetSymHash("myInt"))
	if !ok {
		t.Fatalf("element myInt must be set.")
	}

	if got != obj {
		t.Errorf("wrong value. expected=%s, got=%s",
			obj.Inspect(), got.Inspect())
	}

	if copiedEnv.Items().Inspect() != `{"myInt": 100}` {
		t.Errorf("Items() are wrong. expected=%s, got=%s",
			`{"myInt": 100}`, env.Items().Inspect())
	}

	// env and copiedEnv are independent
	added := NewPanInt(200)
	env.Set(GetSymHash("addedInt"), added)

	_, ok = copiedEnv.Get(GetSymHash("addedInt"))
	if ok {
		t.Fatalf("element added must not be in copiedEnv.")
	}
}

func TestGetInOuter(t *testing.T) {
	outer := NewEnv()
	inner := NewEnclosedEnv(outer)

	obj := NewPanInt(100)
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
