package main

import (
	"bytes"
	"io"
	"strings"
	"syscall/js"

	"github.com/Syuparn/pangaea/di"
	"github.com/Syuparn/pangaea/evaluator"
	"github.com/Syuparn/pangaea/object"
	"github.com/Syuparn/pangaea/parser"
)

// Executor executes pangaea script as js function.
type Executor struct {
	constEnv *object.Env
}

// NewExecutor generates new Executor.
func NewExecutor() *Executor {
	return &Executor{
		// NOTE: prepare consts because they take time to be evaluated
		constEnv: setupEnv(),
	}
}

// Execute executes soruce code.
// args: (src, stdin) => ({res, stdout, errmsg})
func (e *Executor) Execute(this js.Value, args []js.Value) interface{} {
	src := e.setupSrc(args)
	stdin := e.setupStdin(args)
	stdout := &bytes.Buffer{}
	res, errmsg := e.execute(src, stdin, stdout)

	if errmsg != "" {
		return map[string]interface{}{
			"res":    "",
			"stdout": stdout.String(),
			"errmsg": errmsg,
		}
	}

	return map[string]interface{}{
		"res":    res.Repr(),
		"stdout": stdout.String(),
		"errmsg": errmsg,
	}
}

func (e *Executor) execute(
	src,
	in io.Reader,
	out io.Writer,
) (res object.PanObject, errmsg string) {
	// NOTE: IO must be injected globally otherwise built-in objects cannot refer it
	e.constEnv.InjectIO(in, out)
	env := object.NewEnclosedEnv(e.constEnv)

	node, err := parser.Parse(src)
	if err != nil {
		errmsg = err.Error()
		return
	}

	evaluated := evaluator.Eval(node, env)
	if err, ok := evaluated.(*object.PanErr); ok {
		errmsg = err.Inspect() + "\n" + err.StackTrace + "\n"
		return
	}

	res = evaluated
	return
}

func (e *Executor) setupSrc(args []js.Value) io.Reader {
	if len(args) == 0 || args[0].Type() != js.TypeString {
		// empty source code
		return strings.NewReader("")
	}

	return strings.NewReader(args[0].String())
}

func (e *Executor) setupStdin(args []js.Value) io.Reader {
	if len(args) < 2 || args[1].Type() != js.TypeString {
		// empty stdin
		return strings.NewReader("")
	}

	return strings.NewReader(args[1].String())
}

func setupEnv() *object.Env {
	env := object.NewEnvWithConsts()

	// necessary to setup built-in object props
	di.InjectBuiltInProps(env)

	// enable to use Kernel props directly in top-level
	// NOTE: InjectFrom must be called after BuiltInKernelObj is set up
	env.InjectFrom(object.BuiltInKernelObj)

	// NOTE: IO is inject when soruce is evaluated
	return env
}
