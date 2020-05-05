// Source code in this file is inherited and modified from
// "Writing an Interpreter in Go" https://interpreterbook.com/
// MIT License | Copyright (c) 2016-2017 Thorsten Ball
// see https://opensource.org/licenses/MIT

package object

type PanObjType string

type PanObject interface {
	Type() PanObjType
	Inspect() string
	Proto() PanObject
}

type PanScalar interface {
	PanObject
	Hash() HashKey
}
