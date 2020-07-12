package object

import (
	"fmt"
)

const ERR_TYPE = "ERR_TYPE"

type PanErr struct {
	ErrType    ErrType
	Msg        string
	StackTrace string
	proto      PanObject
}

func (e *PanErr) Type() PanObjType {
	return ERR_TYPE
}

func (e *PanErr) Inspect() string {
	return fmt.Sprintf("%s: %s", e.ErrType, e.Msg)
}

func (e *PanErr) Proto() PanObject {
	return e.proto
}

func NewPanErr(msg string) *PanErr {
	return &PanErr{
		ErrType: ERR,
		Msg:     msg,
		proto:   BuiltInErrObj,
	}
}

func NewAssertionErr(msg string) *PanErr {
	return &PanErr{
		ErrType: ASSERTION_ERR,
		Msg:     msg,
		proto:   BuiltInAssertionErr,
	}
}

func NewNameErr(msg string) *PanErr {
	return &PanErr{
		ErrType: NAME_ERR,
		Msg:     msg,
		proto:   BuiltInNameErr,
	}
}

func NewNoPropErr(msg string) *PanErr {
	return &PanErr{
		ErrType: NO_PROP_ERR,
		Msg:     msg,
		proto:   BuiltInNoPropErr,
	}
}

func NewNotImplementedErr(msg string) *PanErr {
	return &PanErr{
		ErrType: NOT_IMPLEMENT_ERR,
		Msg:     msg,
		proto:   BuiltInNotImplementedErr,
	}
}

func NewSyntaxErr(msg string) *PanErr {
	return &PanErr{
		ErrType: SYNTAX_ERR,
		Msg:     msg,
		proto:   BuiltInSyntaxErr,
	}
}

func NewTypeErr(msg string) *PanErr {
	return &PanErr{
		ErrType: TYPE_ERR,
		Msg:     msg,
		proto:   BuiltInTypeErr,
	}
}

func NewValueErr(msg string) *PanErr {
	return &PanErr{
		ErrType: VALUE_ERR,
		Msg:     msg,
		proto:   BuiltInValueErr,
	}
}

func NewZeroDivisionErr(msg string) *PanErr {
	return &PanErr{
		ErrType: ZERO_DIVISION_ERR,
		Msg:     msg,
		proto:   BuiltInZeroDivisionErr,
	}
}

type ErrType string

const (
	ERR               = "Err"
	ASSERTION_ERR     = "AssertionErr"
	NAME_ERR          = "NameErr"
	NO_PROP_ERR       = "NoPropErr"
	NOT_IMPLEMENT_ERR = "NotImplementedErr"
	SYNTAX_ERR        = "SyntaxErr"
	TYPE_ERR          = "TypeErr"
	VALUE_ERR         = "ValueErr"
	ZERO_DIVISION_ERR = "ZeroDivisionErr"
)
