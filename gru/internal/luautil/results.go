package luautil

import "github.com/Shopify/go-lua"

// function(): nil, error
func ErrorResult(l *lua.State, message string) int {
	l.PushNil()
	l.PushString(message)
	return 2
}

// function(): string
func StringResult(l *lua.State, value string) int {
	l.PushString(value)
	return 1
}

// function(): boolean
func BoolResult(l *lua.State, value bool) int {
	l.PushBoolean(value)
	return 1
}

// function(): table
func TableResult(l *lua.State, kvp map[string]any) int {
	PushTable(l, kvp)
	return 1
}
