package gru

import (
	"fmt"
	"path"
	"path/filepath"

	"github.com/Shopify/go-lua"
)

func NewPathModule() GruModule {
	module := NewModule("path", "System Path operations")
	module.FunctionBuilder("basename", "Returns the last portion of the path", pathBasename).
		StringParam("path", "").
		ReturnsString().
		Register()
	module.FunctionBuilder("dirname", "Returns the directory name of the path", pathDirname).
		StringParam("path", "").
		ReturnsString().
		Register()
	module.FunctionBuilder("extname", "Returns the file extension", pathExtname).
		StringParam("path", "").
		ReturnsString().
		Register()
	module.FunctionBuilder("isAbsolute", "Checks if the path is absolute", pathIsAbsolute).
		StringParam("path", "").
		ReturnsBoolean().
		Register()
	module.FunctionBuilder("join", "Joins path elements", pathJoin).
		Vararg("string|number").
		ReturnsString().
		Register()
	module.FunctionBuilder("split",
		"Splits a path separating it into a directory and file name. If there is no slash in path, split returns an empty dir and file set to path.",
		pathSplit).
		StringParam("path", "path to be split").
		Returns("table").
		Register()
	return module
}

func pathBasename(l *lua.State) int {
	value, valid := l.ToString(1)
	if !valid {
		return LuaError("Expected string")
	}

	return LuaStringResult(filepath.Base(value))
}

func pathDirname(l *lua.State) int {
	value, valid := l.ToString(1)
	if !valid {
		return LuaError("Expected string")
	}

	return LuaStringResult(filepath.Dir(value))
}

func pathExtname(l *lua.State) int {
	value, valid := l.ToString(1)
	if !valid {
		return LuaError("Expected string")
	}

	return LuaStringResult(filepath.Ext(value))
}

func pathIsAbsolute(l *lua.State) int {
	value, valid := l.ToString(1)
	if !valid {
		return LuaError("Expected string")
	}

	return LuaBoolResult(filepath.IsAbs(value))
}

func pathJoin(l *lua.State) int {
	count := l.Top()
	parts := make([]string, count)

	for i := 1; i <= count; i++ {
		if !l.IsString(i) {
			return LuaError(fmt.Sprintf("Expected string or number in argument %d", i))
		}
		// Numbers in lua are always convertible to string
		value, _ := l.ToString(i)
		parts[i-1] = value
	}

	return LuaStringResult(filepath.Join(parts...))
}

func pathSplit(l *lua.State) int {
	if !l.IsString(1) || l.IsNumber(1) {
		return LuaError("Expected string")
	}

	value, _ := l.ToString(1)

	dir, file := path.Split(value)

	PushLuaTable(map[string]any{
		"dir":  dir,
		"file": file,
	})

	return 1
}
