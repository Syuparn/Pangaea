// Source code in this file is inherited and modified from
// "Writing an Interpreter in Go" https://interpreterbook.com/
// MIT License | Copyright (c) 2016-2017 Thorsten Ball
// see https://opensource.org/licenses/MIT

package runscript

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"strings"

	"github.com/Syuparn/pangaea/evaluator"
	"github.com/Syuparn/pangaea/object"
	"github.com/Syuparn/pangaea/parser"
)

// StartREPL starts Pangaea interpreter.
func StartREPL(in io.Reader, out io.Writer) {
	io.WriteString(out, fmt.Sprintf("Pangaea %s (alpha)\n", Version))
	io.WriteString(out, fmt.Sprintln("multi : multi-line mode"))
	io.WriteString(out, fmt.Sprintln("single: single-line mode (default)"))
	io.WriteString(out, fmt.Sprintln())

	scanner := newScanner(in)
	env := setup(in, out)

	for {
		io.WriteString(out, scanner.Prompt())
		scanned, ok := scanner.Scan()
		if !ok {
			return
		}

		program, err := parser.Parse(strings.NewReader(scanned))

		if err != nil {
			io.WriteString(out, err.Error())
			continue
		}

		evaluated := evaluator.Eval(program, env)

		io.WriteString(out, object.ReprStr(evaluated)+"\n")
	}
}

func newScanner(in io.Reader) *_Scanner {
	scanner := bufio.NewScanner(in)
	return &_Scanner{
		mode:    newScannerState("single"),
		scanner: scanner,
	}
}

type _Scanner struct {
	mode    _ScannerState
	scanner *bufio.Scanner
}

func (s *_Scanner) Scan() (string, bool) {
	scanned, ok := s.mode.Scan(s.scanner)

	if s.isModeChangeCode(scanned) {
		s.changeMode(scanned)
		return "", ok
	}
	return scanned, ok
}

func (s *_Scanner) Prompt() string { return s.mode.Prompt() }

func (s *_Scanner) isModeChangeCode(scanned string) bool {
	switch scanned {
	case "single":
		return true
	case "multi":
		return true
	}
	return false
}

func (s *_Scanner) changeMode(modeName string) {
	// if current mode is specified, do nothing
	if modeName == s.mode.Name() {
		return
	}

	s.mode = newScannerState(modeName)
}

type _ScannerState interface {
	Prompt() string
	Scan(scanner *bufio.Scanner) (string, bool)
	Name() string
}

func newScannerState(modeName string) _ScannerState {
	switch modeName {
	case "single":
		return &_SingleLineScannerState{}
	case "multi":
		return &_MultiLineScannerState{}
	}
	return nil
}

type _SingleLineScannerState struct{}

func (s *_SingleLineScannerState) Prompt() string { return ">>> " }

func (s *_SingleLineScannerState) Scan(scanner *bufio.Scanner) (string, bool) {
	ok := scanner.Scan()
	if !ok {
		return "", false
	}

	return scanner.Text(), true
}

func (s *_SingleLineScannerState) Name() string { return "single" }

type _MultiLineScannerState struct{}

func (s *_MultiLineScannerState) Prompt() string {
	return "<< multi-line mode (read lines until empty line is found) >>\n"
}

func (s *_MultiLineScannerState) Scan(scanner *bufio.Scanner) (string, bool) {
	// read lines until empty line is found
	var out bytes.Buffer

	for {
		ok := scanner.Scan()
		if !ok {
			return "", false
		}
		line := scanner.Text()
		if line == "" {
			break
		}
		if line == "single" {
			// extract single mode command
			return "single", true
		}
		out.WriteString(line + "\n")
	}

	return out.String(), true
}

func (s *_MultiLineScannerState) Name() string { return "multi" }
