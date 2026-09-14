package gru

import (
	"archive/zip"
	"fmt"
	"io"
	"os"

	"github.com/Shopify/go-lua"
	"github.com/augustofrade/gru-lua/gru/definitions"
	"github.com/augustofrade/gru-lua/gru/internal/luautil"
)

type sysEnt struct {
	basename string
	fullPath string
	path     string
}

func NewZipModule() definitions.GruModule {
	module := definitions.NewModule("zip", "ZIP compression operations")
	module.FunctionBuilder("files", "ZIPs files and directories into a single compressed file specified by the 'files' table", zipFiles).
		StringParam("files", "Files to be compressed").
		StringParam("output", "Output path of the compressed ZIP").
		Register()
	return module
}

func zipFiles(l *lua.State) int {
	if !l.IsTable(1) {
		return luautil.PushError(l, "Expected 'files' arguments to be a table")
	}

	if !luautil.IsString(l, 2) {
		return luautil.PushError(l, "Expected 'output' argument to be a string")
	}

	fileNames, err := luautil.GetTableStringValues(l)
	if err != nil {
		return luautil.PushError(l, "Expected table made only from path strings")
	}

	output, _ := l.ToString(2)

	// buf := new(bytes.Buffer)
	// w := zip.NewWriter(buf)
	// files := make([]sysEnt, 0)
	fmt.Println(fileNames)

	zipFile, _ := os.Create(output)
	zipWriter := zip.NewWriter(zipFile)
	defer zipFile.Close()
	defer zipWriter.Close()
	for _, filename := range *fileNames {
		addFile(zipWriter, filename)
	}

	// for _, file := range *fileNames {
	// 	fileStat, err := os.Stat(file)
	// 	if err != err {
	// 		if errors.Is(err, os.ErrNotExist) {
	// 			return luautil.PushError(l, fmt.Sprintf("File '%s' not found", file))
	// 		}
	// 		return luautil.PushError(l, err.Error())
	// 	}

	// 	fullPath, _ := filepath.Abs(file)
	// 	files = append(files, sysEnt{
	// 		basename: filepath.Base(file),
	// 		fullPath: fullPath,
	// 	})
	// 	fmt.Println(fileStat.Name())
	// }

	return 0
}

func addFile(zipWriter *zip.Writer, filename string) {
	fileToZip, err := os.Open(filename)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer fileToZip.Close()

	w, err := zipWriter.Create(filename)
	if err != nil {
		fmt.Println(err)
		return
	}

	if _, err := io.Copy(w, fileToZip); err != nil {
		fmt.Println(err)
	}
}
