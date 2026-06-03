package definitions

import "github.com/Shopify/go-lua"

type LuaInteropFunc func(l *lua.State) int

// A module that is accessible in Lua with the default gru table: gru.<module-name>
type GruModule struct {
	Name        string
	Description string
	Functions   []GruFunction
	Types       []*GruModuleCustomType
	Alias       []*GruModuleAlias
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

type GruModuleAlias struct {
	Name        string
	Description string
	To          string
}

type GruModuleCustomType struct {
	Name        string
	Description string
	Properties  map[string]GruModuleCustomTypeProperty
	Methods     map[string]*GruModuleCustomTypeMethod
}

type GruModuleCustomTypeMethod struct {
	Description string
	Parameters  []GruFunctionParameter
	ReturnTypes []string
}

type GruModuleCustomTypeProperty struct {
	Description string
	Type        string
}
