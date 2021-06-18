package props

import (
	"encoding/json"
	"fmt"
	"math"

	"github.com/Syuparn/pangaea/object"
)

// JSONProps provides built-in props for JSON.
// NOTE: Some props are defind by native code (not by this function).
func JSONProps(propContainer map[string]object.PanObject) map[string]object.PanObject {
	// NOTE: inject some built-in functions which relate to parser or evaluator
	return map[string]object.PanObject{
		"dec": f(
			func(
				env *object.Env, kwargs *object.PanObj, args ...object.PanObject,
			) object.PanObject {
				if len(args) < 2 {
					return object.NewTypeErr("JSON.dec requires at least 2 args")
				}

				str, ok := object.TraceProtoOfStr(args[1])
				if !ok {
					return object.NewTypeErr(
						fmt.Sprintf("%s cannot be treated as str", args[1].Repr()))
				}

				return decodeJSON(str.Value)
			},
		),
	}
}

func decodeJSON(s string) object.PanObject {
	var elems interface{}

	if err := json.Unmarshal([]byte(s), &elems); err != nil {
		return object.NewValueErr(
			fmt.Sprintf("failed to decode JSON: %v (input `%s`)", err, s))
	}

	return parseJSONElems(elems)
}

func parseJSONElems(elems interface{}) object.PanObject {
	switch e := elems.(type) {
	case string:
		return object.NewPanStr(e)
	case float64:
		return parseJSONNum(e)
	case []interface{}:
		return parseJSONArr(e)
	case map[string]interface{}:
		return parseJSONMap(e)
	case nil:
		return object.BuiltInNil
	default:
		return object.NewValueErr(
			fmt.Sprintf("failed to parse JSON elements unexpectedly: %#v", elems))
	}
}

func parseJSONArr(elems []interface{}) object.PanObject {
	panElems := make([]object.PanObject, len(elems))

	for i, e := range elems {
		elem := parseJSONElems(e)
		// error handling
		if elem.Type() == object.ErrType {
			return elem
		}

		panElems[i] = elem
	}

	return object.NewPanArr(panElems...)
}

func parseJSONMap(elems map[string]interface{}) object.PanObject {
	pairs := map[object.SymHash]object.Pair{}

	for k, v := range elems {
		value := parseJSONElems(v)
		// error handling
		if value.Type() == object.ErrType {
			return value
		}

		pairs[object.GetSymHash(k)] = object.Pair{
			Key:   object.NewPanStr(k),
			Value: value,
		}
	}

	return object.PanObjInstancePtr(&pairs)
}

func parseJSONNum(f float64) object.PanObject {
	// check if f is integer
	if math.Floor(f) == f {
		return object.NewPanInt(int64(f))
	}

	return object.NewPanFloat(float64(f))
}
