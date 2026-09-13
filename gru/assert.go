package gru

import (
	"fmt"
	"math"
	"strings"

	"github.com/Shopify/go-lua"
	"github.com/augustofrade/gru-lua/gru/definitions"
	"github.com/augustofrade/gru-lua/gru/internal/luautil"
)

func NewAssertModule() definitions.GruModule {
	module := definitions.NewModule("assert", "Assertion helpers that throw Lua errors when validations fail.")

	module.FunctionBuilder("type", "Asserts that a value matches the expected Lua type.", assertType).
		Param("value", "any", "Value to validate.").
		StringParam("expected_type", "Expected Lua type: string, number, boolean, table, function or nil.").
		OptionalStringParam("label", "Optional label used in the error message.").
		Returns("any").
		Register()

	module.FunctionBuilder("truthy", "Asserts that a value is truthy in Lua.", assertTruthy).
		Param("value", "any", "Value to validate.").
		OptionalStringParam("label", "Optional label used in the error message.").
		Returns("any").
		Register()

	module.FunctionBuilder("equals", "Asserts that two Lua values are equal.", assertEquals).
		Param("value", "any", "Value to validate.").
		Param("expected", "any", "Expected value.").
		OptionalStringParam("label", "Optional label used in the error message.").
		Returns("any").
		Register()

	module.FunctionBuilder("between", "Asserts that a number is between minimum and maximum, inclusive.", assertBetween).
		NumberParam("value", "Number to validate.").
		NumberParam("minimum", "Minimum accepted value.").
		NumberParam("maximum", "Maximum accepted value.").
		OptionalStringParam("label", "Optional label used in the error message.").
		ReturnsNumber().
		Register()

	module.FunctionBuilder("is_string", "Asserts that a value is a string.", assertIsString).
		Param("value", "any", "Value to validate.").
		OptionalStringParam("label", "Optional label used in the error message.").
		Returns("any").
		Register()

	module.FunctionBuilder("is_number", "Asserts that a value is a number.", assertIsNumber).
		Param("value", "any", "Value to validate.").
		OptionalStringParam("label", "Optional label used in the error message.").
		ReturnsNumber().
		Register()

	module.FunctionBuilder("is_boolean", "Asserts that a value is a boolean.", assertIsBoolean).
		Param("value", "any", "Value to validate.").
		OptionalStringParam("label", "Optional label used in the error message.").
		Returns("any").
		Register()

	module.FunctionBuilder("is_table", "Asserts that a value is a table.", assertIsTable).
		Param("value", "any", "Value to validate.").
		OptionalStringParam("label", "Optional label used in the error message.").
		Returns("any").
		Register()

	module.FunctionBuilder("is_array", "Asserts that a value is an array-like Lua table.", assertIsArray).
		Param("value", "any", "Value to validate.").
		OptionalStringParam("label", "Optional label used in the error message.").
		Returns("any").
		Register()

	module.FunctionBuilder("is_integer", "Asserts that a value is an integer number.", assertIsInteger).
		Param("value", "any", "Value to validate.").
		OptionalStringParam("label", "Optional label used in the error message.").
		ReturnsNumber().
		Register()

	module.FunctionBuilder("is_function", "Asserts that a value is a function.", assertIsFunction).
		Param("value", "any", "Value to validate.").
		OptionalStringParam("label", "Optional label used in the error message.").
		Returns("any").
		Register()

	module.FunctionBuilder("greater_than", "Asserts that a number is greater than another number.", assertGreaterThan).
		NumberParam("value", "Number to validate.").
		NumberParam("minimum", "Comparison value.").
		OptionalStringParam("label", "Optional label used in the error message.").
		ReturnsNumber().
		Register()

	module.FunctionBuilder("lesser_than", "Asserts that a number is less than another number.", assertLesserThan).
		NumberParam("value", "Number to validate.").
		NumberParam("maximum", "Comparison value.").
		OptionalStringParam("label", "Optional label used in the error message.").
		ReturnsNumber().
		Register()

	module.FunctionBuilder("not_nil", "Asserts that a value is not nil.", assertNotNil).
		Param("value", "any", "Value to validate.").
		OptionalStringParam("label", "Optional label used in the error message.").
		Returns("any").
		Register()

	module.FunctionBuilder("not_empty", "Asserts that a string or table is not empty.", assertNotEmpty).
		Param("value", "any", "String or table to validate.").
		OptionalStringParam("label", "Optional label used in the error message.").
		Returns("any").
		Register()

	module.FunctionBuilder("has_keys", "Asserts that a table contains all provided string keys.", assertHasKeys).
		TableParam("value", "Table to validate.").
		TableParam("keys", "Array table of required keys.").
		OptionalStringParam("label", "Optional label used in the error message.").
		Returns("any").
		Register()

	return module
}

