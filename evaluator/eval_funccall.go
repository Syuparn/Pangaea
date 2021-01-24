package evaluator

import (
	"fmt"
	"github.com/Syuparn/pangaea/object"
)

// used for Func#call
func evalFuncCall(
	env *object.Env,
	kwargs *object.PanObj,
	args ...object.PanObject,
) object.PanObject {
	if len(args) < 1 {
		return object.NewTypeErr("Func#call requires at least 1 arg")
	}

	self := args[0]
	// unshift args to ignore func itself
	args = args[1:]

	switch f := self.(type) {
	case *object.PanFunc:
		return evalPanFuncCall(f, env, kwargs, args...)
	case *object.PanBuiltIn:
		return f.Fn(env, kwargs, args...)
	default:
		err := object.NewTypeErr(
			fmt.Sprintf("`%s` is not callable.", self.Inspect()))
		return err
	}
}

func evalPanFuncCall(
	f *object.PanFunc,
	env *object.Env,
	kwargs *object.PanObj,
	args ...object.PanObject,
) object.PanObject {
	assignArgsToEnv(f.Env, f.Args().Elems, f.Kwargs(), args, kwargs)
	retVal := evalStmts(*f.Body(), f.Env)

	if err, ok := retVal.(*object.PanErr); ok {
		return appendStackTrace(err, (*f.Body())[0].Source())
	}

	return retVal
}

func assignArgsToEnv(
	env *object.Env,
	params []object.PanObject,
	kwargParams *object.PanObj,
	args []object.PanObject,
	kwargs *object.PanObj,
) {
	// nil padding if arity of args is fewer than that of params
	args = paddedArgs(args, params)

	for i, param := range params {
		ident, ok := param.(*object.PanStr)

		if ok {
			env.Set(object.GetSymHash(ident.Value), args[i])
		} else {
			// TODO: pattern matching
		}
	}

	// set argvars (like `\1`)
	for i, arg := range args {
		argVar := fmt.Sprintf(`\%d`, i+1)
		env.Set(object.GetSymHash(argVar), arg)
	}
	// `\0`
	env.Set(object.GetSymHash(`\0`), &object.PanArr{Elems: args})
	// `\`
	if len(args) > 0 {
		env.Set(object.GetSymHash(`\`), args[0])
	}

	for symHash, defaultPair := range *kwargParams.Pairs {
		kwargPair, ok := (*kwargs.Pairs)[symHash]

		if ok {
			env.Set(symHash, kwargPair.Value)
		} else {
			env.Set(symHash, defaultPair.Value)
		}
	}

	// set kwargvars (like `\hoge`)
	for _, kwargPair := range *kwargs.Pairs {
		sym := fmt.Sprintf("\\%s", kwargPair.Key.(*object.PanStr).Value)
		env.Set(object.GetSymHash(sym), kwargPair.Value)
	}
	// `\_`
	env.Set(object.GetSymHash("\\_"), kwargs)
}

func paddedArgs(args []object.PanObject, params []object.PanObject) []object.PanObject {
	lackedArityNum := len(params) - len(args)

	if lackedArityNum <= 0 {
		// if arity is sufficient, do nothing
		return args
	}

	paddedArgs := make([]object.PanObject, len(args))
	copy(paddedArgs, args)
	for i := 0; i < lackedArityNum; i++ {
		paddedArgs = append(paddedArgs, object.BuiltInNil)
	}

	return paddedArgs
}
