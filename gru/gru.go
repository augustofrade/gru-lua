package gru

import (
	"fmt"

	"github.com/Shopify/go-lua"
)

const runtimeCurrentVersion string = "0.0.1"

var _l *lua.State

var RegisteredModules = make([]GruModule, 0)

const errorCallbackRegistryKey = "gru.runtime.on_error"

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

// Transforms the map into a Lua table and pushes it onto the stack
func PushLuaTable(kvp map[string]any) {
	_l.CreateTable(0, len(kvp))
	for key, value := range kvp {
		switch v := value.(type) {
		case string:
			_l.PushString(v)
		case int:
			_l.PushInteger(v)
		case int64:
			_l.PushInteger(int(v))
		case float32:
			_l.PushNumber(float64(v))
		case float64:
			_l.PushNumber(v)
		case bool:
			_l.PushBoolean(v)
			// TODO: handle arrays/slices
		case map[string]any:
			PushLuaTable(v)
		default:
			_l.PushString(fmt.Sprint(v))
		}

		_l.SetField(-2, key)
	}
}
