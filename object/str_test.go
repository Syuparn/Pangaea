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

func testStrIsPublic(t *testing.T) {
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
		{NewPanStr("123"), false},
		{NewPanStr("Hello, world!"), false},
		{NewPanStr("+"), false},
		{NewPanStr(`\1`), false},
		{NewPanStr(`\foo`), false},
	}

	for _, tt := range tests {
		if tt.obj.IsPublic != tt.expected {
			t.Errorf("wrong field IsPublic: expected=%t, got=%t",
				tt.expected, tt.obj.IsPublic)
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
