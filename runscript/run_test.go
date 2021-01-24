package runscript

import (
	"os"
	"testing"
)

func BenchmarkSetup(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_ = setup(os.Stdin, os.Stdout)
	}
}
