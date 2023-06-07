package di

import (
	"fmt"
	"os"
	"strings"

	"github.com/Syuparn/pangaea/object"
	"github.com/Syuparn/pangaea/parser"
)

func kernelImport(
	env *object.Env,
	kwargs *object.PanObj,
	args ...object.PanObject,
) object.PanObject {
	if len(args) < 1 {
		return object.NewTypeErr("import requires at least 1 arg")
	}

	fileNameObj, ok := args[0].(*object.PanStr)
	if !ok {
		return object.NewTypeErr("\\1 must be str")
	}

	fileName := fileNameObj.Value
	if !strings.HasSuffix(fileName, ".pangaea") {
		fileName += ".pangaea"
	}

	f, err := os.Open(fileName)
	if err != nil {
		return object.NewFileNotFoundErr(fmt.Sprintf("failed to open %q", fileName))
	}

	// NOTE: object.NewEnv cannot be used because an empty env does not have built-in objects
	newEnv := object.NewEnclosedEnv(env.Global())
	result := eval(parser.NewReader(f, fileName), newEnv)
	if result.Type() == object.ErrType {
		return result
	}

	return newEnv.Items()
}
