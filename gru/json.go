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
