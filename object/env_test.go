package object

import (
	"testing"
)

func TestNewEnv(t *testing.T) {
	env := NewEnv()
	items := env.Items()

	if items.Type() != OBJ_TYPE {
		t.Fatalf("wrong type: expected=%s, got=%s",
			OBJ_TYPE, items.Type())
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
