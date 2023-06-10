package modules

import "github.com/Syuparn/pangaea/object"

func InjectTo(env *object.Env, pairs map[string]object.PanObject) {
	for k, v := range pairs {
		env.Set(object.GetSymHash(k), v)
	}
}
