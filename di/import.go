package di

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/Syuparn/pangaea/native"
	"github.com/Syuparn/pangaea/object"
	"github.com/Syuparn/pangaea/parser"
	"github.com/Syuparn/pangaea/props/modules"
)

func kernelImport(
	env *object.Env,
	kwargs *object.PanObj,
	args ...object.PanObject,
) object.PanObject {
	if len(args) < 1 {
		return object.NewTypeErr("import requires at least 1 arg")
	}

	importPathObj, ok := args[0].(*object.PanStr)
	if !ok {
		return object.NewTypeErr("\\1 must be str")
	}

	return importModule(env, importPathObj.Value)
}

func kernelInvite(
	env *object.Env,
	kwargs *object.PanObj,
	args ...object.PanObject,
) object.PanObject {
	if len(args) < 1 {
		return object.NewTypeErr("invite! requires at least 1 arg")
	}

	importPathObj, ok := args[0].(*object.PanStr)
	if !ok {
		return object.NewTypeErr("\\1 must be str")
	}

	return inviteModule(env, importPathObj.Value)
}

func importModule(env *object.Env, importPath string) object.PanObject {
	// relative path like "./foo/bar"
	if strings.HasPrefix(importPath, ".") {
		return importRelative(env, importPath)
	}

	// NOTE: object.NewEnv cannot be used because an empty env does not have built-in objects
	// NOTE: object.NewEnclosedEnv(env) cannot be used otherwise variables in this env affects the imported module
	newEnv := object.NewEnclosedEnv(env.Global())

	return injectStandardModule(newEnv, importPath)
}

func importRelative(env *object.Env, importPath string) object.PanObject {
	f, importPath, errObj := readSourceFile(env, importPath)
	if errObj != nil {
		return errObj
	}

	// NOTE: object.NewEnv cannot be used because an empty env does not have built-in objects
	// NOTE: object.NewEnclosedEnv(env) cannot be used otherwise variables in this env affects the imported module
	newEnv := object.NewEnclosedEnv(env.Global())

	// set imported file path to SourcePathVar for inner env
	newEnv.SetSourceFilePath(importPath)

	result := eval(parser.NewReader(f, importPath), newEnv)
	if result.Type() == object.ErrType {
		return result
	}

	return newEnv.Items()
}

func inviteModule(env *object.Env, importPath string) object.PanObject {
	// relative path like "./foo/bar"
	if strings.HasPrefix(importPath, ".") {
		return inviteRelative(env, importPath)
	}

	m := injectStandardModule(env, importPath)
	if m.Type() == object.ErrType {
		return m
	}

	return object.BuiltInNil
}

func inviteRelative(env *object.Env, importPath string) object.PanObject {
	f, importPath, errObj := readSourceFile(env, importPath)
	if errObj != nil {
		return errObj
	}

	origPath, existsPath := env.Get(object.GetSymHash(object.SourcePathVar))
	env.SetSourceFilePath(importPath)
	// HACK: set the original value again for the following process
	defer func() {
		if existsPath {
			env.Set(object.GetSymHash(object.SourcePathVar), origPath)
		}
	}()

	result := eval(parser.NewReader(f, importPath), env)
	if result.Type() == object.ErrType {
		return result
	}

	return object.BuiltInNil
}

func readSourceFile(env *object.Env, importPath string) (io.Reader, string, *object.PanErr) {
	// add extension
	if !strings.HasSuffix(importPath, ".pangaea") {
		importPath += ".pangaea"
	}

	// NOTE: if importPath is relative, it is based on the evaluating source file (not based on where pangaea command is executed)
	if p, ok := env.Get(object.GetSymHash(object.SourcePathVar)); ok {
		if p.Type() != object.StrType {
			return nil, "", object.NewTypeErr(fmt.Sprintf("%s %s must be str", object.SourcePathVar, p.Inspect()))
		}
		sourcePath := p.(*object.PanStr).Value
		importPath = filepath.Join(filepath.Dir(sourcePath), importPath)
	}

	f, err := os.Open(importPath)
	if err != nil {
		return nil, "", object.NewFileNotFoundErr(fmt.Sprintf("failed to open %q", importPath))
	}

	return f, importPath, nil
}

func injectStandardModule(env *object.Env, importPath string) object.PanObject {
	// find built-in
	if m, ok := modules.Modules[importPath]; ok {
		modules.InjectTo(env, m())
		return env.Items()
	}

	// find native otherwise
	return injectNativeStandardModule(env, importPath)
}

func injectNativeStandardModule(env *object.Env, importPath string) object.PanObject {
	filePath := fmt.Sprintf("modules/%s.pangaea", importPath)

	fp, err := native.FS.Open(filePath)
	if err != nil {
		return object.NewFileNotFoundErr(fmt.Sprintf("failed to read native module %q: %s", importPath, err))
	}
	defer fp.Close()

	result := eval(parser.NewReader(fp, importPath), env)
	if result.Type() == object.ErrType {
		return result
	}

	return env.Items()
}
