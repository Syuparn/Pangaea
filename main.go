package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/Syuparn/pangaea/repl"
	"github.com/Syuparn/pangaea/runscript"
)

var (
	parse      = flag.Bool("parse", false, "only parse instead of eval")
	testCmdSet = flag.NewFlagSet("test", flag.ExitOnError)
)

func main() {
	// test mode
	if len(os.Args) >= 2 && os.Args[1] == "test" {
		testCmdSet.Parse(os.Args[2:])
		if path := testCmdSet.Arg(0); path != "" {
			exitCode := runTest(path)
			os.Exit(exitCode)
		}
	}

	// normal mode
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

func runTest(path string) int {
	exitCode := runscript.RunTest(path, os.Stdin, os.Stdout)
	return exitCode
}

func runScript(fileName string) int {
	exitCode := runscript.Run(fileName, os.Stdin, os.Stdout)
	return exitCode
}

func runRepl() {
	fmt.Println("Pangaea ver0.1.1 (alpha)")
	repl.Start(os.Stdin, os.Stdout)
}

func runParsing() {
	fmt.Println("Pangaea ver0.1.1 (alpha)")
	fmt.Println("Parsing mode(parsed ast is not evaluated)")
	repl.StartParsing(os.Stdin, os.Stdout)
}