func assertType(l *lua.State) int {
	if !luautil.IsString(l, 2) {
		return luautil.PushError(l, "Expected string on 'expected_type' parameter")
	}

	expected, _ := l.ToString(2)
	label, err := getOptionalAssertLabel(l, 3)
	if err != nil {
		return luautil.PushError(l, err.Error())
	}

	return assertLuaType(l, expected, label)
}

func assertTruthy(l *lua.State) int {
	label, err := getOptionalAssertLabel(l, 2)
	if err != nil {
		return luautil.PushError(l, err.Error())
	}

	if !l.ToBoolean(1) {
		return luautil.PushError(l, fmt.Sprintf("Assertion failed: expected %s to be truthy", assertionSubject(label)))
	}

	return returnAssertedValue(l)
}

func assertEquals(l *lua.State) int {
	label, err := getOptionalAssertLabel(l, 3)
	if err != nil {
		return luautil.PushError(l, err.Error())
	}

	if !l.Compare(1, 2, lua.OpEq) {
		return luautil.PushError(l, fmt.Sprintf("Assertion failed: expected %s to equal the expected value", assertionSubject(label)))
	}

	return returnAssertedValue(l)
}

func assertBetween(l *lua.State) int {
	label, err := getOptionalAssertLabel(l, 4)
	if err != nil {
		return luautil.PushError(l, err.Error())
	}

	value, err := getAssertNumber(l, 1, "value")
	if err != nil {
		return luautil.PushError(l, err.Error())
	}

	minimum, err := getAssertNumber(l, 2, "minimum")
	if err != nil {
		return luautil.PushError(l, err.Error())
	}

	maximum, err := getAssertNumber(l, 3, "maximum")
	if err != nil {
		return luautil.PushError(l, err.Error())
	}

	if minimum > maximum {
		return luautil.PushError(l, "Expected 'minimum' parameter to be less than or equal to 'maximum'")
	}

	if value < minimum || value > maximum {
		return luautil.PushError(l, fmt.Sprintf("Assertion failed: expected %s to be between %v and %v", assertionSubject(label), minimum, maximum))
	}

	return returnAssertedValue(l)
}

func assertIsString(l *lua.State) int {
	return assertFixedLuaType(l, "string")
}

func assertIsNumber(l *lua.State) int {
	return assertFixedLuaType(l, "number")
}

func assertIsBoolean(l *lua.State) int {
	return assertFixedLuaType(l, "boolean")
}

func assertIsTable(l *lua.State) int {
	return assertFixedLuaType(l, "table")
}

func assertIsArray(l *lua.State) int {
	label, err := getOptionalAssertLabel(l, 2)
	if err != nil {
		return luautil.PushError(l, err.Error())
	}

	if l.TypeOf(1) != lua.TypeTable {
		return luautil.PushError(l, fmt.Sprintf("Assertion failed: expected %s to be array-like, found %s", assertionSubject(label), fmt.Sprint(l.TypeOf(1))))
	}

	if luautil.GetTableLength(l, 1) == 0 || luautil.IsArrayTable(l, 1) {
		return returnAssertedValue(l)
	}

	return luautil.PushError(l, fmt.Sprintf("Assertion failed: expected %s to be array-like", assertionSubject(label)))
}

