package gru

import (
	"fmt"

	"github.com/Shopify/go-lua"
)

var _l *lua.State

type LuaInteropFunc func(l *lua.State) int

type GruModule struct {
	Name        string
	Description string
	Functions   []GruFunction
	l           *lua.State
}

type GruFunction struct {
	Name           string
	Description    string
	Implementation LuaInteropFunc
}

func HandleFile(file *string) {
	_l = lua.NewState()

	lua.OpenLibraries(_l)
	RegisterDefaultModules()

	if err := lua.DoFile(_l, *file); err != nil {
		fmt.Println(err)
	}
}

func RegisterDefaultModules() {
	// gru table
	_l.NewTable()

	_l.SetGlobal("gru")
}

func RegisterGruModule(module GruModule) {
	_l.NewTable()
	for _, function := range module.Functions {
		_l.PushGoFunction(lua.Function(function.Implementation))
		_l.SetField(-2, function.Name)
	}
	_l.SetField(-2, module.Name)
}

func NewModule(name string, description string) GruModule {
	return GruModule{
		Name:        name,
		Description: description,
		l:           _l,
		Functions:   make([]GruFunction, 0),
	}
}

func (module *GruModule) Register(funcName string, description string, function LuaInteropFunc) {
	module.Functions = append(module.Functions, GruFunction{
		Name:           funcName,
		Description:    description,
		Implementation: function,
	})
}
