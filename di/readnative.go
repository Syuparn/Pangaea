package di

import (
	"errors"
	"fmt"

	"github.com/rakyll/statik/fs"

	"github.com/Syuparn/pangaea/object"
	// neccessary to read embedded native source files
	_ "github.com/Syuparn/pangaea/statik"
)

func mustReadNativeCode(
	srcName string,
	env *object.Env,
) *map[object.SymHash]object.Pair {
	pairs, err := readNativeCode(srcName, env)
	if err != nil {
		panic(err.Error())
	}

	return pairs
}

func readNativeCode(
	srcName string,
	env *object.Env,
) (*map[object.SymHash]object.Pair, error) {
	fileName := fmt.Sprintf("/%s.pangaea", srcName)

	// NOTE: instead of native source files, open embedded statik file system
	// (in /statik directory)
	statikFS, err := fs.New()
	if err != nil {
		return nil, err
	}

	fp, err := statikFS.Open(fileName)
	if err != nil {
		e := fmt.Errorf("failed to read native %s props in native%s (zipped in statik/)",
			srcName, fileName)
		return nil, e
	}
	defer fp.Close()

	// NOTE: must pass EnclosedEnv otherwise outerenv of func literal cannot work
	// (cannot call top-level consts for example)
	result := eval(fp, object.NewEnclosedEnv(env))
	if result.Type() == object.ErrType {
		return nil, errors.New(result.Inspect())
	}

	obj, ok := result.(*object.PanObj)
	if !ok {
		e := fmt.Errorf("result must be ObjType. got=%s", result.Type())
		return nil, e
	}
	if obj.Pairs == nil {
		return nil, errors.New("Pairs must not be nil")
	}

	return obj.Pairs, nil
}
