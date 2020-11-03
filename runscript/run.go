package runscript

import (
	"fmt"
	"io"
	"os"

	"github.com/Syuparn/pangaea/di"
	"github.com/Syuparn/pangaea/evaluator"
	"github.com/Syuparn/pangaea/object"
	"github.com/Syuparn/pangaea/parser"
)

// Run runs Pangaea script source file.
func Run(fileName string, in io.Reader, out io.Writer) int {
	env := object.NewEnvWithConsts()
	// setup object `IO`
	env.InjectIO(in, out)

	// necessary to setup built-in object props
	di.InjectBuiltInProps()

	// enable to use Kernel props directly in top-level
	// NOTE: InjectFrom must be called after BuiltInKernelObj is set up
	env.InjectFrom(object.BuiltInKernelObj)

	fp, err := os.Open(fileName)
	if err != nil {
		fmt.Fprint(os.Stderr, err.Error()+"\n")
		return 1
	}

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
