package builtin

import (
	"github.com/labstack/echo/v4"

	"github.com/Syuparn/pangaea/object"
)

// serverType is a type of panServer.
const serverType = "ServerType"

// panServer is object of arr literal.
type panServer struct {
	server *echo.Echo
}

// Type returns type of this PanObject.
func (s *panServer) Type() object.PanObjType {
	return serverType
}

// Inspect returns formatted source code of this object.
func (s *panServer) Inspect() string {
	return "[server]"
}

// Repr returns pritty-printed string of this object.
func (s *panServer) Repr() string {
	return "[server]"
}

// Proto returns proto of this object.
func (s *panServer) Proto() object.PanObject {
	return object.BuiltInObjObj
}

// Zero returns zero value of this object.
func (s *panServer) Zero() object.PanObject {
	return s
}

// newPanServer returns new server object.
func NewPanServer() *panServer {
	return &panServer{
		server: echo.New(),
	}
}
