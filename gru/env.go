package gru

import (
	"os"

	"github.com/Shopify/go-lua"
	"github.com/augustofrade/gru-lua/gru/internal/luautil"
)

func NewEnvModule() GruModule {
	module := NewModule("env", "Environment operations.")
	module.FunctionBuilder("set", "Sets an environment variable.", envSet).
		StringParam("key", "Key of the environment variable.").
		StringParam("value", "value to be assigned on the key of the environment variable.").
		Register()
	module.FunctionBuilder("get", " retrieves the value of the environment variable named by the key. It returns the value, which will be empty if the variable is not present.", envGet).
		StringParam("key", "Key of the environment variable.").
		ReturnsString().
		Register()
	module.FunctionBuilder("clear", "Deletes all environment variables.", envClear).
		Register()
	module.FunctionBuilder("unset", "Unsets the environment variable", envUnset).
		StringParam("key", "Key to be unset").
		Register()
	module.FunctionBuilder("lookup", "Looks for the environment variable. Returns the value and true if the key exists, even if the value is empty.", envLookup).
		StringParam("key", "Key to be unset").
		ReturnsString().
		ReturnsBoolean().
		Register()
	module.FunctionBuilder("all", "Gets all environment variables as key as a table made of key=value strings", envAll).
		Returns("table").
		Register()

	return module
}

func envSet(l *lua.State) int {
	if !l.IsString(1) {
		return luautil.PushError(l, "Expected string on 'key' parameter")
	}
	if !l.IsString(2) {
		return luautil.PushError(l, "Expected string on 'value' parameter")
	}

	key, _ := l.ToString(1)
	value, _ := l.ToString(2)

	err := os.Setenv(key, value)
	if err != nil {
		return luautil.PushError(l, err.Error())
	}

	return 0
}

func envGet(l *lua.State) int {
	if !l.IsString(1) {
		return luautil.PushError(l, "Expected string on 'key' parameter")
	}

	key, _ := l.ToString(1)
	value := os.Getenv(key)
	return luautil.StringResult(l, value)
}

func envClear(l *lua.State) int {
	os.Clearenv()
	return 0
}

func envUnset(l *lua.State) int {
	if !l.IsString(1) {
		return luautil.PushError(l, "Expected string on 'key' parameter")
	}

	key, _ := l.ToString(1)
	err := os.Unsetenv(key)
	if err != nil {
		return luautil.PushError(l, err.Error())
	}

	return 0
}

func envLookup(l *lua.State) int {
	if !l.IsString(1) {
		return luautil.PushError(l, "Expected string on 'key' parameter")
	}

	key, _ := l.ToString(1)
	value, exists := os.LookupEnv(key)

	l.PushString(value)
	l.PushBoolean(exists)
	return 2
}

func envAll(l *lua.State) int {
	return luautil.StringTableArrayResult(l, os.Environ())
}
