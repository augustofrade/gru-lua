package cli

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"slices"
	"strings"
)

var gitIgnoreTypesFile string = "#generated with 'gru init'\ngru-types.lua"

func initCommand(path *string) {
	absPath, err := filepath.Abs(*path)
	if err != nil {
		fmt.Println(err)
		return
	}

	err = handlePathGitRepository(&absPath)
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println("Generating type annotations file...")
	typesPath := filepath.Join(absPath, "gru-types.lua")
	generateLuaTypeAnnotations(&typesPath)

	createSampleGruLuaFile(filepath.Join(absPath, "main.lua"))

	fmt.Println("Done.\n\nRun your lua code with 'gru main.lua'")
}

func handlePathGitRepository(path *string) error {
	pathInfo, err := os.Stat(*path)
	if err != nil && !errors.Is(err, os.ErrNotExist) {
		return err
	}

	if err == nil && !pathInfo.IsDir() {
		return fmt.Errorf("Path is a file, expected a directory")
	}

	// it is known that the path is/can be directory
	fmt.Printf("Handling %s...\n", *path)
	files, _ := os.ReadDir(*path)
	gitIgnorePath := filepath.Join(*path, ".gitignore")

	isGitRepo := false
	hasGitIgnore := false

	for _, file := range files {
		if file.IsDir() && file.Name() == ".git" {
			fmt.Println("Found git repository...")
			isGitRepo = true
		}
		if !file.IsDir() && file.Name() == ".gitignore" {
			hasGitIgnore = true
		}
	}

	if !isGitRepo {
		fmt.Println("Creating git repository...")
		err = exec.Command("git", "init", *path).Run()
		if err != nil {
			return err
		}
	}

	if hasGitIgnore {
		fmt.Println("Found .gitignore...")
		return handleExistentGitIgnoreFile(&gitIgnorePath)
	}

	fmt.Println("Creating new .gitignore file...")
	return createGitIgnoreFile(&gitIgnorePath)
}

func handleExistentGitIgnoreFile(gitIgnorePath *string) error {
	fileBytes, err := os.ReadFile(*gitIgnorePath)
	if err != nil {
		return err
	}

	fileContent := string(fileBytes)

	if slices.Contains(strings.Split(fileContent, "\n"), "gru-types.lua") {
		fmt.Println("Found type annotations file name in .gitignore...")
		return nil
	}

	fileContent = fileContent + "\n" + gitIgnoreTypesFile
	err = os.WriteFile(*gitIgnorePath, []byte(fileContent), 0644)
	return err
}

func createGitIgnoreFile(gitIgnorePath *string) error {
	err := os.WriteFile(*gitIgnorePath, []byte(gitIgnoreTypesFile), 0644)
	return err
}

func createSampleGruLuaFile(path string) error {
	_, err := os.Stat(path)
	if err == nil {
		// file already exists
		return nil
	}

	if !errors.Is(err, os.ErrNotExist) {
		return err
	}

	os.WriteFile(path, []byte("print(gru.colors.lightBlue('Hello, world!'))"), 0644)
	return nil
}
