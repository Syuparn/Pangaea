package modules

import (
	"testing"

	"github.com/Syuparn/pangaea/object"
)

func TestInjectTo(t *testing.T) {
	tests := []struct {
		path     string
		expected object.PanObject
	}{
		{
			"dummy",
			object.PanObjInstancePtr(&map[object.SymHash]object.Pair{
				object.GetSymHash("message"): {Key: object.NewPanStr("message"), Value: object.NewPanStr("This is a dummy module.")},
			}),
		},
	}

	for _, tt := range tests {
		env := object.NewEnv()
		InjectTo(env, Modules[tt.path]())

		if tt.expected.Inspect() != env.Items().Inspect() {
			t.Errorf("Expected %q, actual %q", tt.expected.Inspect(), env.Items().Inspect())
		}
	}
}
