package evaluator

import (
	"github.com/Syuparn/pangaea/ast"
	"github.com/Syuparn/pangaea/object"
)

func evalFunc(node *ast.FuncLiteral, env *object.Env) object.PanObject {
	return evalCallable(node.FuncComponent, env, object.FuncFunc)
}

func evalCallable(
	component ast.FuncComponent,
	env *object.Env,
	funcKind object.FuncKind,
) object.PanObject {
	args := []object.PanObject{}
	for _, argNode := range component.Args {
		var arg object.PanObject
		if ident, ok := argNode.(*ast.Ident); ok {
			arg = object.NewPanStr(ident.String())
		} else {
			// TODO: error handling for pattern match exprs
			arg = Eval(argNode, env)
		}

		args = append(args, arg)
	}

	kwargs, err := evalKwargs(component.Kwargs, env)

	if err != nil {
		return err
	}

	wrapper := &FuncWrapperImpl{
		codeStr: component.String(),
		args:    object.NewPanArr(args...),
		kwargs:  kwargs,
		body:    &component.Body,
	}

	if funcKind == object.FuncFunc {
		return object.NewPanFunc(wrapper, object.NewEnclosedEnv(env))
	}

	return object.NewPanIter(wrapper, object.NewEnclosedEnv(env))
}
