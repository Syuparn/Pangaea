package object

import (
	"testing"
)

func TestErrWrapperKind(t *testing.T) {
	errObj := WrapErr(NewPanErr("err"))
	if errObj.Type() != ErrWrapperType {
		t.Fatalf("wrong type: expected=%s, got=%s", ErrWrapperType, errObj.Type())
	}
}

func TestErrWrapperInspect(t *testing.T) {
	tests := []struct {
		obj      PanObject
		expected string
	}{
		{WrapErr(NewPanErr("err")), "Err: err"},
		{WrapErr(NewAssertionErr("err")), "AssertionErr: err"},
		{WrapErr(NewNameErr("err")), "NameErr: err"},
		{WrapErr(NewNoPropErr("err")), "NoPropErr: err"},
		{WrapErr(NewNotImplementedErr("err")), "NotImplementedErr: err"},
		{WrapErr(NewStopIterErr("err")), "StopIterErr: err"},
		{WrapErr(NewSyntaxErr("err")), "SyntaxErr: err"},
		{WrapErr(NewTypeErr("err")), "TypeErr: err"},
		{WrapErr(NewValueErr("err")), "ValueErr: err"},
		{WrapErr(NewZeroDivisionErr("err")), "ZeroDivisionErr: err"},
	}

	for _, tt := range tests {
		if tt.obj.Inspect() != tt.expected {
			t.Errorf("wrong output: expected=%s, got=%s",
				tt.expected, tt.obj.Inspect())
		}
	}
}

func TestErrWrapperProto(t *testing.T) {
	tests := []struct {
		obj              PanObject
		expected         PanObject
		expectedTypeName string
	}{
		{
			WrapErr(NewPanErr("err")),
			BuiltInErrObj,
			"BuiltInErrObj",
		},
		{
			WrapErr(NewAssertionErr("err")),
			BuiltInAssertionErr,
			"BuiltInAssertionErr",
		},
		{
			WrapErr(NewNameErr("err")),
			BuiltInNameErr,
			"BuiltInNameErr",
		},
		{
			WrapErr(NewNoPropErr("err")),
			BuiltInNoPropErr,
			"BuiltInAssertionErr",
		},
		{
			WrapErr(NewNotImplementedErr("err")),
			BuiltInNotImplementedErr,
			"BuiltInNotImplementedErr",
		},
		{
			WrapErr(NewStopIterErr("err")),
			BuiltInStopIterErr,
			"BuiltInStopIterErr",
		},
		{
			WrapErr(NewSyntaxErr("err")),
			BuiltInSyntaxErr,
			"BuiltInSyntaxErr",
		},
		{
			WrapErr(NewTypeErr("err")),
			BuiltInTypeErr,
			"BuiltInTypeErr",
		},
		{
			WrapErr(NewValueErr("err")),
			BuiltInValueErr,
			"BuiltInValueErr",
		},
		{
			WrapErr(NewZeroDivisionErr("err")),
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
