package evaluator

import (
	"bytes"
	"github.com/Syuparn/pangaea/ast"
	"github.com/Syuparn/pangaea/object"
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
		sSym := object.NewPanStr("S")
		evaluatedS := builtInCallProp(env, object.EmptyPanObjPtr(),
			object.EmptyPanObjPtr(), evaluated, sSym)

		evaluatedStr, ok := object.TraceProtoOfStr(evaluatedS)
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

	return object.NewPanStr(out.String())
}
