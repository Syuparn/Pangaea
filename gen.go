package main

// NOTE: this file is used to make auto-generated src files.

// make parser y.go from yacc file parser.go.y
//go:generate go run golang.org/x/tools/cmd/goyacc -o ./parser/y.go -v ./parser/y.output ./parser/parser.go.y
