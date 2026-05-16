package gru

import (
	"fmt"

	"github.com/Shopify/go-lua"
)

var _l *lua.State

var RegisteredModules = make([]GruModule, 0)

type LuaInteropFunc func(l *lua.State) int

// A module that is accessible in Lua with the default gru table: gru.<module-name>
type GruModule struct {
	Name        string
	Description string
	Functions   []GruFunction
}

// A callable function through gru.<module-name>.<function-name>
type GruFunction struct {
	Name           string
	Description    string
	Parameters     []GruFunctionParameter
	Implementation LuaInteropFunc
	returnTypes    []string
}

type GruFunctionParameter struct {
	Name        string
	Description string
	Type        string
	Vararg      bool
}

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
	RegisteredModules = append(RegisteredModules, NewColorsModule())
	RegisteredModules = append(RegisteredModules, NewJsonModule())
	RegisteredModules = append(RegisteredModules, NewTimeModule())
	RegisteredModules = append(RegisteredModules, NewPathModule())
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

// Registers a new Gru module alongside all its functions
func RegisterGruModule(module GruModule) {
	_l.NewTable()
	for _, function := range module.Functions {
		_l.PushGoFunction(lua.Function(function.Implementation))
		_l.SetField(-2, function.Name)
	}
	_l.SetField(-2, module.Name)
}

// Default GruModule factory
func NewModule(name string, description string) GruModule {
	return GruModule{
		Name:        name,
		Description: description,
		Functions:   make([]GruFunction, 0),
	}
}

// Registers a Go function in the Lua GruModule
func (module *GruModule) Register(funcName string, description string, function LuaInteropFunc) {
	module.RegisterGruFunction(GruFunction{
		Name:           funcName,
		Description:    description,
		Implementation: function,
	})
}

// Registers a built GruFunction in the Lua GruModule
func (module *GruModule) RegisterGruFunction(function GruFunction) {
	module.Functions = append(module.Functions, function)
}

// Creates a FunctionBuilder for GruFunctions
func (module *GruModule) FunctionBuilder(name string, description string, function LuaInteropFunc) *GruFunctionBuilder {
	return &GruFunctionBuilder{
		name:        name,
		description: description,
		function:    function,
		parameters:  make([]GruFunctionParameter, 0),
		module:      module,
	}
}

func LuaError(message string) int {
	_l.PushNil()
	_l.PushString(message)
	return 2
}

// TODO: remove this
func LuaVoidError(message string) int {
	_l.PushString(message)
	return 1
}

func LuaStringResult(value string) int {
	_l.PushString(value)
	return 1
}

func LuaBoolResult(value bool) int {
	_l.PushBoolean(value)
	return 1
}
