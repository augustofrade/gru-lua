package gru

import (
	"fmt"

	"github.com/Shopify/go-lua"
)

const runtimeCurrentVersion string = "0.0.2"

var _l *lua.State

var RegisteredModules = make([]GruModule, 0)

const errorCallbackRegistryKey = "gru.runtime.on_error"

type LuaInteropFunc func(l *lua.State) int

// A module that is accessible in Lua with the default gru table: gru.<module-name>
type GruModule struct {
	Name        string
	Description string
	Functions   []GruFunction
	Types       []*GruModuleType
}

type GruModuleType struct {
	Name        string
	Description string
	Properties  map[string]GruModuleTypeProperty
}

type GruModuleTypeProperty struct {
	Description string
	Type        string
}

// A callable function through gru.<module-name>.<function-name>
type GruFunction struct {
	Name           string
	Description    string
	Parameters     []GruFunctionParameter
	Implementation LuaInteropFunc
	ReturnTypes    []string
}

type GruFunctionParameter struct {
	Name        string
	Description string
	Type        string
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
	RegisteredModules = append(RegisteredModules, NewColorsModule())
	RegisteredModules = append(RegisteredModules, NewJsonModule())
	RegisteredModules = append(RegisteredModules, NewTimeModule())
	RegisteredModules = append(RegisteredModules, NewPathModule())
	RegisteredModules = append(RegisteredModules, NewRuntimeModule())
	RegisteredModules = append(RegisteredModules, NewZipModule())
	RegisteredModules = append(RegisteredModules, NewEnvModule())
	RegisteredModules = append(RegisteredModules, NewFsModule())
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
func RegisterGruModule(module GruModule) {
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

// Default GruModule factory
func NewModule(name string, description string) GruModule {
	return GruModule{
		Name:        name,
		Description: description,
		Functions:   make([]GruFunction, 0),
	}
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

func (module *GruModule) HasCustomType(name string, description string) *GruModuleType {
	newType := GruModuleType{
		Name:        name,
		Description: description,
		Properties:  make(map[string]GruModuleTypeProperty),
	}
	module.Types = append(module.Types, &newType)
	return &newType
}

func (cType *GruModuleType) Prop(name string, propType string, description string) *GruModuleType {
	cType.Properties[name] = GruModuleTypeProperty{
		Description: description,
		Type:        propType,
	}
	return cType
}

func (cType *GruModuleType) StringProp(name string, description string) *GruModuleType {
	return cType.Prop(name, "string", description)
}

func (cType *GruModuleType) NumberProp(name string, description string) *GruModuleType {
	return cType.Prop(name, "number", description)
}

func (cType *GruModuleType) BooleanProp(name string, description string) *GruModuleType {
	return cType.Prop(name, "boolean", description)
}
