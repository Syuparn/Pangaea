package object

import (
	"io"
)

const IO_TYPE = "IO_TYPE"

type PanIO struct {
	In  io.Reader
	Out io.Writer
}

func (io *PanIO) Type() PanObjType {
	return IO_TYPE
}

func (io *PanIO) Inspect() string {
	return "IO"
}

func (io *PanIO) Proto() PanObject {
	return BuiltInIOObj
}
