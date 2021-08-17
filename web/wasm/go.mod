module github.com/Syuparn/pangaea/web/wasm

go 1.16

require github.com/Syuparn/pangaea v0.6.1

// bundle to patch lexer
replace github.com/macrat/simplexer v0.0.0-20180110131648-bce8e0661570 => ../../third_party/simplexer
