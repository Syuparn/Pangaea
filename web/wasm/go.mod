module github.com/Syuparn/pangaea/web/wasm

go 1.17

require github.com/Syuparn/pangaea v0.6.2

require (
	github.com/dlclark/regexp2 v1.4.0 // indirect
	github.com/lithammer/dedent v1.1.0 // indirect
	github.com/macrat/simplexer v0.0.0-20180110131648-bce8e0661570 // indirect
	github.com/tanaton/dtoa v0.0.0-20190918101016-f12936c87cdb // indirect
)

// bundle to patch lexer
replace github.com/macrat/simplexer v0.0.0-20180110131648-bce8e0661570 => ../../third_party/simplexer
