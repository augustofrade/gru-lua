package gru

import (
	"errors"
	"os"
	"path/filepath"

	"github.com/Shopify/go-lua"
	"github.com/augustofrade/gru-lua/gru/internal/luautil"
)

type GruDirentry struct {
	Name       string `lua:"name"`
	IsDir      bool   `lua:"is_dir"`
	ParentPath string `lua:"parent_path"`
}

func NewFsModule() GruModule {
	module := NewModule("fs", "File System operations.")
	module.HasCustomType("GruDirentry", "").
		StringProp("name", "Name of the directory entry.").
		BooleanProp("is_dir", "Whether the entry is a directory").
		StringProp("parent_path", "Absolute parent path of the directory entry")

	module.HasCustomType("GruFileInfo", "Information about a file or directory").
		StringProp("name", "Name of the file or directory.").
		StringProp("fullpath", "Fullpath of the file or directory").
		BooleanProp("is_dir", "Whether the entry is a directory").
		NumberProp("last_modification_time", "UNIX date of last modification of the file or directory").
		NumberProp("size", "Size of the file or directory")

	module.FunctionBuilder("read_dir", "Reads the provided directory path and returns its contents.", fsReadDir).
		StringParam("dir", "Dir to be read.").
		ReturnsWithError("GruDirentry[]").
		Register()

	module.FunctionBuilder("read_file", "Reads the file in the provided path and returns its contents as a string.", fsReadFile).
		StringParam("file_path", "File to be read.").
		ReturnsStringWithError().
		Register()

	module.FunctionBuilder("write_file", "Writes data to the provided file path, creating it if necessary with permissions 'permissions'.", fsWriteFile).
		StringParam("file_path", "Path of the file.").
		StringParam("data", "Data to be written.").
		OptionalNumberParam("permissions", "Permission of the file if it doesn't exist.").
		ReturnsError().
		Register()
	module.FunctionBuilder("create", "Creates or truncates the named file. If the file already exists, it is truncated. Creates all neccessary directories of the file if needed.", fsCreate).
		StringParam("file_path", "Path of the file.").
		ReturnsError().
		Register()

	module.FunctionBuilder("mkdir",
		"Creates a directory along with any necessary parents, and returns nil, or else returns an error. The permissions value is used for all directories created. If path is already a directory, nothing happens and nil is returned.",
		fsMkDir).
		StringParam("path", "Path of the directory to be created").
		OptionalNumberParam("permissions", "Permission of the directory and its parents if they don't exist.").
		ReturnsError().
		Register()

	module.FunctionBuilder("remove", "Removes the named file or directory. If you need to remove a directory with children, consider using 'fs.remove_all' function", fsRemove).
		StringParam("path", "Path of the file or directory.").
		ReturnsError().
		Register()

	module.FunctionBuilder("remove_all", "Removes path and any children it contains. If path doesn't exist, nothing happens and nil is returned.", fsRemoveAll).
		StringParam("path", "Path of the directory.").
		ReturnsError().
		Register()
	module.FunctionBuilder("exists", "Returns whether the file or directory at path exists", fsExists).
		StringParam("path", "Path of the file or directory.").
		ReturnsBooleanWithError().
		Register()

	module.FunctionBuilder("stat", "Returns information about the file or directory at path", fsStat).
		StringParam("path", "Path of the file or directory.").
		ReturnsWithError("GruFileInfo").
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

	luautil.PushArrayTable(l, mapped)
	return 1
}

func fsReadFile(l *lua.State) int {
	if !luautil.IsString(l, 1) {
		return luautil.ErrorResult(l, "Expected string type on 'file_path' parameter")
	}

	fPath, _ := l.ToString(1)
	fileContent, err := os.ReadFile(fPath)
	if err != nil {
		return luautil.ErrorResult(l, err.Error())
	}

	return luautil.StringResult(l, string(fileContent))
}

