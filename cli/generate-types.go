package cli

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/augustofrade/gru-lua/gru"
)

func generateTypesCommand(path *string) {
	fullPath, err := resolveTypesFilePath(path)
	if err != nil {
		fmt.Println(err)
		return
	}

	err = ensureTypesPathIsFileOrNotExists(fullPath)
	if err != nil {
		fmt.Println(err)
		return
	}

	generateLuaTypeAnnotations(fullPath)
}

func generateLuaTypeAnnotations(fullPath *string) {
	os.MkdirAll(filepath.Dir(*fullPath), 0755)

	gru.InitDefaultModules()
	fileContent := buildTypeAnnotations()

	err := os.WriteFile(*fullPath, fileContent, 0644)
	if err != nil {
		fmt.Println(err)
		return
	}
}

func buildTypeAnnotations() []byte {
	builder := strings.Builder{}
	builder.WriteString("-- Type annotations generated automatically with 'gru generate-types'\n")
	fmt.Fprintf(&builder, "-- For help about gru use 'gru help'\n---@meta\n\n")
	builder.WriteString("---@alias GruError string\n")

	moduleListBuilder := strings.Builder{}
	moduleListBuilder.WriteString("\n---@class Gru\n")

	for _, module := range gru.RegisteredModules {
		luaModuleName := "Gru" + strings.ToUpper(string(module.Name[0])) + module.Name[1:] + "Module"
		fmt.Fprintf(&moduleListBuilder, "---@field %s %s\n", module.Name, luaModuleName)

		fmt.Fprintf(&builder, "\n---@class %s %s", luaModuleName, module.Description)

		for _, function := range module.Functions {
			fmt.Fprintf(&builder, "\n---@field %s fun(", function.Name)

			paramCount := len(function.Parameters)
			for i, param := range function.Parameters {
				fmt.Fprintf(&builder, "%s: %s", param.Name, param.Type)
				if paramCount > 1 && i < paramCount-1 {
					builder.WriteString(", ")
				}
			}
			builder.WriteString(")")
			if len(function.ReturnTypes) > 0 {
				builder.WriteString(": " + strings.Join(function.ReturnTypes, ", "))
			}
			builder.WriteString(" " + function.Description)
		}
		builder.WriteString("\n")
		if len(module.Types) > 0 {
			builder.WriteString("\n-- Types of this Module\n")
			for _, t := range module.Types {
				fmt.Fprintf(&builder, "---@class %s %s\n", t.Name, t.Description)
				for pName, p := range t.Properties {
					fmt.Fprintf(&builder, "---@field %s %s %s\n", pName, p.Type, p.Description)
				}
			}
		}
	}

	moduleListBuilder.WriteString("\n---@type Gru\ngru = gru")

	builder.WriteString(moduleListBuilder.String())

	return []byte(builder.String())

}

func ensureTypesPathIsFileOrNotExists(path *string) error {
	fileInfo, err := os.Stat(*path)
	if err != nil && errors.Is(err, os.ErrNotExist) {
		return nil
	}

	if fileInfo.IsDir() {
		return fmt.Errorf("Path is a directory")
	}

	return nil
}

func resolveTypesFilePath(path *string) (*string, error) {
	if path == nil {
		p := "."
		path = &p
	}

	pathExt := filepath.Ext(*path)
	filename := "gru-types.lua"

	fullPath, err := filepath.Abs(*path)
	if err != nil {
		return nil, fmt.Errorf("Invalid path")
	}

	if pathExt == "" || pathExt == "." {
		fullPath = filepath.Join(fullPath, filename)
	} else if pathExt != ".lua" {
		fullPath = filepath.Join(filepath.Dir(fullPath), filename)
	}

	return &fullPath, nil
}
