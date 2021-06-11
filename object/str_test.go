package object

import (
	"testing"
)

func TestStrType(t *testing.T) {
	strObj := NewPanStr("hello")
	if strObj.Type() != StrType {
		t.Fatalf("wrong type: expected=%s, got=%s", StrType, strObj.Type())
	}
}

func TestStrInspect(t *testing.T) {
	tests := []struct {
		obj      *PanStr
		expected string
	}{
		{NewPanStr("hello"), `"hello"`},
		{NewPanStr("_foo"), `"_foo"`},
		{NewPanStr("a i u e o"), `"a i u e o"`},
		{NewPanStr(`\a`), `"\a"`},
	}

	for _, tt := range tests {
		if tt.obj.Inspect() != tt.expected {
			t.Errorf("wrong output: expected=%s, got=%s",
				tt.expected, tt.obj.Inspect())
		}
	}
}

func TestStrRepr(t *testing.T) {
	tests := []struct {
		obj      *PanStr
		expected string
	}{
		{NewPanStr("hello"), `"hello"`},
		{NewPanStr("_foo"), `"_foo"`},
		{NewPanStr("a i u e o"), `"a i u e o"`},
		{NewPanStr(`\a`), `"\a"`},
	}

	for _, tt := range tests {
		if tt.obj.Repr() != tt.expected {
			t.Errorf("wrong output: expected=%s, got=%s",
				tt.expected, tt.obj.Inspect())
		}
	}
}

func TestStrProto(t *testing.T) {
	s := NewPanStr("foo")
	if s.Proto() != BuiltInStrObj {
		t.Fatalf("Proto is not BuiltInStrObj. got=%T (%+v)",
			s.Proto(), s.Proto())
	}
}

func TestStrHash(t *testing.T) {
	tests := []struct {
		obj      *PanStr
		expected string
	}{
		{NewPanStr("hello"), "hello"},
		{NewPanStr("a i u e o"), "a i u e o"},
		{NewPanStr(""), ""},
		{NewPanStr("longlonglonglonglonglong"), "longlonglonglonglonglong"},
	}

	for _, tt := range tests {
		// register symbol
		_ = tt.obj.SymHash()

		h := tt.obj.Hash()

		if h.Type != StrType {
			t.Fatalf("hash type must be StrType. got=%s", h.Type)
		}

		if h.Value != symHashTable[tt.expected] {
			t.Errorf("wrong hash key: got=%d, expected=%d",
				h.Value, symHashTable[tt.expected])
		}
	}
}

func TestNewPanStr(t *testing.T) {
	tests := []struct {
		str      string
		isPublic bool
		isSym    bool
	}{
		{"hoge", true, true},
		{"hello world", false, false},
	}

	for _, tt := range tests {
		actual := NewPanStr(tt.str)
		if actual.Value != tt.str {
			t.Errorf("wrong value. expected=%s, got=%s", tt.str, actual.Value)
		}

		if actual.IsPublic != tt.isPublic {
			t.Errorf("wrong isPublic. expected=%t, got=%t", tt.isPublic, actual.IsPublic)
		}

		if actual.IsSym != tt.isSym {
			t.Errorf("wrong isSym. expected=%t, got=%t", tt.isSym, actual.IsSym)
		}
	}
}

func TestStrIsPublic(t *testing.T) {
	tests := []struct {
		obj      *PanStr
		expected bool
	}{
		// public
		{NewPanStr("hello"), true},
		{NewPanStr("hoge1"), true},
		{NewPanStr("Hoge"), true},
		{NewPanStr("Ho_ge"), true},
		{NewPanStr("hoge?"), true},
		{NewPanStr("hoge!"), true},
		// private
		{NewPanStr("_hoge"), false},
		{NewPanStr("+"), false},
		{NewPanStr(`\1`), false},
		{NewPanStr(`\foo`), false},
		// others
		{NewPanStr("Hello, world!"), false},
		{NewPanStr("123"), false},
	}

	for _, tt := range tests {
		if tt.obj.IsPublic != tt.expected {
			t.Errorf("wrong field IsPublic: expected=%t, got=%t",
				tt.expected, tt.obj.IsPublic)
		}
	}
}

func TestStrIsSym(t *testing.T) {
	tests := []struct {
		obj      *PanStr
		expected bool
	}{
		// public
		{NewPanStr("hello"), true},
		{NewPanStr("hoge1"), true},
		{NewPanStr("Hoge"), true},
		{NewPanStr("Ho_ge"), true},
		{NewPanStr("hoge?"), true},
		{NewPanStr("hoge!"), true},
		// private
		{NewPanStr("_hoge"), true},
		{NewPanStr("+"), true},
		{NewPanStr(`\1`), true},
		{NewPanStr(`\foo`), true},
		// cannot be used for prop or var name
		{NewPanStr("Hello, world!"), false},
		{NewPanStr("にほんご"), false},
		{NewPanStr("123"), false},
	}

	for _, tt := range tests {
		if tt.obj.IsSym != tt.expected {
			t.Errorf("wrong field IsSym: expected=%t(%s), got=%t",
				tt.expected, tt.obj.Value, tt.obj.IsSym)
		}
	}
}

// checked by compiler (this function works nothing)
func testStrIsPanObject() {
	var _ PanObject = NewPanStr("FOO")
}

func testStrIsPanScalar() {
	var _ PanScalar = NewPanStr("ABC")
}
