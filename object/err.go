package object

import (
	"fmt"
)

// ErrType is a type of PanErr.
const ErrType = "ErrType"

// PanErr is object of err literal.
type PanErr struct {
	ErrKind    ErrKind
	Msg        string
	StackTrace string
	proto      PanObject
}

// Type returns type of this PanObject.
func (e *PanErr) Type() PanObjType {
	return ErrType
}

// Inspect returns formatted source code of this object.
func (e *PanErr) Inspect() string {
	return fmt.Sprintf("%s: %s", e.ErrKind, e.Msg)
}

// Repr returns pritty-printed string of this object.
func (e *PanErr) Repr() string {
	return e.Inspect()
}

// Proto returns proto of this object.
func (e *PanErr) Proto() PanObject {
	return e.proto
}

// Zero returns zero value of this object.
func (e *PanErr) Zero() PanObject {
	// TODO: implement zero value
	return e
}

// Message returns error message.
func (e *PanErr) Message() string {
	return e.Msg
}

// Kind returns kind of this err.
func (e *PanErr) Kind() string {
	return string(e.ErrKind)
}

// NewPanErr returns new err object.
func NewPanErr(msg string) *PanErr {
	return &PanErr{
		ErrKind: Err,
		Msg:     msg,
		proto:   BuiltInErrObj,
	}
}

// NewAssertionErr returns new assertionErr object.
func NewAssertionErr(msg string) *PanErr {
	return &PanErr{
		ErrKind: AssertionErr,
		Msg:     msg,
		proto:   BuiltInAssertionErr,
	}
}

// NewFileNotFoundErr returns new fileNotFoundErr object.
func NewFileNotFoundErr(msg string) *PanErr {
	return &PanErr{
		ErrKind: FileNotFoundErr,
		Msg:     msg,
		proto:   BuiltInFileNotFoundErr,
	}
}

// NewNameErr returns new nameErr object.
func NewNameErr(msg string) *PanErr {
	return &PanErr{
		ErrKind: NameErr,
		Msg:     msg,
		proto:   BuiltInNameErr,
	}
}

// NewNoPropErr returns new noPropErr object.
func NewNoPropErr(msg string) *PanErr {
	return &PanErr{
		ErrKind: NoPropErr,
		Msg:     msg,
		proto:   BuiltInNoPropErr,
	}
}

// NewNotImplementedErr returns new notImplementedErr object.
func NewNotImplementedErr(msg string) *PanErr {
	return &PanErr{
		ErrKind: NotImplementErr,
		Msg:     msg,
		proto:   BuiltInNotImplementedErr,
	}
}

// NewStopIterErr returns new stopIterErr object.
// NOTE: This error is prepared to make iter simpler, even though
// stopIter is not an error actually.
func NewStopIterErr(msg string) *PanErr {
	return &PanErr{
		ErrKind: StopIterErr,
		Msg:     msg,
		proto:   BuiltInStopIterErr,
	}
}

// NewSyntaxErr returns new syntaxErr object.
func NewSyntaxErr(msg string) *PanErr {
	return &PanErr{
		ErrKind: SyntaxErr,
		Msg:     msg,
		proto:   BuiltInSyntaxErr,
	}
}

// NewTypeErr returns new typeErr object.
func NewTypeErr(msg string) *PanErr {
	return &PanErr{
		ErrKind: TypeErr,
		Msg:     msg,
		proto:   BuiltInTypeErr,
	}
}

// NewValueErr returns new valueErr object.
func NewValueErr(msg string) *PanErr {
	return &PanErr{
		ErrKind: ValueErr,
		Msg:     msg,
		proto:   BuiltInValueErr,
	}
}

// NewZeroDivisionErr returns new zeroDivisionErr object.
func NewZeroDivisionErr(msg string) *PanErr {
	return &PanErr{
		ErrKind: ZeroDivisionErr,
		Msg:     msg,
		proto:   BuiltInZeroDivisionErr,
	}
}

// ErrKind is a information to designate error type.
// NOTE: ErrKind is necessary because all err is implemented by same struct, PanErr.
type ErrKind string

// nolint: comment
const (
	Err             = "Err"
	AssertionErr    = "AssertionErr"
	FileNotFoundErr = "FileNotFoundErr"
	NameErr         = "NameErr"
	NoPropErr       = "NoPropErr"
	NotImplementErr = "NotImplementedErr"
	StopIterErr     = "StopIterErr"
	SyntaxErr       = "SyntaxErr"
	TypeErr         = "TypeErr"
	ValueErr        = "ValueErr"
	ZeroDivisionErr = "ZeroDivisionErr"
)
