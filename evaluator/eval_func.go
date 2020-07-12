package evaluator

import (
	"../ast"
	"../object"
)

func evalFunc(node *ast.FuncLiteral, env *object.Env) object.PanObject {
	return evalCallable(node.FuncComponent, env, object.FUNC_FUNC)
}

func evalCallable(
	component ast.FuncComponent,
	env *object.Env,
	funcType object.FuncType,
) object.PanObject {
	args := []object.PanObject{}
	for _, argNode := range component.Args {
		var arg object.PanObject
		if ident, ok := argNode.(*ast.Ident); ok {
			arg = &object.PanStr{Value: ident.String()}
		} else {
			arg = Eval(argNode, env)
		}

		args = append(args, arg)
	}

	kwargs := evalKwargs(component.Kwargs, env)

	wrapper := &FuncWrapperImpl{
		codeStr: component.String(),
		args:    &object.PanArr{Elems: args},
		kwargs:  kwargs,
		body:    &component.Body,
	}

	return &object.PanFunc{
		FuncWrapper: wrapper,
		FuncType:    funcType,
		Env:         object.NewEnclosedEnv(env),
	}
}
