package evaluator

import (
	"../object"
)

// used for Func#call
func evalFuncCall(
	env *object.Env,
	kwargs *object.PanObj,
	args ...object.PanObject,
) object.PanObject {
	self := args[0]
	f, ok := self.(*object.PanFunc)
	if !ok {
		// TODO: error handling
		panic("1st arg is not PanFunc")
	}

	// unshift args to ignore func itself
	args = args[1:]

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

	for symHash, defaultPair := range *kwargParams.Pairs {
		kwargPair, ok := (*kwargs.Pairs)[symHash]

		if ok {
			env.Set(symHash, kwargPair.Value)
		} else {
			env.Set(symHash, defaultPair.Value)
		}
	}
}

func paddedArgs(args []object.PanObject, params []object.PanObject) []object.PanObject {
	lackedArityNum := len(args) - len(params)

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
