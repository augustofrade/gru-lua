package luautil

import (
	"fmt"

	"github.com/Shopify/go-lua"
)

// Gets all string varargs for the called Lua Go function that are in the stack.
//
// Throws an error if any argument isn't of string type.
func GetStringVarargs(l *lua.State, varargAmount int) ([]string, error) {

	parts := make([]string, varargAmount)

	for i := 1; i <= varargAmount; i++ {
		if !IsString(l, i) {
			return nil, fmt.Errorf("Expected string in argument %d", i)
		}
		value, _ := l.ToString(i)
		parts[i-1] = value
	}

	return parts, nil
}

// Returns whether the value at stack index is strictly a string and not a string convertible number.
//
//   - "text" = true
//   - "1" = true
//   - 1 = false
func IsString(l *lua.State, index int) bool {
	return l.IsString(index) && !l.IsNumber(index)
}
