package di

import (
	"errors"
	"fmt"

	"github.com/Syuparn/pangaea/native"
	"github.com/Syuparn/pangaea/object"
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
	fileName := fmt.Sprintf("%s.pangaea", srcName)

	fp, err := native.FS.Open(fileName)
	if err != nil {
		e := fmt.Errorf("failed to read native %s props in native/%s",
			srcName, fileName)
		return nil, e
	}
	defer fp.Close()

	// NOTE: must pass EnclosedEnv otherwise outerenv of func literal cannot work
	// (cannot call top-level consts for example)
	result := eval(fp, object.NewEnclosedEnv(env))
	if err, ok := result.(*object.PanErr); ok {
		return nil, errors.New(err.Inspect() + "\n" + err.StackTrace)
	}

	obj, ok := result.(*object.PanObj)
	if !ok {
		e := fmt.Errorf("result must be ObjType. got=%s", result.Type())
		return nil, e
	}
	if obj.Pairs == nil {
		return nil, errors.New("pairs must not be nil")
	}

	return obj.Pairs, nil
}
