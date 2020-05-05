package object

// initialize built-in objects like Int, Arr, Str...
// NOTE: Props are inserted in package eval not to make
// package object and package builtin circular reference
var builtInIntObj = &PanObj{}
var builtInFloatObj = &PanObj{}
