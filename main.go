package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/Syuparn/pangaea/repl"
	"github.com/Syuparn/pangaea/runscript"
)

var (
	parse = flag.Bool("parse", false, "only parse instead of eval")
)

func main() {
	flag.Parse()

	if srcFileName := flag.Arg(0); srcFileName != "" {
		exitCode := runScript(srcFileName)
		os.Exit(exitCode)
	}

	if *parse {
		runParsing()
		return
	}

	runRepl()
}

func runScript(fileName string) int {
	exitCode := runscript.Run(fileName, os.Stdin, os.Stdout)
	return exitCode
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
