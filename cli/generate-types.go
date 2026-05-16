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

	moduleListBuilder := strings.Builder{}
	moduleListBuilder.WriteString("\n---@class Gru\n")

	builder.WriteString("-- Type annotations generated automatically with 'gru generate-types'\n")
	builder.WriteString("-- For help about gru use 'gru help'\n")
	builder.WriteString("---@meta\n\n")
	builder.WriteString("---@alias GruError string\n")

	for _, module := range gru.RegisteredModules {
		luaModuleName := "Gru" + strings.ToUpper(string(module.Name[0])) + module.Name[1:] + "Module"
		moduleListBuilder.WriteString("---@field ")
		moduleListBuilder.WriteString(module.Name)
		moduleListBuilder.WriteString(" ")
		moduleListBuilder.WriteString(luaModuleName)
		moduleListBuilder.WriteString("\n")

		builder.WriteString("\n---@class ")
		builder.WriteString(luaModuleName)
		builder.WriteString(" " + module.Description)

		for _, function := range module.Functions {
			builder.WriteString("\n---@field ")
			builder.WriteString(function.Name)
			builder.WriteString(" fun(")
			for _, param := range function.Parameters {
				builder.WriteString(param.Name + ": " + param.Type)
			}
			builder.WriteString(")")
			if len(function.ReturnTypes) > 0 {
				builder.WriteString(": " + strings.Join(function.ReturnTypes, ", "))
			}
			builder.WriteString(" " + function.Description)
		}
		builder.WriteString("\n")
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

	if pathExt == "" {
		fullPath = filepath.Join(fullPath, filename)
	} else if pathExt != ".lua" {
		fullPath = filepath.Join(filepath.Dir(fullPath), filename)
	}

	return &fullPath, nil
}
