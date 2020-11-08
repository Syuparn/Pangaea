package props

import (
	"github.com/Syuparn/pangaea/object"
)

// DiamondProps provides built-in props for <> (Diamond).
// NOTE: Some Diamond props are defind by native code (not by this function).
func DiamondProps(propContainer map[string]object.PanObject) map[string]object.PanObject {
	// NOTE: inject some built-in functions which relate to parser or evaluator
	return map[string]object.PanObject{
		"_iter": f(
			func(
				env *object.Env, kwargs *object.PanObj, args ...object.PanObject,
			) object.PanObject {
				if len(args) < 1 {
					return object.NewTypeErr("Diamond#_iter requires at least 1 arg")
				}

				ioObj, err := getIO(env)
				if err != nil {
					return err
				}

				return &object.PanBuiltInIter{
					Fn:  diamondIter(ioObj),
					Env: env, // not used
				}
			},
		),
		// <> can use all str props (call prop of read line)
		"_missing": f(
			func(
				env *object.Env, kwargs *object.PanObj, args ...object.PanObject,
			) object.PanObject {
				if len(args) < 2 {
					return object.NewTypeErr("_missing requires at least 2 args")
				}
				propName, ok := object.TraceProtoOfStr(args[1])
				if !ok {
					return object.NewTypeErr("\\2 must be prop name str")
				}

				line, err := readLineFromIO(env)
				if err != nil {
					return err
				}

				// call readLine.(prop)
				argsToPass := append(
					[]object.PanObject{object.EmptyPanObjPtr(), line, propName},
					args[2:]...,
				)
				return propContainer["Obj_callProp"].(*object.PanBuiltIn).Fn(
					env, kwargs,
					argsToPass...,
				)
			},
		),
		// S returns read line str
		"S": f(
			func(
				env *object.Env, kwargs *object.PanObj, args ...object.PanObject,
			) object.PanObject {
				if len(args) < 1 {
					return object.NewTypeErr("Diamond#S requires at least 1 arg")
				}

				line, err := readLineFromIO(env)
				if err != nil {
					return err
				}
				return line
			},
		),
	}
}

func readLineFromIO(env *object.Env) (*object.PanStr, *object.PanErr) {
	ioObj, err := getIO(env)
	if err != nil {
		return nil, err
	}
	line, ok := ioObj.ReadLine()
	if !ok {
		// return blank instead
		return object.NewPanStr(""), nil
	}

	return line, nil
}

func getIO(env *object.Env) (*object.PanIO, *object.PanErr) {
	// read line from stdin
	ioVar, ok := env.Get(object.GetSymHash("IO"))
	if !ok {
		return nil, object.NewNameErr("`IO` is not found in env")
	}
	ioObj, ok := object.TraceProtoOfIO(ioVar)
	if !ok {
		return nil, object.NewTypeErr("name `IO` must be IO obj")
	}
	return ioObj, nil
}

func diamondIter(ioObj *object.PanIO) object.BuiltInFunc {
	return func(
		env *object.Env, kwargs *object.PanObj, args ...object.PanObject,
	) object.PanObject {
		line, ok := ioObj.ReadLine()
		if !ok {
			return object.NewStopIterErr("iter stopped")
		}
		return line
	}
}
