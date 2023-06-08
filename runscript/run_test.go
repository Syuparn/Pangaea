package runscript

import (
	"bytes"
	"os"
	"testing"
)

func BenchmarkSetup(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_ = setup(os.Stdin, os.Stdout, "")
	}
}

func TestRunSource(t *testing.T) {
	tests := []struct {
		name     string
		fileName string
		source   string
		expected string
	}{
		{
			"read from file",
			"foo.pangaea",
			`"Hello".p`,
			"Hello\n",
		},
		{
			"read from stdin",
			"<stdin>",
			`"Hello".p`,
			"Hello\n",
		},
		{
			"Kernel props are injected",
			"foo.pangaea",
			`assertEq(1+1, 2)`,
			"",
		},
		{
			"_PANGAEA_SOURCE_PATH is set",
			"foo.pangaea",
			// NOTE: since _PANGAEA_SOURCE_PATH is os-dependent, we test only the file name
			`_PANGAEA_SOURCE_PATH.split(sep: "[\\\\/]")[-1].p`,
			"foo.pangaea\n",
		},
		{
			"_PANGAEA_SOURCE_PATH is not set if filepath is a dummy",
			"<stdin>",
			`nil.try.{ _PANGAEA_SOURCE_PATH.p }`,
			"",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var out bytes.Buffer
			status := RunSource(tt.source, tt.fileName, os.Stdin, &out)

			if status != 0 {
				t.Fatalf("status must be 0. got %v", status)
			}

			actual := out.String()
			if actual != tt.expected {
				t.Errorf("wrong output: expected=%+v, got=%+v", tt.expected, actual)
			}
		})
	}
}
