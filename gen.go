package main

// NOTE: this file is used to make auto-generated src files.

// make parser y.go from yacc file parser.go.y
//go:generate go run golang.org/x/tools/cmd/goyacc -o ./parser/y.go -v ./parser/y.output ./parser/parser.go.y

// make statik/statik.go,
// which contains all native Pangaea source files in /native as binary
//go:generate go run github.com/rakyll/statik -src=native/ -include=*.pangaea
