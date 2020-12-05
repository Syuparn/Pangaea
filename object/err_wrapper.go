package object

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

// Type returns type of this PanObject.
func (w *PanErrWrapper) Type() PanObjType {
	return ErrWrapperType
}
