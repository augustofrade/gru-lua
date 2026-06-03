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
func (module *GruModule) registerGruFunction(function GruFunction) {
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

func (module *GruModule) HasCustomType(name string, description string) *GruModuleCustomType {
	newType := GruModuleCustomType{
		Name:        name,
		Description: description,
		Properties:  make(map[string]GruModuleCustomTypeProperty),
	}
	module.Types = append(module.Types, &newType)
	return &newType
}

func (cType *GruModuleCustomType) Prop(name string, propType string, description string) *GruModuleCustomType {
	cType.Properties[name] = GruModuleCustomTypeProperty{
		Description: description,
		Type:        propType,
	}
	return cType
}

func (cType *GruModuleCustomType) StringProp(name string, description string) *GruModuleCustomType {
	return cType.Prop(name, "string", description)
}

func (cType *GruModuleCustomType) NumberProp(name string, description string) *GruModuleCustomType {
	return cType.Prop(name, "number", description)
}

func (cType *GruModuleCustomType) BooleanProp(name string, description string) *GruModuleCustomType {
	return cType.Prop(name, "boolean", description)
}

func (cType *GruModuleCustomType) Method(name string, description string) *GruModuleCustomTypeMethod {
	cMethod := GruModuleCustomTypeMethod{
		Description: description,
	}

	cType.Methods[name] = &cMethod

	return &cMethod
}
