package gru

import (
	"fmt"
	"path/filepath"

	"github.com/Shopify/go-lua"
)

func NewPathModule() GruModule {
	module := NewModule("path", "System Path operations")
	module.Register("basename", "Returns the last portion of the path", pathBasename)
	module.Register("dirname", "Returns the directory name of the path", pathDirname)
	module.Register("extname", "Returns the file extension", pathExtname)
	module.Register("isAbsolute", "Checks if the path is absolute", pathIsAbsolute)
	module.Register("join", "Joins path elements", pathJoin)
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
		value, valid := l.ToString(i)
		if !valid {
			return LuaError(fmt.Sprintf("Expected string in argument %d", i))
		}
		parts[i-1] = value
	}

	return LuaStringResult(filepath.Join(parts...))
}
