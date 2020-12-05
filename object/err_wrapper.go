package object

import "fmt"

// ErrWrapperType is a type of PanErrWrapper.
const ErrWrapperType = "ErrWrapperType"

// WrapErr wraps PanErr object and enables to treat error info without exception.
func WrapErr(err *PanErr) *PanErrWrapper {
	return &PanErrWrapper{*err}
}

// PanErrWrapper is an error wrapper to treat error info without exception.
type PanErrWrapper struct {
	PanErr
}

// Inspect returns formatted source code of this object.
func (w *PanErrWrapper) Inspect() string {
	// wrapped by [] so that border of err object is easily understood
	// (for example: `{a: [NameErr: err]}`)
	return fmt.Sprintf("[%s]", w.PanErr.Inspect())
}

// Type returns type of this PanObject.
func (w *PanErrWrapper) Type() PanObjType {
	return ErrWrapperType
}
