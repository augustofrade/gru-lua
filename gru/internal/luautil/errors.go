package luautil

import (
	"github.com/Shopify/go-lua"
)

// Pushes to the stack and throws an error
func PushError(l *lua.State, message string) int {
	l.PushString(message)
	l.Error()
	return 1
}
