package luautil

import (
	"fmt"

	"github.com/Shopify/go-lua"
)

// Transforms the map into a Lua table and pushes it onto the stack
func PushTable(l *lua.State, kvp map[string]any) {
	l.CreateTable(0, len(kvp))
	for key, value := range kvp {
		switch v := value.(type) {
		case string:
			l.PushString(v)
		case int:
			l.PushInteger(v)
		case int64:
			l.PushInteger(int(v))
		case float32:
			l.PushNumber(float64(v))
		case float64:
			l.PushNumber(v)
		case bool:
			l.PushBoolean(v)
			// TODO: handle arrays/slices
		case map[string]any:
			PushTable(l, v)
		default:
			l.PushString(fmt.Sprint(v))
		}

		l.SetField(-2, key)
	}
}
