package cli

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/augustofrade/gru-lua/gru"
)

func HandleArgs(args []string) {
	if len(args) <= 1 {
		helpCommand()
		return
	}

	command := args[1]
	subCommand := ""
	if len(args) > 2 {
		subCommand = args[2]
	}

	// TODO: handle params
	switch command {
	case "help":
		helpCommand()
	case "modules":
		modulesCommand(5, subCommand)
	case "eval":
		gru.Evalute(&subCommand)
	default:
		handleFile(&command)
	}
}

func handleFile(file *string) {
	absFile, _ := filepath.Abs(*file)
	if _, err := os.Stat(absFile); err != nil {
		fmt.Println("File not found: " + absFile)
		return
	}

	gru.DoFile(file)
}

func modulesCommand(maxModuleSize int, specifiedModule string) {
	gru.InitDefaultModules()

	if specifiedModule == "" {
		printAllModulesInfo(maxModuleSize)
		return
	}

	var module *gru.GruModule
	for _, m := range gru.RegisteredModules {
		if m.Name == specifiedModule {
			module = &m
		}
	}

	if module == nil {
		fmt.Println("Gru module not found")
		return
	}

	fmt.Printf("[%d function(s)]  %s\n", len(module.Functions), module.Description)
	printModuleInfo(*module, 0)
}

func printAllModulesInfo(maxModuleSize int) {
	for _, module := range gru.RegisteredModules {
		fmt.Printf("%s [%d]  %s\n", module.Name, len(module.Functions), module.Description)

		printModuleInfo(module, maxModuleSize)

		fmt.Println()
	}

	fmt.Println("Use 'gru modules <module-name>' to view all functions of the specified module")
}

func printModuleInfo(module gru.GruModule, maxModuleSize int) {
	functionAmount := len(module.Functions)
	for i := range functionAmount {
		function := module.Functions[i]
		fmt.Printf("  %s()\n    %s\n", function.Name, function.Description)

		if i > 0 && i == maxModuleSize-1 {
			fmt.Printf("  %d functions not shown...\n", functionAmount-maxModuleSize)
			break
		}
	}
}

func helpCommand() {
	fmt.Println("Available commands:")
	fmt.Println("modules <module>  Describes all available modules or the specified module")
	fmt.Println("eval    <code>    Evalutes the passed Lua code")
	fmt.Println("help              This message")
	fmt.Println("\nIf the passed value is not found as a command, it will be treated as a lua file: 'gru main.lua'")
}
