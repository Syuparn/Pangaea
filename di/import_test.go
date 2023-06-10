package di

import (
	"path/filepath"
	"testing"

	"github.com/Syuparn/pangaea/object"
)

func TestEvalKernelImport(t *testing.T) {
	// NOTE: since form of Abspath is OS-dependent, we use fixture to prepare expected string
	abspath := func(path string) string {
		p, _ := filepath.Abs(path)
		return p
	}

	tests := []struct {
		input    string
		expected object.PanObject
	}{
		// success
		{
			`import("./testdata/testSuccess.pangaea")`,
			object.PanObjInstancePtr(&map[object.SymHash]object.Pair{
				object.GetSymHash("_PANGAEA_SOURCE_PATH"): {Key: object.NewPanStr("_PANGAEA_SOURCE_PATH"), Value: object.NewPanStr(abspath("./testdata/testSuccess.pangaea"))},
				object.GetSymHash("a"):                    {Key: object.NewPanStr("a"), Value: object.NewPanInt(1)},
				object.GetSymHash("b"):                    {Key: object.NewPanStr("b"), Value: object.NewPanInt(2)},
			}),
		},
		// extension
		{
			`import("./testdata/testSuccess")`,
			object.PanObjInstancePtr(&map[object.SymHash]object.Pair{
				object.GetSymHash("_PANGAEA_SOURCE_PATH"): {Key: object.NewPanStr("_PANGAEA_SOURCE_PATH"), Value: object.NewPanStr(abspath("./testdata/testSuccess.pangaea"))},
				object.GetSymHash("a"):                    {Key: object.NewPanStr("a"), Value: object.NewPanInt(1)},
				object.GetSymHash("b"):                    {Key: object.NewPanStr("b"), Value: object.NewPanInt(2)},
			}),
		},
		// nested import
		{
			`import("./testdata/importing")`,
			object.PanObjInstancePtr(&map[object.SymHash]object.Pair{
				object.GetSymHash("_PANGAEA_SOURCE_PATH"): {Key: object.NewPanStr("_PANGAEA_SOURCE_PATH"), Value: object.NewPanStr(abspath("./testdata/importing.pangaea"))},
				object.GetSymHash("a"):                    {Key: object.NewPanStr("a"), Value: object.NewPanInt(1)},
			}),
		},
		// standard module
		{
			`import("dummy")`,
			object.PanObjInstancePtr(&map[object.SymHash]object.Pair{
				object.GetSymHash("message"): {Key: object.NewPanStr("message"), Value: object.NewPanStr("This is a dummy module.")},
			}),
		},
		{
			`import("dummy_native")`,
			object.PanObjInstancePtr(&map[object.SymHash]object.Pair{
				object.GetSymHash("message"): {Key: object.NewPanStr("message"), Value: object.NewPanStr("This is a dummy module.")},
			}),
		},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			actual := testEval(t, tt.input)
			testValue(t, actual, tt.expected)
		})
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
		{
			`_PANGAEA_SOURCE_PATH := 1; import("./testdata/testSuccess")`,
			object.NewTypeErr("_PANGAEA_SOURCE_PATH 1 must be str"),
		},
		{
			`import("notfound")`,
			object.NewFileNotFoundErr("failed to read native module \"notfound\": open modules/notfound.pangaea: file does not exist"),
		},
		{
			`import("dummy_native_wrong")`,
			object.NewSyntaxErr("failed to parse"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			actual := testEval(t, tt.input)
			testValue(t, actual, tt.expected)
		})
	}
}

func TestEvalKernelInvite(t *testing.T) {
	tests := []struct {
		input    string
		expected object.PanObject
	}{
		{
			`invite!("./testdata/testSuccess.pangaea"); a`,
			object.NewPanInt(1),
		},
		{
			`invite!("./testdata/testSuccess"); a`,
			object.NewPanInt(1),
		},
		{
			`invite!("./testdata/testSuccess")`,
			object.BuiltInNil,
		},
		// invite! update variables
		{
			`a := "original"; invite!("./testdata/testSuccess"); a`,
			object.NewPanInt(1),
		},
		// nested invite
		{
			`invite!("./testdata/inviting"); a`,
			object.NewPanInt(1),
		},
		// _PANGAEA_SOURCE_PATH is not changed
		{
			`_PANGAEA_SOURCE_PATH := "dummy"; invite!("./testdata/inviting"); _PANGAEA_SOURCE_PATH`,
			object.NewPanStr("dummy"),
		},
		// standard module
		{
			`invite!("dummy"); message`,
			object.NewPanStr("This is a dummy module."),
		},
		{
			`invite!("dummy_native"); message`,
			object.NewPanStr("This is a dummy module."),
		},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			actual := testEval(t, tt.input)
			testValue(t, actual, tt.expected)
		})
	}
}

func TestEvalKernelInviteError(t *testing.T) {
	tests := []struct {
		input    string
		expected object.PanObject
	}{
		{
			`invite!("./testdata/notfound")`,
			object.NewFileNotFoundErr("failed to open \"./testdata/notfound.pangaea\""),
		},
		{
			`invite!("./testdata/syntaxError")`,
			object.NewSyntaxErr("failed to parse"),
		},
		{
			`invite!(1)`,
			object.NewTypeErr("\\1 must be str"),
		},
		{
			`invite!()`,
			object.NewTypeErr("invite! requires at least 1 arg"),
		},
		{
			`_PANGAEA_SOURCE_PATH := 1; invite!("./testdata/testSuccess")`,
			object.NewTypeErr("_PANGAEA_SOURCE_PATH 1 must be str"),
		},
		{
			`invite!("notfound")`,
			object.NewFileNotFoundErr("failed to read native module \"notfound\": open modules/notfound.pangaea: file does not exist"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			actual := testEval(t, tt.input)
			testValue(t, actual, tt.expected)
		})
	}
}
