package envs

import (
	"os"
	"path/filepath"
)

// environment variable names
const (
	// JargonFile is a file path of the jargon script file.
	JargonFileKey = "PANGAEA_JARGON_FILE"
)

// DefaultJargonFile is a default file path ot jargon script file.
func DefaultJargonFile() string {
	dir, _ := os.UserHomeDir()
	return filepath.Join(dir, ".jargon.pangaea")
}
