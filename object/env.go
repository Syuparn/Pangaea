// Source code in this file is inherited and modified from
// "Writing an Interpreter in Go" https://interpreterbook.com/
// MIT License | Copyright (c) 2016-2017 Thorsten Ball
// see https://opensource.org/licenses/MIT

package object

func NewEnv() *Env {
	s := make(map[SymHash]PanObject)
	return &Env{s}
}

type Env struct {
	Store map[SymHash]PanObject
}

func (e *Env) Get(h SymHash) (PanObject, bool) {
	obj, ok := e.Store[h]
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
