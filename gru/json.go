package gru

import (
	"encoding/json"

	"github.com/Shopify/go-lua"
	"github.com/augustofrade/gru-lua/gru/internal/luautil"
)

func NewJsonModule() GruModule {
	module := NewModule("json", "JSON methods")
	module.FunctionBuilder("stringify", "Converts a table to JSON. Properties with 'nil' string values are converted to JSON null. If an invalid table is passed, a runtime error is thrown. Use 'pcall()' if needed.", jsonStringify).
		TableParam("table", "The table to be stringified").
		ReturnsString().
		Register()
	module.FunctionBuilder("parse", "Parses a JSON string into a table.", jsonParse).
		StringParam("json", "The JSON string to be parsed").
		Returns("table").
		Register()
	return module
}

func jsonStringify(l *lua.State) int {
	if !l.IsTable(1) {
		return luautil.ErrorResult(l, "Expected table")
	}
	result := luautil.LuaTableToGo(l, 1)
	stringfiedValue, _ := json.Marshal(result)
	return luautil.StringResult(l, string(stringfiedValue))
}

func jsonParse(l *lua.State) int {
	if !luautil.IsString(l, 1) {
		return luautil.ErrorResult(l, "Expected a JSON string")
	}
	str, _ := l.ToString(1)

	var result any
	json.Unmarshal([]byte(str), &result)

	return luautil.PushValue(l, result)
}
