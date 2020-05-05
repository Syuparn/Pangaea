package object

import (
	"testing"
)

func TestIntType(t *testing.T) {
	intObj := PanInt{10}
	if intObj.Type() != INT_TYPE {
		t.Fatalf("wrong type: expected=%s, got=%s", INT_TYPE, intObj.Type())
	}
}
