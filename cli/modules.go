package cli

import (
	"fmt"
	"strings"

	"github.com/augustofrade/gru-lua/gru"
	"github.com/augustofrade/gru-lua/gru/definitions"
)

func modulesCommand(maxModuleSize int, specifiedModule string) {
	gru.InitDefaultModules()

	if specifiedModule == "" {
		printAllModulesInfo(maxModuleSize)
		return
	}

	var module *definitions.GruModule
	for _, m := range gru.RegisteredModules {
		if m.Name == specifiedModule {
			module = &m
		}
	}

	if module == nil {
		fmt.Println("Gru module not found. Expected one of the below:")
		for _, mod := range gru.RegisteredModules {
			fmt.Printf("  %s\n", mod.Name)
		}
		return
	}

	fmt.Printf("[%d function(s)]  %s\n", len(module.Functions), module.Description)
	printModuleInfo(*module, 0)
}

func printAllModulesInfo(maxModuleSize int) {
	for _, module := range gru.RegisteredModules {
		fmt.Printf("%s [%d]  %s\n", module.Name, len(module.Functions), module.Description)

		printModuleInfo(module, maxModuleSize)
	}

	fmt.Println("Use 'gru modules <module-name>' to view all functions of the specified module")
}

func printModuleInfo(module definitions.GruModule, maxModuleSize int) {
	functionAmount := len(module.Functions)

	builder := strings.Builder{}

	for i := range functionAmount {
		paramBuilder := strings.Builder{}
		function := module.Functions[i]
		fmt.Fprintf(&builder, "  %s(", function.Name)

		for _, param := range function.Parameters {
			fmt.Fprintf(&builder, "%s: %s", param.Name, param.Type)
			fmt.Fprintf(&paramBuilder, "    @param %s: %s   %s\n", param.Name, param.Type, param.Description)
		}

		fmt.Fprintf(&builder, ")\n    %s\n", function.Description)
		builder.WriteString(paramBuilder.String())

		remainingFuncs := functionAmount - maxModuleSize
		if remainingFuncs > 0 && i > 0 && i == maxModuleSize-1 {
			fmt.Fprintf(&builder, "%d functions not shown...\n", remainingFuncs)
			break
		}
	}

	fmt.Println(builder.String())
}
