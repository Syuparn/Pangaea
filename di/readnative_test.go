package di

import (
	"testing"

	"github.com/Syuparn/pangaea/object"
)

func TestReadNativeCodeSuccess(t *testing.T) {
	tests := []struct {
		srcName  string
		expected map[object.SymHash]object.Pair
	}{
		{
			"testSuccess",
			map[object.SymHash]object.Pair{
				object.GetSymHash("a"): {
					Key:   object.NewPanStr("a"),
					Value: object.NewPanInt(1),
				},
				object.GetSymHash("b"): {
					Key:   object.NewPanStr("b"),
					Value: object.NewPanInt(2),
				},
			},
		},
	}

	for _, tt := range tests {
		pairsPtr, err := readNativeCode("testdata/"+tt.srcName, object.NewEnv())
		if err != nil {
			t.Fatalf("err must be nil. got=%s", err.Error())
		}
		if pairsPtr == nil {
			t.Fatalf("pairs must not be nil.")
		}
		pairs := *pairsPtr

		if len(pairs) != len(tt.expected) {
			t.Fatalf("len(pairs) must be %d. got=%d", len(tt.expected), len(pairs))
		}

		for k, expected := range tt.expected {
			actual, ok := pairs[k]
			if !ok {
				t.Fatalf("key %s(%d) not found in pairs.",
					expected.Key.Inspect(), k)
			}

			testValue(t, actual.Key, expected.Key)
			testValue(t, actual.Value, expected.Value)
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
			"NameErr: name `undefinedVar` is not defined\n" +
				"\"testdata/testErr.pangaea\" line: 2, col: 1\n" +
				"undefinedVar\n" +
				"\"testdata/testErr.pangaea\" line: 2, col: 13\n" +
				"undefinedVar",
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
			t.Errorf("err msg must be `\n%s\n`. got=`\n%s\n`", tt.expectedErrMsg, err.Error())
		}
	}
}

func TestMustReadNativeCodeSuccess(t *testing.T) {
	tests := []struct {
		srcName  string
		expected map[object.SymHash]object.Pair
	}{
		{
			"testSuccess",
			map[object.SymHash]object.Pair{
				object.GetSymHash("a"): {
					Key:   object.NewPanStr("a"),
					Value: object.NewPanInt(1),
				},
				object.GetSymHash("b"): {
					Key:   object.NewPanStr("b"),
					Value: object.NewPanInt(2),
				},
			},
		},
	}

	for _, tt := range tests {
		pairsPtr := mustReadNativeCode("testdata/"+tt.srcName, object.NewEnv())
		if pairsPtr == nil {
			t.Fatalf("pairsPtr must not be nil.")
		}
		pairs := *pairsPtr

		if len(pairs) != len(tt.expected) {
			t.Fatalf("len(pairs) must be %d. got=%d", len(tt.expected), len(pairs))
		}

		for k, expected := range tt.expected {
			actual, ok := pairs[k]
			if !ok {
				t.Fatalf("key %s(%d) not found in pairs.",
					expected.Key.Inspect(), k)
			}

			testValue(t, actual.Key, expected.Key)
			testValue(t, actual.Value, expected.Value)
		}
	}
}

func TestMustReadNativeCodePanics(t *testing.T) {
	tests := []struct {
		srcName        string
		expectedErrMsg string
	}{
		{
			"testNotObj",
			"result must be ObjType. got=ArrType",
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
