package gru

import (
	"encoding/json"

	"github.com/Shopify/go-lua"
)

func NewJsonModule() GruModule {
	module := NewModule("json", "JSON methods")
	module.Register("stringify", "Converts a table into a string", jsonStringify)
	return module
}

func jsonStringify(l *lua.State) int {
	if !l.IsTable(1) {
		l.PushNil()
		l.PushString("Expected table")
		return 2
	}

	result := map[string]any{}

	l.PushNil() // first key

	// TODO: handle array tables
	for l.Next(1) {
		// -2 = key, -1 = value
		key, _ := l.ToString(-2)

		var val any

		switch l.TypeOf(-1) {
		case lua.TypeString:
			val, _ = l.ToString(-1)
		case lua.TypeNumber:
			val, _ = l.ToNumber(-1)
		case lua.TypeBoolean:
			val = l.ToBoolean(-1)
		case lua.TypeNil:
			val = nil
		}

		if val == "nil" {
			val = nil
		}

		// TODO: handle table (object/array) props

		result[key] = val
		l.Pop(1) // remove valor, mantém chave para Next()
	}

	out, _ := json.Marshal(result)
	l.PushString(string(out))
	return 1
}
