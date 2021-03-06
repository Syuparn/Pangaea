package runscript

// Version represents Pangaea version.
// NOTE: this version is patched by GoReleaser
var Version = "master (unstable)"

const (
	// ReadStdinLinesTemplate is a template src for one-liner option
	// so that each line of stdin is assigned to \
	ReadStdinLinesTemplate = "<>@{%s}"
	// ReadStdinLinesAndWritesTemplate is a template src for one-liner option
	// similar to ReadStdinLinesTemplate but also prints evaluated values to stdout
	ReadStdinLinesAndWritesTemplate = "<>@{%s}@p"
)
