//go:build tools
// +build tools

package main

// import tools indirectly to control their versions in go.mod
import (
	_ "github.com/Songmu/gocredits/cmd/gocredits"
	_ "golang.org/x/tools/cmd/goyacc"
)
