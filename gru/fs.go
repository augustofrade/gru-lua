package gru

import (
	"os"
	"path/filepath"

	"github.com/Shopify/go-lua"
	"github.com/augustofrade/gru-lua/gru/internal/luautil"
)

type GruDirentry struct {
	Name       string
	IsDir      bool
	ParentPath string
}

func NewFsModule() GruModule {
	module := NewModule("fs", "File System operations.")
	module.HasCustomType("GruDirentry", "").
		StringProp("name", "Name of the directory entry.").
		BooleanProp("is_dir", "Whether the entry is a directory").
		StringProp("parent_path", "Absolute parent path of the directory entry")

	module.FunctionBuilder("read_dir", "Reads the passed dir and returns its contents.", fsReadDir).
		StringParam("dir", "Dir to be read.").
		ReturnsWithError("GruDirentry[]").
		Register()

	return module
}

func fsReadDir(l *lua.State) int {
	if !luautil.IsString(l, 1) {
		return luautil.ErrorResult(l, "Expected string type on 'dir' parameter")
	}
	dir, _ := l.ToString(1)
	contents, err := os.ReadDir(dir)
	if err != nil {
		return luautil.ErrorResult(l, err.Error())
	}

	dirFullPath := dir
	if !filepath.IsAbs(dir) {
		dirFullPath, err = filepath.Abs(dir)
		if err != nil {
			return luautil.ErrorResult(l, err.Error())
		}
	}

	mapped := make([]GruDirentry, 0)
	for _, v := range contents {
		mapped = append(mapped, GruDirentry{
			Name:       v.Name(),
			IsDir:      v.IsDir(),
			ParentPath: dirFullPath,
		})
	}

	arr := make([]any, 0, len(mapped))
	for _, d := range mapped {
		arr = append(arr, map[string]any{
			"name":        d.Name,
			"is_dir":      d.IsDir,
			"parent_path": d.ParentPath,
		})
	}

	luautil.PushArrayTable(l, arr)
	return 1
}
