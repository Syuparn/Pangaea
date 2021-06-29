package main

import "syscall/js"

// RegisterPangaea registers Pangaea modules to global namespace `pangaea.xxx` in js
func RegisterPangaea() {
	ex := NewExecutor()

	js.Global().Set("pangaea", js.ValueOf(
		map[string]interface{}{
			"execute": js.FuncOf(ex.Execute),
		},
	))
}
