package runscript

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/Syuparn/pangaea/di"
	"github.com/Syuparn/pangaea/evaluator"
	"github.com/Syuparn/pangaea/object"
	"github.com/Syuparn/pangaea/parser"
)

// RunTest runs all test script files in path until error is raised.
func RunTest(path string, in io.Reader, out io.Writer) int {
	env := setup(in, out, "")
	exitCode := 0

	filepath.Walk(path, func(path string, info os.FileInfo, err error) error {
		if !strings.HasSuffix(path, ".pangaea") {
			return nil
		}

		out.Write([]byte(fmt.Sprintf("run:  %s\n", path)))

		exitCode = runTest(path, in, out, env)
		if exitCode != 0 {
			return fmt.Errorf("test in %s failed", path)
		}

		out.Write([]byte(fmt.Sprintf("pass: %s\n", path)))
		return nil
	})
	return exitCode
}

// RunSource runs input src.
func RunSource(src string, fileName string, in io.Reader, out io.Writer) int {
	env := setup(in, out, fileName)
	reader := strings.NewReader(src)
	exitCode := runSource(parser.NewReader(reader, fileName), in, out, env)
	return exitCode
}

func ReadFile(fileName string) (string, int) {
	bytes, err := os.ReadFile(fileName)
	if err != nil {
		fmt.Fprint(os.Stderr, err.Error()+"\n")
		return "", 1
	}

	return string(bytes), 0
}

func runTest(fileName string, in io.Reader, out io.Writer, env *object.Env) int {
	fp, err := os.Open(fileName)
	if err != nil {
		fmt.Fprint(os.Stderr, err.Error()+"\n")
		return 1
	}

	inner := object.NewEnclosedEnv(env)
	inner.SetSourceFilePath(fileName)

	exitCode := runSource(parser.NewReader(fp, fileName), in, out, env)
	return exitCode
}

func runSource(fp *parser.Reader, in io.Reader, out io.Writer, env *object.Env) int {
	node, err := parser.Parse(fp)
	if err != nil {
		fmt.Fprint(os.Stderr, err.Error()+"\n")
		return 1
	}

	evaluated := evaluator.Eval(node, env)

	if err, ok := evaluated.(*object.PanErr); ok {
		fmt.Fprint(os.Stderr, err.Inspect()+"\n")
		fmt.Fprint(os.Stderr, err.StackTrace+"\n")
		return 1
	}

	return 0
}

func setup(in io.Reader, out io.Writer, fileName string) *object.Env {
	env := object.NewEnvWithConsts()
	// setup object `IO`
	env.InjectIO(in, out)

	// set current source file path to env
	env.SetSourceFilePath(fileName)

	// necessary to setup built-in object props
	di.InjectBuiltInProps(env)

	// enable to use Kernel props directly in top-level
	// NOTE: InjectFrom must be called after BuiltInKernelObj is set up
	env.InjectFrom(object.BuiltInKernelObj)

	return env
}
