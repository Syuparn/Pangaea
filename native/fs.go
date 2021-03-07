package native

import (
	"embed"
)

// FS exports current directory file system.
// NOTE: FS must not be declared in di package because
// path root of unittest differs from runtime one!
//go:embed *.pangaea testdata/*.pangaea
var FS embed.FS
