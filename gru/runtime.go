package gru

import (
	"github.com/Shopify/go-lua"
)

func NewRuntimeModule() GruModule {
	module := NewModule("runtime", "Provides direct access to the runtime")

	module.FunctionBuilder("version", "Current version of the Gru runtime", runtimeVersion).
		ReturnsString().
		Register()

	module.FunctionBuilder("on_error", "Set function as callback to be triggered on runtime errors", runtimeOnError).
		Register()

	return module
}

func runtimeVersion(l *lua.State) int {
	l.PushString(runtimeCurrentVersion)
	return 1
}

func runtimeOnError(l *lua.State) int {
	if !l.IsFunction(1) {
		return LuaError("Expected a function")
	}

	l.PushValue(1)
	l.SetField(lua.RegistryIndex, errorCallbackRegistryKey)
	return 0
}
