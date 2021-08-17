module github.com/Syuparn/pangaea

go 1.16

require (
	github.com/dlclark/regexp2 v1.4.0
	github.com/lithammer/dedent v1.1.0
	github.com/macrat/simplexer v0.0.0-20180110131648-bce8e0661570
	github.com/tanaton/dtoa v0.0.0-20190918101016-f12936c87cdb
	golang.org/x/tools v0.1.0
)

// bundle to patch lexer
replace github.com/macrat/simplexer v0.0.0-20180110131648-bce8e0661570 => ./third_party/simplexer