func fsWriteFile(l *lua.State) int {
	if !luautil.IsString(l, 1) {
		return luautil.SimpleErrorResult(l, "Expected string on 'file_path' parameter")
	}

	if !luautil.IsString(l, 2) {
		return luautil.SimpleErrorResult(l, "Expected string on 'data' parameter")
	}

	permissions := 0644
	if l.Top() >= 3 {
		// 3 parameters provided
		if !l.IsNumber(3) {
			return luautil.SimpleErrorResult(l, "Expected number or nil on optional 'permissions' parameter")
		}
		fm, _ := l.ToInteger(3)
		permissions = fm
	}

	fPath, _ := l.ToString(1)
	fileData, _ := l.ToString(2)

	err := os.WriteFile(fPath, []byte(fileData), os.FileMode(permissions))

	if err != nil {
		return luautil.SimpleErrorResult(l, err.Error())
	}

	return 0
}

func fsCreate(l *lua.State) int {
	if !luautil.IsString(l, 1) {
		return luautil.SimpleErrorResult(l, "Expected string type on 'path' parameter")
	}

	fPath, _ := l.ToString(1)
	if fPath == "" {
		return luautil.SimpleErrorResult(l, "Expected a valid file path on 'path' parameter")
	}
	dirname := filepath.Dir(fPath)
	err := os.MkdirAll(dirname, 0755)
	if err != nil {
		return luautil.SimpleErrorResult(l, "Error while creating directory of provided path: "+err.Error())
	}

	file, err := os.Create(fPath)
	if err != nil {
		return luautil.SimpleErrorResult(l, err.Error())
	}

	// TODO: return table with file operation such as write, append, read, close
	file.Close()

	return 0
}

func fsMkDir(l *lua.State) int {
	if !luautil.IsString(l, 1) {
		return luautil.SimpleErrorResult(l, "Expected string on 'path' parameter")
	}
	dirPath, _ := l.ToString(1)

	permissions := 0755
	if l.Top() >= 2 {
		// 2 parameters provided
		if !l.IsNumber(2) {
			return luautil.SimpleErrorResult(l, "Expected number or nil on optional 'permissions' parameter")
		}
		fm, _ := l.ToInteger(2)
		permissions = fm
	}

	err := os.MkdirAll(filepath.Clean(dirPath), os.FileMode(permissions))
	if err != nil {
		return luautil.SimpleErrorResult(l, err.Error())
	}

	return 0
}

func fsRemove(l *lua.State) int {
	if !luautil.IsString(l, 1) {
		return luautil.SimpleErrorResult(l, "Expected string on 'path' parameter")
	}
	path, _ := l.ToString(1)
	err := os.Remove(filepath.Clean(path))
	if err != nil {
		return luautil.SimpleErrorResult(l, err.Error())
	}
	return 0
}

func fsRemoveAll(l *lua.State) int {
	if !luautil.IsString(l, 1) {
		return luautil.SimpleErrorResult(l, "Expected string on 'path' parameter")
	}
	path, _ := l.ToString(1)
	err := os.RemoveAll(filepath.Clean(path))
	if err != nil {
		return luautil.SimpleErrorResult(l, err.Error())
	}
	return 0
}

func fsExists(l *lua.State) int {
	if !luautil.IsString(l, 1) {
		return luautil.SimpleErrorResult(l, "Expected string on 'path' parameter")
	}
	path, _ := l.ToString(1)
	_, err := os.Stat(filepath.Clean(path))
	if err == nil {
		return luautil.BoolResult(l, true)
	}

	if errors.Is(err, os.ErrNotExist) {
		return luautil.BoolResult(l, false)
	}
	return luautil.ErrorResult(l, err.Error())

}

func fsStat(l *lua.State) int {
	if !luautil.IsString(l, 1) {
		return luautil.SimpleErrorResult(l, "Expected string on 'path' parameter")
	}
	path, _ := l.ToString(1)
	path = filepath.Clean(path)
	if !filepath.IsAbs(path) {
		fullPath, err := filepath.Abs(path)
		if err != nil {
			return luautil.SimpleErrorResult(l, err.Error())
		}
		path = fullPath
	}
	stat, err := os.Stat(path)
	if err != nil {
		return luautil.SimpleErrorResult(l, err.Error())
	}

	fileInfo := make(map[string]any)

	fileInfo["name"] = stat.Name()
	fileInfo["fullpath"] = path
	fileInfo["is_dir"] = stat.IsDir()
	fileInfo["last_modification_time"] = stat.ModTime().Unix()
	fileInfo["size"] = stat.Size()

	return luautil.TableResult(l, fileInfo)
}
