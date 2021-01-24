package runscript

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"os"
	"strings"
	"testing"
)

func TestStartREPL(t *testing.T) {
	in := strings.NewReader("a := 2\na * 2\n")
	out := &bytes.Buffer{}
	expected := strings.Join([]string{
		fmt.Sprintf("Pangaea %s (alpha)", Version),
		"multi : multi-line mode",
		"single: single-line mode (default)",
		"",
		">>> 2",
		">>> 4",
		">>> ",
	}, "\n")

	StartREPL(in, out)

	actual := out.String()
	if actual != expected {
		t.Errorf("output is wrong: \nexpected: \n%s\nactual: \n%s\n",
			expected, actual)
	}
}

func TestNewScanner(t *testing.T) {
	tests := []struct {
		in       io.Reader
		expected *_Scanner
	}{
		{
			os.Stdin,
			&_Scanner{
				mode:    &_SingleLineScannerState{},
				scanner: bufio.NewScanner(os.Stdin),
			},
		},
	}

	for _, tt := range tests {
		actual := newScanner(tt.in)
		// NOTE: scanners cannot be compared...
		if actual.mode.Name() != tt.expected.mode.Name() {
			t.Errorf("wrong mode: expected=%s, got=%s",
				tt.expected.mode.Name(), actual.mode.Name())
		}
	}
}

func TestSingleLineScan(t *testing.T) {
	tests := []struct {
		in       io.Reader
		expected string
	}{
		{
			strings.NewReader("a := 2\n"),
			"a := 2",
		},
		{
			strings.NewReader("'pangaea\n"),
			"'pangaea",
		},
		{
			strings.NewReader("'first\n'second\n"),
			"'first",
		},
	}

	for _, tt := range tests {
		scanner := newScanner(tt.in)
		actual, ok := scanner.Scan()
		if !ok {
			t.Fatalf("ok must be true (in testcase %s)", tt.expected)
		}

		if actual != tt.expected {
			t.Errorf("wrong scanned result: expected=%s, got=%s",
				tt.expected, actual)
		}
	}
}

func TestPrompt(t *testing.T) {
	tests := []struct {
		in       *_Scanner
		expected string
	}{
		{
			newScanner(os.Stdin),
			">>> ",
		},
		{
			&_Scanner{
				scanner: bufio.NewScanner(os.Stdin),
				mode:    newScannerState("multi"),
			},
			"<< multi-line mode (read lines until empty line is found) >>\n",
		},
	}

	for _, tt := range tests {
		if tt.in.Prompt() != tt.expected {
			t.Errorf("wrong prompt: expected=%s, got=%s",
				tt.expected, tt.in.Prompt())
		}
	}
}

func TestMultiLineScan(t *testing.T) {
	tests := []struct {
		in       io.Reader
		expected string
	}{
		{
			strings.NewReader("a := 2\n\n"),
			"a := 2\n",
		},
		{
			strings.NewReader("'pangaea\n\n"),
			"'pangaea\n",
		},
		{
			strings.NewReader("a := 2\nb := a * 2\n\n"),
			"a := 2\nb := a * 2\n",
		},
	}

	for _, tt := range tests {
		scanner := newScanner(tt.in)
		scanner.changeMode("multi")

		actual, ok := scanner.Scan()
		if !ok {
			t.Fatalf("ok must be true (in testcase %s)", tt.expected)
		}

		if actual != tt.expected {
			t.Errorf("wrong scanned result: expected=%s, got=%s",
				tt.expected, actual)
		}
	}
}

func TestNewScannerState(t *testing.T) {
	tests := []struct {
		actual   _ScannerState
		expected _ScannerState
	}{
		{
			newScannerState("single"),
			&_SingleLineScannerState{},
		},
		{
			newScannerState("multi"),
			&_MultiLineScannerState{},
		},
	}

	for _, tt := range tests {
		if tt.actual != tt.expected {
			t.Errorf("wrong output: expected=%+v, got=%+v",
				tt.expected, tt.actual)
		}
	}
}

func TestScannerChangeMode(t *testing.T) {
	in := strings.NewReader("single\nmulti\nmulti\n\nsingle\n\n")
	scanner := newScanner(in)

	_, ok := scanner.Scan()
	if !ok {
		t.Fatalf("first `single` must be scanned")
	}

	if scanner.mode.Name() != "single" {
		t.Errorf("scanner must be single mode. got=%s", scanner.mode.Name())
	}

	_, ok = scanner.Scan()
	if !ok {
		t.Fatalf("first `multi` must be scanned")
	}

	if scanner.mode.Name() != "multi" {
		t.Errorf("scanner must be multi mode. got=%s", scanner.mode.Name())
	}

	_, ok = scanner.Scan()
	if !ok {
		t.Fatalf("secode `multi` must be scanned")
	}

	if scanner.mode.Name() != "multi" {
		t.Errorf("scanner must be multi mode. got=%s", scanner.mode.Name())
	}

	_, ok = scanner.Scan()
	if !ok {
		t.Fatalf("secode `single` must be scanned")
	}

	if scanner.mode.Name() != "single" {
		t.Errorf("scanner must be single mode. got=%s", scanner.mode.Name())
	}
}
