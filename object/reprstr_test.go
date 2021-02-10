package object

import "testing"

func TestReprStr(t *testing.T) {
	tests := []struct {
		obj      PanObject
		expected string
	}{
		{
			NewPanStr("hello"),
			`"hello"`,
		},
		{
			NewPanInt(2),
			`2`,
		},
		{
			EmptyPanObjPtr(),
			`{}`,
		},
		// if object has prop _name, show it instead
		{
			PanObjInstancePtr(&map[SymHash]Pair{
				(NewPanStr("_name")).SymHash(): {
					NewPanStr("_name"), NewPanStr("Taro"),
				},
			}),
			`Taro`,
		},
		// if _name is not str, ignore it
		{
			PanObjInstancePtr(&map[SymHash]Pair{
				(NewPanStr("_name")).SymHash(): {
					NewPanStr("_name"), NewPanInt(1),
				},
			}),
			`{"_name": 1}`,
		},
	}

	for _, tt := range tests {
		if ReprStr(tt.obj) != tt.expected {
			t.Errorf("wrong output: expected=%s, got=%s",
				tt.expected, ReprStr(tt.obj))
		}
	}
}
