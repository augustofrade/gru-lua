package cli

import (
	"fmt"

	"github.com/augustofrade/gru-lua/gru"
)

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