func assertIsInteger(l *lua.State) int {
	label, err := getOptionalAssertLabel(l, 2)
	if err != nil {
		return luautil.PushError(l, err.Error())
	}

	value, err := getAssertNumber(l, 1, "value")
	if err != nil {
		return luautil.PushError(l, err.Error())
	}

	if math.Trunc(value) != value {
		return luautil.PushError(l, fmt.Sprintf("Assertion failed: expected %s to be an integer, found %v", assertionSubject(label), value))
	}

	return returnAssertedValue(l)
}

func assertIsFunction(l *lua.State) int {
	return assertFixedLuaType(l, "function")
}

func assertGreaterThan(l *lua.State) int {
	label, err := getOptionalAssertLabel(l, 3)
	if err != nil {
		return luautil.PushError(l, err.Error())
	}

	value, err := getAssertNumber(l, 1, "value")
	if err != nil {
		return luautil.PushError(l, err.Error())
	}

	minimum, err := getAssertNumber(l, 2, "minimum")
	if err != nil {
		return luautil.PushError(l, err.Error())
	}

	if value <= minimum {
		return luautil.PushError(l, fmt.Sprintf("Assertion failed: expected %s to be greater than %v", assertionSubject(label), minimum))
	}

	return returnAssertedValue(l)
}

func assertLesserThan(l *lua.State) int {
	label, err := getOptionalAssertLabel(l, 3)
	if err != nil {
		return luautil.PushError(l, err.Error())
	}

	value, err := getAssertNumber(l, 1, "value")
	if err != nil {
		return luautil.PushError(l, err.Error())
	}

	maximum, err := getAssertNumber(l, 2, "maximum")
	if err != nil {
		return luautil.PushError(l, err.Error())
	}

	if value >= maximum {
		return luautil.PushError(l, fmt.Sprintf("Assertion failed: expected %s to be less than %v", assertionSubject(label), maximum))
	}

	return returnAssertedValue(l)
}

func assertNotNil(l *lua.State) int {
	label, err := getOptionalAssertLabel(l, 2)
	if err != nil {
		return luautil.PushError(l, err.Error())
	}

	if l.IsNoneOrNil(1) {
		return luautil.PushError(l, fmt.Sprintf("Assertion failed: expected %s to not be nil", assertionSubject(label)))
	}

	return returnAssertedValue(l)
}

func assertNotEmpty(l *lua.State) int {
	label, err := getOptionalAssertLabel(l, 2)
	if err != nil {
		return luautil.PushError(l, err.Error())
	}

	switch l.TypeOf(1) {
	case lua.TypeString:
		value, _ := l.ToString(1)
		if value == "" {
			return luautil.PushError(l, fmt.Sprintf("Assertion failed: expected %s to not be empty", assertionSubject(label)))
		}
	case lua.TypeTable:
		if luautil.IsTableEmpty(l, 1) {
			return luautil.PushError(l, fmt.Sprintf("Assertion failed: expected %s to not be empty", assertionSubject(label)))
		}
	default:
		return luautil.PushError(l, fmt.Sprintf("Assertion failed: expected %s to be a string or table, found %s", assertionSubject(label), fmt.Sprint(l.TypeOf(1))))
	}

	return returnAssertedValue(l)
}

func assertHasKeys(l *lua.State) int {
	label, err := getOptionalAssertLabel(l, 3)
	if err != nil {
		return luautil.PushError(l, err.Error())
	}

	if l.TypeOf(1) != lua.TypeTable {
		return luautil.PushError(l, "Expected table on 'value' parameter")
	}

	keys, err := getAssertStringArray(l, 2, "keys")
	if err != nil {
		return luautil.PushError(l, err.Error())
	}

	for _, key := range keys {
		if !luautil.TableHasKey(l, 1, key) {
			return luautil.PushError(l, fmt.Sprintf("Assertion failed: expected %s to contain key '%s'", assertionSubject(label), key))
		}
	}

	return returnAssertedValue(l)
}

