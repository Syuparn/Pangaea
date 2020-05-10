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
	return BuiltInObjObj, true
}

func (e *Env) Set(h SymHash, obj PanObject) {
	return
}

func (e *Env) Items() PanObject {
	return BuiltInZeroInt
}
