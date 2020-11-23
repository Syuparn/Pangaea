package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/Syuparn/pangaea/runscript"
)

var (
	oneLiner            = flag.String("e", "", "run one-line script")
	readsLines          = flag.Bool("n", false, "assign stdin each line to \\")
	readsAndWritesLines = flag.Bool("p", false, "similar to -n but also print to evaluated values")
	version             = flag.Bool("v", false, "show version")
	testCmdSet          = flag.NewFlagSet("test", flag.ExitOnError)
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

	// show version
	if *version {
		showVersion()
		os.Exit(0)
	}

	// run one-liner
	if *oneLiner != "" {
		src := wrapSource(*oneLiner, *readsLines, *readsAndWritesLines)
		exitCode := runOneLiner(src)
		os.Exit(exitCode)
	}

	if srcFileName := flag.Arg(0); srcFileName != "" {
		exitCode := runScript(srcFileName)
		os.Exit(exitCode)
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

func runOneLiner(src string) int {
	exitCode := runscript.RunSource(src, os.Stdin, os.Stdout)
	return exitCode
}

func runRepl() {
	runscript.StartREPL(os.Stdin, os.Stdout)
}

func showVersion() {
	fmt.Println(runscript.Version)
}

func wrapSource(original string, readsLines, readsAndWritesLines bool) string {
	if readsAndWritesLines {
		return fmt.Sprintf(runscript.ReadStdinLinesAndWritesTemplate, original)
	}
	if readsLines {
		return fmt.Sprintf(runscript.ReadStdinLinesTemplate, original)
	}
	return original
}
