package di

import (
	"testing"

	"github.com/Syuparn/pangaea/object"
)

func TestReadNativeCodeSuccess(t *testing.T) {
	tests := []struct {
		srcName  string
		expected map[string]object.PanObject
	}{
		{
			"testSuccess",
			map[string]object.PanObject{
				"a": object.NewPanInt(1),
				"b": object.NewPanInt(2),
			},
		},
	}

	for _, tt := range tests {
		ctn, err := readNativeCode("testdata/"+tt.srcName, object.NewEnv())
		if err != nil {
			t.Fatalf("err must be nil. got=%s", err.Error())
		}

		if len(ctn) != len(tt.expected) {
			t.Fatalf("len(ctn) must be %d. got=%d", len(tt.expected), len(ctn))
		}

		for k, expected := range tt.expected {
			actual, ok := ctn[k]
			if !ok {
				t.Fatalf("key %s not found in ctn.", k)
			}

			testValue(t, actual, expected)
		}
	}
}

func TestReadNativeCodeFailed(t *testing.T) {
	tests := []struct {
		srcName        string
		expectedErrMsg string
	}{
		{
			"testNotExist",
			"failed to read native testdata/testNotExist props in " +
				"native/testdata/testNotExist.pangaea",
		},
		{
			"testErr",
			"NameErr: name `undefinedVar` is not defined",
		},
		{
			"testNotObj",
			"result must be ObjType. got=ArrType",
		},
	}

	for _, tt := range tests {
		_, err := readNativeCode("testdata/"+tt.srcName, object.NewEnv())
		if err == nil {
			t.Fatal("err must be occurred.")
		}

		if err.Error() != tt.expectedErrMsg {
			t.Errorf("err msg must be \n%s.\ngot=\n%s", tt.expectedErrMsg, err.Error())
		}
	}
}

func TestMustReadNativeCodeSuccess(t *testing.T) {
	tests := []struct {
		srcName  string
		expected map[string]object.PanObject
	}{
		{
			"testSuccess",
			map[string]object.PanObject{
				"a": object.NewPanInt(1),
				"b": object.NewPanInt(2),
			},
		},
	}

	for _, tt := range tests {
		ctn := mustReadNativeCode("testdata/"+tt.srcName, object.NewEnv())

		if len(ctn) != len(tt.expected) {
			t.Fatalf("len(ctn) must be %d. got=%d", len(tt.expected), len(ctn))
		}

		for k, expected := range tt.expected {
			actual, ok := ctn[k]
			if !ok {
				t.Fatalf("key %s not found in ctn.", k)
			}

			testValue(t, actual, expected)
		}
	}
}

func TestMustReadNativeCodePanics(t *testing.T) {
	tests := []struct {
		srcName        string
		expectedErrMsg string
	}{
		{
			"testErr",
			"NameErr: name `undefinedVar` is not defined",
		},
	}

	for _, tt := range tests {
		func() {
			defer func() {
				err := recover()
				if err != tt.expectedErrMsg {
					t.Errorf("msg must be %v. got=%v", tt.expectedErrMsg, err)
				}
			}()

			_ = mustReadNativeCode("testdata/"+tt.srcName, object.NewEnv())
		}()
	}
}
