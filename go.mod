module github.com/Syuparn/pangaea

go 1.21

require (
	github.com/Songmu/gocredits v0.3.0
	github.com/dlclark/regexp2 v1.4.0
	github.com/labstack/echo/v4 v4.10.2
	github.com/lithammer/dedent v1.1.0
	github.com/macrat/simplexer v0.0.0-20180110131648-bce8e0661570
	github.com/tanaton/dtoa v0.0.0-20190918101016-f12936c87cdb
	golang.org/x/tools v0.6.0
)

require (
	github.com/labstack/gommon v0.4.0 // indirect
	github.com/mattn/go-colorable v0.1.13 // indirect
	github.com/mattn/go-isatty v0.0.17 // indirect
	github.com/valyala/bytebufferpool v1.0.0 // indirect
	github.com/valyala/fasttemplate v1.2.2 // indirect
	golang.org/x/crypto v0.14.0 // indirect
	golang.org/x/net v0.17.0 // indirect
	golang.org/x/sys v0.13.0 // indirect
	golang.org/x/text v0.13.0 // indirect
)

// bundle to patch lexer
replace github.com/macrat/simplexer v0.0.0-20180110131648-bce8e0661570 => ./third_party/simplexer
