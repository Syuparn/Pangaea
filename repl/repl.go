// Source code in this file is inherited and modified from
// "Writing an Interpreter in Go" https://interpreterbook.com/
// MIT License | Copyright (c) 2016-2017 Thorsten Ball
// see https://opensource.org/licenses/MIT

package repl

import (
	"../parser"
	"bufio"
	"fmt"
	"io"
	"strings"
)

const PROMPT = "> "

func Start(in io.Reader, out io.Writer) {
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