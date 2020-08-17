package evaluator

import (
	"../ast"
	"../object"
	"bytes"
)

func evalEmbeddedStr(node *ast.EmbeddedStr, env *object.Env) object.PanObject {
	// cache strs because embeddedstr ast is reverse order of source code
	evaluatedStrs := []string{node.Latter}
	for n := node.Former; n != nil; n = n.Former {
		evaluated := Eval(n.Expr, env)
		if err, ok := evaluated.(*object.PanErr); ok {
			return appendStackTrace(err, node.Source())
		}

		// call .S to convert into str
		sSym := &object.PanStr{Value: "S"}
		evaluatedS := builtInCallProp(env, object.EmptyPanObjPtr(),
			object.EmptyPanObjPtr(), evaluated, sSym)

		// TODO: allow ancestor of str
		evaluatedStr, ok := evaluatedS.(*object.PanStr)
		if !ok {
			err := object.NewValueErr(".S must return str")
			return appendStackTrace(err, node.Source())
		}

		// prepend
		evaluatedStrs = append(
			[]string{n.Str, evaluatedStr.Value}, evaluatedStrs...)
	}

	var out bytes.Buffer

	for _, str := range evaluatedStrs {
		out.WriteString(str)
	}

	return &object.PanStr{Value: out.String()}
}
