package runscript

import (
	"io"
	"os"

	"github.com/Syuparn/pangaea/di"
	"github.com/Syuparn/pangaea/evaluator"
	"github.com/Syuparn/pangaea/object"
	"github.com/Syuparn/pangaea/parser"
)

// Run runs Pangaea script source file.
func Run(fileName string, in io.Reader, out io.Writer) {
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
		io.WriteString(out, err.Error())
		return
	}

	node, err := parser.Parse(fp)
	if err != nil {
		io.WriteString(out, err.Error()+"\n")
		return
	}

	evaluated := evaluator.Eval(node, env)

	if err, ok := evaluated.(*object.PanErr); ok {
		io.WriteString(out, err.Inspect()+"\n")
		io.WriteString(out, err.StackTrace+"\n")
	}
}
