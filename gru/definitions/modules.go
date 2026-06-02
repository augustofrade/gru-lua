package definitions

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

func (module *GruModule) HasCustomAlias(name string, description string, aliasTo string) *GruModuleAlias {
	newAlias := GruModuleAlias{
		Name:        name,
		Description: description,
		To:          aliasTo,
	}
	module.Alias = append(module.Alias, &newAlias)
	return &newAlias
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

func (cType *GruModuleType) Method(name string, description string, returns ...string) *GruModuleType {
	cType.Properties[name] = GruModuleTypeProperty{
		Description: description,
		Type:        "fun():",
	}
	return cType
}
