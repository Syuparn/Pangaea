package main

import (
	"./repl"
	"fmt"
	"os"
)

func main() {
	runRepl()
}

func runRepl() {
	fmt.Println("Pangaea ver0.0.0 (alpha)")
	fmt.Println("Now you can only parsing (eval is under construction...)")
	repl.Start(os.Stdin, os.Stdout)
}
