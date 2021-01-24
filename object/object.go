// Source code in this file is inherited and modified from
// "Writing an Interpreter in Go" https://interpreterbook.com/
// MIT License | Copyright (c) 2016-2017 Thorsten Ball
// see https://opensource.org/licenses/MIT

package object

// PanObjType is type information to designate each PanObject.
type PanObjType string

// PanObject is an interface of all Pangaea object implementations.
type PanObject interface {
	Type() PanObjType
	Inspect() string
	Proto() PanObject
}

// PanScalar is an interface of scalar PanObject, which does not have child components.
type PanScalar interface {
	PanObject
	Hash() HashKey
}
