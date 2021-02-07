// +build tools

package main

// import tools indirectly to control their versions in go.mod
import (
	_ "golang.org/x/tools/cmd/goyacc"
)
