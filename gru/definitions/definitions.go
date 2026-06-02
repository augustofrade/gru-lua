package definitions

import "github.com/Shopify/go-lua"

type LuaInteropFunc func(l *lua.State) int

// A module that is accessible in Lua with the default gru table: gru.<module-name>
type GruModule struct {
	Name        string
	Description string
	Functions   []GruFunction
	Types       []*GruModuleType
	Alias       []*GruModuleAlias
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

type GruModuleTypeMethod struct {
	Description string
	Parameters  []GruFunctionParameter
	ReturnTypes []string
}

type GruModuleTypeMethodParam struct {
}

type GruModuleAlias struct {
	Name        string
	Description string
	To          string
}
