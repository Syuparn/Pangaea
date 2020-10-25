package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/Syuparn/pangaea/repl"
)

var (
	parse = flag.Bool("parse", false, "only parse instead of eval")
)

func main() {
	flag.Parse()

	if *parse {
		runParsing()
		return
	}

	runRepl()
}

func runRepl() {
	fmt.Println("Pangaea ver0.0.0 (alpha)")
	fmt.Println("Under construction...")
	repl.Start(os.Stdin, os.Stdout)
}

func runParsing() {
	fmt.Println("Pangaea ver0.0.0 (alpha)")
	fmt.Println("Parsing mode(parsed ast is not evaluated)")
	repl.StartParsing(os.Stdin, os.Stdout)
}
