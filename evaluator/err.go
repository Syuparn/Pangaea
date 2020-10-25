package evaluator

import (
	"bytes"
	"github.com/Syuparn/pangaea/ast"
	"github.com/Syuparn/pangaea/object"
	"strings"
)

func appendStackTrace(e *object.PanErr, src *ast.Source) *object.PanErr {
	var out bytes.Buffer

	stackTrace := parseSrc(src)
	// NOTE: if stackTrace is same as previous one, just ignore it
	if strings.HasSuffix(e.StackTrace, stackTrace) {
		return e
	}

	// write previous stacktrace first
	if e.StackTrace != "" {
		out.WriteString(e.StackTrace + "\n")
	}
	// append source info of src
	out.WriteString(stackTrace)

	e.StackTrace = out.String()

	return e
}

func parseSrc(src *ast.Source) string {
	var out bytes.Buffer

	out.WriteString(src.Pos.String() + "\n")
	out.WriteString(src.Line)

	return out.String()
}
