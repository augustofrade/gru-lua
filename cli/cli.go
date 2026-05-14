package cli

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/augustofrade/gru-lua/gru"
)

func HandleArgs(args []string) {
	command := args[1]

	// TODO: handle params

	handleFile(&command)
}

func handleFile(file *string) {
	absFile, _ := filepath.Abs(*file)
	if _, err := os.Stat(absFile); err != nil {
		fmt.Println("File not found: " + absFile)
		return
	}

	gru.HandleFile(file)
}
