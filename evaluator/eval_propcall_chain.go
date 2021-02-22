package evaluator

import (
	"github.com/Syuparn/pangaea/ast"
	"github.com/Syuparn/pangaea/object"
)

type _PropCallMiddleware func(next _PropCallMiddlewareHandler) _PropCallMiddlewareHandler

type _PropCallMiddlewareHandler func(
	env *object.Env,
	recv object.PanObject,
	propName string,
	prop object.PanObject,
	chainArg object.PanObject,
	args []object.PanObject,
	kwargs *object.PanObj,
) object.PanObject

func mergePropCallMiddlewares(middlewares ..._PropCallMiddleware) _PropCallMiddleware {
	return func(next _PropCallMiddlewareHandler) _PropCallMiddlewareHandler {
		merged := next
		for i := len(middlewares) - 1; i >= 0; i-- {
			merged = middlewares[i](merged)
		}
		return merged
	}
}

func newChainMiddleware(
	chain ast.Chain,
) _PropCallMiddleware {
	// `=@` chain keeps evaluated nil elements
	if chain.Main == ast.List && chain.Additional == ast.Strict {
		return mergePropCallMiddlewares(
			keepNilPropCallListChainMiddleware,
			findPropMiddleware,
		)
	}
	// `~@` chain replaces evaluated nil/err elements by recv
	if chain.Main == ast.List && chain.Additional == ast.Thoughtful {
		return mergePropCallMiddlewares(
			keepNilPropCallListChainMiddleware,
			propCallThoughtfulChainMiddleware,
			findPropMiddleware,
		)
	}

	// the other chain middlewares simply consists of two chain middlewares
	mainMiddleware := newPropCallMainChainMiddleware(chain.Main)
	additionalMiddleware := newPropCallAdditionalChainMiddleware(chain.Additional)
	return mergePropCallMiddlewares(
		mainMiddleware,
		additionalMiddleware,
		findPropMiddleware,
	)
}

func newPropCallMainChainMiddleware(c ast.MainChain) _PropCallMiddleware {
	switch c {
	case ast.Scalar:
		return propCallNothingMiddleware
	case ast.List:
		return squashNilPropCallListChainMiddleware
	case ast.Reduce:
		return propCallReduceChainMiddleware
	default:
		return propCallNothingMiddleware
	}
}

func findPropMiddleware(next _PropCallMiddlewareHandler) _PropCallMiddlewareHandler {
	return func(
		env *object.Env,
		recv object.PanObject,
		propName string,
		_ object.PanObject,
		chainArg object.PanObject,
		args []object.PanObject,
		kwargs *object.PanObj,
	) object.PanObject {
		// get prop of recv
		prop, isMissing := evalProp(propName, recv)
		if err, ok := prop.(*object.PanErr); ok {
			return err
		}

		// prepend prop name to arg if _missing is called
		argsToPass := args
		if isMissing {
			propName := object.NewPanStr(propName)
			argsToPass = append([]object.PanObject{propName}, argsToPass...)
		}

		return next(env, recv, propName, prop, chainArg, argsToPass, kwargs)
	}
}

func squashNilPropCallListChainMiddleware(next _PropCallMiddlewareHandler) _PropCallMiddlewareHandler {
	return func(
		env *object.Env,
		recv object.PanObject,
		propName string,
		_ object.PanObject,
		chainArg object.PanObject,
		args []object.PanObject,
		kwargs *object.PanObj,
	) object.PanObject {
		iter, err := iterOf(env, recv)
		if err != nil {
			return err
		}

		elems := []object.PanObject{}

		for {
			nextRecv, err := iter.Next(env)
			if err != nil {
				if err.Kind() == object.StopIterErr {
					break
				}
				return err
			}

			elem := next(env, nextRecv, propName, nil, chainArg, args, kwargs)
			if elem.Type() == object.ErrType {
				// raise error
				return elem
			}
			if elem.Type() == object.NilType {
				continue
			}

			elems = append(elems, elem)
		}

		return object.NewPanArr(elems...)
	}
}

func keepNilPropCallListChainMiddleware(next _PropCallMiddlewareHandler) _PropCallMiddlewareHandler {
	return func(
		env *object.Env,
		recv object.PanObject,
		propName string,
		_ object.PanObject,
		chainArg object.PanObject,
		args []object.PanObject,
		kwargs *object.PanObj,
	) object.PanObject {
		iter, err := iterOf(env, recv)
		if err != nil {
			return err
		}

		elems := []object.PanObject{}

		for {
			nextRecv, err := iter.Next(env)
			if err != nil {
				if err.Kind() == object.StopIterErr {
					break
				}
				return err
			}

			elem := next(env, nextRecv, propName, nil, chainArg, args, kwargs)
			elems = append(elems, elem)
		}

		return object.NewPanArr(elems...)
	}
}

func propCallReduceChainMiddleware(next _PropCallMiddlewareHandler) _PropCallMiddlewareHandler {
	return func(
		env *object.Env,
		recv object.PanObject,
		propName string,
		_ object.PanObject,
		chainArg object.PanObject,
		args []object.PanObject,
		kwargs *object.PanObj,
	) object.PanObject {
		iter, err := iterOf(env, recv)
		if err != nil {
			return err
		}

		acc := chainArg

		for {
			nextRet, err := iter.Next(env)
			if err != nil {
				if err.Kind() == object.StopIterErr {
					break
				}
				return err
			}

			// prepend nextRet to args
			argsToPass := append([]object.PanObject{nextRet}, args...)

			evaluated := next(env, acc, propName, nil, chainArg, argsToPass, kwargs)
			if evaluated.Type() == object.ErrType {
				return evaluated
			}
			acc = evaluated
		}

		return acc
	}
}

func newPropCallAdditionalChainMiddleware(c ast.AdditionalChain) _PropCallMiddleware {
	switch c {
	case ast.Lonely:
		return propCallLonelyChainMiddleware
	case ast.Thoughtful:
		return propCallThoughtfulChainMiddleware
	// NOTE: newChainMiddleware deals with combination `=@`
	case ast.Strict:
		return propCallNothingMiddleware
	default:
		return propCallNothingMiddleware
	}
}

func propCallLonelyChainMiddleware(next _PropCallMiddlewareHandler) _PropCallMiddlewareHandler {
	return func(
		env *object.Env,
		recv object.PanObject,
		propName string,
		_ object.PanObject,
		chainArg object.PanObject,
		args []object.PanObject,
		kwargs *object.PanObj,
	) object.PanObject {
		if recv.Type() == object.NilType {
			return recv
		}
		return next(env, recv, propName, nil, chainArg, args, kwargs)
	}
}

func propCallThoughtfulChainMiddleware(next _PropCallMiddlewareHandler) _PropCallMiddlewareHandler {
	return func(
		env *object.Env,
		recv object.PanObject,
		propName string,
		_ object.PanObject,
		chainArg object.PanObject,
		args []object.PanObject,
		kwargs *object.PanObj,
	) object.PanObject {
		ret := next(env, recv, propName, nil, chainArg, args, kwargs)
		if ret.Type() == object.ErrType || ret.Type() == object.NilType {
			return recv
		}
		return ret
	}
}

func propCallNothingMiddleware(next _PropCallMiddlewareHandler) _PropCallMiddlewareHandler {
	return func(
		env *object.Env,
		recv object.PanObject,
		propName string,
		prop object.PanObject,
		chainArg object.PanObject,
		args []object.PanObject,
		kwargs *object.PanObj,
	) object.PanObject {
		return next(env, recv, propName, prop, chainArg, args, kwargs)
	}
}
