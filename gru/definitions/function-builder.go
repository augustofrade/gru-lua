package definitions

type GruFunctionBuilder struct {
	name        string
	description string
	function    LuaInteropFunc
	parameters  []GruFunctionParameter
	returnTypes []string
	module      *GruModule
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

// Builds and registers the function
func (f *GruFunctionBuilder) Register() {
	f.module.RegisterGruFunction(
		GruFunction{
			Name:           f.name,
			Description:    f.description,
			Parameters:     f.parameters,
			ReturnTypes:    f.returnTypes,
			Implementation: f.function,
		})
}

func (f *GruFunctionBuilder) Returns(returnType string) *GruFunctionBuilder {
	f.returnTypes = append(f.returnTypes, returnType)
	return f
}

func (f *GruFunctionBuilder) ReturnsNil() *GruFunctionBuilder {
	f.returnTypes = append(f.returnTypes, "nil")
	return f
}

func (f *GruFunctionBuilder) ReturnsString() *GruFunctionBuilder {
	f.returnTypes = append(f.returnTypes, "string")
	return f
}

func (f *GruFunctionBuilder) ReturnsBoolean() *GruFunctionBuilder {
	f.returnTypes = append(f.returnTypes, "boolean")
	return f
}

func (f *GruFunctionBuilder) ReturnsBooleanWithError() *GruFunctionBuilder {
	f.returnTypes = append(f.returnTypes, "boolean", "GruError")
	return f
}

func (f *GruFunctionBuilder) ReturnsNumber() *GruFunctionBuilder {
	f.returnTypes = append(f.returnTypes, "number")
	return f
}

func (f *GruFunctionBuilder) ReturnsError() *GruFunctionBuilder {
	f.returnTypes = append(f.returnTypes, "GruError")
	return f
}

func (f *GruFunctionBuilder) ReturnsStringWithError() *GruFunctionBuilder {
	f.returnTypes = append(f.returnTypes, "string", "GruError")
	return f
}

func (f *GruFunctionBuilder) ReturnsWithError(returnType string) *GruFunctionBuilder {
	f.returnTypes = append(f.returnTypes, returnType, "GruError")
	return f
}

func (f *GruFunctionBuilder) Vararg(paramType string) *GruFunctionBuilder {
	return f.Param("...", paramType, "")
}

func (f *GruFunctionBuilder) StringParam(name string, description string) *GruFunctionBuilder {
	return f.Param(name, "string", description)
}

func (f *GruFunctionBuilder) OptionalStringParam(name string, description string) *GruFunctionBuilder {
	return f.Param(name, "string?", description)
}

func (f *GruFunctionBuilder) NumberParam(name string, description string) *GruFunctionBuilder {
	return f.Param(name, "number", description)
}

func (f *GruFunctionBuilder) OptionalNumberParam(name string, description string) *GruFunctionBuilder {
	return f.Param(name, "number?", description)
}

func (f *GruFunctionBuilder) TableParam(name string, description string) *GruFunctionBuilder {
	return f.Param(name, "table", description)
}

func (f *GruFunctionBuilder) OptionalTableParam(name string, description string) *GruFunctionBuilder {
	return f.Param(name, "table?", description)
}

func (f *GruFunctionBuilder) BooleanParam(name string, description string) *GruFunctionBuilder {
	return f.Param(name, "boolean", description)
}

func (f *GruFunctionBuilder) Param(name string, paramType string, description string) *GruFunctionBuilder {
	f.parameters = append(f.parameters, GruFunctionParameter{
		Name:        name,
		Description: description,
		Type:        paramType,
	})

	return f
}
