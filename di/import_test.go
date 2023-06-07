package di

import (
	"testing"

	"github.com/Syuparn/pangaea/object"
)

func TestEvalKernelImport(t *testing.T) {
	tests := []struct {
		input    string
		expected object.PanObject
	}{
		// success
		{
			`import("./testdata/testSuccess.pangaea")`,
			object.PanObjInstancePtr(&map[object.SymHash]object.Pair{
				object.GetSymHash("a"): {Key: object.NewPanStr("a"), Value: object.NewPanInt(1)},
				object.GetSymHash("b"): {Key: object.NewPanStr("b"), Value: object.NewPanInt(2)},
			}),
		},
		// extension
		{
			`import("./testdata/testSuccess")`,
			object.PanObjInstancePtr(&map[object.SymHash]object.Pair{
				object.GetSymHash("a"): {Key: object.NewPanStr("a"), Value: object.NewPanInt(1)},
				object.GetSymHash("b"): {Key: object.NewPanStr("b"), Value: object.NewPanInt(2)},
			}),
		},
	}

	for _, tt := range tests {
		actual := testEval(t, tt.input)
		testValue(t, actual, tt.expected)
	}
}

func TestEvalKernelImportError(t *testing.T) {
	tests := []struct {
		input    string
		expected object.PanObject
	}{
		{
			`import("./testdata/notfound")`,
			object.NewFileNotFoundErr("failed to open \"./testdata/notfound.pangaea\""),
		},
		{
			`import("./testdata/syntaxError")`,
			object.NewSyntaxErr("failed to parse"),
		},
		{
			`import(1)`,
			object.NewTypeErr("\\1 must be str"),
		},
		{
			`import()`,
			object.NewTypeErr("import requires at least 1 arg"),
		},
	}

	for _, tt := range tests {
		actual := testEval(t, tt.input)
		testValue(t, actual, tt.expected)
	}
}
