package evaluator

import (
	"../ast"
	"../object"
	"fmt"
)

func evalIdent(ident *ast.Ident, env *object.Env) object.PanObject {
	// check if ident refers keyword
	const_, ok := constVal(ident.Value)

	if ok {
		return const_
	}

	val, ok := env.Get(object.GetSymHash(ident.Value))

	if !ok {
		err := object.NewNameErr(
			fmt.Sprintf("name `%s` is not defined.", ident.String()))
		return appendStackTrace(err, ident.Source())
	}

	return val
}

func constVal(name string) (object.PanObject, bool) {
	switch name {
	case "Int":
		return object.BuiltInIntObj, true
	case "Float":
		return object.BuiltInFloatObj, true
	case "Num":
		return object.BuiltInNumObj, true
	case "Nil":
		return object.BuiltInNilObj, true
	case "Str":
		return object.BuiltInStrObj, true
	case "Arr":
		return object.BuiltInArrObj, true
	case "Range":
		return object.BuiltInRangeObj, true
	case "Func":
		return object.BuiltInFuncObj, true
	case "Match":
		return object.BuiltInMatchObj, true
	case "Obj":
		return object.BuiltInObjObj, true
	case "BaseObj":
		return object.BuiltInBaseObj, true
	case "Map":
		return object.BuiltInMapObj, true
	case "true":
		return object.BuiltInTrue, true
	case "false":
		return object.BuiltInFalse, true
	case "nil":
		return object.BuiltInNil, true
	case "Err":
		return object.BuiltInErrObj, true
	case "AssertionErr":
		return object.BuiltInAssertionErr, true
	case "NameErr":
		return object.BuiltInNameErr, true
	case "NoPropErr":
		return object.BuiltInNoPropErr, true
	case "NotImplementedErr":
		return object.BuiltInNotImplementedErr, true
	case "SyntaxErr":
		return object.BuiltInSyntaxErr, true
	case "TypeErr":
		return object.BuiltInTypeErr, true
	case "ValueErr":
		return object.BuiltInValueErr, true
	case "ZeroDivisionErr":
		return object.BuiltInZeroDivisionErr, true
	case "_":
		return object.NewNotImplementedErr("Not implemented"), true
	default:
		return nil, false
	}
}
