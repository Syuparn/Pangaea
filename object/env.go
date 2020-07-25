// Source code in this file is inherited and modified from
// "Writing an Interpreter in Go" https://interpreterbook.com/
// MIT License | Copyright (c) 2016-2017 Thorsten Ball
// see https://opensource.org/licenses/MIT

package object

import (
	"io"
)

func NewEnv() *Env {
	s := make(map[SymHash]PanObject)
	return &Env{s, nil}
}

func NewEnclosedEnv(e *Env) *Env {
	s := make(map[SymHash]PanObject)
	return &Env{s, e}
}

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
	env.Set(GetSymHash("Match"), BuiltInMatchObj)
	env.Set(GetSymHash("Obj"), BuiltInObjObj)
	env.Set(GetSymHash("BaseObj"), BuiltInBaseObj)
	env.Set(GetSymHash("Map"), BuiltInMapObj)
	env.Set(GetSymHash("true"), BuiltInTrue)
	env.Set(GetSymHash("false"), BuiltInFalse)
	env.Set(GetSymHash("nil"), BuiltInNil)
	env.Set(GetSymHash("Err"), BuiltInErrObj)
	env.Set(GetSymHash("AssertionErr"), BuiltInAssertionErr)
	env.Set(GetSymHash("NameErr"), BuiltInNameErr)
	env.Set(GetSymHash("NoPropErr"), BuiltInNoPropErr)
	env.Set(GetSymHash("NotImplementedErr"), BuiltInNotImplementedErr)
	env.Set(GetSymHash("SyntaxErr"), BuiltInSyntaxErr)
	env.Set(GetSymHash("TypeErr"), BuiltInTypeErr)
	env.Set(GetSymHash("ValueErr"), BuiltInValueErr)
	env.Set(GetSymHash("ZeroDivisionErr"), BuiltInZeroDivisionErr)
	env.Set(GetSymHash("_"), NewNotImplementedErr("Not implemented"))

	return env
}

type Env struct {
	Store map[SymHash]PanObject
	outer *Env
}

func (e *Env) Get(h SymHash) (PanObject, bool) {
	obj, ok := e.Store[h]

	// if not found, search outer scope
	if !ok && e.outer != nil {
		obj, ok = e.outer.Get(h)
	}

	return obj, ok
}

func (e *Env) Set(h SymHash, obj PanObject) {
	e.Store[h] = obj
}

func (e *Env) Items() PanObject {
	pairs := make(map[SymHash]Pair)
	for h, obj := range e.Store {
		strObj, ok := SymHash2Str(h)

		if !ok {
			panic("Failed to fetch PanStr by SymHash2Str().\n" +
				"StrTable may be broken.")
		}

		pairs[h] = Pair{strObj, obj}
	}

	return PanObjInstancePtr(&pairs)
}

func (e *Env) Outer() *Env {
	return e.outer
}

func (e *Env) InjectIO(in io.Reader, out io.Writer) {
	// define const `IO` containing io of args
	ioObj := &PanIO{In: in, Out: out}
	e.Set(GetSymHash("IO"), ioObj)
}

func (e *Env) InjectRecur(recurFunc BuiltInFunc) {
	// define const `recur` builtInFunc, which can only be used inside iter
	recur := &PanBuiltIn{Fn: recurFunc}
	e.Set(GetSymHash("recur"), recur)
}
