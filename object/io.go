package object

import (
	"io"
)

// IOType is a type of PanIO.
const IOType = "IOType"

// PanIO is object of IO literal.
type PanIO struct {
	In  io.Reader
	Out io.Writer
}

// Type returns type of this PanObject.
func (io *PanIO) Type() PanObjType {
	return IOType
}

// Inspect returns formatted source code of this object.
func (io *PanIO) Inspect() string {
	return "IO"
}

// Proto returns proto of this object.
func (io *PanIO) Proto() PanObject {
	return BuiltInIOObj
}
