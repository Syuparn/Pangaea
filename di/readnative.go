package di

import (
	"errors"
	"fmt"
	"os"

	"github.com/Syuparn/pangaea/object"
)

func mustReadNativeCode(srcName string, env *object.Env) map[string]object.PanObject {
	propContainer, err := readNativeCode(srcName, env)
	if err != nil {
		panic(err.Error())
	}

	return propContainer
}

func readNativeCode(srcName string, env *object.Env) (map[string]object.PanObject, error) {
	srcPath := "native"
	fileName := fmt.Sprintf("%s/%s.pangaea", srcPath, srcName)

	fp, err := os.Open(fileName)
	if err != nil {
		e := fmt.Errorf("failed to read native %s props in %s",
			srcName, fileName)
		return map[string]object.PanObject{}, e
	}
	defer fp.Close()

	// NOTE: must pass EnclosedEnv otherwise outerenv of func literal cannot work
	// (cannot call top-level consts for example)
	result := eval(fp, object.NewEnclosedEnv(env))
	if result.Type() == object.ErrType {
		return map[string]object.PanObject{}, errors.New(result.Inspect())
	}

	obj, ok := result.(*object.PanObj)
	if !ok {
		e := fmt.Errorf("result must be ObjType. got=%s", result.Type())
		return map[string]object.PanObject{}, e
	}
	if obj.Pairs == nil {
		return map[string]object.PanObject{}, errors.New("Pairs must not be nil")
	}

	propContainer := map[string]object.PanObject{}
	for _, v := range *obj.Pairs {
		keyStr, ok := v.Key.(*object.PanStr)
		if !ok {
			e := fmt.Errorf("key must be StrType. got=%s", v.Key.Inspect())
			return map[string]object.PanObject{}, e
		}

		propContainer[keyStr.Value] = v.Value
	}

	return propContainer, nil
}
