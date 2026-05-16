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
	case "modules", "mod":
		modulesCommand(5, subCommand)
	case "eval":
		gru.Evalute(&subCommand)
	case "generate-types", "types":
		generateTypesCommand(&subCommand)
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

func helpCommand() {
	fmt.Println("Available commands:")
	fmt.Println("modules <module>  Describes all available modules or the specified module")
	fmt.Println("eval    <code>    Evalutes the passed Lua code")
	fmt.Println("help              This message")
	fmt.Println("\nIf the passed value is not found as a command, it will be treated as a lua file: 'gru main.lua'")
}
