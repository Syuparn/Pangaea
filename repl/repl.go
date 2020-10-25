// Source code in this file is inherited and modified from
// "Writing an Interpreter in Go" https://interpreterbook.com/
// MIT License | Copyright (c) 2016-2017 Thorsten Ball
// see https://opensource.org/licenses/MIT

package repl

import (
	"bufio"
	"fmt"
	"io"
	"strings"

	"github.com/Syuparn/pangaea/evaluator"
	"github.com/Syuparn/pangaea/object"
	"github.com/Syuparn/pangaea/parser"
)

// PROMPT is a prefix string printed in Pangaea interpreter.
const PROMPT = "> "

// Start starts Pangaea interpreter.
func Start(in io.Reader, out io.Writer) {
	scanner := bufio.NewScanner(in)
	env := object.NewEnvWithConsts()
	// setup object `IO`
	env.InjectIO(in, out)

	for {
		fmt.Printf(PROMPT)
		ok := scanner.Scan()
		if !ok {
			return
		}
		// temporally convert to string to parse each line
		line := scanner.Text()

		program, err := parser.Parse(strings.NewReader(line))

		if err != nil {
			io.WriteString(out, err.Error())
			continue
		}

		// necessary to setup built-in object props
		evaluator.InjectBuiltInProps()

		evaluated := evaluator.Eval(program, env)

		io.WriteString(out, evaluated.Inspect()+"\n")
	}
}

// StartParsing starts Pangaea interpreter but only lexing and parsing.
func StartParsing(in io.Reader, out io.Writer) {
	scanner := bufio.NewScanner(in)

	for {
		fmt.Printf(PROMPT)
		ok := scanner.Scan()
		if !ok {
			return
		}
		// temporally convert to string to parse each line
		line := scanner.Text()

		program, err := parser.Parse(strings.NewReader(line))

		if err != nil {
			io.WriteString(out, err.Error())
			continue
		}

		io.WriteString(out, program.String()+"\n")
	}
}
