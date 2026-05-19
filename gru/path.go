package gru

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/Shopify/go-lua"
)

func NewPathModule() GruModule {
	module := NewModule("path", "System Path operations")
	module.FunctionBuilder("basename", "Returns the last portion of the path.", pathBasename).
		StringParam("path", "").
		ReturnsStringWithError().
		Register()
	module.FunctionBuilder("dirname", "Returns the directory name of the path.", pathDirname).
		StringParam("path", "").
		ReturnsStringWithError().
		Register()
	module.FunctionBuilder("extname", "Returns the file extension.", pathExtname).
		StringParam("path", "").
		ReturnsStringWithError().
		Register()
	module.FunctionBuilder("is_absolute", "Checks if the path is absolute.", pathIsAbsolute).
		StringParam("path", "").
		ReturnsBooleanWithError().
		Register()
	module.FunctionBuilder("is_relative", "Checks if the path is relative. It is the opposite of the 'is_absolute' function.", pathIsRelative).
		StringParam("path", "").
		ReturnsBooleanWithError().
		Register()
	module.FunctionBuilder("join", "Joins path elements.", pathJoin).
		Vararg("string").
		ReturnsStringWithError().
		Register()
	module.FunctionBuilder("parse",
		"Parses a path by separating it into a directory, file and extension names. If there is no slash in path, split returns an empty dir and file set to path.",
		pathParse).
		StringParam("path", "Path to be parsed.").
		ReturnsWithError("table").
		Register()
	module.FunctionBuilder("absolute",
		"Returns an absolute representation of a path. If the path is not absolute it will be joined with the current working directory to turn it into an absolute path.",
		pathAbsolute).
		StringParam("path", "").
		ReturnsStringWithError().
		Register()
	module.FunctionBuilder("clean",
		"Normalizes a path by removing redundant separators and resolving . and .. segments when possible.",
		pathClean).
		StringParam("path", "Path to be normalized.").
		ReturnsStringWithError().
		Register()
	module.FunctionBuilder("stem",
		"Returns the file name without the extension",
		pathStem).
		StringParam("path", "Path used to retrieve the file name without extension.").
		ReturnsStringWithError().
		Register()
	module.FunctionBuilder("relative",
		"Returns a relative path from 'base_path' to 'target_path'.",
		pathRelative).
		StringParam("base_path", "Starting path.").
		StringParam("target_path", "Destination path relative to 'base_path'.").
		ReturnsStringWithError().
		Register()
	module.FunctionBuilder("resolve",
		"Resolves all paths passed by joining and cleaning them and returning the resulting absolute path", pathResolve).
		Vararg("string").
		ReturnsStringWithError().
		Register()
	return module
}

func pathBasename(l *lua.State) int {
	if !l.IsString(1) || l.IsNumber(1) {
		return LuaError("Expected string")
	}

	value, _ := l.ToString(1)

	return LuaStringResult(filepath.Base(value))
}

func pathDirname(l *lua.State) int {
	if !l.IsString(1) || l.IsNumber(1) {
		return LuaError("Expected string")
	}

	value, _ := l.ToString(1)

	return LuaStringResult(filepath.Dir(value))
}

func pathExtname(l *lua.State) int {
	if !l.IsString(1) || l.IsNumber(1) {
		return LuaError("Expected string")
	}

	value, _ := l.ToString(1)

	return LuaStringResult(filepath.Ext(value))
}

func pathIsAbsolute(l *lua.State) int {
	if !l.IsString(1) || l.IsNumber(1) {
		return LuaError("Expected string")
	}

	value, _ := l.ToString(1)

	return LuaBoolResult(filepath.IsAbs(value))
}

func pathIsRelative(l *lua.State) int {
	if !l.IsString(1) || l.IsNumber(1) {
		return LuaError("Expected string")
	}

	value, _ := l.ToString(1)

	return LuaBoolResult(!filepath.IsAbs(value))
}

func pathJoin(l *lua.State) int {
	count := l.Top()
	if count == 0 {
		return LuaError("Expected at least 1 argument")
	}

	parts, err := GetLuaStringVarargs(l, count)
	if err != nil {
		return LuaError(err.Error())
	}

	return LuaStringResult(filepath.Join(parts...))
}

func pathParse(l *lua.State) int {
	if !l.IsString(1) || l.IsNumber(1) {
		return LuaError("Expected string")
	}

	value, _ := l.ToString(1)

	dir, file := filepath.Split(value)

	PushLuaTable(map[string]any{
		"dir":  dir,
		"file": file,
		"ext":  filepath.Ext(file),
	})

	return 1
}

func pathAbsolute(l *lua.State) int {
	if !l.IsString(1) || l.IsNumber(1) {
		return LuaError("Expected string")
	}

	value, _ := l.ToString(1)

	absPath, err := filepath.Abs(value)
	if err != nil {
		return LuaError(err.Error())
	}

	return LuaStringResult(absPath)
}

func pathClean(l *lua.State) int {
	if !l.IsString(1) || l.IsNumber(1) {
		return LuaError("Expected string")
	}

	value, _ := l.ToString(1)

	cleanPath := filepath.Clean(value)

	return LuaStringResult(cleanPath)
}

func pathStem(l *lua.State) int {
	if !l.IsString(1) || l.IsNumber(1) {
		return LuaError("Expected string")
	}

	value, _ := l.ToString(1)
	value = filepath.Base(value)

	if value == "." || value == ".." {
		return LuaStringResult(value)
	}

	isHidden := false
	if len(value) > 0 && string(value[0]) == "." {
		isHidden = true
		value = value[1:]
	}

	value = strings.TrimSuffix(value, filepath.Ext(value))
	if isHidden {
		value = "." + value
	}

	return LuaStringResult(value)
}

func pathRelative(l *lua.State) int {
	if !l.IsString(1) || l.IsNumber(1) {
		return LuaError("Expected string in 'base_path' argument")
	}

	if !l.IsString(2) || l.IsNumber(2) {
		return LuaError("Expected string in 'target_path' argument")
	}

	base, _ := l.ToString(1)
	target, _ := l.ToString(2)

	relativePath, err := filepath.Rel(base, target)
	if err != nil {
		return LuaError(err.Error())
	}

	return LuaStringResult(relativePath)
}

func pathResolve(l *lua.State) int {
	count := l.Top()
	parts, err := GetLuaStringVarargs(l, count)
	if err != nil {
		return LuaError(err.Error())
	}

	fmt.Println(filepath.IsAbs(filepath.Join(parts...)))
	resolvedPath, err := filepath.Abs(filepath.Join(parts...))
	if err != nil {
		return LuaError(err.Error())
	}

	return LuaStringResult(resolvedPath)
}
