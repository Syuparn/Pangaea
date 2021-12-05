package object

import (
	"bufio"
	"io"
)

// IOType is a type of PanIO.
const IOType = "IOType"

// NewPanIO makes new io object.
func NewPanIO(in io.Reader, out io.Writer) *PanIO {
	return &PanIO{
		In:      in,
		Out:     out,
		scanner: bufio.NewScanner(in),
	}
}

// PanIO is object of IO literal.
type PanIO struct {
	In      io.Reader
	Out     io.Writer
	scanner *bufio.Scanner
}

// Type returns type of this PanObject.
func (io *PanIO) Type() PanObjType {
	return IOType
}

// Inspect returns formatted source code of this object.
func (io *PanIO) Inspect() string {
	return "IO"
}

// Repr returns pritty-printed string of this object.
func (i *PanIO) Repr() string {
	return i.Inspect()
}

// Proto returns proto of this object.
func (io *PanIO) Proto() PanObject {
	return BuiltInIOObj
}

// Zero returns zero value of this object.
func (io *PanIO) Zero() PanObject {
	// TODO: implement zero value
	return io
}

// ReadLine reads line from in and returns it as PanStr.
func (io *PanIO) ReadLine() (*PanStr, bool) {
	if !io.scanner.Scan() {
		return nil, false
	}
	line := io.scanner.Text()
	return NewPanStr(line), true
}
