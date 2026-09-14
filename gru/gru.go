package gru

import (
	"fmt"

	"github.com/Shopify/go-lua"
	"github.com/augustofrade/gru-lua/gru/definitions"
)

var (
	CurrentBuildVersion = "unknown"
	CurrentBuildDate    = "N/A"
)

var _l *lua.State

var RegisteredModules = make([]definitions.GruModule, 0)

const errorCallbackRegistryKey = "gru.runtime.on_error"

// Inits the VM with all default Lua and Gru libraries
func initGru() {
	_l = lua.NewState()

	lua.OpenLibraries(_l)

	InitDefaultModules()
	RegisterDefaultModules()
}

// Handles a Lua file
func DoFile(file *string) {
	initGru()

	if err := lua.DoFile(_l, *file); err != nil {
		_l.Field(lua.RegistryIndex, errorCallbackRegistryKey)
		if _l.IsFunction(-1) {
			_l.PushString(err.Error())
			_l.ProtectedCall(1, 0, 0)
		} else {
			_l.Pop(1)
		}
		fmt.Println(err)
	}
}

// Evalutes the passed Lua code
func Evalute(code *string) {
	initGru()

	if err := lua.DoString(_l, *code); err != nil {
		fmt.Println(err)
	}
}

// Sets all default Go Gru modules into a module sice.
func InitDefaultModules() {
	RegisteredModules = []definitions.GruModule{
		NewAssertModule(),
		NewColorsModule(),
		NewJsonModule(),
		NewTimeModule(),
		NewPathModule(),
		NewRuntimeModule(),
		NewZipModule(),
		NewEnvModule(),
		NewFsModule(),
		NewHttpModule(),
	}
}

// Registers all default Gru modules into Lua tables accessed through the default "gru" global table.
func RegisterDefaultModules() {
	// gru table
	_l.NewTable()

	for _, module := range RegisteredModules {
		RegisterGruModule(module)
	}

	_l.SetGlobal("gru")
}

// Registers a new Gru module alongside all its functions as Lua tables and functions
func RegisterGruModule(module definitions.GruModule) {
	_l.NewTable()
	for _, function := range module.Functions {
		_l.PushGoFunction(lua.Function(function.Implementation))
		// sets the function (top of stack) as value of function.Name key
		// table[function.Name] = GoFunction
		//
		_l.SetField(-2, function.Name)
	}
	// all functions have been popped from the stack,
	// therefore the top of the stack is the module Lua table
	// gruTable[module.Name] = moduleTable
	_l.SetField(-2, module.Name)
}
