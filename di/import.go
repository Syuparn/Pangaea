package di

import (
	"fmt"
	"os"
	"path/filepath"
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

	// NOTE: if fileName is relative, it is based on the evaluating source file (not based on where pangaea command is executed)
	if p, ok := env.Get(object.GetSymHash(object.SourcePathVar)); ok {
		if p.Type() != object.StrType {
			return object.NewTypeErr(fmt.Sprintf("%s %s must be str", object.SourcePathVar, p.Inspect()))
		}
		sourcePath := p.(*object.PanStr).Value
		fileName = filepath.Join(filepath.Dir(sourcePath), fileName)
	}

	f, err := os.Open(fileName)
	if err != nil {
		return object.NewFileNotFoundErr(fmt.Sprintf("failed to open %q", fileName))
	}

	// NOTE: object.NewEnv cannot be used because an empty env does not have built-in objects
	newEnv := object.NewEnclosedEnv(env.Global())

	// set imported file path to SourcePathVar for inner env
	newEnv.SetSourceFilePath(fileName)

	result := eval(parser.NewReader(f, fileName), newEnv)
	if result.Type() == object.ErrType {
		return result
	}

	return newEnv.Items()
}
