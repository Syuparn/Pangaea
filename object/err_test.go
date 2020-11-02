package object

import (
	"testing"
)

func TestErrKind(t *testing.T) {
	errObj := NewPanErr("err")
	if errObj.Type() != ErrType {
		t.Fatalf("wrong type: expected=%s, got=%s", ErrType, errObj.Type())
	}
}

func TestErrInspect(t *testing.T) {
	tests := []struct {
		obj      PanObject
		expected string
	}{
		{NewPanErr("err"), "Err: err"},
		{NewAssertionErr("err"), "AssertionErr: err"},
		{NewNameErr("err"), "NameErr: err"},
		{NewNoPropErr("err"), "NoPropErr: err"},
		{NewNotImplementedErr("err"), "NotImplementedErr: err"},
		{NewStopIterErr("err"), "StopIterErr: err"},
		{NewSyntaxErr("err"), "SyntaxErr: err"},
		{NewTypeErr("err"), "TypeErr: err"},
		{NewValueErr("err"), "ValueErr: err"},
		{NewZeroDivisionErr("err"), "ZeroDivisionErr: err"},
	}

	for _, tt := range tests {
		if tt.obj.Inspect() != tt.expected {
			t.Errorf("wrong output: expected=%s, got=%s",
				tt.expected, tt.obj.Inspect())
		}
	}
}

func TestErrProto(t *testing.T) {
	tests := []struct {
		obj              PanObject
		expected         PanObject
		expectedTypeName string
	}{
		{
			NewPanErr("err"),
			BuiltInErrObj,
			"BuiltInErrObj",
		},
		{
			NewAssertionErr("err"),
			BuiltInAssertionErr,
			"BuiltInAssertionErr",
		},
		{
			NewNameErr("err"),
			BuiltInNameErr,
			"BuiltInNameErr",
		},
		{
			NewNoPropErr("err"),
			BuiltInNoPropErr,
			"BuiltInAssertionErr",
		},
		{
			NewNotImplementedErr("err"),
			BuiltInNotImplementedErr,
			"BuiltInNotImplementedErr",
		},
		{
			NewStopIterErr("err"),
			BuiltInStopIterErr,
			"BuiltInStopIterErr",
		},
		{
			NewSyntaxErr("err"),
			BuiltInSyntaxErr,
			"BuiltInSyntaxErr",
		},
		{
			NewTypeErr("err"),
			BuiltInTypeErr,
			"BuiltInTypeErr",
		},
		{
			NewValueErr("err"),
			BuiltInValueErr,
			"BuiltInValueErr",
		},
		{
			NewZeroDivisionErr("err"),
			BuiltInZeroDivisionErr,
			"BuiltInZeroDivisionErr",
		},
	}
	for _, tt := range tests {
		if tt.obj.Proto() != tt.expected {
			t.Fatalf("Proto is not %s. got=%T (%+v)",
				tt.expectedTypeName, tt.obj.Proto(), tt.obj.Proto())
		}
	}
}
