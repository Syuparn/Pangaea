// Source code in this file is inherited and modified from
// "Writing an Interpreter in Go" https://interpreterbook.com/
// MIT License | Copyright (c) 2016-2017 Thorsten Ball
// see https://opensource.org/licenses/MIT

package object

func NewEnv() *Env {
	s := make(map[SymHash]PanObject)
	return &Env{s, nil}
}

func NewEnclosedEnv(e *Env) *Env {
	s := make(map[SymHash]PanObject)
	return &Env{s, e}
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