func assertFixedLuaType(l *lua.State, expected string) int {
	label, err := getOptionalAssertLabel(l, 2)
	if err != nil {
		return luautil.PushError(l, err.Error())
	}

	return assertLuaType(l, expected, label)
}

func assertLuaType(l *lua.State, expected string, label string) int {
	normalizedExpected, err := normalizeExpectedLuaType(expected)
	if err != nil {
		return luautil.PushError(l, err.Error())
	}

	if matchesExpectedLuaType(l, 1, normalizedExpected) {
		return returnAssertedValue(l)
	}

	return luautil.PushError(l, fmt.Sprintf(
		"Assertion failed: expected %s to be %s, found %s",
		assertionSubject(label),
		normalizedExpected,
		fmt.Sprint(l.TypeOf(1)),
	))
}

func getOptionalAssertLabel(l *lua.State, index int) (string, error) {
	if l.IsNoneOrNil(index) {
		return "", nil
	}

	if !luautil.IsString(l, index) {
		return "", fmt.Errorf("Expected string on 'label' parameter")
	}

	label, _ := l.ToString(index)
	return label, nil
}

func normalizeExpectedLuaType(expected string) (string, error) {
	switch strings.ToLower(strings.TrimSpace(expected)) {
	case "string":
		return "string", nil
	case "number":
		return "number", nil
	case "bool", "boolean":
		return "boolean", nil
	case "table":
		return "table", nil
	case "func", "function":
		return "function", nil
	case "nil":
		return "nil", nil
	default:
		return "", fmt.Errorf("Unsupported Lua type '%s' on 'expected_type' parameter", expected)
	}
}

func matchesExpectedLuaType(l *lua.State, index int, expected string) bool {
	switch expected {
	case "string":
		return luautil.IsString(l, index)
	case "number":
		return l.TypeOf(index) == lua.TypeNumber
	case "boolean":
		return l.TypeOf(index) == lua.TypeBoolean
	case "table":
		return l.TypeOf(index) == lua.TypeTable
	case "function":
		return l.IsFunction(index)
	case "nil":
		return l.IsNoneOrNil(index)
	default:
		return false
	}
}

func getAssertNumber(l *lua.State, index int, parameter string) (float64, error) {
	if l.TypeOf(index) != lua.TypeNumber {
		return 0, fmt.Errorf("Expected number on '%s' parameter", parameter)
	}

	value, _ := l.ToNumber(index)
	return value, nil
}

func getAssertStringArray(l *lua.State, index int, parameter string) ([]string, error) {
	if l.TypeOf(index) != lua.TypeTable {
		return nil, fmt.Errorf("Expected table on '%s' parameter", parameter)
	}

	if !luautil.IsArrayTable(l, index) {
		return nil, fmt.Errorf("Expected '%s' parameter to be an array table of strings", parameter)
	}

	absIndex := l.AbsIndex(index)
	length := luautil.GetTableLength(l, absIndex)
	values := make([]string, 0, length)

	for i := 1; i <= length; i++ {
		l.RawGetInt(absIndex, i)
		if !luautil.IsString(l, -1) {
			l.Pop(1)
			return nil, fmt.Errorf("Expected '%s' parameter to be an array table of strings", parameter)
		}

		value, _ := l.ToString(-1)
		values = append(values, value)
		l.Pop(1)
	}

	return values, nil
}

func returnAssertedValue(l *lua.State) int {
	// copies the value of the 1st param to the top of the stack
	l.PushValue(1)
	return 1
}

func assertionSubject(label string) string {
	if label == "" {
		return "value"
	}

	return fmt.Sprintf("'%s'", label)
}
