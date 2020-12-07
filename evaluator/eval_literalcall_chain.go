package evaluator

import (
	"github.com/Syuparn/pangaea/ast"
	"github.com/Syuparn/pangaea/object"
)

var (
	nextSym = object.NewPanStr("next")
	iterSym = object.NewPanStr("_iter")
)

type _LiteralCallMiddleware func(next _LiteralCallMiddlewareHandler) _LiteralCallMiddlewareHandler

type _LiteralCallMiddlewareHandler func(
	env *object.Env,
	recv object.PanObject,
	chainArg object.PanObject,
	args []object.PanObject,
	kwargs *object.PanObj,
) object.PanObject

func mergeLiteralCallMiddlewares(middlewares ..._LiteralCallMiddleware) _LiteralCallMiddleware {
	return func(next _LiteralCallMiddlewareHandler) _LiteralCallMiddlewareHandler {
		merged := next
		for i := len(middlewares) - 1; i >= 0; i-- {
			merged = middlewares[i](merged)
		}
		return merged
	}
}

func newLiteralCallChainMiddleware(
	chain ast.Chain,
) _LiteralCallMiddleware {
	// `=@` chain keeps evaluated nil elements
	if chain.Main == ast.List && chain.Additional == ast.Strict {
		return keepNilLiteralCallListChainMiddleware
	}
	// `~@` chain replaces evaluated nil/err elements by recv
	if chain.Main == ast.List && chain.Additional == ast.Thoughtful {
		return mergeLiteralCallMiddlewares(
			keepNilLiteralCallListChainMiddleware,
			literalCallThoughtfulChainMiddleware,
		)
	}
	// `~$` chain replaces evaluated nil/err elements by last acc
	if chain.Main == ast.Reduce && chain.Additional == ast.Thoughtful {
		return literalCallThoughtfulReduceChainMiddleware
	}

	// the other chain middlewares simply consists of two chain middlewares
	mainMiddleware := newLiteralCallMainChainMiddleware(chain.Main)
	additionalMiddleware := newLiteralCallAdditionalChainMiddleware(chain.Additional)
	return mergeLiteralCallMiddlewares(mainMiddleware, additionalMiddleware)
}

func newLiteralCallMainChainMiddleware(c ast.MainChain) _LiteralCallMiddleware {
	switch c {
	case ast.Scalar:
		return literalCallNothingMiddleware
	case ast.List:
		return squashNilLiteralCallListChainMiddleware
	case ast.Reduce:
		return literalCallReduceChainMiddleware
	default:
		return literalCallNothingMiddleware
	}
}

func squashNilLiteralCallListChainMiddleware(next _LiteralCallMiddlewareHandler) _LiteralCallMiddlewareHandler {
	return func(
		env *object.Env,
		recv object.PanObject,
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

			elem := next(env, nextRecv, chainArg, args, kwargs)
			if elem.Type() == object.NilType {
				continue
			}

			elems = append(elems, elem)
		}

		return &object.PanArr{Elems: elems}
	}
}

func keepNilLiteralCallListChainMiddleware(next _LiteralCallMiddlewareHandler) _LiteralCallMiddlewareHandler {
	return func(
		env *object.Env,
		recv object.PanObject,
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

			elem := next(env, nextRecv, chainArg, args, kwargs)
			elems = append(elems, elem)
		}

		return &object.PanArr{Elems: elems}
	}
}

func literalCallReduceChainMiddleware(next _LiteralCallMiddlewareHandler) _LiteralCallMiddlewareHandler {
	return func(
		env *object.Env,
		recv object.PanObject,
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
			nextRecv := &object.PanArr{Elems: []object.PanObject{acc, nextRet}}

			evaluated := next(env, nextRecv, chainArg, args, kwargs)
			if evaluated.Type() == object.ErrType {
				return evaluated
			}
			acc = evaluated
		}

		return acc
	}
}

func literalCallThoughtfulReduceChainMiddleware(next _LiteralCallMiddlewareHandler) _LiteralCallMiddlewareHandler {
	return func(
		env *object.Env,
		recv object.PanObject,
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
			nextRecv := &object.PanArr{Elems: []object.PanObject{acc, nextRet}}

			evaluated := next(env, nextRecv, chainArg, args, kwargs)
			if evaluated.Type() == object.ErrType {
				// replace evaluated value with last acc
				continue
			}
			acc = evaluated
		}

		return acc
	}
}

func newLiteralCallAdditionalChainMiddleware(c ast.AdditionalChain) _LiteralCallMiddleware {
	switch c {
	case ast.Lonely:
		return literalCallLonelyChainMiddleware
	case ast.Thoughtful:
		return literalCallThoughtfulChainMiddleware
	// NOTE: newChainMiddleware deals with combination `=@`
	case ast.Strict:
		return literalCallNothingMiddleware
	default:
		return literalCallNothingMiddleware
	}
}

func literalCallLonelyChainMiddleware(next _LiteralCallMiddlewareHandler) _LiteralCallMiddlewareHandler {
	return func(
		env *object.Env,
		recv object.PanObject,
		chainArg object.PanObject,
		args []object.PanObject,
		kwargs *object.PanObj,
	) object.PanObject {
		if recv.Type() == object.NilType {
			return recv
		}
		return next(env, recv, chainArg, args, kwargs)
	}
}

func literalCallThoughtfulChainMiddleware(next _LiteralCallMiddlewareHandler) _LiteralCallMiddlewareHandler {
	return func(
		env *object.Env,
		recv object.PanObject,
		chainArg object.PanObject,
		args []object.PanObject,
		kwargs *object.PanObj,
	) object.PanObject {
		ret := next(env, recv, chainArg, args, kwargs)
		if ret.Type() == object.ErrType || ret.Type() == object.NilType {
			return recv
		}
		return ret
	}
}

func literalCallNothingMiddleware(next _LiteralCallMiddlewareHandler) _LiteralCallMiddlewareHandler {
	return func(
		env *object.Env,
		recv object.PanObject,
		chainArg object.PanObject,
		args []object.PanObject,
		kwargs *object.PanObj,
	) object.PanObject {
		return next(env, recv, chainArg, args, kwargs)
	}
}

// shorthand for Iter object
type iterHandler struct {
	iter object.PanObject
}

func (h *iterHandler) Next(env *object.Env) (object.PanObject, *object.PanErr) {
	// call `(iter).next`
	nextRet := builtInCallProp(env, object.EmptyPanObjPtr(),
		object.EmptyPanObjPtr(), h.iter, nextSym)

	if err, ok := nextRet.(*object.PanErr); ok {
		return nil, err
	}

	return nextRet, nil
}

func iterOf(
	env *object.Env,
	obj object.PanObject,
) (*iterHandler, *object.PanErr) {
	iter := builtInCallProp(env, object.EmptyPanObjPtr(),
		object.EmptyPanObjPtr(), obj, iterSym)

	if err, ok := iter.(*object.PanErr); ok {
		return nil, err
	}

	if iter == object.BuiltInNil {
		err := object.NewTypeErr("recv must have prop `_iter`")
		return nil, err
	}

	return &iterHandler{iter: iter}, nil
}
