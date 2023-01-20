package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/Syuparn/pangaea/envs"
	"github.com/Syuparn/pangaea/runscript"
)

var (
	oneLiner            = flag.String("e", "", "run one-line script")
	jargon              = flag.Bool("j", false, fmt.Sprintf("read jargon script saved in $%s (`~/.jargon.pangaea` by default)", envs.JargonFileKey))
	readsLines          = flag.Bool("n", false, "assign stdin each line to \\")
	readsAndWritesLines = flag.Bool("p", false, "similar to -n but also print to evaluated values")
	version             = flag.Bool("v", false, "show version")
	testCmdSet          = flag.NewFlagSet("test", flag.ExitOnError)
)

// TODO: refactor handling of each options (by DDD or something else?)
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

	src := ""
	if *jargon {
		src = readJargon()
	}

	// run one-liner
	if *oneLiner != "" {
		src += wrapSource(*oneLiner, *readsLines, *readsAndWritesLines)
		exitCode := run(src)
		os.Exit(exitCode)
	}

	if srcFileName := flag.Arg(0); srcFileName != "" {
		fileSrc, exitCode := runscript.ReadFile(srcFileName)
		if exitCode != 0 {
			os.Exit(exitCode)
		}
		src += fileSrc

		exitCode = run(src)
		os.Exit(exitCode)
	}

	runRepl(src)
}

func runTest(path string) int {
	exitCode := runscript.RunTest(path, os.Stdin, os.Stdout)
	return exitCode
}

func run(src string) int {
	exitCode := runscript.RunSource(src, os.Stdin, os.Stdout)
	return exitCode
}

func runRepl(src string) {
	runscript.StartREPL(src, os.Stdin, os.Stdout)
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

func readJargon() string {
	fileName := envs.DefaultJargonFile()
	if n, ok := os.LookupEnv(envs.JargonFileKey); ok {
		fileName = n
	}

	bytes, err := os.ReadFile(fileName)
	// error does not matter
	if err != nil {
		return ""
	}

	return string(bytes) + "\n"
}
