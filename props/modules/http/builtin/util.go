package builtin

import "github.com/Syuparn/pangaea/object"

func mapToObj(kwargMap map[string]object.PanObject) *object.PanObj {
	p := map[object.SymHash]object.Pair{}
	for k, v := range kwargMap {
		p[object.GetSymHash(k)] = object.Pair{Key: object.NewPanStr(k), Value: v}
	}

	return object.PanObjInstancePtr(&p).(*object.PanObj)
}
