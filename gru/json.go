package gru

import (
	"encoding/json"
	"os"
	"path/filepath"

	"github.com/Shopify/go-lua"
	"github.com/augustofrade/gru-lua/gru/definitions"
	"github.com/augustofrade/gru-lua/gru/internal/luautil"
)

func NewJsonModule() definitions.GruModule {
	module := definitions.NewModule("json", "JSON methods")
	module.FunctionBuilder("stringify", "Converts a table to JSON. Properties with 'nil' string values are converted to JSON null. If an invalid table is passed, a runtime error is thrown. Use 'pcall()' if needed.", jsonStringify).
		TableParam("table", "The table to be stringified").
		ReturnsString().
		Register()
	module.FunctionBuilder("parse", "Parses a JSON string into a table.", jsonParse).
		StringParam("json", "The JSON string to be parsed").
		Returns("table").
		Register()

	module.FunctionBuilder("dump", "Serializes a table to JSON and writes it to a file. Returns nil, error on failure.", jsonDump).
		StringParam("path", "Destination file path").
		TableParam("data", "The table to be serialized").
		ReturnsError().
		Register()
	return module
}

func jsonStringify(l *lua.State) int {
	if !l.IsTable(1) {
		return luautil.ErrorResult(l, "Expected table")
	}
	result := luautil.LuaTableToGo(l, 1)
	stringfiedValue, err := json.Marshal(result)
	if err != nil {
		return luautil.ErrorResult(l, err.Error())
	}

	return luautil.StringResult(l, string(stringfiedValue))
}

func jsonParse(l *lua.State) int {
	if !luautil.IsString(l, 1) {
		return luautil.ErrorResult(l, "Expected a JSON string")
	}
	str, _ := l.ToString(1)

	var result any
	err := json.Unmarshal([]byte(str), &result)
	if err != nil {
		return luautil.ErrorResult(l, err.Error())
	}

	return luautil.PushValue(l, result)
}

func jsonDump(l *lua.State) int {
	if !luautil.IsString(l, 1) {
		return luautil.SimpleErrorResult(l, "Expected string on 'path' parameter")
	}
	if !l.IsTable(2) {
		return luautil.SimpleErrorResult(l, "Expected table on 'data' parameter")
	}

	path, _ := l.ToString(1)
	result := luautil.LuaTableToGo(l, 2)

	jsonBytes, err := json.MarshalIndent(result, "", "  ")
	if err != nil {
		return luautil.SimpleErrorResult(l, err.Error())
	}

	path = filepath.Clean(path)
	dirname := filepath.Dir(path)
	err = os.MkdirAll(dirname, 0755)
	if err != nil {
		return luautil.ErrorResult(l, "Error while creating directory of provided path: "+err.Error())
	}

	if err := os.WriteFile(path, jsonBytes, 0644); err != nil {
		return luautil.ErrorResult(l, err.Error())
	}

	l.PushNil()
	return 1
}
