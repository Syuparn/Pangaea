// Source code in this file is inherited and modified from
// "Writing an Interpreter in Go" https://interpreterbook.com/
// MIT License | Copyright (c) 2016-2017 Thorsten Ball
// see https://opensource.org/licenses/MIT

package object

import (
	"io"
)

// NewEnv makes new environment of variables.
func NewEnv() *Env {
	s := make(map[SymHash]PanObject)
	return &Env{s, nil}
}

// NewEnclosedEnv makes new environment of variables inside e.
// It is used to make closure.
func NewEnclosedEnv(e *Env) *Env {
	s := make(map[SymHash]PanObject)
	return &Env{s, e}
}

// NewEnvWithConsts makes new global environment, which includes all standart objects.
func NewEnvWithConsts() *Env {
	env := NewEnv()
	env.Set(GetSymHash("Int"), BuiltInIntObj)
	env.Set(GetSymHash("Float"), BuiltInFloatObj)
	env.Set(GetSymHash("Num"), BuiltInNumObj)
	env.Set(GetSymHash("Nil"), BuiltInNilObj)
	env.Set(GetSymHash("Str"), BuiltInStrObj)
	env.Set(GetSymHash("Arr"), BuiltInArrObj)
	env.Set(GetSymHash("Range"), BuiltInRangeObj)
	env.Set(GetSymHash("Func"), BuiltInFuncObj)
	env.Set(GetSymHash("Iter"), BuiltInIterObj)
	env.Set(GetSymHash("Iterable"), BuiltInIterableObj)
	env.Set(GetSymHash("Comparable"), BuiltInComparableObj)
	env.Set(GetSymHash("Wrappable"), BuiltInWrappableObj)
	env.Set(GetSymHash("Match"), BuiltInMatchObj)
	env.Set(GetSymHash("Obj"), BuiltInObjObj)
	env.Set(GetSymHash("BaseObj"), BuiltInBaseObj)
	env.Set(GetSymHash("Map"), BuiltInMapObj)
	env.Set(GetSymHash("Diamond"), BuiltInDiamondObj)
	env.Set(GetSymHash("Kernel"), BuiltInKernelObj)
	env.Set(GetSymHash("JSON"), BuiltInJSONObj)
	env.Set(GetSymHash("Either"), BuiltInEitherObj)
	env.Set(GetSymHash("EitherVal"), BuiltInEitherValObj)
	env.Set(GetSymHash("EitherErr"), BuiltInEitherErrObj)
	env.Set(GetSymHash("true"), BuiltInTrue)
	env.Set(GetSymHash("false"), BuiltInFalse)
	env.Set(GetSymHash("nil"), BuiltInNil)
	env.Set(GetSymHash("Err"), BuiltInErrObj)
	env.Set(GetSymHash("AssertionErr"), BuiltInAssertionErr)
	env.Set(GetSymHash("NameErr"), BuiltInNameErr)
	env.Set(GetSymHash("NoPropErr"), BuiltInNoPropErr)
	env.Set(GetSymHash("NotImplementedErr"), BuiltInNotImplementedErr)
	env.Set(GetSymHash("StopIterErr"), BuiltInStopIterErr)
	env.Set(GetSymHash("SyntaxErr"), BuiltInSyntaxErr)
	env.Set(GetSymHash("TypeErr"), BuiltInTypeErr)
	env.Set(GetSymHash("ValueErr"), BuiltInValueErr)
	env.Set(GetSymHash("ZeroDivisionErr"), BuiltInZeroDivisionErr)
	env.Set(GetSymHash("_"), BuiltInNotImplemented)

	return env
}

// NewCopiedEnv makes copied environment of env, which is independent of original one.
func NewCopiedEnv(env *Env) *Env {
	newStore := map[SymHash]PanObject{}
	// copy all variables to new store
	for k, v := range env.Store {
		newStore[k] = v
	}

	return &Env{
		Store: newStore,
		outer: env.outer,
	}
}

// Env is an environment of variables.
type Env struct {
	Store map[SymHash]PanObject
	outer *Env
}

// Get fetches variable value from the environment.
func (e *Env) Get(h SymHash) (PanObject, bool) {
	obj, ok := e.Store[h]

	// if not found, search outer scope
	if !ok && e.outer != nil {
		obj, ok = e.outer.Get(h)
	}

	return obj, ok
}

// Set sets variable to the environment.
func (e *Env) Set(h SymHash, obj PanObject) {
	e.Store[h] = obj
}

// Items returns all variables in the enviroment as obj.
func (e *Env) Items() PanObject {
	pairs := make(map[SymHash]Pair)
	for h, obj := range e.Store {
		strObj, ok := SymHash2Str(h)

		if !ok {
			panic("Failed to fetch PanStr by SymHash2Str().\n" +
				"strTable may be broken.")
		}

		pairs[h] = Pair{strObj, obj}
	}

	return PanObjInstancePtr(&pairs)
}

// Outer returns outer environment.
func (e *Env) Outer() *Env {
	return e.outer
}

// InjectIO injects reader and writer for `IO` object
func (e *Env) InjectIO(in io.Reader, out io.Writer) {
	// define const `IO` containing io of args
	ioObj := NewPanIO(in, out)
	e.Set(GetSymHash("IO"), ioObj)
}

// InjectRecur injects built-in func `recur`.
func (e *Env) InjectRecur(recurFunc BuiltInFunc) {
	// define const `recur` builtInFunc, which can only be used inside iter
	recur := NewPanBuiltInFunc(recurFunc)
	e.Set(GetSymHash("recur"), recur)
}

// InjectFrom injects all obj props to env.
func (e *Env) InjectFrom(obj *PanObj) {
	for sym, pair := range *obj.Pairs {
		e.Set(sym, pair.Value)
	}
}
